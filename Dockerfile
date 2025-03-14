FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o url-shortener ./cmd/url-shortener/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/url-shortener .
COPY .env .env

RUN apk add --no-cache ca-certificates

RUN chmod +x /app/url-shortener

ENTRYPOINT ["./url-shortener"]