# Use golang base image to build the app
FROM golang:1.23-alpine AS builder

LABEL \
    website.name="Forum" \
    description="A web forum application built with Go and SQLite and JavaScript. It allows user communication through posts and comments, supports authentication with sessions and cookies, implements a like/dislike system, and provides filtering options. The application follows best practices, handles HTTP and technical errors, and is containerized with Docker for easy deployment." \
    authors="Fethi Abderrahmane, Aymane Berhili, Mostafa Zakri, Anass Elabsi"

# Install necessary dependencies for go-sqlite3
WORKDIR /app

# Install required packages for CGO
RUN apk add --no-cache gcc g++ musl-dev sqlite-dev

# Enable CGO
ENV CGO_ENABLED=1

# Copy Go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source
COPY . .

# Build the Go application
RUN go build -o main ./app/main.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/app/db/schema.sql ./app/db/schema.sql
COPY --from=builder /app/logs/ ./logs/
COPY --from=builder /app/static/ ./static/
COPY --from=builder /app/templates/ ./templates/

RUN apk add --no-cache sqlite-dev bash

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./main"]
