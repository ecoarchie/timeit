package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Race struct {
	ID       uuid.UUID `json:"race_id"`
	Name     string    `json:"race_name"`
	Timezone string    `json:"timezone"`
}

func NewRace(req *RaceFormData) (*Race, error) {
	if !IsIANATimezone(req.Timezone) {
		return nil, fmt.Errorf("not valid IANA timezone")
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
		Timezone: req.Timezone,
	}, nil
}

func IsIANATimezone(tz string) bool {
	_, err := time.LoadLocation(tz) // tz must correspond to IANA time zones names
	return err == nil
}
