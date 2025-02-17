package entity

import "github.com/google/uuid"

type PhysicalLocation struct {
	ID      uuid.UUID `json:"location_id"`
	RaceID  uuid.UUID `json:"race_id"`
	BoxName string    `json:"box_name"`
}
