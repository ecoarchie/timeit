// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: races.sql

package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const createRace = `-- name: CreateRace :one
INSERT INTO races (id, name, race_date, timezone) VALUES ($1, $2, $3, $4)
RETURNING id, name, race_date, timezone
`

type CreateRaceParams struct {
	ID       uuid.UUID
	Name     string
	RaceDate pgtype.Date
	Timezone string
}

func (q *Queries) CreateRace(ctx context.Context, arg CreateRaceParams) (Race, error) {
	row := q.db.QueryRow(ctx, createRace,
		arg.ID,
		arg.Name,
		arg.RaceDate,
		arg.Timezone,
	)
	var i Race
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.RaceDate,
		&i.Timezone,
	)
	return i, err
}
