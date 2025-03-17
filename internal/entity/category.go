package entity

import (
	"time"

	"github.com/ecoarchie/timeit/internal/controller/httpv1/dto"
	"github.com/ecoarchie/timeit/pkg/validator"
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

func NewCategory(dto *dto.CategoryDTO, v *validator.Validator) *Category {
	v.Check(IsValidGender(CategoryGender(dto.Gender)), "gender", "must be male, female or mixed")
	v.Check(dto.AgeFrom >= 0, "category age from", "must be greater or equal to 0")
	v.Check(dto.AgeTo > 0, "category age to", "must be greater than 0")
	v.Check(dto.AgeFrom < dto.AgeTo, "category age", "upper age limit must be greater than lower age limit")
	if !v.Valid() {
		return nil
	}

	dateFrom, _ := time.Parse(time.DateOnly, dto.DateFrom)
	dateTo, _ := time.Parse(time.DateOnly, dto.DateTo)

	// FIXME rewrite so from ui come only 'on race date' flag true or false, and calculate dateFrom and dateTo depending on these flags
	return &Category{
		ID:       dto.ID,
		RaceID:   dto.RaceID,
		EventID:  dto.EventID,
		Name:     dto.Name,
		Gender:   CategoryGender(dto.Gender),
		AgeFrom:  dto.AgeFrom,
		DateFrom: dateFrom,
		AgeTo:    dto.AgeTo,
		DateTo:   dateTo,
	}
}

// FIXME add method to check valid category, now it fetches from DB for each athlete while creating one
// func (c *Category) ValidForGenderAndDateOfBirth(gender CategoryGender, dob time.Time) bool {
// 	if gender != c.Gender {
// 		return false
// 	}

// 	bdFrom := c.BirthDateFrom(dob)
// 	bdTo := c.BirthDateTo(dob)
// 	return (dob.After(bdFrom) || dob.Equal(bdFrom)) && (dob.Before(bdTo) || dob.Equal(bdTo))
// }

// func (c *Category) BirthDateFrom(dob time.Time) time.Time {
// 	return c.DateFrom.AddDate(-c.AgeFrom, 0, 0)
// }

// func (c *Category) BirthDateTo(dob time.Time) time.Time {
// 	return c.DateTo.AddDate(-c.AgeTo, 0, 0)
// }
