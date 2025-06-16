-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS subscriptions (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    token TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE subscriptions;
-- +goose StatementEnd
