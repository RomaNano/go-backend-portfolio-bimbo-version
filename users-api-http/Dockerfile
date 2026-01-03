# ---------- build stage ----------
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Устанавливаем git (нужен для go mod download)
RUN apk add --no-cache git

# Копируем только файлы зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники
COPY . .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o api ./cmd/api


# ---------- runtime stage ----------
FROM alpine:3.20

WORKDIR /app

# Добавляем ca-certificates (для HTTPS / внешних API)
RUN apk add --no-cache ca-certificates

# Копируем бинарник из build-стейджа
COPY --from=builder /app/api /app/api

# Порт, который слушает приложение
EXPOSE 8080

# Запуск
ENTRYPOINT ["/app/api"]
