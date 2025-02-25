package entity

import "github.com/google/uuid"

type EventAthlete struct {
	RaceID     uuid.UUID `json:"race_id"`
	EventID    uuid.UUID `json:"event_id"`
	AthleteID  uuid.UUID `json:"athlete_id"`
	WaveID     uuid.UUID `json:"wave_id"`
	CategoryID uuid.UUID `json:"category_id"`
	Bib        int       `json:"bib"`
}
