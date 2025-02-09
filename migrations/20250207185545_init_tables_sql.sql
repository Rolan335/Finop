-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    balance INT NOT NULL
);

CREATE TABLE IF NOT EXISTS operations(
    id SERIAL PRIMARY KEY,
    type VARCHAR(255) NOT NULL
);

INSERT INTO operations (type) values ('deposit');
INSERT INTO operations (type) values ('send');

CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY,
    user_id INT NOT NULL,
    receiver_id INT,
    operation_id INT NOT NULL,
    amount INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (receiver_id) REFERENCES users(id),
    FOREIGN KEY (operation_id) REFERENCES operations(id)
);

CREATE INDEX idx_users_name on users USING hash (name);
CREATE INDEX idx_transactions_created_at on transactions(created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_users_name;
DROP INDEX IF EXISTS idx_created_at;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS operations;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
