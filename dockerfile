# Use golang base image to build the app
FROM golang:1.23-alpine AS builder

# Install necessary dependencies for go-sqlite3
RUN apk add --no-cache \
    sqlite-dev \
    gcc \
    g++

WORKDIR /app

# Copy Go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source
COPY . .

# Build the Go application
RUN go build -o main ./app/main.go

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./main"]