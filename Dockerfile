# Step 1: Build the Go binary
FROM golang:1.26.3-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files first to leverage Docker caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the application from the cmd/api directory
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api/main.go

# Step 2: Run the binary in a lightweight container
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/frontend ./frontend
# Copy the .env file into the container
COPY .env .

# Expose the port your Go API listens on (change 8080 if your app uses a different port)
EXPOSE 8080

# Run the application
CMD ["./main"]