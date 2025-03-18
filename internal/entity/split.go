package entity

import (
	"fmt"
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
	ID                 uuid.UUID
	RaceID             uuid.UUID
	EventID            uuid.UUID
	Name               string
	Type               SplitType
	DistanceFromStart  int
	TimeReaderID       uuid.UUID
	MinTime            time.Duration
	MaxTime            time.Duration
	MinLapTime         time.Duration
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

func (s Split) String() string {
	return fmt.Sprintf(
		"Split {\n"+
			"  ID: %s\n"+
			"  RaceID: %s\n"+
			"  EventID: %s\n"+
			"  Name: %q\n"+
			"  Type: %s\n"+
			"  DistanceFromStart: %d meters\n"+
			"  TimeReaderID: %s\n"+
			"  MinTime: %s\n"+
			"  MaxTime: %s\n"+
			"  MinLapTime: %s\n"+
			"  PreviousLapSplitID: %s\n"+
			"}",
		s.ID,
		s.RaceID,
		s.EventID,
		s.Name,
		s.Type,
		s.DistanceFromStart,
		s.TimeReaderID,
		formatDuration(s.MinTime),
		formatDuration(s.MaxTime),
		formatDuration(s.MinLapTime),
		formatNullUUID(s.PreviousLapSplitID),
	)
}

func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func formatNullUUID(n uuid.NullUUID) string {
	if !n.Valid {
		return "null"
	}
	return n.UUID.String()
}
