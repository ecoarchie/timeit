package entity

import (
	"time"

	"github.com/google/uuid"
)

type ReaderRecord struct {
	ID         int       `json:"id"`
	RaceID     uuid.UUID `json:"race_id"`
	Chip       int       `json:"chip"`
	TOD        time.Time `json:"tod"`
	ReaderName string    `json:"reader_name"`
	CanUse     bool      `json:"can_use"`
	// Type uint `json:"type"`
}
