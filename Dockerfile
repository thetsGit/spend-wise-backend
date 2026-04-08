# Stage 1 - Build the app
FROM golang:1.26-alpine3.23 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main ./cmd/server

# Stage 2 - Run the app
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
CMD ["./main"]
