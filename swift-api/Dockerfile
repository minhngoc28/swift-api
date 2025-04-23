# Build stage - chạy trong môi trường Linux
FROM golang:1.24.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# ✅ BẮT BUỘC để Docker build đúng cho Linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o swift-api main.go

# Final stage
FROM alpine:3.19

WORKDIR /app
COPY --from=builder /app/swift-api .

CMD ["./swift-api"]
