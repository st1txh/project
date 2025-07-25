# Этап сборки
FROM golang:1.23.5-alpine AS builder

# Устанавливаем git и SSL (необходимо для зависимостей)
RUN apk add --no-cache git ca-certificates

# Рабочая директория
WORKDIR /app

# Копируем сначала только модули (для кэширования)
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем весь проект
COPY . .

# Собираем приложение (учитываем структуру cmd/main/app.go)
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/restapi ./cmd/main

# Этап запуска
FROM alpine:latest

WORKDIR /app

# Копируем
COPY --from=builder /app/bin/restapi .
COPY config/Config.yml .
COPY docs ./docs
COPY --from=builder /app/.env .

# Создаем папку для логов
RUN mkdir -p /app/logs

EXPOSE "${PORT}"

CMD ["./restapi"]