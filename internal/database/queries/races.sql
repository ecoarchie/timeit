-- name: CreateRace :one
INSERT INTO races (id, name, race_date, timezone) VALUES ($1, $2, $3, $4)
RETURNING *;