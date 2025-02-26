-- name: AddOrUpdateSplit :one
INSERT INTO splits
(id, race_id, event_id, split_name, split_type, distance_from_start, time_reader_id, min_time_sec, max_time_sec, min_lap_time_sec)
VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
ON CONFLICT (race_id, event_id, id)
DO UPDATE
SET split_name=EXCLUDED.split_name, split_type=EXCLUDED. split_type, distance_from_start=EXCLUDED.distance_from_start, time_reader_id=EXCLUDED.time_reader_id, min_time_sec=EXCLUDED.min_time_sec, max_time_sec=EXCLUDED.max_time_sec, min_lap_time_sec=EXCLUDED.min_lap_time_sec
RETURNING *;

-- name: DeleteSplitByID :exec
DELETE FROM splits
WHERE id=$1;

-- name: GetAllSplitsForEvent :many
SELECT 
(id, race_id, event_id, split_name, split_type, distance_from_start, time_reader_id, min_time_sec, max_time_sec, min_lap_time_sec)
FROM splits
WHERE event_id=$1
ORDER BY distance_from_start ASC;