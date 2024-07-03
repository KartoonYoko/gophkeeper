-- +goose Up
-- +goose StatementBegin
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
CREATE UNIQUE INDEX "users_login_idx" ON "users" ("login");

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "user_refresh_token";
DROP TABLE IF EXISTS "users";
-- +goose StatementEnd
