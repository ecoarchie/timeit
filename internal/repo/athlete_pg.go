package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/ecoarchie/timeit/internal/database"
	"github.com/ecoarchie/timeit/internal/entity"
	"github.com/ecoarchie/timeit/pkg/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type ParticipantQuery interface {
	GetAthleteByID(ctx context.Context, id uuid.UUID) (database.GetAthleteByIDRow, error)
	CreateOrUpdateAthlete(ctx context.Context, arg database.CreateOrUpdateAthleteParams) (database.Athlete, error)
	AddChipBib(ctx context.Context, arg database.AddChipBibParams) (database.ChipBib, error)
	AddEventAthlete(ctx context.Context, arg database.AddEventAthleteParams) (database.EventAthlete, error)
	DeleteAthleteByID(ctx context.Context, athleteID uuid.UUID) error
	DeleteAthletesWithRaceID(ctx context.Context, raceID uuid.UUID) error
	DeleteAthletesWithEventID(ctx context.Context, eventID uuid.UUID) error
	DeleteChipBib(ctx context.Context, arg database.DeleteChipBibParams) error
	DeleteChipBibWithEventID(ctx context.Context, arg database.DeleteChipBibWithEventIDParams) error
	DeleteChipBibWithRaceID(ctx context.Context, raceID uuid.UUID) error
	// DeleteEventAthlete(ctx context.Context, arg database.DeleteEventAthleteParams) error
	GetEventAthlete(ctx context.Context, athleteID uuid.UUID) (database.EventAthlete, error)
	GetCategoryForAthlete(ctx context.Context, arg database.GetCategoryForAthleteParams) (database.Category, error)
	GetEventAthleteRecords(ctx context.Context, arg database.GetEventAthleteRecordsParams) ([]database.GetEventAthleteRecordsRow, error)
	GetSplitsForEvent(ctx context.Context, eventID uuid.UUID) ([]database.Split, error)
	WithTx(tx pgx.Tx) *database.Queries
}

type AthleteRepoPG struct {
	q  ParticipantQuery
	pg *postgres.Postgres
}

func NewAthleteRepoPG(q ParticipantQuery, pg *postgres.Postgres) *AthleteRepoPG {
	return &AthleteRepoPG{
		q:  q,
		pg: pg,
	}
}

func (ar *AthleteRepoPG) WithTx(tx pgx.Tx) *AthleteRepoPG {
	return &AthleteRepoPG{
		q:  ar.q.WithTx(tx),
		pg: ar.pg,
	}
}

// func (ar *AthleteRepoPG) WithRollback(ctx context.Context) (pgx.Tx, *AthleteRepoPG, func(ctx context.Context) error, error) {
// 	tx, err := ar.pool.Begin(ctx)
// 	if err != nil {
// 		return nil, nil, nil, err
// 	}
// 	return tx, ar.WithTx(tx), tx.Rollback, nil
// }

