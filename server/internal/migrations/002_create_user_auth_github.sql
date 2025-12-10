-- +goose Up
CREATE TABLE user_auth_github (
    user_id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    github_id BIGINT UNIQUE NOT NULL
);

-- +goose Down
DROP TABLE user_auth_github;
