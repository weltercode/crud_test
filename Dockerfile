# Build stage
FROM golang:1.22.3-alpine as builder

WORKDIR /app

# Copy module files
COPY go.mod go.sum ./

# Install dependencies
RUN go mod tidy && go mod download

# Copy application files
COPY . .

# Set the working directory for the build process
WORKDIR /app/internal/cmd

# Build the application
RUN go build -o main .

# Run stage
FROM alpine:latest

WORKDIR /root/

# Copy the compiled binary
COPY --from=builder /app/internal/cmd/main .

EXPOSE 8080

CMD ["./main"]