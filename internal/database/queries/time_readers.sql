-- name: GetTimeReadersForRace :many
SELECT id, race_id, reader_name
FROM time_readers
WHERE race_id=$1;

-- name: AddOrUpdateTimeReader :one
INSERT INTO time_readers
(id, race_id, reader_name)
VALUES($1, $2, $3)
ON CONFLICT (race_id, reader_name) DO UPDATE
SET race_id=excluded.race_id, reader_name=excluded.reader_name
RETURNING *;

-- name: DeleteTimeReaderByID :exec
DELETE FROM time_readers
WHERE id=$1;