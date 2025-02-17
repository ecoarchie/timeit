package entity

import (
	"github.com/google/uuid"
)

type TPType string

const (
	TPTypeStart    TPType = "start"
	TPTypeStandard TPType = "standard"
	TPTypeFinish   TPType = "finish"
)

type TimingPoint struct {
	ID                uuid.UUID `json:"timing_point_id"`
	RaceID            uuid.UUID `json:"race_id"`
	EventID           uuid.UUID `json:"event_id"`
	Name              string    `json:"timing_point_name"`
	Type              TPType    `json:"timing_point_type"`
	DistanceFromStart int       `json:"distance_from_start"`
	BoxName           string    `json:"box_name"`
	MinTimeSec        int64     `json:"min_time_sec"`
	MaxTimeSec        int64     `json:"max_time_sec"`
	MinLapTimeSec     int64     `json:"min_lap_time_sec"`
}

// func (tp *TimingPoint) SetValidMinMaxTimes(participantWaceStart time.Time) {
// 	if tp.MinTimeSec == 0 {
// 		tp.ValidMinTime = participantWaceStart
// 	} else {
// 		tp.ValidMinTime = participantWaceStart.Add(time.Duration(tp.MinTimeSec) * time.Second)
// 	}
// 	if tp.MaxTimeSec == 0 {
// 		tp.ValidMaxTime = participantWaceStart.Add(time.Duration(time.Hour) * 24)
// 	} else {
// 		tp.ValidMaxTime = participantWaceStart.Add(time.Duration(tp.MaxTimeSec) * time.Second)
// 	}
// }

type NewTPrequest struct {
	RaceID            uuid.UUID `json:"race_id"`
	EventID           uuid.UUID `json:"event_id"`
	Name              string    `json:"name"`
	Type              TPType    `json:"type"`
	DistanceFromStart int       `json:"distance_from_start"`
	BoxName           string    `json:"box_name"`
	MinTimeSec        int64     `json:"min_time_sec"`
	MaxTimeSec        int64     `json:"max_time_sec"`
	MinLapTimeSec     int64     `json:"min_lap_time_sec"`
}

func IsValidTPType(tp TPType) bool {
	switch tp {
	case TPTypeStart, TPTypeFinish, TPTypeStandard:
		return true
	default:
		return false
	}
}

// func NewTimingPoint(raceID uuid.UUID, eventID uuid.UUID, req TimingPointFormData) (*TimingPoint, error) {
// 	if raceID == uuid.Nil {
// 		return nil, fmt.Errorf("empty raceID")
// 	}
// 	if eventID == uuid.Nil {
// 		return nil, fmt.Errorf("empty eventID")
// 	}
// 	if req.Name == "" {
// 		return nil, fmt.Errorf("empty timing point name")
// 	}
// 	if req.DistanceFromStart < 0 {
// 		return nil, fmt.Errorf("distance from start must be equal or greater than 0")
// 	}
// 	if req.BoxName == "" {
// 		return nil, fmt.Errorf("empty box name")
// 	}
// 	if req.MinTimeSec < 0 || req.MaxTimeSec < 0 || req.MinLapTimeSec < 0 {
// 		return nil, fmt.Errorf("min, max and lap times must be equal or greater than 0")
// 	}

// 	id := uuid.New()
// 	return &TimingPoint{
// 		ID:                id,
// 		RaceID:            raceID,
// 		EventID:           eventID,
// 		Name:              req.Name,
// 		Type:              req.Type,
// 		DistanceFromStart: req.DistanceFromStart,
// 		BoxName:           req.BoxName,
// 		MinTimeSec:        req.MinTimeSec,
// 		MaxTimeSec:        req.MaxTimeSec,
// 		MinLapTimeSec:     req.MinLapTimeSec,
// 	}, nil
// }

func RandomTimingPoint(name string, typ TPType, dst int, boxName string, min, max, lap int64) *TimingPoint {
	return &TimingPoint{
		ID:                uuid.New(),
		RaceID:            uuid.New(),
		EventID:           uuid.New(),
		Name:              name,
		Type:              typ,
		DistanceFromStart: dst,
		BoxName:           boxName,
		MinTimeSec:        min,
		MaxTimeSec:        max,
		MinLapTimeSec:     lap,
	}
}
