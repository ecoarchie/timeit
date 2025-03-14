package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ecoarchie/timeit/internal/database"
	"github.com/ecoarchie/timeit/internal/entity"
	"github.com/ecoarchie/timeit/pkg/logger"
	"github.com/google/uuid"
)

type AthleteManager interface {
	GetAthleteByID(ctx context.Context, athleteID uuid.UUID) *entity.Athlete
	CreateAthlete(ctx context.Context, req entity.AthleteCreateRequest) (*entity.Athlete, error)
	UpdateAthlete(ctx context.Context, req entity.AthleteUpdateRequest) (*entity.Athlete, error)
	DeleteAthlete(ctx context.Context, athleteID uuid.UUID) error
	DeleteAthletesForRace(ctx context.Context, raceID, eventID uuid.UUID) error
	FromCSVtoRequestAthlete(raceID uuid.UUID, data []*AthleteCSV) []entity.AthleteCreateRequest
}

type AthleteRepo interface {
	SaveAthlete(ctx context.Context, p *entity.Athlete) error
	GetCategoryFor(ctx context.Context, p *entity.Athlete) (uuid.NullUUID, bool, error)
	GetAthleteWithChip(chip int) (*entity.Athlete, error)
	GetAthleteByID(ctx context.Context, athleteID uuid.UUID) (*entity.Athlete, error)
	DeleteAthlete(ctx context.Context, a *entity.Athlete) error
	DeleteAthletesForRace(ctx context.Context, raceID uuid.UUID) error
	DeleteAthletesForRaceWithEventID(ctx context.Context, raceID, eventID uuid.UUID) error
	GetRecordsAndSplitsForEventAthlete(ctx context.Context, raceID, eventID uuid.UUID) ([]database.GetEventAthleteRecordsRow, []*entity.Split, error)
}

const TimeFormatDDMMYYYY = "02.01.2006"

type AthleteService struct {
	log   *logger.Logger
	repo  AthleteRepo
	cache *RaceCache
}

func NewAthleteService(logger *logger.Logger, repo AthleteRepo, cache *RaceCache) *AthleteService {
	return &AthleteService{
		log:   logger,
		repo:  repo,
		cache: cache,
	}
}

func (ps *AthleteService) GetAthleteByID(ctx context.Context, athleteID uuid.UUID) *entity.Athlete {
	p, err := ps.repo.GetAthleteByID(ctx, athleteID)
	if err != nil {
		return nil
	}
	return p
}

func (ps *AthleteService) CreateAthlete(ctx context.Context, req entity.AthleteCreateRequest) (*entity.Athlete, error) {
	p, err := entity.NewAthlete(req)
	if err != nil {
		return nil, err
	}

	// TODO check if category with this ID is exists. Complete rewrite here
	if !req.CategoryID.Valid {
		err := ps.assignCategory(ctx, p)
		fmt.Println("assign category for athlete: ", req.Bib, p.CategoryID)
		if err != nil {
			fmt.Println("error assigning category", err)
		}
	}

	err = ps.repo.SaveAthlete(ctx, p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (ps *AthleteService) assignCategory(ctx context.Context, p *entity.Athlete) error {
	catID, _, err := ps.repo.GetCategoryFor(ctx, p)
	if err != nil {
		return fmt.Errorf("error assigning category for athlete with bib %d", p.Bib)
	}
	p.CategoryID = catID
	return nil
}

func (ps *AthleteService) UpdateAthlete(ctx context.Context, req entity.AthleteUpdateRequest) (*entity.Athlete, error) {
	p, err := ps.repo.GetAthleteByID(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("updateAthlete: athlete with ID %s not found", req.ID)
	}
	newP, err := entity.NewAthlete(req.AthleteCreateRequest)
	if err != nil {
		return nil, err
	}
	newP.ID = p.ID

	err = ps.repo.SaveAthlete(ctx, newP)
	if err != nil {
		return nil, err
	}
	return newP, nil
}

func (ps *AthleteService) DeleteAthlete(ctx context.Context, athleteID uuid.UUID) error {
	a, err := ps.repo.GetAthleteByID(ctx, athleteID)
	if err != nil {
		return err
	}
	if a == nil {
		return fmt.Errorf("athlete with ID %s not found", athleteID)
	}
	err = ps.repo.DeleteAthlete(ctx, a)
	if err != nil {
		return fmt.Errorf("delete athlete: error deleting athlete %s from DB", athleteID)
	}
	return nil
}

func (ps *AthleteService) DeleteAthleteBulk(ctx context.Context, raceID uuid.UUID, ids []uuid.UUID) []error {
	var errors []error
	for _, id := range ids {
		err := ps.DeleteAthlete(ctx, id)
		if err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

func (as *AthleteService) DeleteAthletesForRace(ctx context.Context, raceID, eventID uuid.UUID) error {
	if eventID == uuid.Nil {
		err := as.repo.DeleteAthletesForRace(ctx, raceID)
		if err != nil {
			return fmt.Errorf("error deleting athletes for raceID = %s", raceID)
		}
	} else {
		err := as.repo.DeleteAthletesForRaceWithEventID(ctx, raceID, eventID)
		if err != nil {
			return fmt.Errorf("error deleting athletes for raceID = %s, eventID = %s", raceID, eventID)
		}
	}
	return nil
}

func (as *AthleteService) FromCSVtoRequestAthlete(raceID uuid.UUID, data []*AthleteCSV) []entity.AthleteCreateRequest {
	eventsMap := as.cache.GetEventNameIDforRace(raceID)
	fmt.Println("events Map: ", eventsMap)
	waves := as.cache.GetWavesForRace(raceID)
	fmt.Println("waves ", waves)
	var res []entity.AthleteCreateRequest
	for _, a := range data {
		eID := eventsMap[a.Event]
		var wID uuid.UUID
		for _, w := range waves {
			if w.EventID == eID {
				if a.Wave == "" {
					wID = w.ID
					break
				} else if w.Name == a.Wave {
					wID = w.ID
				}
			}
		}
		dob, _ := time.Parse(TimeFormatDDMMYYYY, a.DateOfBirth)
		r := entity.AthleteCreateRequest{
			RaceID:      raceID,
			EventID:     eID,
			WaveID:      wID,
			Bib:         a.Bib,
			Chip:        a.Chip,
			FirstName:   a.FirstName,
			LastName:    a.LastName,
			Gender:      entity.CategoryGender(a.Gender),
			DateOfBirth: dob,
			CategoryID: uuid.NullUUID{
				UUID:  uuid.UUID{},
				Valid: false,
			},
			Phone:    a.Phone,
			Comments: a.Comments,
		}
		res = append(res, r)
	}
	return res
}
