package service

import (
	"fmt"
	"time"

	"github.com/ecoarchie/timeit/internal/entity"
	"github.com/ecoarchie/timeit/pkg/logger"
	"github.com/google/uuid"
)

type ParticipantManager interface {
	CreateParticipant(req entity.ParticipantCreateRequest) (*entity.Participant, error)
	UpdateParticipant(req entity.ParticipantUpdateRequest) (*entity.Participant, error)
	DeleteParticipant(raceID, id uuid.UUID) error
}

type ParticipantRepo interface {
	SaveParticipant(p *entity.Participant) error
	GetCategoryFor(eventID uuid.UUID, gender entity.CategoryGender, dob time.Time) (uuid.NullUUID, error)
	GetParticipantWithChip(chip int) (*entity.Participant, error)
	GetParticipantByID(id uuid.UUID) (*entity.Participant, error)
	DeleteParticipant(raceID uuid.UUID, id uuid.UUID) error
}

type ParticipantService struct {
	l    logger.Interface
	repo ParticipantRepo
}

func NewParticipantService(logger logger.Interface, repo ParticipantRepo) *ParticipantService {
	return &ParticipantService{
		l:    logger,
		repo: repo,
	}
}

func (ps ParticipantService) CreateParticipant(req entity.ParticipantCreateRequest) (*entity.Participant, error) {
	p, err := entity.NewParticipant(req)
	if err != nil {
		return nil, err
	}

	if !req.CategoryID.Valid {
		ps.assignCategory(p)
	}

	err = ps.repo.SaveParticipant(p)
	if err != nil {
		return nil, err
	}

	// TODO create and store entity EventParticipant. DO it in repo layer in transaction
	return p, nil
}

func (ps ParticipantService) assignCategory(p *entity.Participant) error {
	catID, err := ps.repo.GetCategoryFor(p.EventID, p.Gender, p.DateOfBirth)
	if err != nil {
		return fmt.Errorf("error assigning category for participant with bib %d", p.Bib)
	}
	p.CategoryID = catID
	return nil
}

func (ps ParticipantService) UpdateParticipant(req entity.ParticipantUpdateRequest) (*entity.Participant, error) {
	p, err := ps.repo.GetParticipantByID(req.ID)
	if err != nil {
		return nil, fmt.Errorf("updateParticipant: participant with ID %s not found", req.ID)
	}
	newP, err := entity.NewParticipant(req.ParticipantCreateRequest)
	if err != nil {
		return nil, err
	}
	newP.ID = p.ID

	err = ps.repo.SaveParticipant(newP)
	if err != nil {
		return nil, err
	}
	return newP, nil
}

func (ps ParticipantService) DeleteParticipant(raceID, id uuid.UUID) error {
	_, err := ps.repo.GetParticipantByID(id)
	if err != nil {
		return fmt.Errorf("deleteParticipant: participant with ID %s not found", id)
	}
	err = ps.repo.DeleteParticipant(raceID, id)
	if err != nil {
		return fmt.Errorf("delete participant: error deleting participant from DB")
	}
	return nil
}
