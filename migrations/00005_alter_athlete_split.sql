-- +goose Up
-- +goose StatementBegin
ALTER TABLE athlete_split
ADD COLUMN gun_rank_gender INTEGER,
ADD COLUMN gun_rank_category INTEGER,
ADD COLUMN gun_rank_overall INTEGER,
ADD COLUMN net_rank_gender INTEGER,
ADD COLUMN net_rank_category INTEGER,
ADD COLUMN net_rank_overall INTEGER;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE athlete_split
DROP COLUMN IF EXISTS gun_rank_gender, 
DROP COLUMN IF EXITST gun_rank_category, 
DROP COLUMN IF EXITST gun_rank_overall,
DROP COLUMN IF EXITST net_rank_gender,
DROP COLUMN IF EXITST net_rank_category,
DROP COLUMN IF EXITST net_rank_overall;
-- +goose StatementEnd