# 🚀 URL Shortener API using Go Fiber, MongoDB & JWT

This is a simple and secure URL shortener service built with [Fiber](https://gofiber.io/), [MongoDB](https://www.mongodb.com/), and JWT authentication. Users can register, log in, create short URLs, and access redirection through the generated codes.

---

## 📁 Project Structure

short-url/
├── cmd/
│ └── main.go # Entry point
├── internal/
│ ├── user/ # User model, auth handlers
│ ├── shorturl/ # Short URL model and logic
│ ├── middleware/ # JWT middleware
├── pkg/
│ └── db/
│ └── mongo.go # MongoDB connection setup
├── routes/
│ └── router.go # Route definitions
├── go.mod / go.sum
├── .env

---

## 🔧 Prerequisites

- Go 1.20+
- MongoDB instance running (local or cloud)
- Git

---

## ⚙️ Setup Instructions

1. **Clone the repository:**
   ```bash
   git clone https://github.com/yourusername/short-url.git
   cd short-url
Create a .env file:

MONGO_URI=mongodb://localhost:27017
DB_NAME=shorturl
JWT_SECRET=your_super_secret_key
Install dependencies:

go mod tidy
Run the project:

go run cmd/main.go
Your server will start on http://localhost:3000.

📡 API Endpoints
🔐 Auth Routes
Method	Endpoint	Description
POST	/api/auth/register	Register a new user
POST	/api/auth/login	Log in and get token

📝 Example Login Response

{
  "token": "eyJhbGciOiJIUzI1NiIsInR..."
}
🔗 URL Routes (Protected)
Method	Endpoint	Description
POST	/api/url/create	Create short URL
GET	/:code	Redirect to original

Note: Add Authorization: Bearer <token> in headers for protected routes.

📦 Sample Request: Create Short URL
Endpoint: POST /api/url/create
Headers: Authorization: Bearer <token>
Body:

{
  "original": "https://example.com/very/long/url"
}
Response:

{
  "id": "64b123...",
  "original": "https://example.com/very/long/url",
  "short_code": "abc123",
  "created_by": "userId"
}
🌐 Redirect
Once created, access the short URL like:

http://localhost:3000/abc123
It will redirect to the original URL.

🧪 Testing
Use Postman or curl to test all the endpoints. Remember to register and login to get your token.

🧰 Built With
Fiber

MongoDB Go Driver

JWT

📝 License
MIT — feel free to use and modify!

---

Let me know if you want this included automatically in your project folder once the ZIP generation feature is back.

