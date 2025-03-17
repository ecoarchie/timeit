-- +goose Up
-- +goose StatementBegin
CREATE TYPE rr_tod AS (
	reader_id uuid,
	tod timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TYPE rr_tod;
-- +goose StatementEnd