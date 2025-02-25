// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: races.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const addRace = `-- name: AddRace :one
INSERT INTO races (id, name, timezone) VALUES ($1, $2, $3)
ON CONFLICT (id) DO UPDATE
SET name=EXCLUDED.name, timezone=EXCLUDED.timezone
RETURNING id, name, timezone
`

type AddRaceParams struct {
	ID       uuid.UUID
	Name     string
	Timezone string
}

func (q *Queries) AddRace(ctx context.Context, arg AddRaceParams) (Race, error) {
	row := q.db.QueryRow(ctx, addRace, arg.ID, arg.Name, arg.Timezone)
	var i Race
	err := row.Scan(&i.ID, &i.Name, &i.Timezone)
	return i, err
}

const deleteRace = `-- name: DeleteRace :exec
DELETE FROM races
WHERE id=$1
`

func (q *Queries) DeleteRace(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteRace, id)
	return err
}
