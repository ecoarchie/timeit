package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Race struct {
	ID       uuid.UUID `json:"race_id"`
	Name     string    `json:"race_name"`
	RaceDate time.Time `json:"race_date"`
	Timezone string    `json:"timezone"`
}

func NewRace(req RaceFormData) (*Race, error) {
	if err := IsValidTimezone(req.Timezone); err != nil {
		return nil, err
	}
	if req.Name == "" {
		return nil, fmt.Errorf("empty race name")
	}
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, fmt.Errorf("invalid race id")
	}
	return &Race{
		ID:       id,
		Name:     req.Name,
		RaceDate: req.RaceDate,
		Timezone: req.Timezone,
	}, nil
}

func IsValidTimezone(tz string) error {
	if tz == "" {
		return fmt.Errorf("empty timezone in race config")
	}
	_, err := time.LoadLocation(tz) // tz must correspond to IANA time zones names
	return err
}
