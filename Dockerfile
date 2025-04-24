# Stage 1: Build binary
FROM golang:1.24.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o swift-api main.go

# Stage 2: Use tiny image to run
FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/swift-api .
COPY swift.csv .

EXPOSE 8080

CMD ["./swift-api"]
