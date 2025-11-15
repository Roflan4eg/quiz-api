-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS answer (
    id SERIAL PRIMARY KEY,
    question_id INTEGER NOT NULL REFERENCES question(id) ON DELETE CASCADE,
    user_id VARCHAR(36) NOT NULL,
    text TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS answer;
-- +goose StatementEnd
