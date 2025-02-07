package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ecoarchie/timeit/internal/entity"
	"github.com/ecoarchie/timeit/pkg/logger"
	"github.com/google/uuid"
)

type CreateRaceRequest struct {
	Name     string    `json:"name"`
	RaceDate time.Time `json:"race_date"`
	Timezone string    `json:"timezone"`
}

type RaceConfigurator interface {
	Create(ctx context.Context, req CreateRaceRequest) (uuid.UUID, error)
}

type RaceRepo interface {
	Create(ctx context.Context, r entity.Race) (uuid.UUID, error)
}

type RaceService struct {
	l    logger.Interface
	repo RaceRepo
}

func NewRaceService(logger logger.Interface, repo RaceRepo) *RaceService {
	return &RaceService{
		l:    logger,
		repo: repo,
	}
}

func (rs RaceService) Create(ctx context.Context, req CreateRaceRequest) (uuid.UUID, error) {
	race, err := entity.NewRace(req.Name, req.RaceDate, req.Timezone)
	if err != nil {
		msg := "create race validation error"
		rs.l.Error(msg, "error", err)
		return uuid.Nil, fmt.Errorf("error creating race. Invalid fields: %w", err)
	}
	// TODO initialize pg pool properly with DB_URL
	id, err := rs.repo.Create(ctx, *race)
	if err != nil {
		const msg = "error saving race to repo"
		rs.l.Error(msg, err)
		return uuid.Nil, errors.New(msg)
	}
	return id, nil
}
