-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA "store";

CREATE TABLE "users" (
  "id" varchar PRIMARY KEY,
  "login" varchar,
  "password" varchar,
  "secret_key" varchar
);

CREATE TABLE "user_refresh_token" (
  "token_id" varchar PRIMARY KEY,
  "user_id" varchar,
  "created_at" timestamp without time zone default (now() at time zone 'utc'),
  "expired_at" timestamp
);

ALTER TABLE "user_refresh_token" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "user_refresh_token" IF EXISTS;
DROP TABLE "users" IF EXISTS;
-- +goose StatementEnd
