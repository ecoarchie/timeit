-- +goose Up
-- +goose StatementBegin
ALTER TABLE athlete_split
ADD COLUMN is_manual BOOLEAN DEFAULT FALSE;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE athlete_split
DROP COLUMN IF EXISTS is_manual, 
-- +goose StatementEnd