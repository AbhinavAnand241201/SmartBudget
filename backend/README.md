# SmartBudget Backend

The backend service for the SmartBudget application, built with Go and Gin.

## Features

- RESTful API endpoints for users, transactions, and budgets
- PostgreSQL database integration
- Docker support for development and deployment
- Comprehensive test coverage
- Environment-based configuration

## Prerequisites

- Go 1.21 or later
- PostgreSQL 15 or later
- Docker and Docker Compose (for containerized development)
- Make (optional, for using Makefile commands)

## Getting Started

### Local Development

1. Clone the repository
2. Copy `.env.example` to `.env` and update the values
3. Install dependencies:
   ```bash
   make deps
   ```
4. Run the application:
   ```bash
   make run
   ```

### Using Docker

1. Build and start the containers:
   ```bash
   docker-compose up --build
   ```
2. The API will be available at `http://localhost:8080`

### Running Tests

```bash
make test
```

For test coverage:
```bash
make test-coverage
```

### Database Setup

Create the test database:
```bash
make test-db
```

Drop the test database:
```bash
make drop-test-db
```

## API Endpoints

### Health Check
- `GET /health` - Check API health

### Users
- `POST /api/users` - Create a new user
- `GET /api/users/:id` - Get user by ID
- `PUT /api/users/:id` - Update user
- `DELETE /api/users/:id` - Delete user

## Project Structure

```
backend/
├── api/           # API handlers
├── config/        # Configuration
├── db/            # Database models and connection
├── services/      # Business logic
└── utils/         # Utility functions
```

## Development

### Adding New Dependencies

```bash
go get github.com/example/package
go mod tidy
```

### Code Style

Run the linter:
```bash
make lint
```

## Deployment

1. Build the Docker image:
   ```bash
   docker build -t smartbudget-backend .
   ```
2. Run the container:
   ```bash
   docker run -p 8080:8080 smartbudget-backend
   ```

## License

MIT License 