-- name: AddOrUpdateWave :one
INSERT INTO waves
(id, race_id, event_id, "name", start_time, is_launched)
VALUES($1, $2, $3, $4, $5, $6)
ON CONFLICT (race_id, event_id, id)
DO UPDATE
SET "name"=EXCLUDED."name", start_time=EXCLUDED.start_time, is_launched=EXCLUDED.is_launched
RETURNING *;

-- name: DeleteWaveByID :exec
DELETE FROM waves
WHERE id=$1;

-- name: GetAllWavesForEvent :many
SELECT id, race_id, event_id, "name", start_time, is_launched
FROM waves
WHERE event_id=$1
ORDER BY start_time ASC;

-- name: StartWave :exec
UPDATE waves
SET is_launched=true
WHERE id=$1; 