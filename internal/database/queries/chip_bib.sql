-- name: AddChipBib :one
INSERT INTO chip_bib
(race_id, event_id, chip, bib)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: DeleteChipBib :exec
DELETE FROM chip_bib
WHERE race_id=$1 AND chip=$2 and bib=$3;