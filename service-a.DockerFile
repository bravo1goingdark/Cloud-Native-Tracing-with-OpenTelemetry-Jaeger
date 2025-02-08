# Build Stage
FROM golang:1.23.2 AS builder

WORKDIR /app

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy entire project
COPY . .

# Build the service-a binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o service-a ./cmd/service-a/main.go

# Final minimal image
FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/service-a .
RUN chmod +x service-a

EXPOSE 8080 50051

CMD ["./service-a"]
