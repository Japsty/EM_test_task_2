-- +goose Up
CREATE TABLE IF NOT EXISTS tasks
(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    user_id INT NOT NULL,
    start_time TIMESTAMP,
    end_time TIMESTAMP,
    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id)
    );

-- +goose Down
DROP TABLE IF EXISTS tasks;