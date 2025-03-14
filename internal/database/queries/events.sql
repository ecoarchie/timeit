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

-- name: GetEventsForRace :many
SELECT id, race_id, event_name, distance_in_meters, event_date
FROM events
WHERE race_id=$1
ORDER BY event_date ASC;

-- name: GetEventIDsWithWavesStarted :many
select distinct e.id
from events e
join waves w on w.event_id = e.id
where e.race_id = $1 and w.is_launched is true;
