package entity

import "github.com/google/uuid"

type PhysicalLocation struct {
	RaceID  uuid.UUID `json:"race_id"`
	BoxName string    `json:"box_name"`
}
