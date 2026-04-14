CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS chats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    from_id UUID NOT NULL,
    to_id UUID NOT NULL,
    message TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at_unix BIGINT NOT NULL
);

CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    from_id UUID NOT NULL,
    to_id UUID NOT NULL,

    message TEXT NOT NULL,

    -- readable timestamp
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    -- for sorting / fast queries
    created_at_unix BIGINT NOT NULL
);

-- Indexes for performance
CREATE INDEX idx_chats_from_id ON chats(from_id);
CREATE INDEX idx_chats_to_id ON chats(to_id);
CREATE INDEX idx_chats_created_at ON chats(created_at);