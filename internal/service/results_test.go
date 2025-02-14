package service

import (
	"testing"
	"time"

	"github.com/ecoarchie/timeit/internal/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var w = entity.Wave{
	ID:        uuid.New(),
	Name:      "Elite Runners",
	StartTime: time.Date(2025, 6, 1, 8, 0, 0, 0, time.UTC),
}

var tpsForStandartEvent = []*entity.TimingPoint{
	entity.RandomTimingPoint("Start Line", entity.TPTypeStart, 0, "box_start", 0, 0, 0),
	entity.RandomTimingPoint("Checkpoint 1", entity.TPTypeStandard, 1000, "box_cp1", 180, 0, 0),
	entity.RandomTimingPoint("Finish Line", entity.TPTypeFinish, 2000, "box_finish", 300, 0, 0),
}

var (
	p1 = entity.RandomParticipant("John", "Doe", "male", 100, 100)
	p2 = entity.RandomParticipant("Jane", "Doe", "female", 101, 101)
	p3 = entity.RandomParticipant("Mike", "Smith", "male", 102, 102)
)

var recs1 = []entity.BoxRecord{
	{
		ID:      1,
		Chip:    100,
		TOD:     time.Date(2025, 6, 1, 8, 0, 0, 0, time.UTC),
		BoxName: "box_start",
		CanUse:  true,
	},
	{
		ID:      2,
		Chip:    100,
		TOD:     time.Date(2025, 6, 1, 8, 5, 0, 0, time.UTC),
		BoxName: "box_cp1",
		CanUse:  true,
	},
	{
		ID:      3,
		Chip:    100,
		TOD:     time.Date(2025, 6, 1, 8, 10, 0, 0, time.UTC),
		BoxName: "box_finish",
		CanUse:  true,
	},
}

var recs2 = []entity.BoxRecord{
	{
		ID:      0, // this start rec must be skipped
		Chip:    101,
		TOD:     time.Date(2025, 6, 1, 8, 0, 30, 0, time.UTC),
		BoxName: "box_start",
		CanUse:  true,
	},
	{
		ID:      1,
		Chip:    101,
		TOD:     time.Date(2025, 6, 1, 8, 1, 0, 0, time.UTC),
		BoxName: "box_start",
		CanUse:  true,
	},
	{
		ID:      2,
		Chip:    101,
		TOD:     time.Date(2025, 6, 1, 8, 4, 59, 0, time.UTC),
		BoxName: "box_cp1",
		CanUse:  true,
	},
	{
		ID:      3,
		Chip:    101,
		TOD:     time.Date(2025, 6, 1, 8, 5, 0, 0, time.UTC),
		BoxName: "box_cp1",
		CanUse:  true,
	},
	{
		ID:      4,
		Chip:    101,
		TOD:     time.Date(2025, 6, 1, 8, 10, 0, 0, time.UTC),
		BoxName: "box_finish",
		CanUse:  true,
	},
	{
		ID:      5, // this finish rec must be skipped
		Chip:    101,
		TOD:     time.Date(2025, 6, 1, 8, 10, 1, 0, time.UTC),
		BoxName: "box_finish",
		CanUse:  true,
	},
}

var recs3 = []entity.BoxRecord{
	{
		ID:      0, // this start rec must be skipped
		Chip:    102,
		TOD:     time.Date(2025, 6, 1, 8, 0, 30, 0, time.UTC),
		BoxName: "box_start",
		CanUse:  true,
	},
	{
		ID:      0, // this start rec must be skipped
		Chip:    102,
		TOD:     time.Date(2025, 6, 1, 8, 0, 35, 0, time.UTC),
		BoxName: "box_start",
		CanUse:  true,
	},
	{
		ID:      1,
		Chip:    102,
		TOD:     time.Date(2025, 6, 1, 8, 1, 1, 0, time.UTC),
		BoxName: "box_start",
		CanUse:  true,
	},
	{
		ID:      2,
		Chip:    102,
		TOD:     time.Date(2025, 6, 1, 8, 4, 59, 0, time.UTC),
		BoxName: "box_cp1",
		CanUse:  true,
	},
	{
		ID:      3,
		Chip:    102,
		TOD:     time.Date(2025, 6, 1, 8, 5, 0, 0, time.UTC),
		BoxName: "box_cp1",
		CanUse:  true,
	},
	{
		ID:      4,
		Chip:    102,
		TOD:     time.Date(2025, 6, 1, 8, 10, 0, 0, time.UTC),
		BoxName: "box_finish",
		CanUse:  true,
	},
	{
		ID:      5, // this finish rec must be skipped
		Chip:    102,
		TOD:     time.Date(2025, 6, 1, 8, 10, 0, 1, time.UTC),
		BoxName: "box_finish",
		CanUse:  true,
	},
}