func (ar *AthleteRepoPG) SaveAthlete(ctx context.Context, p *entity.Athlete) error {
	tx, err := ar.pg.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	qtx := ar.WithTx(tx)
	aParams := database.CreateOrUpdateAthleteParams{
		ID:     p.ID,
		RaceID: p.RaceID,
		FirstName: pgtype.Text{
			String: p.FirstName,
			Valid:  true,
		},
		LastName: pgtype.Text{
			String: p.LastName,
			Valid:  true,
		},
		Gender: database.CategoryGender(p.Gender),
		DateOfBirth: pgtype.Date{
			Time:             p.DateOfBirth,
			InfinityModifier: 0,
			Valid:            true,
		},
		Phone: pgtype.Text{
			String: p.Phone,
			Valid:  true,
		},
		AthleteComments: pgtype.Text{
			String: p.Comments,
			Valid:  true,
		},
	}

	_, err = qtx.q.CreateOrUpdateAthlete(ctx, aParams)
	if err != nil {
		return err
	}
	cParams := database.AddChipBibParams{
		RaceID:  p.RaceID,
		EventID: p.EventID,
		Chip:    int32(p.Chip),
		Bib:     int32(p.Bib),
	}
	_, err = qtx.q.AddChipBib(ctx, cParams)
	if err != nil {
		return err
	}

	eaParams := database.AddEventAthleteParams{
		RaceID:     p.RaceID,
		EventID:    p.EventID,
		AthleteID:  p.ID,
		WaveID:     p.WaveID,
		CategoryID: p.CategoryID,
		Bib:        int32(p.Bib),
	}

	_, err = qtx.q.AddEventAthlete(ctx, eaParams)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (ar *AthleteRepoPG) GetCategoryFor(ctx context.Context, p *entity.Athlete) (uuid.NullUUID, bool, error) {
	params := database.GetCategoryForAthleteParams{
		EventID: p.EventID,
		Gender:  database.CategoryGender(p.Gender),
		DateFrom: pgtype.Timestamp{
			Time:             p.DateOfBirth,
			InfinityModifier: 0,
			Valid:            true,
		},
	}

	c, err := ar.q.GetCategoryForAthlete(ctx, params)
	if err != nil {
		if c.ID == uuid.Nil {
			return uuid.NullUUID{}, false, nil
		}
		return uuid.NullUUID{}, false, err
	}
	return uuid.NullUUID{
		UUID:  c.ID,
		Valid: true,
	}, true, nil
}

func (ar *AthleteRepoPG) GetAthleteWithChip(chip int) (*entity.Athlete, error) {
	return nil, nil
}

func (ar *AthleteRepoPG) GetAthleteByID(ctx context.Context, athleteID uuid.UUID) (*entity.Athlete, error) {
	a, err := ar.q.GetAthleteByID(ctx, athleteID)
	if err != nil {
		return nil, err
	}

	athlete := &entity.Athlete{
		ID:          a.ID,
		RaceID:      a.RaceID,
		EventID:     a.EventID,
		WaveID:      a.WaveID,
		Bib:         int(a.Bib),
		Chip:        int(a.Chip),
		FirstName:   a.FirstName.String,
		LastName:    a.LastName.String,
		Gender:      entity.CategoryGender(a.Gender),
		DateOfBirth: a.DateOfBirth.Time,
		CategoryID:  a.CategoryID,
		Phone:       a.Phone.String,
		Comments:    a.AthleteComments.String,
	}
	return athlete, nil
}

func (ar *AthleteRepoPG) DeleteAthletesForRace(ctx context.Context, raceID uuid.UUID) error {
	tx, err := ar.pg.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	qtx := ar.WithTx(tx)

	err = qtx.q.DeleteChipBibWithRaceID(ctx, raceID)
	if err != nil {
		return fmt.Errorf("error deleting chipbib for race = %s", raceID)
	}

	err = qtx.q.DeleteAthletesWithRaceID(ctx, raceID)
	if err != nil {
		return fmt.Errorf("error deleting athletes for race = %s", raceID)
	}

	return tx.Commit(ctx)
}

func (ar *AthleteRepoPG) DeleteAthletesForRaceWithEventID(ctx context.Context, raceID, eventID uuid.UUID) error {
	tx, err := ar.pg.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	qtx := ar.WithTx(tx)

	dParams := database.DeleteChipBibWithEventIDParams{
		RaceID:  raceID,
		EventID: eventID,
	}
	err = qtx.q.DeleteChipBibWithEventID(ctx, dParams)
	if err != nil {
		return fmt.Errorf("error deleting chipbib for eventID = %s", eventID)
	}

	err = qtx.q.DeleteAthletesWithEventID(ctx, eventID)
	if err != nil {
		return fmt.Errorf("error deleting athletes for race = %s", raceID)
	}

	return tx.Commit(ctx)
}

func (ar *AthleteRepoPG) DeleteAthlete(ctx context.Context, a *entity.Athlete) error {
	tx, err := ar.pg.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	qtx := ar.WithTx(tx)
	cbParams := database.DeleteChipBibParams{
		RaceID: a.RaceID,
		Chip:   int32(a.Chip),
		Bib:    int32(a.Bib),
	}
	err = qtx.q.DeleteChipBib(ctx, cbParams)
	if err != nil {
		return err
	}

	err = qtx.q.DeleteAthleteByID(ctx, a.ID)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (ar *AthleteRepoPG) GetRecordsAndSplitsForEventAthlete(ctx context.Context, raceID, eventID uuid.UUID) ([]database.GetEventAthleteRecordsRow, []*entity.Split, error) {
	ss, err := ar.q.GetSplitsForEvent(ctx, eventID)
	if err != nil {
		return nil, nil, err
	}

	eaParams := database.GetEventAthleteRecordsParams{
		RaceID:  raceID,
		EventID: eventID,
	}
	records, err := ar.q.GetEventAthleteRecords(ctx, eaParams)
	if err != nil {
		return nil, nil, err
	}

	splits := []*entity.Split{}
	for _, s := range ss {
		split := &entity.Split{
			ID:                 s.ID,
			RaceID:             s.RaceID,
			EventID:            s.EventID,
			Name:               s.SplitName,
			Type:               entity.SplitType(s.SplitType),
			DistanceFromStart:  int(s.DistanceFromStart),
			TimeReaderID:       s.TimeReaderID,
			MinTime:            time.Duration(s.MinTime.Microseconds * 1000),
			MaxTime:            time.Duration(s.MaxTime.Microseconds * 1000),
			MinLapTime:         time.Duration(s.MinLapTime.Microseconds * 1000),
			PreviousLapSplitID: s.PreviousLapSplitID,
		}
		splits = append(splits, split)
	}

	return records, splits, nil
}
