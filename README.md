# backend-go (портфолио API)

Полноценный REST API на Go (net/http) для портфолио с двумя режимами хранения: JSON (read-only) и PostgreSQL (полный CRUD) + аутентификация на JWT/refresh.

## Возможности
- JSON режим (только чтение) или PostgreSQL (полный CRUD)
- JWT: access + refresh c учётом JTI, Argon2id для паролей
- Авторизация: GET публично, POST/PUT/DELETE только с Bearer
- Сущности: Projects, Skills, Contacts, Posts (tags как text[])
- CORS: конфиг через env, с credentials
- Docker Compose для прод-режима с volume и init.sql

## Быстрый старт

### Локально, JSON (read-only)
```bash
go run ./cmd/server
```
Эндпоинты: GET /health, /api/projects, /api/skills, /api/contacts, /api/posts, /api/{entity}/{id}
Запись (POST/PUT/DELETE) вернёт 403 "write operations not supported in JSON mode".

### PostgreSQL (полный CRUD + Auth)
1) Скопируйте env:
```bash
cp .env.example .env
```
2) Заполните .env (пример):
```env
PORT=8080
DATABASE_URL=postgres://postgres:postgres@localhost:5432/portfolio?sslmode=disable
JWT_SECRET=заменить-на-секрет
JWT_REFRESH_SECRET=заменить-на-refresh-секрет
ACCESS_TTL=15m
REFRESH_TTL=168h
CORS_ORIGINS=http://localhost:3000
```
Секреты: `openssl rand -base64 32`

3) Инициализация БД:
```bash
docker compose up -d db
psql -U postgres -d portfolio -f init.sql
```
Таблицы: users, refresh_tokens, projects, skills, contacts, posts (tags text[]).

4) Запуск API:
```bash
export $(cat .env | xargs)
go run ./cmd/server
```
Или docker compose up (см. ниже).

## Аутентификация (JWT/refresh)
- Регистрация: POST /api/auth/register {email,password}
- Логин: POST /api/auth/login {email,password}
- Обновление: POST /api/auth/refresh {refresh_token}
- Логаут: POST /api/auth/logout {refresh_token} (ревокация по jti)
- Доступ: Bearer access обязателен для POST/PUT/DELETE.
- TTL по умолчанию: access 15m, refresh 7d.

## Примеры cURL
Регистрация:
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"u@test.com","password":"pass"}'
```
Создание проекта (нужен токен):
```bash
curl -X POST http://localhost:8080/api/projects \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"id":"p1","title":"Demo","tags":["go"],"url":"https://example.com"}'
```
Обновление:
```bash
curl -X PUT http://localhost:8080/api/projects/p1 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Updated"}'
```
Удаление:
```bash
curl -X DELETE http://localhost:8080/api/projects/p1 \
  -H "Authorization: Bearer $TOKEN"
```
Публичные GET:
```bash
curl http://localhost:8080/api/projects
curl http://localhost:8080/api/projects/p1
```

## Интеграция с Next.js
Базовый URL: `http://localhost:8080`

### Хранение токенов
- access в памяти (React state) или httpOnly cookie (более безопасно через бэкенд-прокси), refresh в httpOnly cookie или secure storage; для простоты примера — в localStorage.

### Минимальный клиент
```ts
const API = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

const getJson = async <T>(url: string, token?: string): Promise<T> => {
  const res = await fetch(`${API}${url}`, {
    headers: token ? { Authorization: `Bearer ${token}` } : {},
    credentials: 'include',
  });
  if (!res.ok) throw new Error(await res.text());
  return res.json();
};

export const login = async (email: string, password: string) => {
  const res = await fetch(`${API}/api/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password }),
    credentials: 'include',
  });
  if (!res.ok) throw new Error(await res.text());
  const data = await res.json();
  localStorage.setItem('access_token', data.access_token);
  localStorage.setItem('refresh_token', data.refresh_token);
  return data;
};

export const fetchProjects = (token?: string) => getJson('/api/projects', token);
export const fetchProject = (id: string, token?: string) => getJson(`/api/projects/${id}`, token);

export const createProject = async (project: any) => {
  const token = localStorage.getItem('access_token');
  const res = await fetch(`${API}/api/projects`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(project),
    credentials: 'include',
  });
  if (res.status === 401 || res.status === 403) throw new Error('unauthorized');
  if (!res.ok) throw new Error(await res.text());
  return res.json();
};
```

### Обновление access по refresh (пример хука)
```ts
import { useEffect, useState, useCallback } from 'react';

export const useAuthToken = () => {
  const [token, setToken] = useState<string | null>(null);

  const refresh = useCallback(async () => {
    const rt = localStorage.getItem('refresh_token');
    if (!rt) return null;
    const res = await fetch(`${API}/api/auth/refresh`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ refresh_token: rt }),
      credentials: 'include',
    });
    if (!res.ok) return null;
    const data = await res.json();
    localStorage.setItem('access_token', data.access_token);
    setToken(data.access_token);
    return data.access_token;
  }, []);

  useEffect(() => {
    const t = localStorage.getItem('access_token');
    if (t) setToken(t);
  }, []);

  return { token, refresh };
};
```
Используйте refresh при 401/403, затем повторяйте запрос с новым access.

### CORS для Next.js
В .env API установите `CORS_ORIGINS=http://localhost:3000` (или домен фронта). В fetch используйте `credentials: 'include'`, если refresh/куки.

### SSR/Next.js Route Handlers
Для серверных роутов можно проксировать запросы к API, выставляя Authorization заголовок из cookie сессии. Пример:
```ts
// app/api/projects/route.ts
import { NextRequest, NextResponse } from 'next/server';

export async function GET(req: NextRequest) {
  const token = req.cookies.get('access_token')?.value;
  const res = await fetch(`${process.env.API_URL}/api/projects`, {
    headers: token ? { Authorization: `Bearer ${token}` } : {},
  });
  const data = await res.json();
  return NextResponse.json(data, { status: res.status });
}
```

## Docker Compose (prod)
```bash
cp .env.example .env
# заполните .env

# запуск
docker compose up -d
# логи
docker compose logs -f
# остановка
docker compose down
# очистка с volume (удалит БД)
docker compose down -v
```
Сервисы: api (8080), postgres (5432), volume для данных, автоприменение init.sql.

## Модели
Project: id, title, description, tags[], url
Skill: id, name, level, category
Contact: id, email, telegram, linkedin, github
Post: id, title, content, tags[], published_at

## Политика доступа
- GET — публично
- POST/PUT/DELETE — только с валидным Bearer access
- /api/auth/* и /health — без авторизации

## Диагностика
- Подключение к БД: `psql -U postgres -h localhost -d portfolio`
- Переменные: `echo $DATABASE_URL`
- JWT не работает: проверьте секреты и TTL, формат `Authorization: Bearer <token>`
- Ошибки записи в JSON-режиме: ожидаемо, переходите на PG (DATABASE_URL)

## Безопасность
- В проде ставьте сильные секреты и HTTPS
- Ограничьте CORS точными доменами
- Добавьте rate limiting в реальных сценариях
