# Chat App

Real-time chat application built with Go, Next.js, PostgreSQL, and Redis.

---

## Stack

- **Backend** — Go, Gorilla WebSocket, Redis Stack, PostgreSQL + GORM
- **Frontend** — Next.js 14, Tailwind v4, TypeScript
- **Auth** — Google OAuth 2.0 + JWT
- **Deployment** — Render (backend), Vercel (frontend)

---

## Live

**Link** -  https://banterrr.vercel.app

## Environment Variables

Create a `.env` file in the project root:

```env
# Google OAuth
CLIENT_ID=        # from Google API Console
CLIENT_SECRET=
REDIRECT_URL=http://localhost:8080/auth/google/callback
TOKEN_SECRET=     # any random alphanumeric string
JWT_SECRET=       # any random string, min 32 chars

# Database
DB_HOST=127.0.0.1
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=chat
DB_SSLMODE=disable

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# Server
SERVER_PORT=8080
ENV=development

# Frontend
FRONTEND_URL=http://localhost:3000
```

---

## Running Locally

```bash
# Start backend, PostgreSQL, and Redis
docker compose up

# Start frontend (separate terminal)
cd chat-frontend && npm run dev
```

The backend runs on `:8080` and frontend on `:3000`.

---

## Database Migrations

Migrations live in `internal/database/migrations/` and run automatically on startup.

**Create a new migration**
```bash
migrate create -ext sql -dir internal/database/migrations -seq <name>
# example
migrate create -ext sql -dir internal/database/migrations -seq add_read_receipts
```

**Run manually**
```bash
migrate -path internal/database/migrations \
  -database "postgres://postgres:password@localhost:5432/chat?sslmode=disable" up
```

**Rollback**
```bash
# One step
migrate -path internal/database/migrations \
  -database "postgres://postgres:password@localhost:5432/chat?sslmode=disable" down 1

# All
migrate -path internal/database/migrations \
  -database "postgres://postgres:password@localhost:5432/chat?sslmode=disable" down
```

**Check version**
```bash
migrate -path internal/database/migrations \
  -database "postgres://postgres:password@localhost:5432/chat?sslmode=disable" version
```

**Force version** (emergency only — use if migration state is dirty)
```bash
migrate -path internal/database/migrations \
  -database "postgres://postgres:password@localhost:5432/chat?sslmode=disable" force 1
```

---
### Delete Database

```bash
# allow script to be executable
chmod +x scripts/reset-db-dev.sh
```

```bash
# run the script
./scripts/reset-db-dev.sh
```

---

## HTTP Routes

| Method | Path | Description |
|--------|------|-------------|
| GET | `/health` | Health check |
| GET | `/auth/google-sso` | Initiate Google OAuth |
| GET | `/auth/google/callback` | OAuth callback |
| GET | `/auth/success` | Successful Login |