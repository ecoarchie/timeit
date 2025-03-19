package service

import (
	"context"
	"testing"
	"time"

	"github.com/ecoarchie/timeit/internal/entity"
	"github.com/ecoarchie/timeit/internal/testdata"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var athleteRepoMock testdata.AthleteRepoMock

func TestCreateAthlete(t *testing.T) {
	service := AthleteService{
		l:    nil,
		repo: athleteRepoMock,
	}
	raceID := uuid.New()
	eventID := uuid.New()
	waveID := uuid.New()
	categoryID := uuid.NullUUID{
		UUID:  uuid.New(),
		Valid: true,
	}
	bib := 100
	chip := 1
	req := entity.AthleteCreateRequest{
		RaceID:      raceID,
		EventID:     eventID,
		WaveID:      waveID,
		Bib:         bib,
		Chip:        chip,
		FirstName:   "Jack",
		LastName:    "Smith",
		Gender:      "male",
		DateOfBirth: time.Time{},
		CategoryID:  categoryID,
		Phone:       "5555555",
		Comments:    "some comments",
	}

	got, err := service.CreateAthlete(context.Background(), req)
	if assert.NoError(t, err) {
		zeroDOB, _ := time.Parse(time.DateOnly, "1900-01-01")
		want := &entity.Athlete{
			ID:          got.ID,
			RaceID:      raceID,
			EventID:     eventID,
			WaveID:      waveID,
			Bib:         100,
			Chip:        1,
			FirstName:   "Jack",
			LastName:    "Smith",
			Gender:      entity.CategoryGenderMale,
			DateOfBirth: zeroDOB,
			CategoryID:  categoryID,
			Phone:       "5555555",
			Comments:    "some comments",
		}
		assert.Equal(t, want, got)
	}
}
