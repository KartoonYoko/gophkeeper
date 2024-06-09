-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS data_store (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    description TEXT,
    data_type TEXT,
    hash TEXT,
    modification_timestamp INTEGER,
    is_deleted INTEGER DEFAULT 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS data_store;
-- +goose StatementEnd
