-- +goose Up
    CREATE TABLE user_auth_email (
        user_id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE
        email TEXT UNIQUE NOT NULL,
        password_hash TEXT NOT NULL
    );

-- +goose Up
DROP TABLE user_auth_email;
