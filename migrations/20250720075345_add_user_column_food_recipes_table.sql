-- +goose Up
-- +goose StatementBegin
ALTER TABLE food_recipes ADD IF NOT EXISTS user_id VARCHAR(100) REFERENCES users;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE food_recipes
DROP COLUMN user_id;

-- +goose StatementEnd