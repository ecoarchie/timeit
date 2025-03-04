-- name: AddRace :one
INSERT INTO races (id, race_name, timezone) VALUES ($1, $2, $3)
ON CONFLICT (id) DO UPDATE
SET race_name=EXCLUDED.race_name, timezone=EXCLUDED.timezone
RETURNING *;

-- name: DeleteRace :exec
DELETE FROM races
WHERE id=$1;

-- name: GetRaceInfo :one
SELECT id, race_name, timezone
FROM races
WHERE id = $1;

-- name: GetRaces :many
SELECT id, race_name, timezone FROM races;