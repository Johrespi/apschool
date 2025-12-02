-- +goose Up
CREATE TABLE IF NOT EXISTS user_auth_email (
    user_id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    password_hash TEXT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS user_auth_email;
