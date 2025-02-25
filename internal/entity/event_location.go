package entity

import "github.com/google/uuid"

type EventTimeReader struct {
	EventID    uuid.UUID `json:"event_id"`
	RaceID     uuid.UUID `json:"race_id"`
	ReaderName string    `json:"reader_name"`
}
