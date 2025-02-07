package entity

import (
	"time"

	"github.com/google/uuid"
)

type BoxRecord struct {
	ID      int       `json:"id"`
	RaceID  uuid.UUID `json:"race_id"`
	Chip    int       `json:"chip"`
	TOD     time.Time `json:"tod"`
	BoxName string    `json:"box_name"`
	CanUse  bool      `json:"can_use"`
}
