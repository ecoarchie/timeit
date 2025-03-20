package service

import (
	"context"
	"fmt"
	"slices"
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
	FromCSVtoRequestAthlete(ctx context.Context, raceID uuid.UUID, data []*AthleteCSV) ([]entity.AthleteCreateRequest, error)
}

type AthleteRepo interface {
	SaveAthlete(ctx context.Context, p *entity.Athlete) error
	GetCategoryFor(ctx context.Context, p *entity.Athlete) (uuid.NullUUID, bool, error)
	GetAthleteWithChip(chip int) (*entity.Athlete, error)
	GetAthleteByID(ctx context.Context, athleteID uuid.UUID) (*entity.Athlete, error)
	DeleteAthlete(ctx context.Context, a *entity.Athlete) error
	DeleteAthletesForRace(ctx context.Context, raceID uuid.UUID) error
	DeleteAthletesForRaceWithEventID(ctx context.Context, raceID, eventID uuid.UUID) error
	GetRecordsAndSplitsForEventAthlete(ctx context.Context, raceID, eventID uuid.UUID) ([]database.GetEventAthleteRecordsCRow, []*entity.Split, error)
	SaveAthleteSplits(ctx context.Context, as []database.CreateAthleteSplitsParams) error
	GetEventIDsWithWavesStarted(ctx context.Context, raceID uuid.UUID) ([]uuid.UUID, error)
}

const TimeFormatDDMMYYYY = "02.01.2006"

type AthleteService struct {
	log         *logger.Logger
	athleteRepo AthleteRepo
	raceRepo    RaceConfigurator
}

func NewAthleteService(logger *logger.Logger, athleteRepo AthleteRepo, raceRepo RaceConfigurator) *AthleteService {
	return &AthleteService{
		log:         logger,
		athleteRepo: athleteRepo,
		raceRepo:    raceRepo,
	}
}

func (ps *AthleteService) GetAthleteByID(ctx context.Context, athleteID uuid.UUID) *entity.Athlete {
	p, err := ps.athleteRepo.GetAthleteByID(ctx, athleteID)
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
		if err != nil {
			fmt.Println("error assigning category", err)
		}
	}

	err = ps.athleteRepo.SaveAthlete(ctx, p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (ps *AthleteService) assignCategory(ctx context.Context, p *entity.Athlete) error {
	catID, _, err := ps.athleteRepo.GetCategoryFor(ctx, p)
	if err != nil {
		return fmt.Errorf("error assigning category for athlete with bib %d", p.Bib)
	}
	p.CategoryID = catID
	return nil
}

func (ps *AthleteService) UpdateAthlete(ctx context.Context, req entity.AthleteUpdateRequest) (*entity.Athlete, error) {
	p, err := ps.athleteRepo.GetAthleteByID(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("updateAthlete: athlete with ID %s not found", req.ID)
	}
	newP, err := entity.NewAthlete(req.AthleteCreateRequest)
	if err != nil {
		return nil, err
	}
	newP.ID = p.ID

	err = ps.athleteRepo.SaveAthlete(ctx, newP)
	if err != nil {
		return nil, err
	}
	return newP, nil
}

func (ps *AthleteService) DeleteAthlete(ctx context.Context, athleteID uuid.UUID) error {
	a, err := ps.athleteRepo.GetAthleteByID(ctx, athleteID)
	if err != nil {
		return err
	}
	if a == nil {
		return fmt.Errorf("athlete with ID %s not found", athleteID)
	}
	err = ps.athleteRepo.DeleteAthlete(ctx, a)
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
		err := as.athleteRepo.DeleteAthletesForRace(ctx, raceID)
		if err != nil {
			return fmt.Errorf("error deleting athletes for raceID = %s", raceID)
		}
	} else {
		err := as.athleteRepo.DeleteAthletesForRaceWithEventID(ctx, raceID, eventID)
		if err != nil {
			return fmt.Errorf("error deleting athletes for raceID = %s, eventID = %s", raceID, eventID)
		}
	}
	return nil
}

func (as *AthleteService) FromCSVtoRequestAthlete(ctx context.Context, raceID uuid.UUID, data []*AthleteCSV) ([]entity.AthleteCreateRequest, error) {
	raceModel, err := as.raceRepo.GetRaceConfig(ctx, raceID)
	if err != nil {
		return nil, err
	}
	var res []entity.AthleteCreateRequest
	for _, a := range data {
		eventIndex := slices.IndexFunc(raceModel.Events, func(e *entity.Event) bool {
			return e.Name == a.Event
		})
		if eventIndex == -1 {
			return nil, fmt.Errorf("event with name %s does not exists. Import aborted", a.Event)
		}
		eventID := raceModel.Events[eventIndex].ID
		var waveID uuid.UUID
		if a.Wave == "" {
			// if wave is not provided in CSV, asign to athlete the first wave of event by default
			waveID = raceModel.Events[eventIndex].Waves[0].ID
		} else {
			waveIndex := slices.IndexFunc(raceModel.Events[eventIndex].Waves, func(w *entity.Wave) bool {
				return w.Name == a.Wave
			})
			if waveIndex == -1 {
				return nil, fmt.Errorf("wave with name %s does not exists. Import aborted", a.Wave)
			}
			waveID = raceModel.Events[eventIndex].Waves[waveIndex].ID
		}
		dob, err := time.Parse(TimeFormatDDMMYYYY, a.DateOfBirth)
		if err != nil {
			dob = time.Date(1900, time.January, 1, 0, 0, 0, 0, time.UTC)
		}
		r := entity.AthleteCreateRequest{
			RaceID:      raceID,
			EventID:     eventID,
			WaveID:      waveID,
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
	return res, nil
}
