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

type RaceQuery interface {
	AddRace(ctx context.Context, arg database.AddRaceParams) (database.Race, error)
	AddOrUpdateTimeReader(ctx context.Context, arg database.AddOrUpdateTimeReaderParams) (database.TimeReader, error)
	AddOrUpdateEvent(ctx context.Context, arg database.AddOrUpdateEventParams) (database.Event, error)
	AddOrUpdateSplit(ctx context.Context, arg database.AddOrUpdateSplitParams) (database.Split, error)
	AddOrUpdateWave(ctx context.Context, arg database.AddOrUpdateWaveParams) (database.Wave, error)
	AddOrUpdateCategory(ctx context.Context, arg database.AddOrUpdateCategoryParams) (database.Category, error)
	WithTx(tx pgx.Tx) *database.Queries
}

type RaceRepoPG struct {
	q    RaceQuery
	pool *pgxpool.Pool
}

func NewRaceRepoPG(q RaceQuery, pool *pgxpool.Pool) *RaceRepoPG {
	return &RaceRepoPG{
		q:    q,
		pool: pool,
	}
}

func (rr *RaceRepoPG) WithTx(tx pgx.Tx) *RaceRepoPG {
	return &RaceRepoPG{
		q:    rr.q.WithTx(tx),
		pool: rr.pool,
	}
}

func (rr *RaceRepoPG) SaveRaceConfig(ctx context.Context, r *entity.RaceConfig) error {
	tx, err := rr.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	qtx := rr.WithTx(tx)

	// Save race
	addRaceParams := database.AddRaceParams{
		ID:       r.ID,
		RaceName: r.Name,
		Timezone: r.Timezone,
	}
	race, err := qtx.q.AddRace(ctx, addRaceParams)
	if err != nil {
		return err
	}

	fmt.Printf("race: %v\n", race)

	// Save physical time_readers
	for _, l := range r.TimeReaders {
		locParam := database.AddOrUpdateTimeReaderParams{
			ID:         l.ID,
			RaceID:     l.RaceID,
			ReaderName: l.ReaderName,
		}
		_, err := qtx.q.AddOrUpdateTimeReader(ctx, locParam)
		if err != nil {
			return fmt.Errorf("error adding time_reader with ID %s: %v", locParam.ID, err)
		}
	}

	// Save events
	for _, e := range r.Events {
		eParams := database.AddOrUpdateEventParams{
			ID:               e.ID,
			RaceID:           e.RaceID,
			EventName:        e.Name,
			DistanceInMeters: int32(e.DistanceInMeters),
			EventDate: pgtype.Timestamptz{
				Time:             e.EventDate,
				InfinityModifier: 0,
				Valid:            true,
			},
		}
		_, err := qtx.q.AddOrUpdateEvent(ctx, eParams)
		if err != nil {
			return fmt.Errorf("error adding event with ID %s: %v", eParams.ID, err)
		}

		// Save splits for event
		for _, tp := range e.Splits {
			tpParams := database.AddOrUpdateSplitParams{
				ID:                tp.ID,
				RaceID:            tp.RaceID,
				EventID:           tp.EventID,
				SplitName:         tp.Name,
				SplitType:         database.TpType(tp.Type),
				DistanceFromStart: int32(tp.DistanceFromStart),
				TimeReaderID:      tp.TimeReaderID,
				MinTimeSec: pgtype.Int4{
					Int32: tp.MinTimeSec,
					Valid: true,
				},
				MaxTimeSec: pgtype.Int4{
					Int32: tp.MaxTimeSec,
					Valid: true,
				},
				MinLapTimeSec: pgtype.Int4{
					Int32: tp.MinLapTimeSec,
					Valid: true,
				},
			}
			_, err := qtx.q.AddOrUpdateSplit(ctx, tpParams)
			if err != nil {
				return fmt.Errorf("error adding split with ID %s: %v", tpParams.ID, err)
			}
		}

		// Save waves
		for _, w := range e.Waves {
			wParams := database.AddOrUpdateWaveParams{
				ID:       w.ID,
				RaceID:   w.RaceID,
				EventID:  w.EventID,
				WaveName: w.Name,
				StartTime: pgtype.Timestamptz{
					Time:             w.StartTime,
					InfinityModifier: 0,
					Valid:            true,
				},
				IsLaunched: w.IsLaunched,
			}
			_, err := qtx.q.AddOrUpdateWave(ctx, wParams)
			if err != nil {
				return fmt.Errorf("error adding wave with ID %s: %v", wParams.ID, err)
			}
		}

		// Save categories
		for _, c := range e.Categories {
			cParams := database.AddOrUpdateCategoryParams{
				ID:           c.ID,
				RaceID:       c.RaceID,
				EventID:      c.EventID,
				CategoryName: c.Name,
				Gender:       database.CategoryGender(c.Gender),
				FromAge:      int32(c.FromAge),
				FromRaceDate: c.FromRaceDate,
				ToAge:        int32(c.ToAge),
				ToRaceDate:   c.ToRaceDate,
			}
			_, err := qtx.q.AddOrUpdateCategory(ctx, cParams)
			if err != nil {
				return fmt.Errorf("error adding category with ID %s: %v", cParams.ID, err)
			}
		}
	}
	return tx.Commit(ctx)
}

func (rr *RaceRepoPG) GetRaceConfig(ctx context.Context, raceID uuid.UUID) (*entity.RaceConfig, error) {
	return nil, nil
}
