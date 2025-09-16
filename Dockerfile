# Multi-stage build for smaller image
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o whu_class .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/whu_class .

# Expose port
EXPOSE 8080

# Run the web server by default
CMD ["./whu_class", "-web"]