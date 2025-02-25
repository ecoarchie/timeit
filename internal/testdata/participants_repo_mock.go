package testdata

import (
	"time"

	"github.com/ecoarchie/timeit/internal/entity"
	"github.com/google/uuid"
)

type AthleteRepoMock struct{}

func (prm AthleteRepoMock) SaveAthlete(p *entity.Athlete) error {
	return nil
}

func (prm AthleteRepoMock) GetCategoryFor(eventID uuid.UUID, gender entity.CategoryGender, dob time.Time) (uuid.NullUUID, error) {
	return uuid.NullUUID{}, nil
}

func (prm AthleteRepoMock) GetAthleteWithChip(chip int) (*entity.Athlete, error) {
	return nil, nil
}

func (prm AthleteRepoMock) GetAthleteByID(raceID, athleteID uuid.UUID) (*entity.Athlete, error) {
	return nil, nil
}

func (prm AthleteRepoMock) DeleteAthlete(raceID uuid.UUID, id uuid.UUID) error {
	return nil
}
