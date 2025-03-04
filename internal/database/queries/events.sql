-- name: DeleteEvent :exec
DELETE FROM events
WHERE id=$1;

-- name: GetEventByID :one
SELECT id, race_id, event_name, distance_in_meters, event_date
FROM events
WHERE id=$1;

-- name: AddOrUpdateEvent :one
INSERT INTO events
(id, race_id, event_name, distance_in_meters, event_date)
VALUES($1, $2, $3, $4, $5)
ON CONFLICT (race_id, id) DO UPDATE
SET id=excluded.id, event_name=excluded.event_name, distance_in_meters=excluded.distance_in_meters, event_date=excluded.event_date
RETURNING *;

-- name: GetAllEventsForRace :many
SELECT id, race_id, event_name, distance_in_meters, event_date
FROM events
WHERE race_id=$1
ORDER BY event_date ASC;
