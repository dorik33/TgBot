-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS portfolio (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    token TEXT NOT NULL,
    amount FLOAT,
    price FLOAT,
    created_at TIMESTAMP DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE portfolio;
-- +goose StatementEnd
