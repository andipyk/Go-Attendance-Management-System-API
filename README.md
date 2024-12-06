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