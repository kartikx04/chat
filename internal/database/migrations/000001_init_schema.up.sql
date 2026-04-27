CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    auth_o_id   VARCHAR(255) NOT NULL UNIQUE,
    email       VARCHAR(255) NOT NULL UNIQUE,
    username    VARCHAR(255) NOT NULL UNIQUE,
    picture     VARCHAR(255),
    role        VARCHAR(50)  NOT NULL DEFAULT 'user',
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Queried on every login
CREATE INDEX idx_users_auth_o_id ON users(auth_o_id);
CREATE INDEX idx_users_email     ON users(email);

CREATE TABLE chats (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    from_id         UUID        NOT NULL,
    to_id           UUID        NOT NULL,
    message         TEXT        NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at_unix BIGINT      NOT NULL
);

-- Queried heavily for chat history
CREATE INDEX idx_chats_from_id    ON chats(from_id);
CREATE INDEX idx_chats_to_id      ON chats(to_id);
CREATE INDEX idx_chats_created_at ON chats(created_at);
-- Composite index for fetching conversation between two users
CREATE INDEX idx_chats_conversation ON chats(from_id, to_id, created_at_unix DESC);