package entity

import (
	"time"

	"github.com/ecoarchie/timeit/internal/controller/httpv1/dto"
	"github.com/ecoarchie/timeit/pkg/validator"
	"github.com/google/uuid"
)

type Wave struct {
	ID         uuid.UUID `json:"wave_id"`
	RaceID     uuid.UUID `json:"race_id"`
	EventID    uuid.UUID `json:"event_id"`
	Name       string    `json:"wave_name"`
	StartTime  time.Time `json:"wave_start_time"`
	IsLaunched bool      `json:"is_launched"`
}

type WaveStart struct {
	WaveID    uuid.UUID `json:"wave_id"`
	StartTime time.Time `json:"wave_start_time"`
}

func NewWave(dto *dto.WaveDTO, v *validator.Validator) *Wave {
	startTime, _ := time.Parse(time.RFC3339, dto.StartTime)
	return &Wave{
		ID:         dto.ID,
		RaceID:     dto.RaceID,
		EventID:    dto.EventID,
		Name:       dto.Name,
		StartTime:  startTime,
		IsLaunched: dto.IsLaunched,
	}
}
