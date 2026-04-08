package shortUrl

import (
	"errors"
	"math/rand"
	"time"

	"short-url/internal/user"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var rng *rand.Rand

func init() {
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func generateCode(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rng.Intn(len(charset))]
	}
	return string(b)
}

type Service interface {
	CreateShortURL(original string, customAlias string, userID string) (*ShortURL, error)
	GetShortURL(code string) (*ShortURL, error)
}

type shortUrlService struct {
	repo    Repository
	userSvc user.Service
}

func NewService(repo Repository, userSvc user.Service) Service {
	return &shortUrlService{repo: repo, userSvc: userSvc}
}

func (s *shortUrlService) CreateShortURL(original string, customAlias string, userID string) (*ShortURL, error) {
	// Check user quota
	count, err := s.userSvc.GetRemainingQuota(userID)
	if err != nil {
		return nil, errors.New("could not retrieve user count")
	}
	if count <= 0 {
		return nil, errors.New("url creation limit reached")
	}

	var code string
	if customAlias != "" {
		if len(customAlias) < 3 || len(customAlias) > 20 {
			return nil, errors.New("custom alias must be between 3 and 20 characters")
		}
		// Check if alias already exists
		_, err := s.repo.FindByCode(customAlias)
		if err == nil {
			return nil, errors.New("custom alias already in use")
		}
		code = customAlias
	} else {
		// Generate random code and ensure uniqueness
		for {
			code = generateCode(6)
			_, err := s.repo.FindByCode(code)
			if err != nil {
				// We expect an error (not found)
				break
			}
		}
	}

	short := &ShortURL{
		Original:  original,
		ShortCode: code,
		CreatedBy: userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.CreateShortURL(short); err != nil {
		return nil, errors.New("could not save short url")
	}

	// Decrement quota
	_ = s.userSvc.DecrementQuota(userID)

	return short, nil
}

func (s *shortUrlService) GetShortURL(code string) (*ShortURL, error) {
	return s.repo.FindByCode(code)
}
