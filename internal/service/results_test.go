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
	p4 = entity.RandomParticipant("Aaron", "Paul", "male", 103, 103)
)

var recs1 = []entity.BoxRecord{
	{
		Chip:    100,
		TOD:     time.Date(2025, 6, 1, 8, 0, 0, 0, time.UTC),
		BoxName: "box_start",
		CanUse:  true,
	},
	{
		Chip:    100,
		TOD:     time.Date(2025, 6, 1, 8, 5, 0, 0, time.UTC),
		BoxName: "box_cp1",
		CanUse:  true,
	},
	{
		Chip:    100,
		TOD:     time.Date(2025, 6, 1, 8, 10, 0, 0, time.UTC),
		BoxName: "box_finish",
		CanUse:  true,
	},
}

var recs2 = []entity.BoxRecord{
	{
		Chip:    101, // this start rec must be skipped
		TOD:     time.Date(2025, 6, 1, 8, 0, 30, 0, time.UTC),
		BoxName: "box_start",
		CanUse:  true,
	},
	{
		Chip:    101,
		TOD:     time.Date(2025, 6, 1, 8, 1, 0, 0, time.UTC),
		BoxName: "box_start",
		CanUse:  true,
	},
	{
		Chip:    101,
		TOD:     time.Date(2025, 6, 1, 8, 4, 59, 0, time.UTC),
		BoxName: "box_cp1",
		CanUse:  true,
	},
	{
		Chip:    101,
		TOD:     time.Date(2025, 6, 1, 8, 5, 0, 0, time.UTC),
		BoxName: "box_cp1",
		CanUse:  true,
	},
	{
		Chip:    101,
		TOD:     time.Date(2025, 6, 1, 8, 10, 0, 0, time.UTC),
		BoxName: "box_finish",
		CanUse:  true,
	},
	{
		Chip:    101, // this finish rec must be skipped
		TOD:     time.Date(2025, 6, 1, 8, 10, 1, 0, time.UTC),
		BoxName: "box_finish",
		CanUse:  true,
	},
}

var recs3 = []entity.BoxRecord{
	{
		Chip:    102, // this start rec must be skipped
		TOD:     time.Date(2025, 6, 1, 8, 0, 30, 0, time.UTC),
		BoxName: "box_start",
		CanUse:  true,
	},
	{
		Chip:    102, // this start rec must be skipped
		TOD:     time.Date(2025, 6, 1, 8, 0, 35, 0, time.UTC),
		BoxName: "box_start",
		CanUse:  true,
	},
	{
		Chip:    102,
		TOD:     time.Date(2025, 6, 1, 8, 1, 1, 0, time.UTC),
		BoxName: "box_start",
		CanUse:  true,
	},
	{
		Chip:    102,
		TOD:     time.Date(2025, 6, 1, 8, 4, 59, 0, time.UTC),
		BoxName: "box_cp1",
		CanUse:  true,
	},
	{
		Chip:    102,
		TOD:     time.Date(2025, 6, 1, 8, 5, 0, 0, time.UTC),
		BoxName: "box_cp1",
		CanUse:  true,
	},
	{
		Chip:    102,
		TOD:     time.Date(2025, 6, 1, 8, 10, 0, 0, time.UTC),
		BoxName: "box_finish",
		CanUse:  true,
	},
	{
		Chip:    102, // this finish rec must be skipped
		TOD:     time.Date(2025, 6, 1, 8, 10, 0, 1, time.UTC),
		BoxName: "box_finish",
		CanUse:  true,
	},
}

var recs4 = []entity.BoxRecord{ // Missing starting record
	{},
	{
		Chip:    103,
		TOD:     time.Date(2025, 6, 1, 8, 5, 0, 0, time.UTC),
		BoxName: "box_cp1",
		CanUse:  true,
	},
	{
		Chip:    103,
		TOD:     time.Date(2025, 6, 1, 8, 10, 0, 0, time.UTC),
		BoxName: "box_finish",
		CanUse:  true,
	},
}

