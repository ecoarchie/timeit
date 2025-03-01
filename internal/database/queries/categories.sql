-- name: AddOrUpdateCategory :one
INSERT INTO categories
(id, race_id, event_id, category_name, gender, age_from, date_from, age_to, date_to)
VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)
ON CONFLICT (id)
DO UPDATE
SET category_name=EXCLUDED. category_name, gender=EXCLUDED.gender, age_from=EXCLUDED.age_from, date_from=EXCLUDED.date_from, age_to=EXCLUDED.age_to, date_to=EXCLUDED.date_to
RETURNING *;

-- name: DeleteCategoryByID :exec
DELETE FROM categories
WHERE id=$1;

-- name: GetCategoriesForEvent :many
SELECT (id, race_id, event_id, category_name, gender, age_from, date_from, age_to, date_to)
FROM categories
WHERE event_id=$1
ORDER BY age_from ASC;

-- name: GetCategoryForAthlete :one
SELECT id, race_id, event_id, category_name, gender, age_from, date_from, age_to, date_to
FROM categories
WHERE 
event_id = $1 
AND gender = $2 
AND $3 BETWEEN (date_to - (age_to || ' years')::INTERVAL) AND (date_from - (age_from || ' years')::INTERVAL);