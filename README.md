# Sistem Internal Backend

A Go backend API with MySQL database integration using Gin framework and GORM.

## 🗄️ Database Setup

### 1. Install MySQL
- Download and install MySQL from [mysql.com](https://dev.mysql.com/downloads/)
- Or use XAMPP/WAMP for easier setup

### 2. Create Database
```sql
CREATE DATABASE sistem_internal;
```

### 3. Configure Environment Variables
Create a `.env` file in the backend directory:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password_here
DB_NAME=sistem_internal

# Server Configuration
GIN_MODE=debug
```

### 4. Update Database Configuration
If you don't want to use environment variables, you can directly update the database configuration in `database/database.go`:

```go
dbHost := "localhost"
dbPort := "3306"
dbUser := "root"
dbPassword := "your_password_here"
dbName := "sistem_internal"
```

## 🚀 Running the Application

1. **Install dependencies:**
   ```bash
   go mod tidy
   ```

2. **Run the server:**
   ```bash
   go run main.go
   ```

3. **The server will start on:** `http://localhost:8080`

## 📋 API Endpoints

### Health Check
- `GET /api/health` - Check API status

### Users
- `GET /api/users` - Get all users
- `GET /api/users/:id` - Get user by ID
- `POST /api/users` - Create new user
- `PUT /api/users/:id` - Update user
- `DELETE /api/users/:id` - Delete user
- `GET /api/users/count` - Get user count

### Example API Usage

#### Create User
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com"}'
```

#### Get All Users
```bash
curl http://localhost:8080/api/users
```

## 🛠️ Features

- ✅ MySQL database integration
- ✅ GORM ORM for database operations
- ✅ RESTful API endpoints
- ✅ CORS configuration
- ✅ Auto database migration
- ✅ Initial data seeding
- ✅ CRUD operations for users

## 📁 Project Structure

```
backend/
├── main.go              # Main application file
├── models/
│   └── user.go         # User model
├── database/
│   └── database.go     # Database connection
├── handlers/
│   └── user_handler.go # User API handlers
└── README.md           # This file
```
