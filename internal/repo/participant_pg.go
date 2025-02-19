package repo

import (
	"time"

	"github.com/ecoarchie/timeit/internal/entity"
	"github.com/google/uuid"
)

type ParticipantRepoPG struct{}

func (pr ParticipantRepoPG) SaveParticipant(p *entity.Participant) error {
	return nil
}

func (pr ParticipantRepoPG) GetCategoryFor(eventID uuid.UUID, gender entity.CategoryGender, dob time.Time) (uuid.NullUUID, error) {
	return uuid.NullUUID{}, nil
}

func (pr ParticipantRepoPG) GetParticipantWithChip(chip int) (*entity.Participant, error) {
	return nil, nil
}

func (pr ParticipantRepoPG) GetParticipantByID(id uuid.UUID) (*entity.Participant, error) {
	return nil, nil
}

func (pr ParticipantRepoPG) DeleteParticipant(raceID uuid.UUID, id uuid.UUID) error {
	return nil
}

func NewParticipantPGRepo() *ParticipantRepoPG {
	return &ParticipantRepoPG{}
}
