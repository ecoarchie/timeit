package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ecoarchie/timeit/internal/database"
	"github.com/ecoarchie/timeit/internal/entity"
	"github.com/ecoarchie/timeit/pkg/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type RaceQuery interface {
	GetRaces(ctx context.Context) ([]database.Race, error)
	GetRaceInfo(ctx context.Context, id uuid.UUID) (database.Race, error)
	AddRace(ctx context.Context, arg database.AddRaceParams) (database.Race, error)
	DeleteRace(ctx context.Context, id uuid.UUID) error
	AddOrUpdateTimeReader(ctx context.Context, arg database.AddOrUpdateTimeReaderParams) (database.TimeReader, error)
	AddOrUpdateEvent(ctx context.Context, arg database.AddOrUpdateEventParams) (database.Event, error)
	AddOrUpdateSplit(ctx context.Context, arg database.AddOrUpdateSplitParams) (database.Split, error)
	AddOrUpdateWave(ctx context.Context, arg database.AddOrUpdateWaveParams) (database.Wave, error)
	AddOrUpdateCategory(ctx context.Context, arg database.AddOrUpdateCategoryParams) (database.Category, error)
	GetTimeReadersForRace(ctx context.Context, raceID uuid.UUID) ([]database.TimeReader, error)
	GetEventsForRace(ctx context.Context, raceID uuid.UUID) ([]database.Event, error)
	GetSplitsForEvent(ctx context.Context, eventID uuid.UUID) ([]database.Split, error)
	GetWavesForEvent(ctx context.Context, eventID uuid.UUID) ([]database.Wave, error)
	GetWavesForRace(ctx context.Context, raceID uuid.UUID) ([]database.Wave, error)
	GetCategoriesForEvent(ctx context.Context, eventID uuid.UUID) ([]database.Category, error)
	GetWaveByID(ctx context.Context, id uuid.UUID) (database.Wave, error)
	GetEventIDsWithWavesStarted(ctx context.Context, raceID uuid.UUID) ([]uuid.UUID, error)
	WithTx(tx pgx.Tx) *database.Queries
}

type RaceRepoPG struct {
	q  RaceQuery
	pg *postgres.Postgres
}

func NewRaceRepoPG(q RaceQuery, pg *postgres.Postgres) *RaceRepoPG {
	return &RaceRepoPG{
		q:  q,
		pg: pg,
	}
}

func (rr *RaceRepoPG) WithTx(tx pgx.Tx) *RaceRepoPG {
	return &RaceRepoPG{
		q:  rr.q.WithTx(tx),
		pg: rr.pg,
	}
}

func (rr *RaceRepoPG) SaveRaceConfig(ctx context.Context, r *entity.RaceConfig) error {
	tx, err := rr.pg.Pool.Begin(ctx)
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
				MinTime: pgtype.Interval{
					Microseconds: time.Duration(tp.MinTime).Microseconds(),
					Days:         0,
					Months:       0,
					Valid:        true,
				},
				MaxTime: pgtype.Interval{
					Microseconds: time.Duration(tp.MaxTime).Microseconds(),
					Days:         0,
					Months:       0,
					Valid:        true,
				},
				MinLapTime: pgtype.Interval{
					Microseconds: time.Duration(tp.MinLapTime).Microseconds(),
					Days:         0,
					Months:       0,
					Valid:        true,
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
	trs, err := rr.q.GetTimeReadersForRace(ctx, raceID)
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
	events, err := rr.q.GetEventsForRace(ctx, raceID)
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
			Splits:     []*entity.SplitConfig{},
			Waves:      []*entity.Wave{},
			Categories: []*entity.Category{},
		}

		// get splits for event
		splits, err := rr.q.GetSplitsForEvent(ctx, e.ID)
		if err != nil {
			return nil, err
		}
		for _, s := range splits {
			split := &entity.SplitConfig{
				ID:                s.ID,
				RaceID:            s.RaceID,
				EventID:           s.EventID,
				Name:              s.SplitName,
				Type:              entity.SplitType(s.SplitType),
				DistanceFromStart: int(s.DistanceFromStart),
				TimeReaderID:      s.TimeReaderID,
				MinTime:           entity.Duration(s.MinTime.Microseconds * 1000),
				MaxTime:           entity.Duration(s.MaxTime.Microseconds * 1000),
				MinLapTime:        entity.Duration(s.MinLapTime.Microseconds * 1000),
			}
			event.Splits = append(event.Splits, split)
		}

		// get waves for event
		waves, err := rr.q.GetWavesForEvent(ctx, e.ID)
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

func (rr *RaceRepoPG) DeleteRace(ctx context.Context, raceID uuid.UUID) error {
	err := rr.q.DeleteRace(ctx, raceID)
	return err
}

func (rr *RaceRepoPG) GetWavesForRace(ctx context.Context, raceID uuid.UUID) ([]*entity.Wave, error) {
	ws, err := rr.q.GetWavesForRace(ctx, raceID)
	if err != nil {
		return nil, err
	}
	if ws == nil {
		return nil, nil
	}
	waves := []*entity.Wave{}
	for _, w := range ws {
		wave := &entity.Wave{
			ID:         w.ID,
			RaceID:     w.RaceID,
			EventID:    w.EventID,
			Name:       w.WaveName,
			StartTime:  w.StartTime.Time,
			IsLaunched: w.IsLaunched,
		}
		waves = append(waves, wave)
	}
	return waves, nil
}

func (rr *RaceRepoPG) GetWaveByID(ctx context.Context, waveID uuid.UUID) (*entity.Wave, error) {
	w, err := rr.q.GetWaveByID(ctx, waveID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	wave := &entity.Wave{
		ID:         w.ID,
		RaceID:     w.RaceID,
		EventID:    w.EventID,
		Name:       w.WaveName,
		StartTime:  w.StartTime.Time,
		IsLaunched: w.IsLaunched,
	}
	return wave, nil
}

func (rr *RaceRepoPG) SaveWave(ctx context.Context, wave *entity.Wave) error {
	wParams := database.AddOrUpdateWaveParams{
		ID:       wave.ID,
		RaceID:   wave.RaceID,
		EventID:  wave.EventID,
		WaveName: wave.Name,
		StartTime: pgtype.Timestamp{
			Time:             wave.StartTime,
			InfinityModifier: 0,
			Valid:            true,
		},
		IsLaunched: wave.IsLaunched,
	}
	_, err := rr.q.AddOrUpdateWave(ctx, wParams)
	if err != nil {
		return err
	}
	return nil
}

func (rr *RaceRepoPG) GetEventIDsWithWavesStarted(ctx context.Context, raceID uuid.UUID) ([]uuid.UUID, error) {
	return rr.q.GetEventIDsWithWavesStarted(ctx, raceID)
}
