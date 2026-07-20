# Build stage
FROM golang:1.24-alpine AS builder

# SQLite requires CGO
RUN apk add --no-cache gcc musl-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main .

# Run stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8082
CMD ["./main"]
