package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ecoarchie/timeit/internal/entity"
	"github.com/ecoarchie/timeit/pkg/logger"
	"github.com/google/uuid"
)

// TODO separate errors to ErrorToSave (preventing from saving) and Warning (can save, just pay attention)
// TODO add category boundaries validation

type RaceConfigurator interface {
	Save(ctx context.Context, rc *entity.RaceConfig) []error
	CreateRace(ctx context.Context, req *entity.RaceFormData) (*entity.Race, error)
	GetRaceConfig(ctx context.Context, raceID string) (*entity.RaceConfig, error)
}

type RaceRepo interface {
	SaveRaceConfig(ctx context.Context, r *entity.RaceConfig) error
	GetRaceConfig(ctx context.Context, raceID uuid.UUID) (*entity.RaceConfig, error)
}

type RaceService struct {
	l         logger.Interface
	raceCache *RaceCache
	repo      RaceRepo
}

func NewRaceService(logger logger.Interface, rc *RaceCache, repo RaceRepo) *RaceService {
	return &RaceService{
		l:         logger,
		raceCache: rc,
		repo:      repo,
	}
}

func (rc RaceService) CreateRace(ctx context.Context, req *entity.RaceFormData) (*entity.Race, error) {
	r, err := entity.NewRace(req)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (rc RaceService) GetRaceConfig(ctx context.Context, raceID string) (*entity.RaceConfig, error) {
	uuID, err := uuid.Parse(raceID)
	if err != nil {
		return nil, err
	}
	rconfig, err := rc.repo.GetRaceConfig(ctx, uuID)
	if err != nil {
		return nil, err
	}
	if rconfig == nil {
		return nil, fmt.Errorf("race with id %s not found", raceID)
	}

	return rconfig, nil
}

func (rs RaceService) Save(ctx context.Context, rc *entity.RaceConfig) []error {
	errors := rs.validate(rc)
	if len(errors) != 0 {
		return errors
	}
	err := rs.repo.SaveRaceConfig(ctx, rc)
	if err != nil {
		const msg = "error saving race to repo"
		rs.l.Error(msg, err)
		errors = append(errors, err)
	}
	rs.raceCache.StoreRaceConfig(rc)
	rs.l.Info("race cache updated")
	// fmt.Println("Race config cats are: ", rc.Events[0].Categories[0], rc.Events[0].Categories[1])
	return errors
}

func (rs RaceService) validate(rc *entity.RaceConfig) []error {
	fmt.Println(rc)
	errors := []error{}
	if err := validateRace(rc.Race); err != nil {
		errors = append(errors, err)
	}

	if err := validateTimeReaders(rc.TimeReaders); err != nil {
		errors = append(errors, err)
		return errors
	}

	if len(rc.Events) == 0 {
		errors = append(errors, fmt.Errorf("race must have at least one event"))
		return errors
	}
	for _, ec := range rc.Events {
		if err := validateEventConfig(rc.Race, rc.TimeReaders, ec); err != nil {
			errors = append(errors, err...)
		}
	}
	return errors
}

func validateRace(race *entity.Race) error {
	if err := entity.IsValidTimezone(race.Timezone); err != nil {
		return err
	}
	if race.Name == "" {
		return fmt.Errorf("empty race name")
	}
	return nil
}

func validateTimeReaders(locs []*entity.TimeReader) error {
	if len(locs) == 0 {
		return fmt.Errorf("race must have at least one physical time_reader")
	}
	boxNames := make(map[string]struct{})

	// Loop through time_readers and check for uniqueness
	for _, l := range locs {
		if _, exists := boxNames[l.ReaderName]; exists {
			return fmt.Errorf("duplicate box name found: %s", l.ReaderName)
		}
		boxNames[l.ReaderName] = struct{}{}
	}
	return nil
}

func validateEventConfig(race *entity.Race, locs []*entity.TimeReader, ec *entity.EventConfig) []error {
	errors := []error{}
	if race.ID != ec.RaceID {
		errors = append(errors, fmt.Errorf("wrong race id for event %s", ec.Name))
	}
	if ec.ID == uuid.Nil {
		errors = append(errors, fmt.Errorf("empty event id for event %s", ec.Name))
	}
	if ec.Name == "" {
		errors = append(errors, fmt.Errorf("empty event name"))
	}
	if ec.DistanceInMeters <= 0 {
		errors = append(errors, fmt.Errorf("distance for event %s must be greater than 0", ec.Name))
	}

	if len(ec.Splits) == 0 {
		errors = append(errors, fmt.Errorf("event must have at least one split"))
	}

	if len(errors) != 0 {
		return errors
	}

	for _, tp := range ec.Splits {
		if err := validateSplit(ec.RaceID, ec.ID, locs, tp); err != nil {
			errors = append(errors, err)
		}
	}

	for _, w := range ec.Waves {
		if err := validateWave(race.ID, ec.ID, w); err != nil {
			errors = append(errors, err)
		}
	}

	// pass addr of Category to update its BirthDate fields
	for _, c := range ec.Categories {
		if err := validateCategory(race.ID, ec.ID, ec.EventDate, c); err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

func validateCategory(raceID, eventID uuid.UUID, eventDate time.Time, c *entity.Category) error {
	if raceID == uuid.Nil || raceID != c.RaceID {
		return fmt.Errorf("empty or invalid raceID for category")
	}
	if eventID == uuid.Nil || eventID != c.EventID {
		return fmt.Errorf("empty or invalid eventID for category")
	}
	if c.Name == "" {
		return fmt.Errorf("empty split name")
	}
	if !entity.IsValidGender(c.Gender) {
		return fmt.Errorf("invalid gender")
	}
	if c.AgeFrom < 0 {
		return fmt.Errorf("from age must be greater or equal to 0")
	}
	if c.AgeTo < 0 {
		return fmt.Errorf("to age must be greater or equal to 0")
	}
	if c.AgeFrom > c.AgeTo {
		return fmt.Errorf("upper age limit must be greater than lower age limit")
	}
	return nil
}

func validateWave(raceID, eventID uuid.UUID, w *entity.Wave) error {
	if raceID == uuid.Nil || raceID != w.RaceID {
		return fmt.Errorf("empty or invalid raceID for wave")
	}
	if eventID == uuid.Nil || eventID != w.EventID {
		return fmt.Errorf("empty or invalid eventID for wave")
	}
	if w.Name == "" {
		return fmt.Errorf("empty wave name")
	}
	return nil
}

func validateSplit(raceID, eventID uuid.UUID, locs []*entity.TimeReader, tp *entity.Split) error {
	if raceID == uuid.Nil {
		return fmt.Errorf("empty raceID")
	}
	if eventID != tp.EventID {
		return fmt.Errorf("wrong event id for timint point")
	}
	if tp.Name == "" {
		return fmt.Errorf("empty split name")
	}
	if tp.Type == "" || !entity.IsValidSplitType(tp.Type) {
		return fmt.Errorf("empty or invalid split type")
	}
	if tp.DistanceFromStart < 0 {
		return fmt.Errorf("distance from start must be equal or greater than 0")
	}

	// check box name for split
	if tp.TimeReaderID.String() == "" {
		return fmt.Errorf("empty time_reader ID")
	}
	unknownBoxID := true
	for _, l := range locs {
		if l.ID == tp.TimeReaderID {
			unknownBoxID = false
		}
	}
	if unknownBoxID {
		return fmt.Errorf("unknown box ID for split")
	}

	// check time restrictions
	if tp.MinTime < 0 || tp.MaxTime < 0 || tp.MinLapTime < 0 {
		return fmt.Errorf("min, max and lap times must be equal or greater than 0")
	}
	return nil
}
