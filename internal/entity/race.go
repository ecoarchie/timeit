package entity

import (
	"time"

	"github.com/ecoarchie/timeit/internal/controller/httpv1/dto"
	"github.com/ecoarchie/timeit/pkg/validator"
	"github.com/google/uuid"
)

type Race struct {
	ID       uuid.UUID `json:"race_id"`
	Name     string    `json:"race_name"`
	Timezone string    `json:"timezone"`
}

func NewRace(req *dto.RaceDTO, v *validator.Validator) *Race {
	v.Check(IsIANATimezone(req.Timezone), "timezone", "must be valid IANA timezone")
	if !v.Valid() {
		return nil
	}

	return &Race{
		ID:       req.ID,
		Name:     req.Name,
		Timezone: req.Timezone,
	}
}

func IsIANATimezone(tz string) bool {
	_, err := time.LoadLocation(tz) // tz must correspond to IANA time zones names
	return err == nil
}

type RaceModel struct {
	*Race
	TimeReaders []*TimeReader
	Events      []*Event
}
