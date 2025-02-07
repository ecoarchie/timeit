package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Race struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	RaceDate time.Time `json:"race_date"`
	Timezone string    `json:"timezone"`
}

func NewRace(name string, raceDate time.Time, tz string) (*Race, error) {
	if err := isValidTimezone(tz); err != nil {
		return nil, err
	}
	if name == "" {
		return nil, fmt.Errorf("empty race name")
	}
	id := uuid.New()
	return &Race{
		ID:       id,
		Name:     name,
		RaceDate: raceDate,
		Timezone: tz,
	}, nil
}

func isValidTimezone(tz string) error {
	if tz == "" {
		return fmt.Errorf("empty timezone")
	}
	_, err := time.LoadLocation(tz) // tz must correspond to IANA time zones names
	return err
}
