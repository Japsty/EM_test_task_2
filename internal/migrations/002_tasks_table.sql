-- +goose Up
CREATE TABLE IF NOT EXISTS tasks
(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    user_id INT NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP
    );

-- +goose Down
DROP TABLE IF EXISTS tasks;