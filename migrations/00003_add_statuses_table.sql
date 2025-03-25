-- +goose Up
-- +goose StatementBegin
CREATE TABLE statuses (
	status_id SERIAL PRIMARY KEY,
  status_full TEXT NOT NULL,
  status_code TEXT NOT NULL,
  can_get_rank BOOLEAN NOT NULL,
  can_assign_at_split BOOLEAN NOT NULL
);

INSERT INTO statuses (status_full, status_code, can_get_rank, can_assign_at_split) VALUES
  ('not yet started', 'NYS', false, false),
  ('running', 'RUN', true, true),
  ('finished', 'FIN', true, true),
  ('disqualified', 'DSQ', false, false),
  ('quarantine', 'QRT', false, false),
  ('pre-race withdrawal', 'DNS', false, false),
  ('withdrawn during race', 'DNF', false, false);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE statuses;
-- +goose StatementEnd