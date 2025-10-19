-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    IF NOT EXISTS difficulties (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        deleted_at TIMESTAMP
    );

INSERT INTO
    difficulties (name, created_at, updated_at)
VALUES
    ('Easy', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('Medium', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('Hard', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS difficulties;

-- +goose StatementEnd