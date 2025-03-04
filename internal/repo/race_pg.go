package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

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
	GetRaces(ctx context.Context) ([]database.Race, error)
	GetRaceInfo(ctx context.Context, id uuid.UUID) (database.Race, error)
	GetAllTimeReadersForRace(ctx context.Context, raceID uuid.UUID) ([]database.TimeReader, error)
	GetAllEventsForRace(ctx context.Context, raceID uuid.UUID) ([]database.Event, error)
	GetAllSplitsForEvent(ctx context.Context, eventID uuid.UUID) ([]database.Split, error)
	GetAllWavesForEvent(ctx context.Context, eventID uuid.UUID) ([]database.Wave, error)
	GetCategoriesForEvent(ctx context.Context, eventID uuid.UUID) ([]database.Category, error)
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
	_, err = qtx.q.AddRace(ctx, addRaceParams)
	if err != nil {
		return err
	}

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
			EventDate: pgtype.Timestamp{
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
				MinTime: pgtype.Int8{
					Int64: int64(tp.MinTime),
					Valid: true,
				},
				MaxTime: pgtype.Int8{
					Int64: int64(tp.MaxTime),
					Valid: true,
				},
				MinLapTime: pgtype.Int8{
					Int64: int64(tp.MinLapTime),
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
				StartTime: pgtype.Timestamp{
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
				AgeFrom:      int32(c.AgeFrom),
				DateFrom: pgtype.Timestamp{
					Time:             c.DateFrom,
					InfinityModifier: 0,
					Valid:            true,
				},
				AgeTo: int32(c.AgeTo),
				DateTo: pgtype.Timestamp{
					Time:             c.DateTo,
					InfinityModifier: 0,
					Valid:            true,
				},
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
	r, err := rr.q.GetRaceInfo(ctx, raceID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	raceCfg := &entity.RaceConfig{
		Race: &entity.Race{
			ID:       r.ID,
			Name:     r.RaceName,
			Timezone: r.Timezone,
		},
		TimeReaders: []*entity.TimeReader{},
		Events:      []*entity.EventConfig{},
	}

	// get all time readers for race
	trs, err := rr.q.GetAllTimeReadersForRace(ctx, raceID)
	if err != nil {
		return nil, err
	}
	for _, tr := range trs {
		reader := &entity.TimeReader{
			ID:         tr.ID,
			RaceID:     tr.RaceID,
			ReaderName: tr.ReaderName,
		}
		raceCfg.TimeReaders = append(raceCfg.TimeReaders, reader)
	}

	// get all events for race
	events, err := rr.q.GetAllEventsForRace(ctx, raceID)
	if err != nil {
		return nil, err
	}
	for _, e := range events {
		event := &entity.EventConfig{
			Event: &entity.Event{
				ID:               e.ID,
				RaceID:           e.RaceID,
				Name:             e.EventName,
				DistanceInMeters: int(e.DistanceInMeters),
				EventDate:        e.EventDate.Time,
			},
			Splits:     []*entity.Split{},
			Waves:      []*entity.Wave{},
			Categories: []*entity.Category{},
		}

		// get splits for event
		splits, err := rr.q.GetAllSplitsForEvent(ctx, e.ID)
		if err != nil {
			return nil, err
		}
		for _, s := range splits {
			split := &entity.Split{
				ID:                s.ID,
				RaceID:            s.RaceID,
				EventID:           s.EventID,
				Name:              s.SplitName,
				Type:              entity.SplitType(s.SplitType),
				DistanceFromStart: int(s.DistanceFromStart),
				TimeReaderID:      s.TimeReaderID,
				MinTime:           time.Duration(s.MinTime.Int64),
				MaxTime:           time.Duration(s.MaxTime.Int64),
				MinLapTime:        time.Duration(s.MinLapTime.Int64),
			}
			event.Splits = append(event.Splits, split)
		}

		// get waves for event
		waves, err := rr.q.GetAllWavesForEvent(ctx, e.ID)
		if err != nil {
			return nil, err
		}
		for _, w := range waves {
			wave := &entity.Wave{
				ID:         w.ID,
				RaceID:     w.RaceID,
				EventID:    w.EventID,
				Name:       w.WaveName,
				StartTime:  w.StartTime.Time,
				IsLaunched: w.IsLaunched,
			}
			event.Waves = append(event.Waves, wave)
		}

		// get categories for event
		cats, err := rr.q.GetCategoriesForEvent(ctx, e.ID)
		if err != nil {
			return nil, err
		}
		for _, c := range cats {
			category := &entity.Category{
				ID:       c.ID,
				RaceID:   c.RaceID,
				EventID:  c.EventID,
				Name:     c.CategoryName,
				Gender:   entity.CategoryGender(c.Gender),
				AgeFrom:  int(c.AgeFrom),
				DateFrom: c.DateFrom.Time,
				AgeTo:    int(c.AgeTo),
				DateTo:   c.DateTo.Time,
			}
			event.Categories = append(event.Categories, category)
		}

		raceCfg.Events = append(raceCfg.Events, event)
	}

	return raceCfg, nil
}

func (rr *RaceRepoPG) GetRaces(ctx context.Context) ([]*entity.Race, error) {
	races, err := rr.q.GetRaces(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	var res []*entity.Race
	for _, r := range races {
		race := &entity.Race{
			ID:       r.ID,
			Name:     r.RaceName,
			Timezone: r.Timezone,
		}
		res = append(res, race)
	}
	return res, nil
}

func (rr *RaceRepoPG) SaveRaceInfo(ctx context.Context, race *entity.Race) error {
	params := database.AddRaceParams{
		ID:       race.ID,
		RaceName: race.Name,
		Timezone: race.Timezone,
	}
	_, err := rr.q.AddRace(ctx, params)
	if err != nil {
		return err
	}
	return nil
}
