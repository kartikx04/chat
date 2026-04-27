DROP INDEX IF EXISTS idx_chats_conversation;
DROP INDEX IF EXISTS idx_chats_created_at;
DROP INDEX IF EXISTS idx_chats_to_id;
DROP INDEX IF EXISTS idx_chats_from_id;
DROP TABLE IF EXISTS chats;

DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_auth_o_id;
DROP TABLE IF EXISTS users;

DROP EXTENSION IF EXISTS "pgcrypto";