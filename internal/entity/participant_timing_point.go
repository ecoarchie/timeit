package entity

import (
	"time"

	"github.com/google/uuid"
)

type ParticipantTimingPoint struct {
	// RaceID        uuid.UUID `json:"race_id"`
	// EventID       uuid.UUID `json:"event_id"`
	// ParticipantID uuid.UUID `json:"participant_id"`
	TimingPointID uuid.UUID `json:"timing_point_id"`
	TOD           time.Time `json:"tod"`
	GunTime       int64     `json:"gun_time"`
	NetTime       int64     `json:"net_time"`
}
