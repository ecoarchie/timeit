// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: time_readers.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const addOrUpdateTimeReader = `-- name: AddOrUpdateTimeReader :one
INSERT INTO time_readers
(id, race_id, reader_name)
VALUES($1, $2, $3)
ON CONFLICT (race_id, reader_name) DO UPDATE
SET race_id=excluded.race_id, reader_name=excluded.reader_name
RETURNING id, race_id, reader_name
`

type AddOrUpdateTimeReaderParams struct {
	ID         uuid.UUID
	RaceID     uuid.UUID
	ReaderName string
}

func (q *Queries) AddOrUpdateTimeReader(ctx context.Context, arg AddOrUpdateTimeReaderParams) (TimeReader, error) {
	row := q.db.QueryRow(ctx, addOrUpdateTimeReader, arg.ID, arg.RaceID, arg.ReaderName)
	var i TimeReader
	err := row.Scan(&i.ID, &i.RaceID, &i.ReaderName)
	return i, err
}

const deleteTimeReaderByID = `-- name: DeleteTimeReaderByID :exec
DELETE FROM time_readers
WHERE id=$1
`

func (q *Queries) DeleteTimeReaderByID(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteTimeReaderByID, id)
	return err
}

const getTimeReadersForRace = `-- name: GetTimeReadersForRace :many
SELECT id, race_id, reader_name
FROM time_readers
WHERE race_id=$1
`

func (q *Queries) GetTimeReadersForRace(ctx context.Context, raceID uuid.UUID) ([]TimeReader, error) {
	rows, err := q.db.Query(ctx, getTimeReadersForRace, raceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []TimeReader
	for rows.Next() {
		var i TimeReader
		if err := rows.Scan(&i.ID, &i.RaceID, &i.ReaderName); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
