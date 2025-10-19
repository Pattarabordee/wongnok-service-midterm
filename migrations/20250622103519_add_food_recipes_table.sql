-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    IF NOT EXISTS food_recipes (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        description TEXT NOT NULL,
        ingredient TEXT NOT NULL,
        instruction TEXT NOT NULL,
        image_url TEXT NULL,
        cooking_duration_id INT NOT NULL REFERENCES cooking_durations,
        difficulty_id INT NOT NULL REFERENCES difficulties,
        created_at TIMESTAMP,
        updated_at TIMESTAMP,
        deleted_at TIMESTAMP
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS food_recipes;

-- +goose StatementEnd