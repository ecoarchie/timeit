package entity

import (
	"fmt"
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
	ID       uuid.UUID
	RaceID   uuid.UUID
	EventID  uuid.UUID
	Name     string
	Gender   CategoryGender
	AgeFrom  int
	DateFrom time.Time
	AgeTo    int
	DateTo   time.Time
}

func NewCategory(dto *dto.CategoryDTO, eventDate time.Time, v *validator.Validator) *Category {
	v.Check(IsValidGender(CategoryGender(dto.Gender)), "gender", "must be male, female or mixed")
	v.Check(dto.AgeFrom >= 0, "category age from", "must be greater or equal to 0")
	v.Check(dto.AgeTo > 0, "category age to", "must be greater than 0")
	v.Check(dto.AgeFrom < dto.AgeTo, "category age", "upper age limit must be greater than lower age limit")
	if !v.Valid() {
		return nil
	}

	dateFrom, dateTo := GetDateRange(dto, eventDate)
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

// TODO Write test for it
func (c *Category) Valid(dob time.Time) bool {
	return (dob.Before(c.DateTo) || dob.Equal(c.DateTo)) && (dob.After(c.DateFrom) || dob.Equal(c.DateFrom))
}

func GetDateRange(dto *dto.CategoryDTO, eventDate time.Time) (dateFrom, dateTo time.Time) {
	if dto.FromRaceDate {
		dateTo = time.Date(eventDate.Year()-dto.AgeFrom, eventDate.Month(), eventDate.Day(), 23, 59, 59, 0, eventDate.Location())
	} else {
		dateTo = time.Date(eventDate.Year()-dto.AgeFrom, time.December, 31, 23, 59, 59, 0, eventDate.Location())
	}
	if dto.ToRaceDate {
		dateFrom = time.Date(eventDate.Year()-dto.AgeTo-1, eventDate.Month(), eventDate.Day()+1, 0, 0, 0, 0, eventDate.Location())
	} else {
		dateFrom = time.Date(eventDate.Year()-dto.AgeTo, time.January, 1, 0, 0, 0, 0, eventDate.Location())
	}
	return dateFrom, dateTo
}

func (c Category) String() string {
	return fmt.Sprintf(
		"Category {\n"+
			"  ID: %s\n"+
			"  RaceID: %s\n"+
			"  EventID: %s\n"+
			"  Name: %q\n"+
			"  Gender: %s\n"+
			"  Age Range: %d - %d years\n"+
			"  Date Range: %s - %s\n"+
			"}",
		c.ID,
		c.RaceID,
		c.EventID,
		c.Name,
		c.Gender,
		c.AgeFrom, c.AgeTo,
		c.DateFrom.Format(time.DateOnly), c.DateTo.Format(time.DateOnly),
	)
}
