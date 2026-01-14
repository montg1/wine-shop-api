# Wine Shop API 🍷

A RESTful backend API for a Wine Shop, built with Go, Gin, GORM, and PostgreSQL.

## 🚀 Tech Stack
- **Language**: Go (Golang)
- **Framework**: Gin Web Framework
- **Database**: PostgreSQL
- **ORM**: GORM
- **Authentication**: JWT & BCrypt
- **Containerization**: Docker

## 🛠️ Setup & Installation

### Prerequisites
- Go 1.20+
- Docker & Docker Compose

### 1. Clone the repository
```bash
git clone https://github.com/montg1/wine-shop-api.git
cd wine-shop-api
```

### 2. Environment Variables
The project comes with a default `.env` file for development.
```env
DB_HOST=localhost
DB_PORT=5433
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=wine_shop
API_SECRET=mysecretkey
TOKEN_HOUR_LIFESPAN=24
```

### 3. Start Database (Docker)
This project uses port **5433** for PostgreSQL to avoid conflicts with local instances.
```bash
docker compose up -d
```

### 4. Run the API
```bash
go mod tidy
go run cmd/server/main.go
```
The server will start on `http://localhost:8080`.

## 📚 API Endpoints

### Public
- `GET /api/health` - Check API status
- `POST /api/register` - Create a new user account
  - Body: `{"email": "user@example.com", "password": "password"}`
- `POST /api/login` - Login and receive JWT
  - Body: `{"email": "user@example.com", "password": "password"}`
- `GET /api/products` - List all wines
- `GET /api/products/:id` - Get wine details

### Protected (Requires Bearer Token)
**Admin**
- `GET /api/admin/profile` - Verify token validity
- `POST /api/admin/products` - Create new wine
- `PUT /api/admin/products/:id` - Update wine details
- `DELETE /api/admin/products/:id` - Delete wine

**User (Shopping)**
- `GET /api/cart` - View current cart
- `POST /api/cart` - Add item to cart
  - Body: `{"product_id": 1, "quantity": 2}`
- `POST /api/orders` - Checkout
- `GET /api/orders` - View order history

## 🗂️ Project Structure
```text
wine-shop-api/
├── cmd/server/      # Entry point
├── internal/
│   ├── domain/      # Data models (User, Product, Order)
│   ├── handler/     # HTTP Controllers
│   ├── middleware/  # Auth Middleware
│   ├── service/     # Business Logic
├── pkg/
│   ├── config/      # DB Connection
│   ├── utils/       # JWT Helper
```
