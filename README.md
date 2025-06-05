# Student-Teacher Management API

A Go-based REST API for managing students and teachers, built with hexagonal architecture and modern Go practices.

## Features

- Student and Teacher CRUD operations
- Authentication with JWT
- Student-Teacher relationship management
- PostgreSQL database with GORM
- Database migrations
- API documentation with Swagger
- CORS middleware
- Dependency Injection with Google Wire
- Automatic timestamps for records
- Standardized API responses
- Input validation
- Error handling

## Prerequisites

- Go 1.22 or higher
- PostgreSQL 12 or higher
- Docker (optional)
- Make (for using Makefile commands)

## Setup

1. Clone the repository:
```bash
git clone https://github.com/OICjangirrahul/students.git
cd students
```

2. Set up the environment variables (copy .env.example to .env):
```bash
cp .env.example .env
```

Edit the .env file with your PostgreSQL credentials:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_user
DB_PASSWORD=your_password
DB_NAME=students_db
JWT_SECRET=your_jwt_secret
CONFIG_PATH=config/local.yaml
```

3. Install dependencies:
```bash
go mod download
```

4. Start the database:
```bash
make docker-up
```

5. Build and run the application:
```bash
make run
```

The server will start at http://localhost:8082 (or the port specified in your config).

## Available Make Commands

- `make build` - Build the application
- `make run` - Run the application
- `make test` - Run tests
- `make docs` - Generate Swagger documentation
- `make docker-up` - Start Docker containers
- `make docker-down` - Stop Docker containers
- `make clean` - Clean build artifacts
- `make test-with-docker` - Run tests with Docker

## API Documentation

Once the server is running, you can access the Swagger documentation at:
http://localhost:8082/swagger/index.html

## API Endpoints

### Student Endpoints
- `POST /api/v1/students` - Create a new student
- `GET /api/v1/students/{id}` - Get a student by ID
- `POST /api/v1/students/login` - Login a student

### Teacher Endpoints
- `POST /api/v1/teachers` - Create a new teacher
- `GET /api/v1/teachers/{id}` - Get a teacher by ID
- `PUT /api/v1/teachers/{id}` - Update a teacher
- `DELETE /api/v1/teachers/{id}` - Delete a teacher
- `POST /api/v1/teachers/login` - Login a teacher
- `POST /api/v1/teachers/{teacherId}/students/{studentId}` - Assign a student to a teacher
- `GET /api/v1/teachers/{teacherId}/students` - Get all students assigned to a teacher

## Project Structure

```
.
├── cmd/
│   ├── main.go           # Application entry point
│   └── migrate/          # Database migration tool
├── internal/
│   ├── core/            # Domain layer
│   │   ├── domain/      # Domain entities
│   │   ├── ports/       # Interfaces
│   │   └── services/    # Business logic
│   ├── adapters/        # Adapters layer
│   │   ├── http/        # HTTP handlers
│   │   └── repositories/# Database repositories
│   ├── middleware/      # HTTP middleware
│   └── config/         # Configuration
├── migrations/         # Database migrations
├── docs/              # Swagger documentation
├── docker/            # Docker configuration
└── Makefile          # Build and development commands
```

## API Response Format

All API responses follow a standard format:

### Success Response
```json
{
  "success": true,
  "data": {
    // Response data here
  }
}
```

### Error Response
```json
{
  "status": "Error",
  "error": "Error message here"
}
```

## Development

### Adding New Migrations

```bash
go run cmd/migrate/main.go create add_new_table
```

### Running Tests

```bash
make test
```

### Updating Swagger Documentation

```bash
make docs
```

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request 