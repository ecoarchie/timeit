package service

import (
	"fmt"
	"time"

	"github.com/ecoarchie/timeit/internal/entity"
	"github.com/ecoarchie/timeit/pkg/logger"
	"github.com/google/uuid"
)

type ParticipantEditor interface{}

type ParticipantRepo interface {
	SaveParticipant(p *entity.Participant) error
	GetCategoryFor(eventID uuid.UUID, gender entity.CategoryGender, dob time.Time) (uuid.UUID, error)
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

func (ps ParticipantService) CreateParticipant(req entity.ParticipantCreateRequest) (uuid.UUID, error) {
	p, err := entity.NewParticipant(req)
	if err != nil {
		return uuid.Nil, err
	}

	if !req.CategoryID.Valid {
		ps.assignCategory(p)
	}

	err = ps.repo.SaveParticipant(p)
	if err != nil {
		return uuid.Nil, err
	}

	// TODO create and store entity EventParticipant. DO it in repo layer in transaction
	return p.ID, nil
}

func (ps ParticipantService) assignCategory(p *entity.Participant) error {
	catID, err := ps.repo.GetCategoryFor(p.EventID, p.Gender, p.DateOfBirth)
	if err != nil {
		return fmt.Errorf("error assigning category for participant with bib %d", p.Bib)
	}
	p.CategoryID = uuid.NullUUID{
		UUID:  catID,
		Valid: true,
	}
	return nil
}
