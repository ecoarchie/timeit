-- +goose Up
-- +goose StatementBegin
ALTER TABLE event_athlete ADD COLUMN status_id INTEGER;

ALTER TABLE event_athlete ALTER COLUMN status_id SET DEFAULT 1;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE event_athlete DROP COLUMN IF EXISTS status_id;
-- +goose StatementEnd