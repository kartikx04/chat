# .env config

CLIENT_ID= #get from google api console
CLIENT_SECRET=
REDIRECT_URL=http://localhost:8007/callback
TOKEN_SECRET= #any alphanumeric string

DB_HOST="127.0.0.1"
DB_PORT=5432
DB_USER="postgres"
DB_PASSWORD="password"
DB_NAME="chat"
DB_SSLMODE="disable"

REDIS_HOST=
REDIS_PASSWORD=

# migrate database
migrate -path ./database/migrations -database "postgresql://postgres:password@localhost:5432/chat?sslmode=disable" up