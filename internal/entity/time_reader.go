package entity

import "github.com/google/uuid"

type TimeReader struct {
	ID         uuid.UUID `json:"time_reader_id"`
	RaceID     uuid.UUID `json:"race_id"`
	ReaderName string    `json:"reader_name"`
}
