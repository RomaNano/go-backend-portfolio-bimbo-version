# Users API (Gin)

Учебный backend-проект на Go, демонстрирующий базовую production-архитектуру:
слоистую структуру, работу с PostgreSQL, middleware, логирование и конфигурацию
через переменные окружения.

## Стек
- Go
- Gin
- PostgreSQL
- Docker
- slog (logging)

## Архитектура
HTTP (Gin)
  → Handler (DTO, HTTP)
    → Service (business logic)
      → Repository (PostgreSQL)

## Функциональность
- Создание пользователя
- Получение списка пользователей
- Получение пользователя по ID
- Валидация входных данных
- Корректные HTTP-коды ошибок
- Request ID middleware
- Логирование запросов

------------------------------------------------
## Запуск проекта

### 1. Запуск PostgreSQL (Docker)
docker compose up postgres


PostgreSQL будет доступен на `localhost:5433`.

### 2. Запуск API локально
[DB_DSN раскрыто для учебных целей] 
DB_DSN="postgres://users:users@localhost:5433/users?sslmode=disable" \
go run ./cmd/api

API стартует на `http://localhost:8080`.

------------------------------------------------
## Эндпоинты

### Health check
GET /health
Ответ:
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

### Получение списка пользователей
GET /users

### Ошибки
- `422 Unprocessable Entity` — некорректные данные
- `404 Not Found` — пользователь не найден
- `500 Internal Server Error` — внутренняя ошибка

## Примечания
- Бизнес-логика изолирована в service layer
- Repository реализован через интерфейс
- Уникальность email обеспечивается на уровне БД
