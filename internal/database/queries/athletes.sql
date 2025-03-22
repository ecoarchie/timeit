-- name: GetAthleteByID :one
SELECT a.id, a.race_id, a.first_name, a.last_name, a.gender, a.date_of_birth, a.phone, a.athlete_comments, ea.event_id, ea.wave_id, ea.category_id,
cb.bib, cb.chip
FROM athletes a
join event_athlete ea 
on ea.athlete_id = a.id
join chip_bib cb 
on cb.bib = ea.bib
where a.id = $1;

-- name: CreateOrUpdateAthlete :one
INSERT INTO athletes
(id, race_id, first_name, last_name, gender, date_of_birth, phone, athlete_comments)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (id)
DO UPDATE
SET first_name=EXCLUDED.first_name, last_name=EXCLUDED.last_name, gender=EXCLUDED.gender, date_of_birth=EXCLUDED.date_of_birth, phone=EXCLUDED.phone, athlete_comments=EXCLUDED.athlete_comments, updated_at=EXCLUDED.updated_at
RETURNING *;

-- name: CreateAthleteBulk :copyfrom
INSERT INTO athletes
(id, race_id, first_name, last_name, gender, date_of_birth, phone, athlete_comments)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: DeleteAthleteByID :exec
DELETE FROM athletes
WHERE id=$1;

-- name: DeleteAthletesWithRaceID :exec
DELETE FROM athletes
WHERE race_id=$1;

-- name: DeleteAthletesWithEventID :exec
DELETE from athletes a
WHERE a.id IN (
  SELECT ea.athlete_id 
  FROM event_athlete ea WHERE ea.event_id = $1);
