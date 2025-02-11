package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type CategoryGender string

const (
	CategoryGenderMale    CategoryGender = "male"
	CategoryGenderFemale  CategoryGender = "female"
	CategoryGenderMixed   CategoryGender = "mixed"
	CategoryGenderUnknown CategoryGender = "unknown"
)

type Category struct {
	ID            uuid.UUID      `json:"category_id"`
	RaceID        uuid.UUID      `json:"race_id"`
	EventID       uuid.UUID      `json:"event_id"`
	Name          string         `json:"category_name"`
	Gender        CategoryGender `json:"category_gender"`
	FromAge       int            `json:"from_age"`
	FromRaceDate  bool           `json:"from_race_date"`
	ToAge         int            `json:"to_age"`
	ToRaceDate    bool           `json:"to_race_date"`
	BirthDateFrom time.Time
	BirthDateTo   time.Time
}

type CategoryFormData struct {
	RaceID       uuid.UUID      `json:"race_id"`
	EventID      uuid.UUID      `json:"event_id"`
	Name         string         `json:"name"`
	Gender       CategoryGender `json:"gender"`
	FromAge      int            `json:"from_age"`
	FromRaceDate bool           `json:"from_race_date"`
	ToAge        int            `json:"to_age"`
	ToRaceDate   bool           `json:"to_race_date"`
}

func NewCategory(raceID uuid.UUID, eventID uuid.UUID, req CategoryFormData) (*Category, error) {
	if raceID == uuid.Nil {
		return nil, fmt.Errorf("empty raceID")
	}
	if eventID == uuid.Nil {
		return nil, fmt.Errorf("empty eventID")
	}
	if req.Name == "" {
		return nil, fmt.Errorf("empty timing point name")
	}
	if req.FromAge < 0 {
		return nil, fmt.Errorf("from age must be greater or equal to 0")
	}
	if req.ToAge < 0 {
		return nil, fmt.Errorf("to age must be greater or equal to 0")
	}
	if req.FromAge > req.ToAge {
		return nil, fmt.Errorf("from age must be less than to age")
	}

	id := uuid.New()

	return &Category{
		ID:           id,
		RaceID:       raceID,
		EventID:      eventID,
		Name:         req.Name,
		Gender:       req.Gender,
		FromAge:      req.FromAge,
		FromRaceDate: req.FromRaceDate,
		ToAge:        req.ToAge,
		ToRaceDate:   req.ToRaceDate,
	}, nil
}

func (c Category) ValidForGenderAndDateOfBirth(gender CategoryGender, dob time.Time, eventDate time.Time) bool {
	if gender != c.Gender {
		return false
	}
	return (dob.After(c.BirthDateFrom) || dob.Equal(c.BirthDateFrom)) && (dob.Before(c.BirthDateTo) || dob.Equal(c.BirthDateTo))
}
