package entity

import (
	"github.com/ecoarchie/timeit/internal/controller/httpv1/dto"
	"github.com/ecoarchie/timeit/pkg/validator"
	"github.com/google/uuid"
)

type TimeReader struct {
	ID         uuid.UUID `json:"time_reader_id"`
	RaceID     uuid.UUID `json:"race_id"`
	ReaderName string    `json:"reader_name"`
}

func NewTimeReader(dto *dto.TimeReaderDTO, v *validator.Validator) *TimeReader {
	v.Check(dto.ID != uuid.Nil, "time_reader_id", "must be valid UUID")
	v.Check(dto.RaceID != uuid.Nil, "reader_race_id", "must be valid UUID")
	v.Check(dto.ReaderName != "", "reader_name", "must not be empty")
	if !v.Valid() {
		return nil
	}
	return &TimeReader{
		ID:         dto.ID,
		RaceID:     dto.RaceID,
		ReaderName: dto.ReaderName,
	}
}
