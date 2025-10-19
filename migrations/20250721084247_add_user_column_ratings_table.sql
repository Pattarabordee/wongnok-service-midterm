-- +goose Up
-- +goose StatementBegin
ALTER TABLE ratings ADD IF NOT EXISTS user_id VARCHAR(100)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE ratings DROP COLUMN user_id
-- +goose StatementEnd