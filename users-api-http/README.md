# Users API (net/http)

Учебный backend-проект на Go, реализованный на чистом `net/http`.
Проект демонстрирует слоистую архитектуру, middleware, работу с PostgreSQL
и отказ от фреймворков в пользу стандартной библиотеки.

## Стек
- Go
- net/http
- PostgreSQL
- database/sql
- Docker
- slog / log

## Архитектура

HTTP (net/http)
  → Middleware (Request ID, Logging)
    → Handler (HTTP / DTO)
      → Service (business logic)
        → Repository (PostgreSQL)

## Функциональность
- Health check
- Создание пользователя
- Получение списка пользователей
- Получение пользователя по ID
- Валидация входных данных
- Domain → HTTP error mapping
- Request ID propagation
- Logging middleware (status, duration)

## Запуск

### 1. Запуск PostgreSQL
docker compose up postgres

PostgreSQL будет доступен на `localhost:5433`.

### 2. Запуск API
(!!!) DB_DSN тут оставлено в учебных целях, в реальном проекте скрываю от паблик поля.
DB_DSN="postgres://users:users@localhost:5433/users?sslmode=disable" \
go run ./cmd/api

API стартует на `http://localhost:8080`.

## Эндпоинты

### Health
GET /health
{"status":"ok"}

### Создание пользователя
POST /users
Content-Type: application/json

{
  "email": "test@example.com"
}

Ответ:
{
  "id": 1,
  "email": "test@example.com",
  "created_at": "2025-12-20T13:28:33Z"
}


## Примечания
- Transport слой реализован без фреймворков
- Middleware написаны вручную через `http.Handler`
- Repository изолирует работу с БД
- Service слой не зависит от HTTP

