-- name: AddRace :one
INSERT INTO races (id, name, timezone) VALUES ($1, $2, $3)
ON CONFLICT (id) DO UPDATE
SET name=EXCLUDED.name, timezone=EXCLUDED.timezone
RETURNING *;

-- name: DeleteRace :exec
DELETE FROM races
WHERE id=$1;