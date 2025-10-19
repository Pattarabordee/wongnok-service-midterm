-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    IF NOT EXISTS cooking_durations (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        deleted_at TIMESTAMP
    );

INSERT INTO
    cooking_durations (name, created_at, updated_at)
VALUES
    ('5 - 10', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('11 - 30', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('31 - 60', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('60+', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS cooking_durations;

-- +goose StatementEnd