.PHONY: build test run clean

# Build the application
build:
	go build -o bin/smartbudget

# Run tests
test:
	go test -v ./...

# Run the application
run:
	go run main.go

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Install dependencies
deps:
	go mod tidy

# Run linter
lint:
	golangci-lint run

# Create test database
test-db:
	createdb smartbudget_test || true
	psql -d smartbudget_test -f db/schema.sql

# Drop test database
drop-test-db:
	dropdb smartbudget_test || true

# Run all tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out 