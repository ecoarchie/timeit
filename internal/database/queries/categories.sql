-- name: AddOrUpdateCategory :one
INSERT INTO categories
(id, race_id, event_id, "name", gender, from_age, from_race_date, to_age, to_race_date)
VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)
ON CONFLICT (id)
DO UPDATE
SET "name"=EXCLUDED."name", gender=EXCLUDED.gender, from_age=EXCLUDED.from_age, from_race_date=EXCLUDED.from_race_date, to_age=EXCLUDED.to_age, to_race_date=EXCLUDED.to_race_date
RETURNING *;

-- name: DeleteCategoryByID :exec
DELETE FROM categories
WHERE id=$1;

-- name: GetCategoriesForEvent :many
SELECT (id, race_id, event_id, "name", gender, from_age, from_race_date, to_age, to_race_date)
FROM categories
WHERE event_id=$1
ORDER BY from_age ASC;