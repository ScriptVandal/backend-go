# backend-go (portfolio API)

Full-featured REST API for portfolio management built with Go (net/http), supporting both JSON (read-only) and PostgreSQL (full CRUD) storage modes, with JWT-based authentication.

## Features

- **Dual Storage Modes**: JSON files (read-only) or PostgreSQL (full CRUD)
- **JWT Authentication**: Access tokens + refresh tokens with JTI tracking
- **Password Security**: Argon2id hashing
- **Authorization**: Public GET endpoints, authenticated POST/PUT/DELETE
- **CORS**: Configurable origins with credentials support
- **Entity Management**: Projects, Skills, Contacts, Posts with tags support
- **Docker Support**: Production-ready compose with persistent volumes

## Quick Start

### Local Development (JSON mode - read-only)

```bash
go run ./cmd/server
```

Endpoints available:
- `GET /health` - Health check
- `GET /api/projects` - List projects
- `GET /api/skills` - List skills  
- `GET /api/contacts` - List contacts
- `GET /api/posts` - List posts
- `GET /api/{entity}/{id}` - Get single item

Note: POST/PUT/DELETE will return `403 Forbidden` with message "write operations not supported in JSON mode"

### PostgreSQL Mode (Full CRUD + Auth)

#### 1. Environment Setup

Copy `.env.example` to `.env` and configure:

```bash
cp .env.example .env
# Edit .env with your configuration
```

Required environment variables for full features:
```env
PORT=8080
DATABASE_URL=postgres://postgres:postgres@localhost:5432/portfolio?sslmode=disable
JWT_SECRET=your-secret-key-change-in-production
JWT_REFRESH_SECRET=your-refresh-secret-change-in-production
ACCESS_TTL=15m
REFRESH_TTL=168h
CORS_ORIGINS=http://localhost:3000
```

Generate secure secrets:
```bash
openssl rand -base64 32
```

#### 2. Database Setup

Start PostgreSQL and initialize:

```bash
# Using Docker Compose
docker compose up -d db

# Or use your local PostgreSQL instance
# Then initialize the schema:
psql -U postgres -d portfolio -f init.sql
```

The `init.sql` creates:
- `users` table (id, email, password_hash)
- `refresh_tokens` table (jti, user_id, expires_at, revoked_at)
- `projects` table (id, title, description, tags[], url)
- `skills` table (id, name, level, category)
- `contacts` table (id, email, telegram, linkedin, github)
- `posts` table (id, title, content, tags[], published_at)

#### 3. Run the API

```bash
# Load environment variables
export $(cat .env | xargs)

# Start server
go run ./cmd/server
```

## Authentication Flow

### Register

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"securepassword"}'
```

Response:
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "abc123",
    "email": "user@example.com",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

### Login

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"securepassword"}'
```

Returns same response as register.

### Refresh Token

```bash
curl -X POST http://localhost:8080/api/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"eyJhbGciOiJIUzI1NiIs..."}'
```

Response:
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

### Logout

```bash
curl -X POST http://localhost:8080/api/auth/logout \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"eyJhbGciOiJIUzI1NiIs..."}'
```

Returns `204 No Content` on success.

## CRUD Operations

### List Items (Public)

```bash
curl http://localhost:8080/api/projects
curl http://localhost:8080/api/skills
curl http://localhost:8080/api/contacts
curl http://localhost:8080/api/posts
```

### Get Single Item (Public)

```bash
curl http://localhost:8080/api/projects/{id}
curl http://localhost:8080/api/skills/{id}
```

### Create Item (Requires Auth)

```bash
curl -X POST http://localhost:8080/api/projects \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "id": "proj1",
    "title": "My Project",
    "description": "A cool project",
    "tags": ["go", "api"],
    "url": "https://github.com/user/project"
  }'
```

### Update Item (Requires Auth)

```bash
curl -X PUT http://localhost:8080/api/projects/proj1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "title": "Updated Project",
    "description": "Updated description",
    "tags": ["go", "api", "rest"],
    "url": "https://github.com/user/project"
  }'
```

### Delete Item (Requires Auth)

