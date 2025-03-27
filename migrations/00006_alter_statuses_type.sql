-- +goose Up
-- +goose StatementBegin
ALTER TABLE statuses
ALTER COLUMN status_id DROP DEFAULT;
ALTER TABLE statuses
ALTER COLUMN status_id TYPE SMALLINT;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE statuses
ALTER COLUMN status_id TYPE INTEGER;
-- +goose StatementEnd