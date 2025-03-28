package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ecoarchie/timeit/internal/controller/httpv1/dto"
	"github.com/ecoarchie/timeit/internal/entity"
	"github.com/ecoarchie/timeit/pkg/logger"
	"github.com/ecoarchie/timeit/pkg/validator"
	"github.com/google/uuid"
)

type ValidationErrors map[string]string

type RaceConfigurator interface {
	SaveRaceConfig(ctx context.Context, rc *dto.RaceModelDTO, v *validator.Validator) error
	GetRaces(ctx context.Context) ([]*entity.Race, error)
	CreateRace(ctx context.Context, req *dto.RaceDTO, v *validator.Validator) (*entity.Race, error)
	DeleteRace(ctx context.Context, raceID uuid.UUID) error
	GetRaceConfig(ctx context.Context, raceID uuid.UUID) (*entity.RaceModel, error)
	GetWavesForRace(ctx context.Context, raceID uuid.UUID) ([]*entity.Wave, error)
	StartWave(ctx context.Context, raceID uuid.UUID, startInfo entity.WaveStart) (time.Time, bool, error)
	GetEventIDsWithWaveStarted(ctx context.Context, raceID uuid.UUID) ([]uuid.UUID, error)
}

type RaceRepo interface {
	SaveRaceConfig(ctx context.Context, r *entity.Race, trs []*entity.TimeReader, ee []*entity.Event) error
	GetRaceConfig(ctx context.Context, raceID uuid.UUID) (*entity.RaceModel, error)
	GetRaces(ctx context.Context) ([]*entity.Race, error)
	SaveRaceInfo(ctx context.Context, race *entity.Race) error
	SaveWave(ctx context.Context, wave *entity.Wave) error
	DeleteRace(ctx context.Context, raceID uuid.UUID) error
	GetWavesForRace(ctx context.Context, raceID uuid.UUID) ([]*entity.Wave, error)
	GetWaveByID(ctx context.Context, waveID uuid.UUID) (*entity.Wave, error)
	GetEventIDsWithWavesStarted(ctx context.Context, raceID uuid.UUID) ([]uuid.UUID, error)
}

type RaceService struct {
	repo RaceRepo
	log  *logger.Logger
}

func NewRaceService(logger *logger.Logger, repo RaceRepo) *RaceService {
	return &RaceService{
		log:  logger,
		repo: repo,
	}
}

func (rs RaceService) GetRaces(ctx context.Context) ([]*entity.Race, error) {
	return rs.repo.GetRaces(ctx)
}

func (rs RaceService) GetEventIDsWithWaveStarted(ctx context.Context, raceID uuid.UUID) ([]uuid.UUID, error) {
	return rs.repo.GetEventIDsWithWavesStarted(ctx, raceID)
}

func (rs RaceService) CreateRace(ctx context.Context, req *dto.RaceDTO, v *validator.Validator) (*entity.Race, error) {
	r := entity.NewRace(req, v)
	if !v.Valid() {
		return nil, validator.ErrValidation
	}
	err := rs.repo.SaveRaceInfo(ctx, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (rs RaceService) GetRaceConfig(ctx context.Context, raceID uuid.UUID) (*entity.RaceModel, error) {
	rconfig, err := rs.repo.GetRaceConfig(ctx, raceID)
	if err != nil {
		return nil, err
	}
	if rconfig == nil {
		return nil, nil
	}

	return rconfig, nil
}

// FIXME return ErrValidation when v is not Valid
func (rs RaceService) SaveRaceConfig(ctx context.Context, rc *dto.RaceModelDTO, v *validator.Validator) error {
	race := entity.NewRace(rc.RaceDTO, v)
	if !v.Valid() {
		return validator.ErrValidation
	}

	v.Check(len(rc.TimeReaders) != 0, "time readers", "race must have at least one time reader")
	// no point for further validation since there are no time readers
	if !v.Valid() {
		return validator.ErrValidation
	}

	timeReaders := make([]*entity.TimeReader, 0, len(rc.TimeReaders))
	var timeReadersNames []string
	for _, tr := range rc.TimeReaders {
		timeReader := entity.NewTimeReader(tr, v)
		if !v.Valid() {
			return validator.ErrValidation
		}
		timeReadersNames = append(timeReadersNames, tr.ReaderName)
		timeReaders = append(timeReaders, timeReader)
	}
	v.Check(validator.Unique(timeReadersNames), "time readers names", "must be unique")

	events := make([]*entity.Event, 0, len(rc.Events))
	for _, e := range rc.Events {
		event := entity.NewEvent(e.EventDTO, e.Splits, rc.TimeReaders, e.Waves, e.Categories, v)
		if !v.Valid() {
			return validator.ErrValidation
		}
		events = append(events, event)
	}

	err := rs.repo.SaveRaceConfig(ctx, race, timeReaders, events)
	if err != nil {
		const msg = "error saving race to repo"
		rs.log.Error(msg, "error", err)
		return err
	}
	// BUG if event already exists then reassign categories and recalculate results
	return nil
}

func (rs RaceService) DeleteRace(ctx context.Context, raceID uuid.UUID) error {
	err := rs.repo.DeleteRace(ctx, raceID)
	if err != nil {
		return fmt.Errorf("error deleting race: %w", err)
	}
	return nil
}

func (rs RaceService) GetWavesForRace(ctx context.Context, raceID uuid.UUID) ([]*entity.Wave, error) {
	waves, err := rs.repo.GetWavesForRace(ctx, raceID)
	if err != nil {
		return nil, err
	}
	if waves == nil {
		return nil, nil
	}
	return waves, nil
}

func (rs RaceService) StartWave(ctx context.Context, raceID uuid.UUID, startInfo entity.WaveStart) (time.Time, bool, error) {
	w, err := rs.repo.GetWaveByID(ctx, startInfo.WaveID)
	if err != nil {
		return time.Time{}, false, err
	}
	if w == nil {
		return time.Time{}, false, nil
	}

	if startInfo.StartTime.IsZero() {
		w.StartTime = time.Now()
	} else {
		w.StartTime = startInfo.StartTime
	}
	w.IsLaunched = true

	err = rs.repo.SaveWave(ctx, w)
	if err != nil {
		return time.Time{}, true, fmt.Errorf("error saving wave: %w", err)
	}

	return w.StartTime, true, nil
}
