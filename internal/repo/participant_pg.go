package repo

import (
	"time"

	"github.com/ecoarchie/timeit/internal/entity"
	"github.com/google/uuid"
)

type AthleteRepoPG struct{}

func (pr AthleteRepoPG) SaveAthlete(p *entity.Athlete) error {
	return nil
}

func (pr AthleteRepoPG) GetCategoryFor(eventID uuid.UUID, gender entity.CategoryGender, dob time.Time) (uuid.NullUUID, error) {
	return uuid.NullUUID{}, nil
}

func (pr AthleteRepoPG) GetAthleteWithChip(chip int) (*entity.Athlete, error) {
	return nil, nil
}

func (pr AthleteRepoPG) GetAthleteByID(raceID, athleteID uuid.UUID) (*entity.Athlete, error) {
	return nil, nil
}

func (pr AthleteRepoPG) DeleteAthlete(raceID uuid.UUID, id uuid.UUID) error {
	return nil
}

func NewAthletePGRepo() *AthleteRepoPG {
	return &AthleteRepoPG{}
}
