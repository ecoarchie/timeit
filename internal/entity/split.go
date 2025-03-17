package entity

import (
	"time"

	"github.com/ecoarchie/timeit/internal/controller/httpv1/dto"
	"github.com/ecoarchie/timeit/pkg/validator"
	"github.com/google/uuid"
)

type SplitType string

const (
	SplitTypeStart    SplitType = "start"
	SplitTypeStandard SplitType = "standard"
	SplitTypeFinish   SplitType = "finish"
)

type Split struct {
	ID                 uuid.UUID     `json:"split_id"`
	RaceID             uuid.UUID     `json:"race_id"`
	EventID            uuid.UUID     `json:"event_id"`
	Name               string        `json:"split_name"`
	Type               SplitType     `json:"split_type"`
	DistanceFromStart  int           `json:"distance_from_start"`
	TimeReaderID       uuid.UUID     `json:"time_reader_id"`
	MinTime            time.Duration `json:"min_time_sec"`
	MaxTime            time.Duration `json:"max_time_sec"`
	MinLapTime         time.Duration `json:"min_lap_time_sec"`
	PreviousLapSplitID uuid.NullUUID
}

func NewSplit(dto *dto.SplitDTO, trs []*dto.TimeReaderDTO, v *validator.Validator) *Split {
	v.Check(IsValidSplitType(SplitType(dto.Type)), "split type", "must be start, standard or finish")
	var tpIDsForLocs []uuid.UUID
	for _, l := range trs {
		tpIDsForLocs = append(tpIDsForLocs, l.ID)
	}

	v.Check(validator.PermittedValue(dto.TimeReaderID, tpIDsForLocs...), "split ID", "must have valid corresponded time reader")

	v.Check(dto.DistanceFromStart >= 0, "split distance from start", "must be greater or equal to 0")
	minTime, _ := time.ParseDuration(dto.MinTime)
	maxTime, _ := time.ParseDuration(dto.MaxTime)
	minLapTime, _ := time.ParseDuration(dto.MinLapTime)
	v.Check(minTime >= 0, "split min time", "must be greater or equal to 0")
	v.Check(maxTime >= 0, "split max time", "must be greater or equal to 0")
	v.Check(minLapTime >= 0, "split min lap time", "must be greater or equal to 0")

	if !v.Valid() {
		return nil
	}
	return &Split{
		ID:                 dto.ID,
		RaceID:             dto.RaceID,
		EventID:            dto.EventID,
		Name:               dto.Name,
		Type:               SplitType(dto.Type),
		DistanceFromStart:  dto.DistanceFromStart,
		TimeReaderID:       dto.TimeReaderID,
		MinTime:            minTime,
		MaxTime:            maxTime,
		MinLapTime:         minLapTime,
		PreviousLapSplitID: uuid.NullUUID{},
	}
}

func IsValidSplitType(tp SplitType) bool {
	switch tp {
	case SplitTypeStart, SplitTypeFinish, SplitTypeStandard:
		return true
	default:
		return false
	}
}

func (s *Split) IsValidForRecord(waveStart time.Time, tod time.Time, prev *AthleteSplit) bool {
	if tod.Before(waveStart) {
		return false
	}
	validMinTime := waveStart.Add(time.Duration(s.MinTime))

	if !(tod.After(validMinTime) || tod.Equal(validMinTime)) {
		return false
	}
	if s.MaxTime != 0 && !(tod.Before(waveStart.Add(time.Duration(s.MaxTime))) || tod.Equal(waveStart.Add(time.Duration(s.MaxTime)))) {
		return false
	}
	if prev != nil {
		if prev.TOD.Add(time.Duration(s.MinLapTime)).After(tod) {
			return false
		}
	}
	return true
}
