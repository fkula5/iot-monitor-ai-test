# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Run stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8083
CMD ["./main"]
