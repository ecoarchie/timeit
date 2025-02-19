package testdata

import (
	"time"

	"github.com/ecoarchie/timeit/internal/entity"
	"github.com/google/uuid"
)

type ParticipantRepoMock struct{}

func (prm ParticipantRepoMock) SaveParticipant(p *entity.Participant) error {
	return nil
}

func (prm ParticipantRepoMock) GetCategoryFor(eventID uuid.UUID, gender entity.CategoryGender, dob time.Time) (uuid.NullUUID, error) {
	return uuid.NullUUID{}, nil
}

func (prm ParticipantRepoMock) GetParticipantWithChip(chip int) (*entity.Participant, error) {
	return nil, nil
}

func (prm ParticipantRepoMock) GetParticipantByID(id uuid.UUID) (*entity.Participant, error) {
	return nil, nil
}

func (prm ParticipantRepoMock) DeleteParticipant(raceID uuid.UUID, id uuid.UUID) error {
	return nil
}
