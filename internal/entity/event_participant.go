package entity

import "github.com/google/uuid"

type EventParticipant struct {
	RaceID        uuid.UUID `json:"race_id"`
	EventID       uuid.UUID `json:"event_id"`
	ParticipantID uuid.UUID `json:"participant_id"`
	WaveID        uuid.UUID `json:"wave_id"`
	CategoryID    uuid.UUID `json:"category_id"`
	Bib           int       `json:"bib"`
}
