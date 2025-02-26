package repo

import (
	"context"
	"fmt"

	"github.com/ecoarchie/timeit/internal/database"
	"github.com/ecoarchie/timeit/internal/entity"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ParticipantQuery interface {
	GetAthleteByID(ctx context.Context, id uuid.UUID) (database.GetAthleteByIDRow, error)
	CreateOrUpdateAthlete(ctx context.Context, arg database.CreateOrUpdateAthleteParams) (database.Athlete, error)
	AddChipBib(ctx context.Context, arg database.AddChipBibParams) (database.ChipBib, error)
	AddEventAthlete(ctx context.Context, arg database.AddEventAthleteParams) (database.EventAthlete, error)
	DeleteAthleteByID(ctx context.Context, athleteID uuid.UUID) error
	DeleteChipBib(ctx context.Context, arg database.DeleteChipBibParams) error
	DeleteEventAthlete(ctx context.Context, arg database.DeleteEventAthleteParams) error
	GetEventAthlete(ctx context.Context, athleteID uuid.UUID) (database.EventAthlete, error)
	WithTx(tx pgx.Tx) *database.Queries
}

type AthleteRepoPG struct {
	q    ParticipantQuery
	pool *pgxpool.Pool
}

func NewAthleteRepoPG(q ParticipantQuery, pool *pgxpool.Pool) *AthleteRepoPG {
	return &AthleteRepoPG{
		q:    q,
		pool: pool,
	}
}

func (ar *AthleteRepoPG) WithTx(tx pgx.Tx) *AthleteRepoPG {
	return &AthleteRepoPG{
		q:    ar.q.WithTx(tx),
		pool: ar.pool,
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
	tx, err := ar.pool.Begin(ctx)
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
		Bib: pgtype.Int4{
			Int32: int32(p.Bib),
			Valid: true,
		},
	}

	_, err = qtx.q.AddEventAthlete(ctx, eaParams)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// func athelteFromDBtoEntity(a database.Athlete, cb database.ChipBib) *entity.Athlete {
// 	return &entity.Athlete{
// 		ID:          a.ID,
// 		RaceID:      a.RaceID,
// 		EventID:     a.EventID,
// 		WaveID:      uuid.UUID{},
// 		Bib:         0,
// 		Chip:        0,
// 		FirstName:   "",
// 		LastName:    "",
// 		Gender:      "",
// 		DateOfBirth: time.Time{},
// 		CategoryID:  uuid.NullUUID{},
// 		Phone:       "",
// 		Comments:    "",
// 	}
// }

func (ar *AthleteRepoPG) GetCategoryFor(p *entity.Athlete) (uuid.NullUUID, error) {
	return uuid.NullUUID{}, nil
}

func (ar *AthleteRepoPG) GetAthleteWithChip(chip int) (*entity.Athlete, error) {
	return nil, nil
}

func (ar *AthleteRepoPG) GetAthleteByID(ctx context.Context, athleteID uuid.UUID) (*entity.Athlete, error) {
	a, err := ar.q.GetAthleteByID(ctx, athleteID)
	if err != nil {
		return nil, err
	}

	// TODO query full data for athlete from db
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

func (ar *AthleteRepoPG) DeleteAthlete(ctx context.Context, a *entity.Athlete) error {
	tx, err := ar.pool.Begin(ctx)
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
		fmt.Println("error here 2")
		return err
	}

	eaParams := database.DeleteEventAthleteParams{
		RaceID:    a.RaceID,
		AthleteID: a.ID,
	}
	err = qtx.q.DeleteEventAthlete(ctx, eaParams)
	if err != nil {
		fmt.Println("error here 3")
		return err
	}
	err = qtx.q.DeleteAthleteByID(ctx, a.ID)
	if err != nil {
		fmt.Println("error here")
		fmt.Println(a)
		return err
	}
	return tx.Commit(ctx)
}
