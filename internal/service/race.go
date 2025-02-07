package service

import (
	"context"

	"github.com/ecoarchie/timeit/internal/entity"
	"github.com/google/uuid"
)

type RaceRepo interface {
	Create(ctx context.Context, r entity.Race) (uuid.UUID, error)
}

type RaceService struct {
	repo RaceRepo
}

func NewRaceService(repo RaceRepo) *RaceService {
	return &RaceService{
		repo: repo,
	}
}
