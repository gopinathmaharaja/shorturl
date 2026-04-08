package user

import (
	"errors"
	"log"
	"net/mail"
	"os"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(u *User) error
	Login(email, password string) (string, *User, error)
	DecrementQuota(userID string) error
	GetRemainingQuota(userID string) (int, error)
}

type userService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &userService{repo: repo}
}

func (s *userService) Register(u *User) error {
	if !ValidateEmail(u.Email) {
		return errors.New("invalid email format")
	}

	if valid, msg := ValidatePassword(u.Password); !valid {
		return errors.New(msg)
	}

	_, err := s.repo.FindOne(bson.M{"email": u.Email})
	if err == nil {
		return errors.New("email already registered")
	}
	if err != mongo.ErrNoDocuments {
		return errors.New("database error")
	}

	hash, err := HashPassword(u.Password)
	if err != nil {
		return errors.New("error processing password")
	}

	u.Password = hash
	u.TotalCount = 10
	u.RemainingCount = 10
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

	return s.repo.CreateUser(u)
}

func (s *userService) Login(email, password string) (string, *User, error) {
	if email == "" || password == "" {
		return "", nil, errors.New("email and password required")
	}

	u, err := s.repo.FindOne(bson.M{"email": email})
	if err == mongo.ErrNoDocuments {
		return "", nil, errors.New("invalid credentials")
	}
	if err != nil {
		return "", nil, errors.New("database error")
	}

	if !CheckPassword(u.Password, password) {
		return "", nil, errors.New("invalid credentials")
	}

	token, err := GenerateToken(u.ID)
	if err != nil {
		return "", nil, errors.New("could not generate token")
	}

	return token, u, nil
}

func (s *userService) DecrementQuota(userID string) error {
	log.Printf("[USER-SERVICE] Decrementing quota for user: %s", userID)
    
    // In mongodb string ID needs to be carefully handled.
    // Assuming the repo implementation handles the string vs ObjectID logic if needed.
    // Our existing users have string ID in the Handler decrement logic. Wait, let's look back at shortUrl handler decrement:
    // It used primitive.ObjectIDFromHex(userID) for GetUserShortURLCount and raw string `userID` for Decrement.
    // That means it was inconsistent. Let's fix this here by always converting it or using the ID as stored.
    
    // For safety, let's try ObjectID first, if it fails, fallback to string.
    var idFilter interface{} = userID
    if objID, err := primitive.ObjectIDFromHex(userID); err == nil {
        idFilter = objID
    }

	_, err := s.repo.UpdateOne(
		bson.M{"_id": idFilter},
		bson.M{
			"$inc": bson.M{"remaining_count": -1},
			"$set": bson.M{"updated_at": time.Now()},
		},
	)
	return err
}

func (s *userService) GetRemainingQuota(userID string) (int, error) {
    var idFilter interface{} = userID
    if objID, err := primitive.ObjectIDFromHex(userID); err == nil {
        idFilter = objID
    }

	u, err := s.repo.FindOne(bson.M{"_id": idFilter})
	if err != nil {
		return 0, err
	}
	return u.RemainingCount, nil
}

func HashPassword(pw string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pw), 14)
	return string(bytes), err
}

func CheckPassword(hash, pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw)) == nil
}

func GenerateToken(id string) (string, error) {
	claims := jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

// ValidateEmail checks if email is valid
func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	valid := err == nil
	if !valid {
		log.Printf("[USER-VALIDATION] Invalid email format: %s", email)
	}
	return valid
}

// ValidatePassword checks password strength
func ValidatePassword(password string) (bool, string) {
	if len(password) < 8 {
		log.Printf("[USER-VALIDATION] Password too short (length: %d)", len(password))
		return false, "Password must be at least 8 characters"
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

	if !hasUpper || !hasLower || !hasNumber {
		log.Printf("[USER-VALIDATION] Password missing required complexity. Upper: %v, Lower: %v, Number: %v", hasUpper, hasLower, hasNumber)
		return false, "Password must contain uppercase, lowercase, and numbers"
	}

	return true, ""
}
