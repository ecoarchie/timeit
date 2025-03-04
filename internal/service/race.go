package service

import (
	"context"
	"fmt"

	"github.com/ecoarchie/timeit/internal/entity"
	"github.com/ecoarchie/timeit/pkg/logger"
	"github.com/google/uuid"
)

type ValidationErrors map[string]string

// TODO separate errors to ErrorToSave (preventing from saving) and Warning (can save, just pay attention)
// TODO add category boundaries validation

type RaceConfigurator interface {
	Save(ctx context.Context, rc *entity.RaceConfig) error
	CreateRace(ctx context.Context, req *entity.RaceFormData) (*entity.Race, error)
	GetRaceConfig(ctx context.Context, raceID string) (*entity.RaceConfig, error)
}

type RaceRepo interface {
	SaveRaceConfig(ctx context.Context, r *entity.RaceConfig) error
	GetRaceConfig(ctx context.Context, raceID uuid.UUID) (*entity.RaceConfig, error)
}

type RaceService struct {
	l         logger.Interface
	raceCache *RaceCache
	repo      RaceRepo
}

func NewRaceService(logger logger.Interface, rc *RaceCache, repo RaceRepo) *RaceService {
	return &RaceService{
		l:         logger,
		raceCache: rc,
		repo:      repo,
	}
}

func (rc RaceService) CreateRace(ctx context.Context, req *entity.RaceFormData) (*entity.Race, error) {
	r, err := entity.NewRace(req)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (rc RaceService) GetRaceConfig(ctx context.Context, raceID string) (*entity.RaceConfig, error) {
	uuID, err := uuid.Parse(raceID)
	if err != nil {
		return nil, fmt.Errorf("error parsing race UUID")
	}
	rconfig, err := rc.repo.GetRaceConfig(ctx, uuID)
	if err != nil {
		return nil, err
	}
	if rconfig == nil {
		return nil, nil
	}

	// FIXME update RaceCache

	return rconfig, nil
}

func (rs RaceService) Save(ctx context.Context, rc *entity.RaceConfig) error {
	err := rs.repo.SaveRaceConfig(ctx, rc)
	if err != nil {
		const msg = "error saving race to repo"
		rs.l.Error(msg, err)
		return err
	}
	rs.raceCache.StoreRaceConfig(rc)
	rs.l.Info("race cache updated")
	return nil
}
