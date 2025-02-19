package service

import (
	"testing"
	"time"

	"github.com/ecoarchie/timeit/internal/entity"
	"github.com/ecoarchie/timeit/internal/testdata"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var participantRepoMock testdata.ParticipantRepoMock

func TestCreateParticipant(t *testing.T) {
	service := ParticipantService{
		l:    nil,
		repo: participantRepoMock,
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
	req := entity.ParticipantCreateRequest{
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

	got, err := service.CreateParticipant(req)
	if assert.NoError(t, err) {
		zeroDOB, _ := time.Parse(time.DateOnly, "1900-01-01")
		want := &entity.Participant{
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
