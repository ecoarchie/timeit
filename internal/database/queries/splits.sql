-- name: AddOrUpdateSplit :one
INSERT INTO splits
(id, race_id, event_id, split_name, split_type, distance_from_start, time_reader_id, min_time, max_time, min_lap_time)
VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
ON CONFLICT (race_id, event_id, id)
DO UPDATE
SET split_name=EXCLUDED.split_name, split_type=EXCLUDED. split_type, distance_from_start=EXCLUDED.distance_from_start, time_reader_id=EXCLUDED.time_reader_id, min_time=EXCLUDED.min_time, max_time=EXCLUDED.max_time, min_lap_time=EXCLUDED.min_lap_time
RETURNING *;

-- name: DeleteSplitByID :exec
DELETE FROM splits
WHERE id=$1;

-- name: GetAllSplitsForEvent :many
SELECT id, race_id, event_id, split_name, split_type, distance_from_start, time_reader_id, min_time, max_time, min_lap_time
FROM splits
WHERE event_id=$1
ORDER BY distance_from_start ASC;

-- name: GetAllSplitsForRace :many
SELECT id, race_id, event_id, split_name, split_type, distance_from_start, time_reader_id, min_time, max_time, min_lap_time
FROM splits
WHERE race_id=$1
ORDER BY distance_from_start ASC;