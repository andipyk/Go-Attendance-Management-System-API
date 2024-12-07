# Attendance Management System API

A robust REST API built with Go (Golang) that implements clean architecture principles for managing employee attendance. This system provides a secure and efficient way to track employee attendance with features like marking attendance, viewing attendance history, and user management.

## Overview

This project is designed to provide a modern solution for attendance management with the following key features:
- User authentication and authorization using JWT
- Attendance marking with status (present, absent, late)
- Daily attendance reports
- User profile management
- Clean and maintainable codebase using clean architecture

## Technology Stack

- **Language**: Go 1.16+
- **Framework**: Gin Web Framework
- **Database**: MySQL 8.0+
- **Authentication**: JWT (JSON Web Tokens)
- **Architecture**: Clean Architecture

## Project Structure

```
.
├── cmd/
│   └── server/             # Application entry point
├── config/                 # Configuration management
├── internal/
│   ├── domain/            # Business entities and interfaces
│   ├── repository/        # Data persistence implementations
│   ├── usecase/           # Business logic implementations
│   ├── delivery/
│   │   └── http/         # HTTP handlers and routes
│   ├── middleware/        # HTTP middlewares
│   └── utils/            # Shared utilities
└── pkg/
    └── db/               # Database connection management
```

## Prerequisites

- Go 1.16 or higher
- MySQL 8.0 or higher
- Git

## Installation

1. Clone the repository
```bash
git clone <repository-url>
cd attendance-system
```

2. Install dependencies
```bash
go mod tidy
```

3. Set up the database
```sql
-- Create database
CREATE DATABASE attendance_db;

-- Create users table
CREATE TABLE users (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Create attendances table
CREATE TABLE attendances (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    attendance_date DATE NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'present',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE KEY unique_user_date (user_id, attendance_date)
);
```

4. Configure environment variables
```bash
cp .env.example .env
```

Edit the `.env` file with your configuration:
```env
DB_DRIVER=mysql
DB_SOURCE=root:password@tcp(localhost:3306)/attendance_db?parseTime=true
SERVER_ADDRESS=:8080
JWT_SECRET=your-super-secret-key-change-this-in-production
```

5. Run the application
```bash
go run cmd/server/main.go
```

## API Endpoints

### Authentication Endpoints
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | /api/users/register | Register new user | No |
| POST | /api/users/login | User login | No |

### User Endpoints
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | /api/users/profile | Get user profile | Yes |
| PUT | /api/users/profile | Update user profile | Yes |

### Attendance Endpoints
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | /api/attendance | Mark attendance | Yes |
| GET | /api/attendance | Get attendance by date | Yes |
| GET | /api/attendance/user | Get user's attendance history | Yes |

## API Usage Examples

### Register User
```bash
curl -X POST http://localhost:8080/api/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/api/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

### Mark Attendance
```bash
curl -X POST http://localhost:8080/api/attendance \
  -H "Authorization: Bearer <your-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "present"
  }'
```

## Future Development Plans

1. **Enhanced Features**
   - Mobile application integration
   - Geolocation-based attendance marking
   - Face recognition for attendance verification
   - Leave management system
   - Overtime tracking

2. **Technical Improvements**
   - Implement caching layer (Redis)
   - Add GraphQL support
   - Implement real-time notifications
   - Add comprehensive logging and monitoring
   - Containerization with Docker
   - Implement CI/CD pipeline

3. **Scalability**
   - Microservices architecture
   - Message queue implementation
   - Horizontal scaling capabilities
   - Multi-tenant support

4. **Security Enhancements**
   - Two-factor authentication
   - Role-based access control (RBAC)
   - API rate limiting
   - Enhanced audit logging


## Testing Guide

### Test Structure
Our tests follow a consistent table-driven test pattern for better readability and maintainability. Here's the standard structure we use: 

### Usecase Tests

The usecase layer contains the core business logic tests. We have two main test files:

1. **User Usecase Tests** (`user_usecase_test.go`)
   - Tests for user registration
   - Tests for user login
   - Tests for profile management

2. **Attendance Usecase Tests** (`attendance_usecase_test.go`)
   - Tests for marking attendance
   - Tests for retrieving attendance records
   - Tests for user attendance history

### Implemented Test Cases

#### User Tests
1. **Registration Tests**
   - Success registration with valid data
   - Failure when email already exists
   - Database errors during registration
   - Password hashing errors

2. **Login Tests**
   - Success login with correct credentials
   - Failure with wrong password
   - Failure with non-existent email
   - Database errors during login

3. **Profile Tests**
   - Success profile retrieval
   - Success profile update with/without password change
   - User not found scenarios
   - Database errors during profile operations

#### Attendance Tests
1. **Mark Attendance Tests**
   - Success marking attendance
   - Failure when already marked for the day
   - User not found scenarios
   - Default status handling
   - Database errors during marking

2. **Attendance Retrieval Tests**
   - Success getting attendance by date
   - Success getting user attendance history
   - No records found scenarios
   - Database errors during retrieval

### Running the Tests

1. **Run All Tests**
```bash
go test ./internal/usecase/...
```

2. **Run Specific Test Files**
```bash
# User tests
go test ./internal/usecase/user_usecase_test.go

# Attendance tests
go test ./internal/usecase/attendance_usecase_test.go
```

3. **Run with Coverage**
```bash
# Get coverage percentage
go test -cover ./internal/usecase/...

# Generate detailed coverage report
go test -coverprofile=coverage.out ./internal/usecase/...
go tool cover -html=coverage.out
```

4. **Run Specific Test Cases**
```bash
# Run specific test function
go test -run TestUserUsecase_Register ./internal/usecase/...

# Run specific test case
go test -run TestUserUsecase_Register/Success ./internal/usecase/...
```

### Test Coverage Results
Current test coverage for usecase layer:
- User Usecase: 100% coverage
- Attendance Usecase: 100% coverage

Key areas covered:
- Success scenarios
- Error handling
- Edge cases
- Database interactions
- Input validation
- Business logic validation

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For support, please open an issue in the GitHub repository or contact the development team.