func TestGetResult(t *testing.T) {
	t.Run("Valid records return correct results", func(t *testing.T) {
		rs := ResultsService{} // Assuming ResultsService is already defined

		result, err := rs.GetResult(p1, recs1, w, tpsForStandartEvent)

		assert.NoError(t, err)

		assert.NotNil(t, result)
		assert.Equal(t, 100, result.Chip)
		assert.Equal(t, int64(0), result.ResultsForTPs["Start Line"].GunTime)
		assert.Equal(t, int64(0), result.ResultsForTPs["Start Line"].NetTime)
		assert.Equal(t, w.StartTime, result.ResultsForTPs["Start Line"].TOD)
		assert.Equal(t, (5 * time.Minute).Microseconds(), result.ResultsForTPs["Checkpoint 1"].GunTime)
		assert.Equal(t, (5 * time.Minute).Microseconds(), result.ResultsForTPs["Checkpoint 1"].NetTime)
		assert.Equal(t, w.StartTime.Add(5*time.Minute), result.ResultsForTPs["Checkpoint 1"].TOD)
		assert.Equal(t, (10 * time.Minute).Microseconds(), result.ResultsForTPs["Finish Line"].GunTime)
		assert.Equal(t, (10 * time.Minute).Microseconds(), result.ResultsForTPs["Finish Line"].NetTime)
		assert.Equal(t, w.StartTime.Add(10*time.Minute), result.ResultsForTPs["Finish Line"].TOD)
	})

	t.Run("Valid records with 2 recs for intermediate point should skip second one", func(t *testing.T) {
		rs := ResultsService{}
		checkpoint1for2participantNet, _ := time.ParseDuration("3m59s")
		checkpoint1for2participantGun, _ := time.ParseDuration("4m59s")

		result, err := rs.GetResult(p2, recs2, w, tpsForStandartEvent)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 101, result.Chip)
		assert.Equal(t, (time.Minute * 1).Microseconds(), result.ResultsForTPs["Start Line"].GunTime)
		assert.Equal(t, (time.Minute * 1).Microseconds(), result.ResultsForTPs["Start Line"].NetTime)
		assert.Equal(t, w.StartTime.Add(time.Minute*1), result.ResultsForTPs["Start Line"].TOD)
		assert.Equal(t, checkpoint1for2participantGun.Microseconds(), result.ResultsForTPs["Checkpoint 1"].GunTime)
		assert.Equal(t, checkpoint1for2participantNet.Microseconds(), result.ResultsForTPs["Checkpoint 1"].NetTime)
		assert.Equal(t, time.Date(2025, 6, 1, 8, 4, 59, 0, time.UTC), result.ResultsForTPs["Checkpoint 1"].TOD)
		assert.Equal(t, (10 * time.Minute).Microseconds(), result.ResultsForTPs["Finish Line"].GunTime)
		assert.Equal(t, (9 * time.Minute).Microseconds(), result.ResultsForTPs["Finish Line"].NetTime)
		assert.Equal(t, w.StartTime.Add(10*time.Minute), result.ResultsForTPs["Finish Line"].TOD)
	})

	t.Run("Valid records. 2 Start recs should be skipped", func(t *testing.T) {
		rs := ResultsService{}
		checkpoint1for3participantGun, _ := time.ParseDuration("4m59s")
		checkpoint1for3participantNet, _ := time.ParseDuration("3m58s")

		result, err := rs.GetResult(p3, recs3, w, tpsForStandartEvent)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 102, result.Chip)
		assert.Equal(t, (time.Minute*1 + time.Second*1).Microseconds(), result.ResultsForTPs["Start Line"].GunTime)
		assert.Equal(t, (time.Minute*1 + time.Second*1).Microseconds(), result.ResultsForTPs["Start Line"].NetTime)
		assert.Equal(t, w.StartTime.Add(time.Minute*1+time.Second*1), result.ResultsForTPs["Start Line"].TOD)
		assert.Equal(t, checkpoint1for3participantGun.Microseconds(), result.ResultsForTPs["Checkpoint 1"].GunTime)
		assert.Equal(t, checkpoint1for3participantNet.Microseconds(), result.ResultsForTPs["Checkpoint 1"].NetTime)
		assert.Equal(t, time.Date(2025, 6, 1, 8, 4, 59, 0, time.UTC), result.ResultsForTPs["Checkpoint 1"].TOD)
		assert.Equal(t, (10 * time.Minute).Microseconds(), result.ResultsForTPs["Finish Line"].GunTime)
		assert.Equal(t, (8*time.Minute + time.Second*59).Microseconds(), result.ResultsForTPs["Finish Line"].NetTime)
		assert.Equal(t, w.StartTime.Add(10*time.Minute), result.ResultsForTPs["Finish Line"].TOD)
	})
}
