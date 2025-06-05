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

## Prerequisites

- Go 1.22 or higher
- PostgreSQL 12 or higher
- Docker (optional)

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

4. Run database migrations:
```bash
go run cmd/migrate/main.go up
```

5. Run the application:
```bash
go run cmd/main.go
```

The server will start at http://localhost:8082 (or the port specified in your config).

## Docker Setup

1. Build the Docker image:
```bash
docker-compose build
```

2. Start the services:
```bash
docker-compose up -d
```

## API Documentation

Once the server is running, you can access the Swagger documentation at:
http://localhost:8082/swagger/index.html

## API Endpoints

### Student Endpoints
- `POST /students` - Create a new student
- `GET /students/{id}` - Get a student by ID
- `POST /students/login` - Login a student

### Teacher Endpoints
- `POST /teachers` - Create a new teacher
- `GET /teachers/{id}` - Get a teacher by ID
- `PUT /teachers/{id}` - Update a teacher
- `DELETE /teachers/{id}` - Delete a teacher
- `POST /teachers/login` - Login a teacher
- `POST /teachers/{teacherId}/students/{studentId}` - Assign a student to a teacher
- `GET /teachers/{teacherId}/students` - Get all students assigned to a teacher

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
└── docker/            # Docker configuration
```

## Development

### Adding New Migrations

```bash
go run cmd/migrate/main.go create add_new_table
```

### Running Tests

```bash
go test ./...
```

### Generating Swagger Documentation

```bash
swag init -g cmd/main.go
```

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request 