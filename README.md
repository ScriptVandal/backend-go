# backend-go (portfolio API)

API для портфолио на чистом Go (net/http) с JSON-хранилищем и инструкциями по миграции на PostgreSQL.

## Запуск локально (JSON-хранилище)
1. `go run ./cmd/server`
2. Эндпоинты:
   - `GET /health`
   - `GET /api/projects`
   - `GET /api/skills`
   - `GET /api/contacts`
   - `GET /api/posts`

## CORS
Разрешен `http://localhost:3000`. При необходимости поменяйте в `internal/middleware/cors.go`.

## Тесты curl
```
curl -i http://localhost:8080/health
curl -i http://localhost:8080/api/projects
```

## Миграция на PostgreSQL (пример для projects)
1. Запустите Postgres (например, `docker-compose up db`).
2. Создайте таблицу:
```sql
CREATE TABLE projects (
  id TEXT PRIMARY KEY,
  title TEXT NOT NULL,
  description TEXT,
  tags TEXT[],
  url TEXT
);
```
3. Установите `DATABASE_URL`, например:
```
export DATABASE_URL=postgres://postgres:postgres@localhost:5432/portfolio?sslmode=disable
```
4. В `cmd/server/main.go` замените JSON-репозиторий на Postgres:
```go
// import "database/sql"
// import _ "github.com/lib/pq"
// dsn := os.Getenv("DATABASE_URL")
// db, _ := sql.Open("postgres", dsn)
// projectRepo := repositories.NewPGProjectRepository(db)
```
5. Аналогично можно сделать PG-репозитории для skills/contacts/posts.

## Docker
```
docker build -t backend-go .
docker run -p 8080:8080 backend-go
```

## Docker Compose (API + Postgres)
```
docker compose up --build
```
API будет на `localhost:8080`, Postgres на `localhost:5432`.

## Интеграция с Next.js
```ts
useEffect(() => {
  fetch("http://localhost:8080/api/projects")
    .then((r) => r.json())
    .then(setProjects)
}, [])
```
Или `axios.get("http://localhost:8080/api/projects")`.

## Дальнейшие шаги
- Добавить аутентификацию (JWT / cookies) для POST/PUT/DELETE.
- Вынести конфиг CORS и PORT в env.
- Добавить go.sum: `go mod tidy`.
- Написать репозитории для Postgres для остальных сущностей.
- Добавить авто-миграции или SQL-скрипты.
