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

**Inside Docker**
```bash
docker exec -it chat_app sh scripts/migrate.sh up
docker exec -it chat_app sh scripts/migrate.sh down 1
docker exec -it chat_app sh scripts/migrate.sh version
```

---

## Useful Commands

**Access containers**
```bash
# App shell
docker exec -it chat_app sh

# Postgres shell
docker exec -it chat_postgres psql -U postgres -d chat
```

**Test WebSocket manually**
```bash
npm install -g wscat
wscat -c ws://localhost:8080/ws
```

Open two terminals and paste:

Terminal 1 (user 1)
```json
{"type":"bootup","user":"user1","user_id":"11111111-1111-1111-1111-111111111111"}
{"type":"message","chat":{"from_id":"11111111-1111-1111-1111-111111111111","to_id":"22222222-2222-2222-2222-222222222222","message":"hello"}}
```

Terminal 2 (user 2)
```json
{"type":"bootup","user":"user2","user_id":"22222222-2222-2222-2222-222222222222"}
```

---

## HTTP Routes

| Method | Path | Description |
|--------|------|-------------|
| GET | `/health` | Health check |
| GET | `/google-sso` | Initiate Google OAuth |
| GET | `/auth/google/callback` | OAuth callback |
| GET | `/me` | Get current user from session |
| GET | `/contacts?id=<uuid>` | Fetch contact list |
| GET | `/chat-history?id=<uuid>&contact=<uuid>` | Fetch chat history |
| POST | `/add-contact?id=<uuid>&contact=<username>` | Add a contact |
| WS | `/ws` | WebSocket connection |