# Stage 1: Build the Go application
FROM golang:1.24-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
# CGO_ENABLED=0: Disable CGO for static linking
# -ldflags="-s -w": Strip debugging information
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" ./cmd/protomock

# ----------------------------------------------------------------------------

# Stage 2: Create a minimal runtime image
FROM alpine:3.21

# Set the working directory
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/protomock .

# Expose the ports the application will run on
EXPOSE 8000
EXPOSE 8010

# Command to run the application
CMD ["./protomock"]
