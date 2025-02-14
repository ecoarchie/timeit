package entity

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSetMinMaxTime(t *testing.T) {
	tp := &TimingPoint{
		ID:                uuid.New(),
		RaceID:            uuid.New(),
		EventID:           uuid.New(),
		Name:              "Test TP",
		Type:              "standard",
		DistanceFromStart: 0,
		BoxName:           "START",
		MinTimeSec:        0,
		MaxTimeSec:        100,
		MinLapTimeSec:     0,
	}

	startTime := time.Date(2025, time.January, 1, 8, 0, 0, 0, time.UTC)
	tp.SetValidMinMaxTimes(startTime)
	assert.Equal(t, startTime, tp.ValidMinTime)
	assert.Equal(t, startTime.Add(time.Duration(tp.MaxTimeSec)*time.Second), tp.ValidMaxTime)

	tp = &TimingPoint{
		ID:                uuid.New(),
		RaceID:            uuid.New(),
		EventID:           uuid.New(),
		Name:              "Test TP",
		Type:              "standard",
		DistanceFromStart: 0,
		BoxName:           "START",
		MinTimeSec:        10,
		MaxTimeSec:        0,
		MinLapTimeSec:     0,
	}
	tp.SetValidMinMaxTimes(startTime)
	assert.Equal(t, startTime.Add(time.Duration(tp.MinTimeSec)*time.Second), tp.ValidMinTime)
	assert.Equal(t, startTime.Add(time.Duration(time.Hour)*24), tp.ValidMaxTime)
}
