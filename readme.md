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

# OAuth Configuration
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
GOOGLE_REDIRECT_URI=http://127.0.0.1:3002/api/auth/google/callback

# Frontend URL for OAuth redirects
FRONTEND_URL=http://localhost:3000
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

| Command                | Description                                               |
| ---------------------- | --------------------------------------------------------- |
| `go run main.go`       | Run with auto migration (preserves existing data)         |
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

## 📊 Team Capacity Feature

The project creation process includes intelligent team capacity management:

### Stage 2: Set Team Capacity

- Define `total_team` - the maximum number of team members for your project
- Example: If you set `total_team` to 5, you can have a maximum of 5 people

### Stage 4: Allocate Team Positions

- **Add Members**: Directly invite specific users to join your project
- **Create Roles**: Define open positions for recruitment
- **Automatic Calculation**: The system tracks:
  - `filled_team`: Number of members already added
  - `remaining_team`: Available positions for new roles
  - Total allocation cannot exceed `total_team`

### Example Workflow

```
Stage 2: total_team = 5
Stage 4:
  - Add 2 members directly → filled_team = 2
  - Create 2 roles with 1 slot each → role_slots = 2
  - Remaining capacity = 5 - 2 - 2 = 1 slot available
```

## 🔐 OAuth Configuration

The project supports OAuth authentication with Google. After successful authentication, users are redirected to the frontend with tokens in query parameters.

### OAuth Flow

1. **Login Initiation**: `GET /api/auth/google/login`

- Redirects user to Google OAuth consent screen

2. **OAuth Callback**: `GET /api/auth/google/callback`

- Handles Google OAuth response
- Redirects to frontend with authentication data

### Frontend Redirect Format

**Success Redirect:**

```
{FRONTEND_URL}/auth/callback?success=true&token={jwt_token}&user_id={id}&user_name={name}&user_email={email}
```

**Error Redirect:**

```
{FRONTEND_URL}/auth/callback?error={error_type}
```

### Environment Variables

| Variable               | Description                              | Example                                          |
| ---------------------- | ---------------------------------------- | ------------------------------------------------ |
| `FRONTEND_URL`         | Primary frontend URL for OAuth redirects | `http://localhost:3000`                          |
| `CLIENT_URL`           | Alternative frontend URL (fallback)      | `http://localhost:3000`                          |
| `GOOGLE_CLIENT_ID`     | Google OAuth Client ID                   | Your Google OAuth client ID                      |
| `GOOGLE_CLIENT_SECRET` | Google OAuth Client Secret               | Your Google OAuth client secret                  |
| `GOOGLE_REDIRECT_URI`  | Google OAuth callback URL                | `http://127.0.0.1:3002/api/auth/google/callback` |

### URL Priority

The system uses the following priority for determining the frontend redirect URL:

1. `FRONTEND_URL` environment variable
2. `CLIENT_URL` environment variable
3. Derived from `APP_URL` (changes port to 3000)
4. Default fallback: `http://localhost:3000`

### Additional OAuth Endpoints

- `GET /api/auth/success` - Handles OAuth success data extraction
- `GET /api/auth/error` - Handles OAuth error information

````

### API Endpoints

- `GET /api/projects/:id/capacity` - Get current team capacity information
- Response includes:
  ```json
  {
    "total_team": 5,
    "filled_team": 2,
    "total_role_slots": 2,
    "remaining_team": 1,
    "members": 2,
    "roles": 2
  }
````
