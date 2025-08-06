# Build Stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod .

COPY go.sum .

# Download Go modules
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /app/myapp .

# Runtime Stage
FROM alpine:latest

# Install ca-certificates for HTTPS communication if needed
RUN apk update --no-cache && apk add --no-cache ca-certificates

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/myapp /app/myapp

# Expose the port your application listens on
# EXPOSE 8081

EXPOSE 50051
# Command to run when the container starts
CMD ["/app/myapp"]