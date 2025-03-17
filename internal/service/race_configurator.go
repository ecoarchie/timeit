package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ecoarchie/timeit/internal/entity"
	"github.com/ecoarchie/timeit/pkg/logger"
	"github.com/google/uuid"
)

type ValidationErrors map[string]string

// TODO separate errors to ErrorToSave (preventing from saving) and Warning (can save, just pay attention)
// TODO add category boundaries validation

type RaceConfigurator interface {
	SaveRaceConfig(ctx context.Context, rc *entity.RaceConfig) error
	GetRaces(ctx context.Context) ([]*entity.Race, error)
	CreateRace(ctx context.Context, req *entity.RaceFormData) (*entity.Race, error)
	DeleteRace(ctx context.Context, raceID string) error
	GetRaceConfig(ctx context.Context, raceID string) (*entity.RaceConfig, error)
	GetWavesForRace(ctx context.Context, raceID string) ([]*entity.Wave, error)
	StartWave(ctx context.Context, raceID string, startInfo entity.WaveStart) (time.Time, bool, error)
	GetEventIDsWithWaveStarted(ctx context.Context, raceID uuid.UUID) ([]uuid.UUID, error)
}

type RaceRepo interface {
	SaveRaceConfig(ctx context.Context, r *entity.RaceConfig) error
	GetRaceConfig(ctx context.Context, raceID uuid.UUID) (*entity.RaceConfig, error)
	GetRaces(ctx context.Context) ([]*entity.Race, error)
	SaveRaceInfo(ctx context.Context, race *entity.Race) error
	SaveWave(ctx context.Context, wave *entity.Wave) error
	DeleteRace(ctx context.Context, raceID uuid.UUID) error
	GetWavesForRace(ctx context.Context, raceID uuid.UUID) ([]*entity.Wave, error)
	GetWaveByID(ctx context.Context, waveID uuid.UUID) (*entity.Wave, error)
	GetEventIDsWithWavesStarted(ctx context.Context, raceID uuid.UUID) ([]uuid.UUID, error)
}

type RaceService struct {
	raceCache *RaceCache
	repo      RaceRepo
	log       *logger.Logger
}

func NewRaceService(logger *logger.Logger, rc *RaceCache, repo RaceRepo) *RaceService {
	return &RaceService{
		log:       logger,
		raceCache: rc,
		repo:      repo,
	}
}

func (rs RaceService) GetRaces(ctx context.Context) ([]*entity.Race, error) {
	return rs.repo.GetRaces(ctx)
}

func (rs RaceService) GetEventIDsWithWaveStarted(ctx context.Context, raceID uuid.UUID) ([]uuid.UUID, error) {
	return rs.repo.GetEventIDsWithWavesStarted(ctx, raceID)
}

func (rs RaceService) CreateRace(ctx context.Context, req *entity.RaceFormData) (*entity.Race, error) {
	r, err := entity.NewRace(req)
	if err != nil {
		return nil, err
	}
	err = rs.repo.SaveRaceInfo(ctx, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (rs RaceService) GetRaceConfig(ctx context.Context, raceID string) (*entity.RaceConfig, error) {
	uuID, err := uuid.Parse(raceID)
	if err != nil {
		return nil, fmt.Errorf("error parsing race UUID")
	}
	rconfig, err := rs.repo.GetRaceConfig(ctx, uuID)
	if err != nil {
		return nil, err
	}
	if rconfig == nil {
		return nil, nil
	}

	rs.raceCache.UpdateWith(rconfig)
	return rconfig, nil
}

func (rs RaceService) SaveRaceConfig(ctx context.Context, rc *entity.RaceConfig) error {
	err := rs.repo.SaveRaceConfig(ctx, rc)
	if err != nil {
		const msg = "error saving race to repo"
		rs.log.Error(msg, err)
		return err
	}
	rs.raceCache.UpdateWith(rc)
	rs.log.Info("race cache updated")
	return nil
}

func (rs RaceService) DeleteRace(ctx context.Context, raceID string) error {
	id, err := uuid.Parse(raceID)
	if err != nil {
		return fmt.Errorf("error parsing raceID")
	}
	err = rs.repo.DeleteRace(ctx, id)
	if err != nil {
		return fmt.Errorf("error deleting race: %w", err)
	}
	return nil
}

func (rs RaceService) GetWavesForRace(ctx context.Context, raceID string) ([]*entity.Wave, error) {
	id := uuid.MustParse(raceID)
	waves, err := rs.repo.GetWavesForRace(ctx, id)
	if err != nil {
		return nil, err
	}
	if waves == nil {
		return nil, nil
	}
	return waves, nil
}

func (rs RaceService) StartWave(ctx context.Context, raceID string, startInfo entity.WaveStart) (time.Time, bool, error) {
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
