package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID               uuid.UUID `json:"event_id"`
	RaceID           uuid.UUID `json:"race_id"`
	Name             string    `json:"event_name"`
	DistanceInMeters int       `json:"distance_in_meters"`
	EventDate        time.Time `json:"event_date"`
}

func NewEvent(raceID uuid.UUID, name string, distanse int, eventDate time.Time) (*Event, error) {
	if raceID == uuid.Nil {
		return nil, fmt.Errorf("empty raceID")
	}
	if name == "" {
		return nil, fmt.Errorf("empty event name")
	}
	if distanse <= 0 {
		return nil, fmt.Errorf("distance must be greater than 0")
	}
	id := uuid.New()
	return &Event{
		ID:               id,
		RaceID:           raceID,
		Name:             name,
		DistanceInMeters: distanse,
		EventDate:        eventDate,
	}, nil
}
