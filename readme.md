# .env config

#Google Authorization Configuration
CLIENT_ID= #get from google api console
CLIENT_SECRET=
REDIRECT_URL=http://localhost:8007/callback
TOKEN_SECRET= #any alphanumeric string

#Database Configuration
DB_HOST="127.0.0.1"
DB_PORT=5432
DB_USER="postgres"
DB_PASSWORD="password"
DB_NAME="chat"
DB_SSLMODE="disable"

#Server Configuration
SERVER_PORT=8007
ENV=development

#Redis Configuration
REDIS_HOST=
REDIS_PASSWORD=

# migrate database
migrate -path ./internal/database/migrations -database "postgresql://postgres:password@localhost:5432/chat?sslmode=disable" up

npm install -g wscat

docker exec -it chat_app sh

docker exec -it chat_postgres psql -U postgres -d chat 

# Check if user chat working via terminal
wscat -c ws://localhost:8081/ws

terminal #1

> {"type":"bootup","user":"user1","user_id":"11111111-1111-1111-1111-111111111111"}
> {"type":"message","chat":{"from_id":"11111111-1111-1111-1111-111111111111","to_id":"22222222-2222-2222-2222-222222222222","message":"hello bro"}}

terminal #2

> {"type":"bootup","user":"user2","user_id":"22222222-2222-2222-2222-222222222222"}