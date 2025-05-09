# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o gopet .

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/gopet .
# Copy the frontend files
COPY ./dist ./dist

# Expose the port the app runs on
EXPOSE 8080

# Command to run the application
CMD ["./gopet", "-port=8080"]