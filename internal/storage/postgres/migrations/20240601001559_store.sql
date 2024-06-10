-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA "store";

CREATE TABLE "store"."data" (
  "id" varchar PRIMARY KEY,
  "user_id" varchar,
  "binary_id" varchar,
  "description" varchar,
  "is_deleted" boolean,
  "hash" varchar,
  "modification_timestamp" bigint,
  "data_type" varchar
);

ALTER TABLE "store"."data" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
CREATE UNIQUE INDEX IF NOT EXISTS idx_store_data_binary_id ON "store"."data"("binary_id");

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "store"."data";
DROP SCHEMA IF EXISTS "store";
-- +goose StatementEnd
