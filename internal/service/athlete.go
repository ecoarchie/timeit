package service

import (
	"fmt"
	"time"

	"github.com/ecoarchie/timeit/internal/entity"
	"github.com/ecoarchie/timeit/pkg/logger"
	"github.com/google/uuid"
)

type AthleteManager interface {
	GetAthleteByID(raceID, athleteID uuid.UUID) *entity.Athlete
	CreateAthlete(req entity.AthleteCreateRequest) (*entity.Athlete, error)
	UpdateAthlete(req entity.AthleteUpdateRequest) (*entity.Athlete, error)
	DeleteAthlete(raceID, id uuid.UUID) error
}

type AthleteRepo interface {
	SaveAthlete(p *entity.Athlete) error
	GetCategoryFor(eventID uuid.UUID, gender entity.CategoryGender, dob time.Time) (uuid.NullUUID, error)
	GetAthleteWithChip(chip int) (*entity.Athlete, error)
	GetAthleteByID(raceID, athleteID uuid.UUID) (*entity.Athlete, error)
	DeleteAthlete(raceID uuid.UUID, id uuid.UUID) error
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

func (ps AthleteService) GetAthleteByID(raceID, athleteID uuid.UUID) *entity.Athlete {
	p, err := ps.repo.GetAthleteByID(raceID, athleteID)
	if err != nil {
		return nil
	}
	return p
}

func (ps AthleteService) CreateAthlete(req entity.AthleteCreateRequest) (*entity.Athlete, error) {
	p, err := entity.NewAthlete(req)
	if err != nil {
		return nil, err
	}

	if !req.CategoryID.Valid {
		ps.assignCategory(p)
	}

	err = ps.repo.SaveAthlete(p)
	if err != nil {
		return nil, err
	}

	// TODO create and store entity EventAthlete. DO it in repo layer in transaction
	return p, nil
}

func (ps AthleteService) assignCategory(p *entity.Athlete) error {
	catID, err := ps.repo.GetCategoryFor(p.EventID, p.Gender, p.DateOfBirth)
	if err != nil {
		return fmt.Errorf("error assigning category for athlete with bib %d", p.Bib)
	}
	p.CategoryID = catID
	return nil
}

func (ps AthleteService) UpdateAthlete(req entity.AthleteUpdateRequest) (*entity.Athlete, error) {
	p, err := ps.repo.GetAthleteByID(req.RaceID, req.ID)
	if err != nil {
		return nil, fmt.Errorf("updateAthlete: athlete with ID %s not found", req.ID)
	}
	newP, err := entity.NewAthlete(req.AthleteCreateRequest)
	if err != nil {
		return nil, err
	}
	newP.ID = p.ID

	err = ps.repo.SaveAthlete(newP)
	if err != nil {
		return nil, err
	}
	return newP, nil
}

func (ps AthleteService) DeleteAthlete(raceID, id uuid.UUID) error {
	_, err := ps.repo.GetAthleteByID(raceID, id)
	if err != nil {
		return fmt.Errorf("deleteAthlete: athlete with ID %s not found", id)
	}
	err = ps.repo.DeleteAthlete(raceID, id)
	if err != nil {
		return fmt.Errorf("delete athlete: error deleting athlete %s from DB", id)
	}
	return nil
}

func (ps AthleteService) DeleteAthleteBulk(raceID uuid.UUID, ids []uuid.UUID) []error {
	var errors []error
	for _, id := range ids {
		err := ps.DeleteAthlete(raceID, id)
		if err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}
