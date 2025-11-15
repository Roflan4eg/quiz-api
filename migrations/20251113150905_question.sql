-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS question (
    id SERIAL PRIMARY KEY,
    text TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS question;
-- +goose StatementEnd
