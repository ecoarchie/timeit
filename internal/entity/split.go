package entity

import (
	"time"

	"github.com/google/uuid"
)

type SplitType string

const (
	SplitTypeStart    SplitType = "start"
	SplitTypeStandard SplitType = "standard"
	SplitTypeFinish   SplitType = "finish"
)

type Split struct {
	ID                uuid.UUID     `json:"split_id"`
	RaceID            uuid.UUID     `json:"race_id"`
	EventID           uuid.UUID     `json:"event_id"`
	Name              string        `json:"split_name"`
	Type              SplitType     `json:"split_type"`
	DistanceFromStart int           `json:"distance_from_start"`
	TimeReaderID      uuid.UUID     `json:"time_reader_id"`
	MinTime           time.Duration `json:"min_time_sec"`
	MaxTime           time.Duration `json:"max_time_sec"`
	MinLapTime        time.Duration `json:"min_lap_time_sec"`
}

// func (tp *Split) SetValidMinMaxTimes(athleteWaceStart time.Time) {
// 	if tp.MinTime == 0 {
// 		tp.ValidMinTime = athleteWaceStart
// 	} else {
// 		tp.ValidMinTime = athleteWaceStart.Add(time.Duration(tp.MinTime) * time.Second)
// 	}
// 	if tp.MaxTime == 0 {
// 		tp.ValidMaxTime = athleteWaceStart.Add(time.Duration(time.Hour) * 24)
// 	} else {
// 		tp.ValidMaxTime = athleteWaceStart.Add(time.Duration(tp.MaxTime) * time.Second)
// 	}
// }

type NewSplitrequest struct {
	RaceID            uuid.UUID `json:"race_id"`
	EventID           uuid.UUID `json:"event_id"`
	Name              string    `json:"name"`
	Type              SplitType `json:"type"`
	DistanceFromStart int       `json:"distance_from_start"`
	ReaderName        string    `json:"reader_name"`
	MinTime           int64     `json:"min_time_sec"`
	MaxTime           int64     `json:"max_time_sec"`
	MinLapTime        int64     `json:"min_lap_time_sec"`
}

func IsValidSplitType(tp SplitType) bool {
	switch tp {
	case SplitTypeStart, SplitTypeFinish, SplitTypeStandard:
		return true
	default:
		return false
	}
}

// func NewSplit(raceID uuid.UUID, eventID uuid.UUID, req SplitFormData) (*Split, error) {
// 	if raceID == uuid.Nil {
// 		return nil, fmt.Errorf("empty raceID")
// 	}
// 	if eventID == uuid.Nil {
// 		return nil, fmt.Errorf("empty eventID")
// 	}
// 	if req.Name == "" {
// 		return nil, fmt.Errorf("empty split name")
// 	}
// 	if req.DistanceFromStart < 0 {
// 		return nil, fmt.Errorf("distance from start must be equal or greater than 0")
// 	}
// 	if req.ReaderName == "" {
// 		return nil, fmt.Errorf("empty box name")
// 	}
// 	if req.MinTime < 0 || req.MaxTime < 0 || req.MinLapTime < 0 {
// 		return nil, fmt.Errorf("min, max and lap times must be equal or greater than 0")
// 	}

// 	id := uuid.New()
// 	return &Split{
// 		ID:                id,
// 		RaceID:            raceID,
// 		EventID:           eventID,
// 		Name:              req.Name,
// 		Type:              req.Type,
// 		DistanceFromStart: req.DistanceFromStart,
// 		ReaderName:           req.ReaderName,
// 		MinTime:        req.MinTime,
// 		MaxTime:        req.MaxTime,
// 		MinLapTime:     req.MinLapTime,
// 	}, nil
// }

func RandomSplit(name string, typ SplitType, dst int, timeReaderID uuid.UUID, min, max, lap time.Duration) *Split {
	return &Split{
		ID:                uuid.New(),
		RaceID:            uuid.New(),
		EventID:           uuid.New(),
		Name:              name,
		Type:              typ,
		DistanceFromStart: dst,
		TimeReaderID:      timeReaderID,
		MinTime:           min,
		MaxTime:           max,
		MinLapTime:        lap,
	}
}