```bash
curl -X DELETE http://localhost:8080/api/projects/proj1 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## Authorization Rules

- **GET requests**: Public, no authentication required
- **POST/PUT/DELETE requests**: Require valid Bearer access token
- **Auth endpoints**: No authentication required (`/api/auth/*`)
- **Health check**: No authentication required (`/health`)

## Data Models

### Project
```json
{
  "id": "string",
  "title": "string",
  "description": "string",
  "tags": ["string"],
  "url": "string"
}
```

### Skill
```json
{
  "id": "string",
  "name": "string",
  "level": "string",
  "category": "string"
}
```

### Contact
```json
{
  "id": "string",
  "email": "string",
  "telegram": "string",
  "linkedin": "string",
  "github": "string"
}
```

### Post
```json
{
  "id": "string",
  "title": "string",
  "content": "string",
  "tags": ["string"],
  "published_at": "string"
}
```

## Docker Deployment

### Production with Docker Compose

```bash
# Set environment variables
cp .env.example .env
# Edit .env with production values

# Build and start services
docker compose up -d

# View logs
docker compose logs -f

# Stop services
docker compose down

# Stop and remove volumes (CAUTION: deletes data)
docker compose down -v
```

The compose setup includes:
- API service (Go app) on port 8080
- PostgreSQL 16 on port 5432
- Persistent volume for database
- Auto-initialization with init.sql

### Build Docker Image Only

```bash
docker build -t backend-go .
docker run -p 8080:8080 \
  -e DATABASE_URL="postgres://..." \
  -e JWT_SECRET="..." \
  -e JWT_REFRESH_SECRET="..." \
  backend-go
```

## CORS Configuration

Configure allowed origins via `CORS_ORIGINS` environment variable:

```env
# Single origin
CORS_ORIGINS=http://localhost:3000

# Multiple origins (comma-separated)
CORS_ORIGINS=http://localhost:3000,https://myapp.com,https://staging.myapp.com
```

Default: `http://localhost:3000`

CORS headers set:
- `Access-Control-Allow-Origin`: Matched origin from allowed list
- `Access-Control-Allow-Methods`: GET, POST, PUT, DELETE, OPTIONS
- `Access-Control-Allow-Headers`: Content-Type, Authorization
- `Access-Control-Allow-Credentials`: true

## Development

### Dependencies

```bash
go mod download
```

Required packages:
- `github.com/lib/pq` - PostgreSQL driver
- `github.com/golang-jwt/jwt/v5` - JWT implementation
- `golang.org/x/crypto/argon2` - Password hashing

### Build

```bash
go build -o server ./cmd/server
./server
```

### Project Structure

```
.
├── cmd/server/          # Application entry point
├── internal/
│   ├── config/         # Environment configuration
│   ├── handlers/       # HTTP handlers
│   ├── middleware/     # Auth, CORS, logging
│   ├── models/         # Data models
│   ├── repositories/   # Data access (JSON & PG)
│   └── services/       # Business logic
├── data/               # JSON storage files
├── init.sql            # Database schema
├── docker-compose.yml  # Docker orchestration
├── Dockerfile          # Container build
└── .env.example        # Environment template
```

## Integration with Frontend

### Next.js / React Example

```typescript
// API client
const API_URL = 'http://localhost:8080';

// Login
const login = async (email: string, password: string) => {
  const res = await fetch(`${API_URL}/api/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password }),
  });
  const data = await res.json();
  localStorage.setItem('access_token', data.access_token);
  localStorage.setItem('refresh_token', data.refresh_token);
  return data;
};

// Fetch with auth
const fetchProjects = async () => {
  const token = localStorage.getItem('access_token');
  const res = await fetch(`${API_URL}/api/projects`, {
    headers: {
      'Authorization': `Bearer ${token}`,
    },
  });
  return res.json();
};

// Create with auth
const createProject = async (project: Project) => {
  const token = localStorage.getItem('access_token');
  const res = await fetch(`${API_URL}/api/projects`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`,
    },
    body: JSON.stringify(project),
  });
  return res.json();
};
```

## Security Notes

- **Production Secrets**: Always use strong, randomly generated secrets
- **HTTPS**: Use HTTPS in production (configure reverse proxy)
- **Token TTLs**: Adjust based on security requirements
- **CORS**: Set specific origins, avoid wildcards in production
- **Password Policy**: Enforce strong passwords at application level
- **Rate Limiting**: Add rate limiting middleware for production

## Troubleshooting

### Database Connection Issues

```bash
# Test database connection
psql -U postgres -h localhost -d portfolio

# Check if tables exist
\dt

# View environment
echo $DATABASE_URL
```

### Authentication Not Working

- Verify `JWT_SECRET` and `JWT_REFRESH_SECRET` are set
- Check token hasn't expired (default access: 15m, refresh: 7d)
- Ensure Bearer token format: `Authorization: Bearer <token>`

### Write Operations Failing

- In JSON mode: Expected behavior, switch to PostgreSQL
- In PG mode: Check authentication token is valid
- Verify HTTP method is correct (POST/PUT/DELETE)

## License

MIT