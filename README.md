# backend-go (portfolio API)

API для портфолио на чистом Go (net/http) с JSON-хранилищем и опциональной работой с PostgreSQL.

## Запуск локально (JSON-хранилище по умолчанию)
1. `go run ./cmd/server`
2. Эндпоинты:
   - `GET /health`
   - `GET /api/projects`
   - `GET /api/skills`
   - `GET /api/contacts`
   - `GET /api/posts`

## Переключение на PostgreSQL
- Приложение автоматически использует PG, если установлена переменная `DATABASE_URL`.
- В `cmd/server/main.go` уже есть переключение: при наличии `DATABASE_URL` выбираются PG-репозитории, иначе JSON.

### Подготовка БД
1. Запустите Postgres (например, `docker-compose up db`).
2. Создайте таблицы:
```sql
CREATE TABLE projects (
  id TEXT PRIMARY KEY,
  title TEXT NOT NULL,
  description TEXT,
  tags TEXT[],
  url TEXT
);

CREATE TABLE skills (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  level TEXT,
  category TEXT
);

CREATE TABLE contacts (
  email TEXT,
  telegram TEXT,
  linkedin TEXT,
  github TEXT
);

CREATE TABLE posts (
  id TEXT PRIMARY KEY,
  title TEXT NOT NULL,
  content TEXT,
  tags TEXT[],
  published_at TEXT
);
```
3. Установите `DATABASE_URL`, например:
```
export DATABASE_URL=postgres://postgres:postgres@localhost:5432/portfolio?sslmode=disable
```
4. Запустите: `go run ./cmd/server` — при наличии `DATABASE_URL` будут использованы PG-репозитории для всех сущностей.

### Пример работы с массивами (tags)
- В PG-репозиториях для projects и posts поле `tags` читается через `pq.Array(&slice)`. 
- При вставке используйте `pq.Array(slice)` в `Exec`/`Query`.

## CORS
Разрешен `http://localhost:3000`. При необходимости поменяйте в `internal/middleware/cors.go`.

## Тесты curl
```
curl -i http://localhost:8080/health
curl -i http://localhost:8080/api/projects
```

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
Или `axios.get("http://localhost:8080/api/projects");`.

## Дальнейшие шаги
- Добавить аутентификацию (JWT / cookies) для POST/PUT/DELETE.
- Вынести конфиг CORS и PORT в env.
- Добавить авто-миграции или SQL-скрипты.