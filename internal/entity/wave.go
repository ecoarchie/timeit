package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Wave struct {
	ID         uuid.UUID `json:"id"`
	RaceID     uuid.UUID `json:"race_id"`
	EventID    uuid.UUID `json:"event_id"`
	Name       string    `json:"name"`
	StartTime  time.Time `json:"start_time"`
	IsLaunched bool      `json:"is_launched"`
}

func NewWave(raceID uuid.UUID, eventID uuid.UUID, name string, startTime time.Time) (*Wave, error) {
	if raceID == uuid.Nil {
		return nil, fmt.Errorf("empty raceID")
	}
	if eventID == uuid.Nil {
		return nil, fmt.Errorf("empty eventID")
	}
	if name == "" {
		return nil, fmt.Errorf("empty wave name")
	}
	id := uuid.New()
	return &Wave{
		ID:         id,
		RaceID:     raceID,
		EventID:    eventID,
		Name:       name,
		StartTime:  startTime,
		IsLaunched: false,
	}, nil
}
