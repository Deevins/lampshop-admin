FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o lampshop-admin-service ./cmd/lampshop_admin/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/lampshop-admin-service .

EXPOSE 8083
ENTRYPOINT ["./lampshop-admin-service"]