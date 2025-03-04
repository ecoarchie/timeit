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
	SaveRaceConfig(ctx context.Context, rc *entity.RaceConfig) error
	GetRaces(ctx context.Context) ([]*entity.Race, error)
	CreateRace(ctx context.Context, req *entity.RaceFormData) (*entity.Race, error)
	GetRaceConfig(ctx context.Context, raceID string) (*entity.RaceConfig, error)
}

type RaceRepo interface {
	SaveRaceConfig(ctx context.Context, r *entity.RaceConfig) error
	GetRaceConfig(ctx context.Context, raceID uuid.UUID) (*entity.RaceConfig, error)
	GetRaces(ctx context.Context) ([]*entity.Race, error)
	SaveRaceInfo(ctx context.Context, race *entity.Race) error
}

type RaceService struct {
	raceCache *RaceCache
	repo      RaceRepo
	log       logger.Interface
}

func NewRaceService(logger logger.Interface, rc *RaceCache, repo RaceRepo) *RaceService {
	return &RaceService{
		log:       logger,
		raceCache: rc,
		repo:      repo,
	}
}

func (rs RaceService) GetRaces(ctx context.Context) ([]*entity.Race, error) {
	return rs.repo.GetRaces(ctx)
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
