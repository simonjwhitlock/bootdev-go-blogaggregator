-- +goose Up
ALTER TABLE feeds
ADD last_fetched_at TIMESTAMP;

-- +goose Down
ALTER TABLE feeds
DELETE IF EXISTS last_fetched_at;