
# 🚀 URL Shortener API using Go Fiber, MongoDB & JWT

This is a simple and secure URL shortener service built with [Fiber](https://gofiber.io/), [MongoDB](https://www.mongodb.com/), and JWT authentication. Users can register, log in, create short URLs, and access redirection through the generated codes.

---

## 📁 Project Structure

```
short-url/
├── cmd/
│   └── main.go              # Entry point
├── internal/
│   ├── user/                # User model, auth handlers
│   ├── shorturl/            # Short URL model and logic
│   ├── middleware/          # JWT middleware
├── pkg/
│   └── db/
│       └── mongo.go         # MongoDB connection setup
├── routes/
│   └── router.go            # Route definitions
├── go.mod / go.sum
├── .env
```

---

## 🔧 Prerequisites

- Go 1.20+
- MongoDB instance running (local or cloud)
- Git

---

## ⚙️ Setup Instructions

1. **Clone the repository:**
   ```bash
   git clone https://github.com/gopinathmaharaja/shorturl.git
   cd short-url
   ```

2. **Create a `.env` file:**
   ```env
   MONGO_URI=mongodb://localhost:27017
   DB_NAME=shorturl
   JWT_SECRET=your_super_secret_key
   ```

3. **Install dependencies:**
   ```bash
   go mod tidy
   ```

4. **Run the project:**
   ```bash
   go run cmd/main.go
   ```

   Your server will start on `http://localhost:3000`.

---

## 📡 API Endpoints

### 🔐 Auth Routes

| Method | Endpoint            | Description          |
|--------|---------------------|----------------------|
| POST   | `/api/auth/register` | Register a new user  |
| POST   | `/api/auth/login`    | Log in and get token |

#### 📝 Example Login Response
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR..."
}
```

---

### 🔗 URL Routes (Protected)

| Method | Endpoint             | Description            |
|--------|----------------------|------------------------|
| POST   | `/api/url/create`    | Create short URL       |
| GET    | `/:code`             | Redirect to original   |

> **Note:** Add `Authorization: Bearer <token>` in headers for protected routes.

---

## 📦 Sample Request: Create Short URL

**Endpoint:** `POST /api/url/create`  
**Headers:** `Authorization: Bearer <token>`  
**Body:**
```json
{
  "original": "https://example.com/very/long/url"
}
```

**Response:**
```json
{
  "id": "64b123...",
  "original": "https://example.com/very/long/url",
  "short_code": "abc123",
  "created_by": "userId"
}
```

---

## 🌐 Redirect

Once created, access the short URL like:
```
http://localhost:3000/abc123
```
It will redirect to the original URL.

---

## 🧪 Testing

Use [Postman](https://www.postman.com/) or curl to test all the endpoints. Remember to register and login to get your token.

---

## 🧰 Built With

- [Fiber](https://gofiber.io/)
- [MongoDB Go Driver](https://go.mongodb.org/mongo-driver)
- [JWT](https://pkg.go.dev/github.com/golang-jwt/jwt/v4)

---

## 📝 License

MIT — feel free to use and modify!
