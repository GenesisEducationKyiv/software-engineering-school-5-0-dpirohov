include .env

DB_URL=$(DATABASE_URL)

# Build the application
all: build test

build:
	@echo "Building..."
	
	
	@go build -o main cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi

.PHONY: all build run test clean watch

# Generate new migration: make migration name=create_users
migration:
	migrate create -ext sql -dir migrations -seq $(name)

# Apply all migrations
migrate-up:
	migrate -path migrations -database "$(DB_URL)?sslmode=disable" up

# Rollback
migrate-down:
	migrate -path migrations -database "$(DB_URL)" down 1

# Show current migration version
migrate-version:
	migrate -path migrations -database "$(DB_URL)" version

# Rollback all migrations
migrate-reset:
	migrate -path migrations -database "$(DB_URL)" drop -f