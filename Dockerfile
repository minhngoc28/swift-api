# Stage 1: Build binary
FROM golang:1.24.2-alpine AS builder

WORKDIR /app

# Copy only go.mod and go.sum first to leverage caching
COPY go.mod go.sum ./
RUN go mod download

# Copy rest of the code
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o swift-api main.go

# Stage 2: Use tiny image to run
FROM alpine:3.19

WORKDIR /app

# Copy binary đã build từ builder
COPY --from=builder /app/swift-api .

# ✅ Copy thêm file dữ liệu CSV
COPY swift.csv .

# Expose port
EXPOSE 8080

CMD ["./swift-api"]
