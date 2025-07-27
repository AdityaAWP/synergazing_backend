# Backend Go Project (Synergazing API)

## ⚙️ Installation & Configuration

Follow the following steps to run this project.

### 1. Clone Repository
Open your terminal and clone repository:
```bash
git clone https://github.com/xx/synergazing_backend.git
cd synergazing_backend
```

### 2. Environment Setup
Copy the environment example file and configure your settings:
```bash
cp .env.example .env
```

Edit the `.env` file with your database credentials and other configurations:
```bash
DB_HOST=127.0.0.1
DB_PORT=5432
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=your_database
DB_SSLMODE=disable

JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
APP_URL=http://127.0.0.1:3002
```

### 3. Install Dependencies
```bash
go mod tidy
```

### 4. Database Migration

#### Auto Migration (Default - Preserves Data)
For normal development and production use:
```bash
go run main.go
```

#### Fresh Migration (Drops & Recreates All Tables)
For development when you want to reset the database:
```bash
go run main.go fresh
```

### 5. Run the Application
The server will start automatically after migration. You can access the API at:
```
http://127.0.0.1:3002
```

## 🛠️ Migration Commands

| Command | Description |
|---------|-------------|
| `go run main.go` | Run with auto migration (preserves existing data) |
| `go run main.go fresh` | Run with fresh migration (drops all tables and recreates) |

## 📁 Project Structure

```
synergazing_backend/
├── config/          # Database configuration
├── controller/      # Request controllers
├── handler/         # Request handlers
├── helper/          # Utility functions
├── middleware/      # Custom middleware
├── migrations/      # Database migrations
├── model/          # Database models
├── routes/         # API routes
├── service/        # Business logic
├── storage/        # File uploads
├── .env.example    # Environment example
├── go.mod          # Go modules
├── go.sum          # Go dependencies
└── main.go         # Application entry point
```