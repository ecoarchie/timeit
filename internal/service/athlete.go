package service

import (
	"context"
	"fmt"

	"github.com/ecoarchie/timeit/internal/entity"
	"github.com/ecoarchie/timeit/pkg/logger"
	"github.com/google/uuid"
)

type AthleteManager interface {
	GetAthleteByID(ctx context.Context, athleteID uuid.UUID) *entity.Athlete
	CreateAthlete(ctx context.Context, req entity.AthleteCreateRequest) (*entity.Athlete, error)
	UpdateAthlete(ctx context.Context, req entity.AthleteUpdateRequest) (*entity.Athlete, error)
	DeleteAthlete(ctx context.Context, athleteID uuid.UUID) error
}

type AthleteRepo interface {
	SaveAthlete(ctx context.Context, p *entity.Athlete) error
	GetCategoryFor(p *entity.Athlete) (uuid.NullUUID, error)
	GetAthleteWithChip(chip int) (*entity.Athlete, error)
	GetAthleteByID(ctx context.Context, athleteID uuid.UUID) (*entity.Athlete, error)
	DeleteAthlete(ctx context.Context, a *entity.Athlete) error
}

type AthleteService struct {
	l    logger.Interface
	repo AthleteRepo
}

func NewAthleteService(logger logger.Interface, repo AthleteRepo) *AthleteService {
	return &AthleteService{
		l:    logger,
		repo: repo,
	}
}

func (ps AthleteService) GetAthleteByID(ctx context.Context, athleteID uuid.UUID) *entity.Athlete {
	p, err := ps.repo.GetAthleteByID(ctx, athleteID)
	if err != nil {
		return nil
	}
	return p
}

func (ps AthleteService) CreateAthlete(ctx context.Context, req entity.AthleteCreateRequest) (*entity.Athlete, error) {
	p, err := entity.NewAthlete(req)
	if err != nil {
		return nil, err
	}

	if !req.CategoryID.Valid {
		ps.assignCategory(p)
	}

	err = ps.repo.SaveAthlete(ctx, p)
	if err != nil {
		return nil, err
	}

	// TODO create and store entity EventAthlete. DO it in repo layer in transaction
	return p, nil
}

func (ps AthleteService) assignCategory(p *entity.Athlete) error {
	catID, err := ps.repo.GetCategoryFor(p)
	if err != nil {
		return fmt.Errorf("error assigning category for athlete with bib %d", p.Bib)
	}
	p.CategoryID = catID
	return nil
}

func (ps AthleteService) UpdateAthlete(ctx context.Context, req entity.AthleteUpdateRequest) (*entity.Athlete, error) {
	p, err := ps.repo.GetAthleteByID(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("updateAthlete: athlete with ID %s not found", req.ID)
	}
	newP, err := entity.NewAthlete(req.AthleteCreateRequest)
	if err != nil {
		return nil, err
	}
	newP.ID = p.ID

	err = ps.repo.SaveAthlete(ctx, newP)
	if err != nil {
		return nil, err
	}
	return newP, nil
}

func (ps AthleteService) DeleteAthlete(ctx context.Context, athleteID uuid.UUID) error {
	a, err := ps.repo.GetAthleteByID(ctx, athleteID)
	if err != nil {
		return err
	}
	if a == nil {
		return fmt.Errorf("athlete with ID %s not found", athleteID)
	}
	err = ps.repo.DeleteAthlete(ctx, a)
	if err != nil {
		return fmt.Errorf("delete athlete: error deleting athlete %s from DB", athleteID)
	}
	return nil
}

func (ps AthleteService) DeleteAthleteBulk(ctx context.Context, raceID uuid.UUID, ids []uuid.UUID) []error {
	var errors []error
	for _, id := range ids {
		err := ps.DeleteAthlete(ctx, id)
		if err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}
