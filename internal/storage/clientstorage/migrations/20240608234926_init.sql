-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS data_store (
    id INTEGER PRIMARY KEY,
    user_id TEXT NOT NULL,
    description TEXT,
    data_type TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS data_store;
-- +goose StatementEnd
