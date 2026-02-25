-- +goose Up
-- +goose StatementBegin

DROP TABLE IF EXISTS messages CASCADE;

DROP TABLE IF EXISTS chat_users CASCADE;

DROP TABLE IF EXISTS chats CASCADE;

CREATE TABLE chats (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE chat_users (
    chat_id BIGINT NOT NULL REFERENCES chats (id) ON DELETE CASCADE,
    username TEXT NOT NULL,
    PRIMARY KEY (chat_id, username)
);

CREATE TABLE messages (
    id BIGSERIAL PRIMARY KEY,
    chat_id BIGINT NOT NULL REFERENCES chats (id) ON DELETE CASCADE,
    sender TEXT NOT NULL,
    text TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX messages_chat_created_at_idx ON messages (chat_id, created_at);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS messages_chat_created_at_idx;

DROP TABLE IF EXISTS messages;

DROP TABLE IF EXISTS chat_users;

DROP TABLE IF EXISTS chats;

-- +goose StatementEnd