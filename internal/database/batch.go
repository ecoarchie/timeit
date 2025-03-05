// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: batch.go

package database

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrBatchAlreadyClosed = errors.New("batch already closed")
)

const createAthlete = `-- name: CreateAthlete :batchexec
INSERT INTO athletes
(id, race_id, first_name, last_name, gender, date_of_birth, phone, athlete_comments)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
`

type CreateAthleteBatchResults struct {
	br     pgx.BatchResults
	tot    int
	closed bool
}

type CreateAthleteParams struct {
	ID              uuid.UUID
	RaceID          uuid.UUID
	FirstName       pgtype.Text
	LastName        pgtype.Text
	Gender          CategoryGender
	DateOfBirth     pgtype.Date
	Phone           pgtype.Text
	AthleteComments pgtype.Text
}

func (q *Queries) CreateAthlete(ctx context.Context, arg []CreateAthleteParams) *CreateAthleteBatchResults {
	batch := &pgx.Batch{}
	for _, a := range arg {
		vals := []interface{}{
			a.ID,
			a.RaceID,
			a.FirstName,
			a.LastName,
			a.Gender,
			a.DateOfBirth,
			a.Phone,
			a.AthleteComments,
		}
		batch.Queue(createAthlete, vals...)
	}
	br := q.db.SendBatch(ctx, batch)
	return &CreateAthleteBatchResults{br, len(arg), false}
}

func (b *CreateAthleteBatchResults) Exec(f func(int, error)) {
	defer b.br.Close()
	for t := 0; t < b.tot; t++ {
		if b.closed {
			if f != nil {
				f(t, ErrBatchAlreadyClosed)
			}
			continue
		}
		_, err := b.br.Exec()
		if f != nil {
			f(t, err)
		}
	}
}

func (b *CreateAthleteBatchResults) Close() error {
	b.closed = true
	return b.br.Close()
}
