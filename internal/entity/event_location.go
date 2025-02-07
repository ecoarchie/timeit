package entity

import "github.com/google/uuid"

type EventLocation struct {
	EventID uuid.UUID `json:"event_id"`
	RaceID  uuid.UUID `json:"race_id"`
	BoxName string    `json:"box_name"`
}
