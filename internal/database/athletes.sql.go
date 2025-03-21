// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: athletes.sql

package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateAthleteBulkParams struct {
	ID              uuid.UUID
	RaceID          uuid.UUID
	FirstName       pgtype.Text
	LastName        pgtype.Text
	Gender          CategoryGender
	DateOfBirth     pgtype.Date
	Phone           pgtype.Text
	AthleteComments pgtype.Text
}

const createOrUpdateAthlete = `-- name: CreateOrUpdateAthlete :one
INSERT INTO athletes
(id, race_id, first_name, last_name, gender, date_of_birth, phone, athlete_comments)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (id)
DO UPDATE
SET first_name=EXCLUDED.first_name, last_name=EXCLUDED.last_name, gender=EXCLUDED.gender, date_of_birth=EXCLUDED.date_of_birth, phone=EXCLUDED.phone, athlete_comments=EXCLUDED.athlete_comments, updated_at=EXCLUDED.updated_at
RETURNING id, race_id, first_name, last_name, gender, date_of_birth, phone, athlete_comments, created_at, updated_at
`

type CreateOrUpdateAthleteParams struct {
	ID              uuid.UUID
	RaceID          uuid.UUID
	FirstName       pgtype.Text
	LastName        pgtype.Text
	Gender          CategoryGender
	DateOfBirth     pgtype.Date
	Phone           pgtype.Text
	AthleteComments pgtype.Text
}

func (q *Queries) CreateOrUpdateAthlete(ctx context.Context, arg CreateOrUpdateAthleteParams) (Athlete, error) {
	row := q.db.QueryRow(ctx, createOrUpdateAthlete,
		arg.ID,
		arg.RaceID,
		arg.FirstName,
		arg.LastName,
		arg.Gender,
		arg.DateOfBirth,
		arg.Phone,
		arg.AthleteComments,
	)
	var i Athlete
	err := row.Scan(
		&i.ID,
		&i.RaceID,
		&i.FirstName,
		&i.LastName,
		&i.Gender,
		&i.DateOfBirth,
		&i.Phone,
		&i.AthleteComments,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteAthleteByID = `-- name: DeleteAthleteByID :exec
DELETE FROM athletes
WHERE id=$1
`

func (q *Queries) DeleteAthleteByID(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteAthleteByID, id)
	return err
}

const deleteAthletesWithEventID = `-- name: DeleteAthletesWithEventID :exec
DELETE from athletes a
WHERE a.id IN (
  SELECT ea.athlete_id 
  FROM event_athlete ea WHERE ea.event_id = $1)
`

func (q *Queries) DeleteAthletesWithEventID(ctx context.Context, eventID uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteAthletesWithEventID, eventID)
	return err
}

const deleteAthletesWithRaceID = `-- name: DeleteAthletesWithRaceID :exec
DELETE FROM athletes
WHERE race_id=$1
`

func (q *Queries) DeleteAthletesWithRaceID(ctx context.Context, raceID uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteAthletesWithRaceID, raceID)
	return err
}

const getAthleteByID = `-- name: GetAthleteByID :one
SELECT a.id, a.race_id, a.first_name, a.last_name, a.gender, a.date_of_birth, a.phone, a.athlete_comments, ea.event_id, ea.wave_id, ea.category_id,
cb.bib, cb.chip
FROM athletes a
join event_athlete ea 
on ea.athlete_id = a.id
join chip_bib cb 
on cb.bib = ea.bib
where a.id = $1
`

type GetAthleteByIDRow struct {
	ID              uuid.UUID
	RaceID          uuid.UUID
	FirstName       pgtype.Text
	LastName        pgtype.Text
	Gender          CategoryGender
	DateOfBirth     pgtype.Date
	Phone           pgtype.Text
	AthleteComments pgtype.Text
	EventID         uuid.UUID
	WaveID          uuid.UUID
	CategoryID      uuid.NullUUID
	Bib             int32
	Chip            int32
}

func (q *Queries) GetAthleteByID(ctx context.Context, id uuid.UUID) (GetAthleteByIDRow, error) {
	row := q.db.QueryRow(ctx, getAthleteByID, id)
	var i GetAthleteByIDRow
	err := row.Scan(
		&i.ID,
		&i.RaceID,
		&i.FirstName,
		&i.LastName,
		&i.Gender,
		&i.DateOfBirth,
		&i.Phone,
		&i.AthleteComments,
		&i.EventID,
		&i.WaveID,
		&i.CategoryID,
		&i.Bib,
		&i.Chip,
	)
	return i, err
}
