-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    IF NOT EXISTS ratings (
        id SERIAL PRIMARY KEY,
        score INT NOT NULL,
        food_recipe_id INT NOT NULL REFERENCES food_recipes,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        deleted_at TIMESTAMP
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS ratings;
-- +goose StatementEnd