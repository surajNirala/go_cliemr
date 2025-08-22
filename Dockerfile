# ---------- Build Stage ----------
FROM golang:1.24.4-alpine AS builder

# Install git & build tools
RUN apk add --no-cache git build-base

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go binary (entrypoint: cmd/server/main.go)
RUN go build -o server ./cmd/server

# ---------- Run Stage ----------
FROM alpine:latest

# Install timezone data & certificates
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy compiled binary from builder
COPY --from=builder /app/server .

# Expose Gin default port (change if you use a different one)
EXPOSE 7070

# Run the binary
CMD ["./server"]
