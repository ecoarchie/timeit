// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: event_athlete.sql

package database

import (
	"context"

	"github.com/ecoarchie/timeit/internal/entity"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const addEventAthlete = `-- name: AddEventAthlete :one
INSERT INTO event_athlete
(race_id, event_id, athlete_id, wave_id, category_id, bib)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (race_id, event_id, athlete_id)
DO UPDATE
SET event_id=EXCLUDED.event_id, wave_id=EXCLUDED.wave_id, category_id=EXCLUDED.category_id, bib=EXCLUDED.bib 
RETURNING race_id, event_id, athlete_id, wave_id, category_id, bib
`

type AddEventAthleteParams struct {
	RaceID     uuid.UUID
	EventID    uuid.UUID
	AthleteID  uuid.UUID
	WaveID     uuid.UUID
	CategoryID uuid.NullUUID
	Bib        int32
}

func (q *Queries) AddEventAthlete(ctx context.Context, arg AddEventAthleteParams) (EventAthlete, error) {
	row := q.db.QueryRow(ctx, addEventAthlete,
		arg.RaceID,
		arg.EventID,
		arg.AthleteID,
		arg.WaveID,
		arg.CategoryID,
		arg.Bib,
	)
	var i EventAthlete
	err := row.Scan(
		&i.RaceID,
		&i.EventID,
		&i.AthleteID,
		&i.WaveID,
		&i.CategoryID,
		&i.Bib,
	)
	return i, err
}

type AddEventAthleteBulkParams struct {
	RaceID     uuid.UUID
	EventID    uuid.UUID
	AthleteID  uuid.UUID
	WaveID     uuid.UUID
	CategoryID uuid.NullUUID
	Bib        int32
}

const getEventAthlete = `-- name: GetEventAthlete :one
SELECT race_id, event_id, athlete_id, wave_id, category_id, bib
FROM event_athlete
WHERE athlete_id=$1
`

func (q *Queries) GetEventAthlete(ctx context.Context, athleteID uuid.UUID) (EventAthlete, error) {
	row := q.db.QueryRow(ctx, getEventAthlete, athleteID)
	var i EventAthlete
	err := row.Scan(
		&i.RaceID,
		&i.EventID,
		&i.AthleteID,
		&i.WaveID,
		&i.CategoryID,
		&i.Bib,
	)
	return i, err
}

const getEventAthleteRecords = `-- name: GetEventAthleteRecords :many
select 
	ea.athlete_id,
	ea.bib,
	cb.chip,
	ea.category_id,
	a.gender,
	w.start_time as wave_start,
	(select array_agg(rr.tod order by rr.tod)::timestamp[]
	from reader_records rr
	where rr.race_id = ea.race_id and rr.chip = cb.chip and rr.can_use is true) as records,
	(select array_agg(tr.id order by rr.tod)::uuid[]
	from reader_records rr
	join time_readers tr on tr.reader_name = rr.reader_name and tr.race_id = rr.race_id
	where rr.race_id = ea.race_id and rr.chip = cb.chip and rr.can_use is true) as reader_ids
	from event_athlete ea
	join waves w on w.race_id = ea.race_id and w.event_id = ea.event_id  and w.id = ea.wave_id
	join chip_bib cb on cb.race_id = ea.race_id and cb.event_id = ea.event_id and cb.bib = ea.bib
	join athletes a on a.id = ea.athlete_id and a.race_id = ea.race_id
	where ea.race_id = $1 
		and ea.event_id = $2 
		and w.is_launched is true
`

type GetEventAthleteRecordsParams struct {
	RaceID  uuid.UUID
	EventID uuid.UUID
}

type GetEventAthleteRecordsRow struct {
	AthleteID  uuid.UUID
	Bib        int32
	Chip       int32
	CategoryID uuid.NullUUID
	Gender     CategoryGender
	WaveStart  pgtype.Timestamp
	Records    []pgtype.Timestamp
	ReaderIds  []uuid.UUID
}

func (q *Queries) GetEventAthleteRecords(ctx context.Context, arg GetEventAthleteRecordsParams) ([]GetEventAthleteRecordsRow, error) {
	rows, err := q.db.Query(ctx, getEventAthleteRecords, arg.RaceID, arg.EventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetEventAthleteRecordsRow
	for rows.Next() {
		var i GetEventAthleteRecordsRow
		if err := rows.Scan(
			&i.AthleteID,
			&i.Bib,
			&i.Chip,
			&i.CategoryID,
			&i.Gender,
			&i.WaveStart,
			&i.Records,
			&i.ReaderIds,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getEventAthleteRecordsC = `-- name: GetEventAthleteRecordsC :many
select 
	ea.athlete_id,
	ea.bib,
	cb.chip,
	ea.category_id,
	a.gender,
	w.start_time as wave_start,
	(select array_agg(row(tr.id, rr.tod)::rr_tod order by rr.tod)::rr_tod[]
	from
		reader_records rr
	join time_readers tr on
		tr.reader_name = rr.reader_name
		and tr.race_id = rr.race_id
	where
		rr.race_id = ea.race_id
		and rr.chip = cb.chip
		and rr.can_use is true) as rr_tod
	from event_athlete ea
	join waves w on w.race_id = ea.race_id and w.event_id = ea.event_id  and w.id = ea.wave_id
	join chip_bib cb on cb.race_id = ea.race_id and cb.event_id = ea.event_id and cb.bib = ea.bib
	join athletes a on a.id = ea.athlete_id and a.race_id = ea.race_id
	where ea.race_id = $1 
		and ea.event_id = $2 
		and w.is_launched is true
`

type GetEventAthleteRecordsCParams struct {
	RaceID  uuid.UUID
	EventID uuid.UUID
}

type GetEventAthleteRecordsCRow struct {
	AthleteID  uuid.UUID
	Bib        int32
	Chip       int32
	CategoryID uuid.NullUUID
	Gender     CategoryGender
	WaveStart  pgtype.Timestamp
	RrTod      []entity.RecordTOD
}

func (q *Queries) GetEventAthleteRecordsC(ctx context.Context, arg GetEventAthleteRecordsCParams) ([]GetEventAthleteRecordsCRow, error) {
	rows, err := q.db.Query(ctx, getEventAthleteRecordsC, arg.RaceID, arg.EventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetEventAthleteRecordsCRow
	for rows.Next() {
		var i GetEventAthleteRecordsCRow
		if err := rows.Scan(
			&i.AthleteID,
			&i.Bib,
			&i.Chip,
			&i.CategoryID,
			&i.Gender,
			&i.WaveStart,
			&i.RrTod,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getEventAthleteRecordsJSON = `-- name: GetEventAthleteRecordsJSON :many
select 
	ea.race_id,
	ea.event_id,
	ea.wave_id,
	ea.athlete_id,
	ea.bib,
	cb.chip,
	w.start_time as wave_start,
	(select json_agg(json_build_object('tod', rr.tod, 'reader', tr.id)) from reader_records rr 
	join time_readers tr on tr.reader_name = rr.reader_name and tr.race_id = rr.race_id
	where rr.race_id = ea.race_id and rr.chip = cb.chip and rr.can_use is true) as j_recs
from event_athlete ea
join waves w on w.race_id = ea.race_id and w.event_id = ea.event_id  and w.id = ea.wave_id
join chip_bib cb on cb.race_id = ea.race_id and cb.event_id = ea.event_id and cb.bib = ea.bib
where ea.race_id = $1
	and ea.event_id = $2
	and w.is_launched is true
`

type GetEventAthleteRecordsJSONParams struct {
	RaceID  uuid.UUID
	EventID uuid.UUID
}

type GetEventAthleteRecordsJSONRow struct {
	RaceID    uuid.UUID
	EventID   uuid.UUID
	WaveID    uuid.UUID
	AthleteID uuid.UUID
	Bib       int32
	Chip      int32
	WaveStart pgtype.Timestamp
	JRecs     []byte
}

func (q *Queries) GetEventAthleteRecordsJSON(ctx context.Context, arg GetEventAthleteRecordsJSONParams) ([]GetEventAthleteRecordsJSONRow, error) {
	rows, err := q.db.Query(ctx, getEventAthleteRecordsJSON, arg.RaceID, arg.EventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetEventAthleteRecordsJSONRow
	for rows.Next() {
		var i GetEventAthleteRecordsJSONRow
		if err := rows.Scan(
			&i.RaceID,
			&i.EventID,
			&i.WaveID,
			&i.AthleteID,
			&i.Bib,
			&i.Chip,
			&i.WaveStart,
			&i.JRecs,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
