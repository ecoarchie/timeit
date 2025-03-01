-- name: AddChipBib :one
INSERT INTO chip_bib
(race_id, event_id, chip, bib)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: DeleteChipBib :exec
DELETE FROM chip_bib
WHERE race_id=$1 AND chip=$2 and bib=$3;

-- name: DeleteChipBibWithRaceID :exec
DELETE FROM chip_bib
WHERE race_id=$1;

-- name: DeleteChipBibWithEventID :exec
DELETE FROM chip_bib
WHERE race_id=$1 and event_id=$2;