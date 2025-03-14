package entity

import (
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
	ID       uuid.UUID      `json:"category_id"`
	RaceID   uuid.UUID      `json:"race_id"`
	EventID  uuid.UUID      `json:"event_id"`
	Name     string         `json:"category_name"`
	Gender   CategoryGender `json:"category_gender"`
	AgeFrom  int            `json:"age_from"`
	DateFrom time.Time      `json:"date_from"`
	AgeTo    int            `json:"age_to"`
	DateTo   time.Time      `json:"date_to"`
}

// type CategoryFormData struct {
// 	RaceID   uuid.UUID      `json:"race_id"`
// 	EventID  uuid.UUID      `json:"event_id"`
// 	Name     string         `json:"name"`
// 	Gender   CategoryGender `json:"gender"`
// 	AgeFrom  int            `json:"age_from"`
// 	DateFrom time.Time      `json:"date_from"`
// 	AgeTo    int            `json:"age_to"`
// 	DateTo   time.Time      `json:"date_to"`
// }

// func NewCategory(raceID uuid.UUID, eventID uuid.UUID, req CategoryFormData) (*Category, error) {
// 	if raceID == uuid.Nil {
// 		return nil, fmt.Errorf("empty raceID")
// 	}
// 	if eventID == uuid.Nil {
// 		return nil, fmt.Errorf("empty eventID")
// 	}
// 	if req.Name == "" {
// 		return nil, fmt.Errorf("empty split name")
// 	}
// 	if req.AgeFrom < 0 {
// 		return nil, fmt.Errorf("from age must be greater or equal to 0")
// 	}
// 	if req.AgeTo < 0 {
// 		return nil, fmt.Errorf("to age must be greater or equal to 0")
// 	}
// 	if req.AgeFrom > req.AgeTo {
// 		return nil, fmt.Errorf("from age must be less than to age")
// 	}

// 	id := uuid.New()

// 	return &Category{
// 		ID:       id,
// 		RaceID:   raceID,
// 		EventID:  eventID,
// 		Name:     req.Name,
// 		Gender:   req.Gender,
// 		AgeFrom:  req.AgeFrom,
// 		DateFrom: req.DateFrom,
// 		AgeTo:    req.AgeTo,
// 		DateTo:   req.DateTo,
// 	}, nil
// }

func (c Category) ValidForGenderAndDateOfBirth(gender CategoryGender, dob time.Time) bool {
	if gender != c.Gender {
		return false
	}

	bdFrom := c.BirthDateFrom(dob)
	bdTo := c.BirthDateTo(dob)
	return (dob.After(bdFrom) || dob.Equal(bdFrom)) && (dob.Before(bdTo) || dob.Equal(bdTo))
}

func (c Category) BirthDateFrom(dob time.Time) time.Time {
	return c.DateFrom.AddDate(-c.AgeFrom, 0, 0)
}

func (c Category) BirthDateTo(dob time.Time) time.Time {
	return c.DateTo.AddDate(-c.AgeTo, 0, 0)
}
