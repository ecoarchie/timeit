package entity

import "github.com/google/uuid"

type ChipBib struct {
	RaceID  uuid.UUID `json:"race_id"`
	EventID uuid.UUID `json:"event_id"`
	Chip    int       `json:"chip"`
	Bib     int       `json:"bib"`
}
