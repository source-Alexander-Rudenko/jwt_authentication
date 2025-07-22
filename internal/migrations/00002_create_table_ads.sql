-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "ADS" (
                       id           UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
                       author_id    UUID        NOT NULL REFERENCES "USER"(id) ON DELETE CASCADE,
                       title        TEXT        NOT NULL,
                       description  TEXT        NOT NULL,
                       price        NUMERIC(10,2) NOT NULL CHECK (price >= 0),
                       image_key    TEXT        NOT NULL,
                       created_at   TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
                       updated_at   TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "ADS";
-- +goose StatementEnd
