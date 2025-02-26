-- name: AddEventAthlete :one
INSERT INTO event_athlete
(race_id, event_id, athlete_id, wave_id, category_id, bib)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetEventAthlete :one
SELECT race_id, event_id, athlete_id, wave_id, category_id, bib
FROM event_athlete
WHERE athlete_id=$1;

-- name: DeleteEventAthlete :exec
DELETE FROM event_athlete
WHERE race_id = $1 AND athlete_id = $2;