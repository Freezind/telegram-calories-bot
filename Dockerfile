# Simple single-stage Dockerfile for homework/demo purposes
FROM golang:1.24-alpine

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy source code
COPY . .

# Download dependencies and build
RUN go mod download
RUN go build -o unified cmd/unified/main.go

# Set timezone
ENV TZ=Asia/Seoul

# Expose port
EXPOSE 8080

# Run the unified backend
CMD ["./unified"]