func TestGetResult(t *testing.T) {
	pr1 := entity.NewParticipantResults(p1)
	pr2 := entity.NewParticipantResults(p2)
	pr3 := entity.NewParticipantResults(p3)
	pr4 := entity.NewParticipantResults(p4)
	t.Run("Valid records return correct results", func(t *testing.T) {
		rs := ResultsService{} // Assuming ResultsService is already defined

		result, err := rs.GetResults(pr1, recs1, w.StartTime, tpsForStandartEvent)

		assert.NoError(t, err)

		assert.NotNil(t, result)
		assert.Equal(t, p1.Chip, result.Chip)
		assert.Equal(t, time.Duration(0), result.ResultsForTPs[tpsForStandartEvent[0].ID].GunTime)
		assert.Equal(t, time.Duration(0), result.ResultsForTPs[tpsForStandartEvent[0].ID].NetTime)
		assert.Equal(t, w.StartTime, result.ResultsForTPs[tpsForStandartEvent[0].ID].TOD)
		assert.Equal(t, (5 * time.Minute), result.ResultsForTPs[tpsForStandartEvent[1].ID].GunTime)
		assert.Equal(t, (5 * time.Minute), result.ResultsForTPs[tpsForStandartEvent[1].ID].NetTime)
		assert.Equal(t, w.StartTime.Add(5*time.Minute), result.ResultsForTPs[tpsForStandartEvent[1].ID].TOD)
		assert.Equal(t, (10 * time.Minute), result.ResultsForTPs[tpsForStandartEvent[2].ID].GunTime)
		assert.Equal(t, (10 * time.Minute), result.ResultsForTPs[tpsForStandartEvent[2].ID].NetTime)
		assert.Equal(t, w.StartTime.Add(10*time.Minute), result.ResultsForTPs[tpsForStandartEvent[2].ID].TOD)
	})

	t.Run("Valid records with 2 recs for intermediate point should skip second one", func(t *testing.T) {
		rs := ResultsService{}
		checkpoint1for2participantNet, _ := time.ParseDuration("3m59s")
		checkpoint1for2participantGun, _ := time.ParseDuration("4m59s")

		result, err := rs.GetResults(pr2, recs2, w.StartTime, tpsForStandartEvent)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, p2.Chip, result.Chip)
		assert.Equal(t, (time.Minute * 1), result.ResultsForTPs[tpsForStandartEvent[0].ID].GunTime)
		assert.Equal(t, (time.Minute * 1), result.ResultsForTPs[tpsForStandartEvent[0].ID].NetTime)
		assert.Equal(t, w.StartTime.Add(time.Minute*1), result.ResultsForTPs[tpsForStandartEvent[0].ID].TOD)
		assert.Equal(t, checkpoint1for2participantGun, result.ResultsForTPs[tpsForStandartEvent[1].ID].GunTime)
		assert.Equal(t, checkpoint1for2participantNet, result.ResultsForTPs[tpsForStandartEvent[1].ID].NetTime)
		assert.Equal(t, time.Date(2025, 6, 1, 8, 4, 59, 0, time.UTC), result.ResultsForTPs[tpsForStandartEvent[1].ID].TOD)
		assert.Equal(t, (10 * time.Minute), result.ResultsForTPs[tpsForStandartEvent[2].ID].GunTime)
		assert.Equal(t, (9 * time.Minute), result.ResultsForTPs[tpsForStandartEvent[2].ID].NetTime)
		assert.Equal(t, w.StartTime.Add(10*time.Minute), result.ResultsForTPs[tpsForStandartEvent[2].ID].TOD)
	})

	t.Run("Valid records. 2 Start recs should be skipped", func(t *testing.T) {
		rs := ResultsService{}
		checkpoint1for3participantGun, _ := time.ParseDuration("4m59s")
		checkpoint1for3participantNet, _ := time.ParseDuration("3m58s")

		result, err := rs.GetResults(pr3, recs3, w.StartTime, tpsForStandartEvent)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, p3.Chip, result.Chip)
		assert.Equal(t, (time.Minute*1 + time.Second*1), result.ResultsForTPs[tpsForStandartEvent[0].ID].GunTime)
		assert.Equal(t, (time.Minute*1 + time.Second*1), result.ResultsForTPs[tpsForStandartEvent[0].ID].NetTime)
		assert.Equal(t, w.StartTime.Add(time.Minute*1+time.Second*1), result.ResultsForTPs[tpsForStandartEvent[0].ID].TOD)
		assert.Equal(t, checkpoint1for3participantGun, result.ResultsForTPs[tpsForStandartEvent[1].ID].GunTime)
		assert.Equal(t, checkpoint1for3participantNet, result.ResultsForTPs[tpsForStandartEvent[1].ID].NetTime)
		assert.Equal(t, time.Date(2025, 6, 1, 8, 4, 59, 0, time.UTC), result.ResultsForTPs[tpsForStandartEvent[1].ID].TOD)
		assert.Equal(t, (10 * time.Minute), result.ResultsForTPs[tpsForStandartEvent[2].ID].GunTime)
		assert.Equal(t, (8*time.Minute + time.Second*59), result.ResultsForTPs[tpsForStandartEvent[2].ID].NetTime)
		assert.Equal(t, w.StartTime.Add(10*time.Minute), result.ResultsForTPs[tpsForStandartEvent[2].ID].TOD)
	})

	t.Run("Missing starting record", func(t *testing.T) {
		rs := ResultsService{}
		checkpoint1for4participantGun, _ := time.ParseDuration("5m00s")
		checkpoint1for4participantNet, _ := time.ParseDuration("5m00s")

		result, err := rs.GetResults(pr4, recs4, w.StartTime, tpsForStandartEvent)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, p4.Chip, result.Chip)
		assert.Nil(t, result.ResultsForTPs[tpsForStandartEvent[0].ID])

		checkpointGun := result.ResultsForTPs[tpsForStandartEvent[1].ID].GunTime
		checkpointNet := result.ResultsForTPs[tpsForStandartEvent[1].ID].NetTime
		checkpointTOD := result.ResultsForTPs[tpsForStandartEvent[1].ID].TOD

		assert.Equal(t, checkpoint1for4participantGun, checkpointGun)
		assert.Equal(t, checkpoint1for4participantNet, checkpointNet)
		assert.Equal(t, time.Date(2025, 6, 1, 8, 5, 0, 0, time.UTC), checkpointTOD)

		finishGun := result.ResultsForTPs[tpsForStandartEvent[2].ID].GunTime
		finishNet := result.ResultsForTPs[tpsForStandartEvent[2].ID].NetTime
		finishTOD := result.ResultsForTPs[tpsForStandartEvent[2].ID].TOD
		assert.Equal(t, (10 * time.Minute), finishGun)
		assert.Equal(t, (10 * time.Minute), finishNet)
		assert.Equal(t, finishGun, finishNet, "Guntime and Net time must be equal")
		assert.Equal(t, w.StartTime.Add(10*time.Minute), finishTOD)
	})
}
