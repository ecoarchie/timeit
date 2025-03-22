package entity

import (
	"cmp"
	"fmt"
	"slices"
	"strings"
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
	Gender   CategoryGender `json:"gender"`
	AgeFrom  int            `json:"age_from"`
	DateFrom time.Time      `json:"date_from"`
	AgeTo    int            `json:"age_to"`
	DateTo   time.Time      `json:"date_to"`
}

func GenderFrom(g string) CategoryGender {
	switch strings.ToLower(g) {
	case "male":
		return CategoryGenderMale
	case "female":
		return CategoryGenderFemale
	case "mixed":
		return CategoryGenderMixed
	case "unknown":
		return CategoryGenderUnknown
	default:
		return CategoryGenderUnknown
	}
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
func (c *Category) Valid(gender CategoryGender, dob time.Time) bool {
	return c.Gender == gender && ((dob.Before(c.DateTo) || dob.Equal(c.DateTo)) && (dob.After(c.DateFrom) || dob.Equal(c.DateFrom)))
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

func CheckCategoriesBoundary(cats []*Category, v *validator.Validator) {
	slices.SortFunc(cats, func(a, b *Category) int {
		return cmp.Compare(a.AgeFrom, b.AgeFrom)
	})
	genderMap := map[CategoryGender][]*Category{}
	for _, c := range cats {
		genderMap[c.Gender] = append(genderMap[c.Gender], c)
	}
	for _, cc := range genderMap {
		if len(cc) > 1 {
			for i := 1; i < len(cc)-1; i++ {
				if cc[i].DateTo.After(cc[i-1].DateFrom) || cc[i-1].DateFrom.Sub(cc[i].DateTo) > time.Duration(time.Hour*24) {
					v.AddError(cc[i-1].Name+" and "+cc[i].Name, "not consequent dates")
					fmt.Println("DateTo, DateFrom", cc[i].DateTo, cc[i-1].DateFrom)
				} else if cc[i].DateFrom.Before(cc[i+1].DateTo) || cc[i].DateFrom.Sub(cc[i+1].DateTo) > time.Duration(time.Hour*24) {
					v.AddError(cc[i].Name+" and "+cc[i+1].Name, "not consequent dates")
					fmt.Println("DateTo, DateFrom", cc[i].DateTo, cc[i+1].DateFrom)
				}
			}
		}
	}
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
