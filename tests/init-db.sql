-- users table
CREATE TABLE
    IF NOT EXISTS users (
        id VARCHAR(100) PRIMARY KEY,
        first_name VARCHAR(100) NOT NULL,
        last_name VARCHAR(100) NOT NULL,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        deleted_at TIMESTAMP
    );

INSERT INTO
    users (id, first_name, last_name, created_at, updated_at)
VALUES
    (
        '38fa4e9e-27de-42d5-a70f-9f01d41f32c2',
        'Demo',
        'Tester',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    );

-- difficulties table
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

-- cooking_durations table
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

-- food_recipes table
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
        user_id VARCHAR(100) REFERENCES users,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        deleted_at TIMESTAMP
    );

INSERT INTO
    food_recipes (
        name,
        description,
        ingredient,
        instruction,
        cooking_duration_id,
        difficulty_id,
        user_id,
        created_at,
        updated_at
    )
VALUES
    (
        'Omlet',
        'Eggs fried?',
        'Eggs',
        'Cooking',
        1,
        1,
        '38fa4e9e-27de-42d5-a70f-9f01d41f32c2',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    );

-- ratings table
CREATE TABLE
    IF NOT EXISTS ratings (
        id SERIAL PRIMARY KEY,
        food_recipe_id INT NOT NULL REFERENCES food_recipes,
        score INT NOT NULL,
        user_id VARCHAR(100) NOT NULL REFERENCES users,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        deleted_at TIMESTAMP
    );

INSERT INTO
    ratings (
        food_recipe_id,
        score,
        user_id,
        created_at,
        updated_at
    )
VALUES
    (
        1,
        5,
        '38fa4e9e-27de-42d5-a70f-9f01d41f32c2',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        1,
        3,
        '38fa4e9e-27de-42d5-a70f-9f01d41f32c2',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    );