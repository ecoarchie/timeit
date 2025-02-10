package service

import (
	"context"
	"fmt"

	"github.com/ecoarchie/timeit/internal/entity"
	"github.com/ecoarchie/timeit/pkg/logger"
	"github.com/google/uuid"
)

// TODO separate errors to ErrorToSave (preventing from saving) and Warning (can save, just pay attention)
// TODO add category boundaries validation

type RaceConfigurator interface {
	Save(ctx context.Context, rc entity.RaceConfig) []error
}

type RaceRepo interface {
	SaveRaceConfig(ctx context.Context, r entity.RaceConfig) error
}

type RaceService struct {
	l    logger.Interface
	repo RaceRepo
}

func NewRaceService(logger logger.Interface, repo RaceRepo) *RaceService {
	return &RaceService{
		l:    logger,
		repo: repo,
	}
}

func (rs RaceService) Save(ctx context.Context, rc entity.RaceConfig) []error {
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
	return errors
}

func (rs RaceService) validate(rc entity.RaceConfig) []error {
	errors := []error{}
	if err := validateRace(rc.Race); err != nil {
		errors = append(errors, err)
	}

	if err := validatePhysicalLocations(rc.PhysicalLocations); err != nil {
		errors = append(errors, err)
		return errors
	}

	if len(rc.Events) == 0 {
		errors = append(errors, fmt.Errorf("race must have at least one event"))
		return errors
	}
	for _, ec := range rc.Events {
		if err := validateEventConfig(rc.Race, rc.PhysicalLocations, ec); err != nil {
			errors = append(errors, err...)
		}
	}
	return errors
}

func validateRace(race entity.Race) error {
	if err := entity.IsValidTimezone(race.Timezone); err != nil {
		return err
	}
	if race.Name == "" {
		return fmt.Errorf("empty race name")
	}
	return nil
}

func validatePhysicalLocations(locs []entity.PhysicalLocation) error {
	if len(locs) == 0 {
		return fmt.Errorf("race must have at least one physical location")
	}
	boxNames := make(map[string]struct{})

	// Loop through locations and check for uniqueness
	for _, l := range locs {
		if _, exists := boxNames[l.BoxName]; exists {
			return fmt.Errorf("duplicate box name found: %s", l.BoxName)
		}
		boxNames[l.BoxName] = struct{}{}
	}
	return nil
}

func validateEventConfig(race entity.Race, locs []entity.PhysicalLocation, ec entity.EventConfig) []error {
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

	if len(ec.TimingPoints) == 0 {
		errors = append(errors, fmt.Errorf("event must have at least one timing point"))
	}

	if len(errors) != 0 {
		return errors
	}

	for _, tp := range ec.TimingPoints {
		if err := validateTimingPoint(ec.RaceID, ec.ID, locs, tp); err != nil {
			errors = append(errors, err)
		}
	}

	for _, w := range ec.Waves {
		if err := validateWave(race.ID, ec.ID, w); err != nil {
			errors = append(errors, err)
		}
	}

	for _, c := range ec.Categories {
		if err := validateCategory(race.ID, ec.ID, c); err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

func validateCategory(raceID, eventID uuid.UUID, c entity.Category) error {
	if raceID == uuid.Nil || raceID != c.RaceID {
		return fmt.Errorf("empty or invalid raceID for category")
	}
	if eventID == uuid.Nil || eventID != c.EventID {
		return fmt.Errorf("empty or invalid eventID for category")
	}
	if c.Name == "" {
		return fmt.Errorf("empty timing point name")
	}
	if c.FromAge < 0 {
		return fmt.Errorf("from age must be greater or equal to 0")
	}
	if c.ToAge < 0 {
		return fmt.Errorf("to age must be greater or equal to 0")
	}
	if c.FromAge > c.ToAge {
		return fmt.Errorf("upper age limit must be greater than lower age limit")
	}
	return nil
}

func validateWave(raceID, eventID uuid.UUID, w entity.Wave) error {
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

func validateTimingPoint(raceID, eventID uuid.UUID, locs []entity.PhysicalLocation, tp entity.TimingPoint) error {
	if raceID == uuid.Nil {
		return fmt.Errorf("empty raceID")
	}
	if eventID != tp.EventID {
		return fmt.Errorf("wrong event id for timint point")
	}
	if tp.Name == "" {
		return fmt.Errorf("empty timing point name")
	}
	if tp.Type == "" || !entity.IsValidTPType(tp.Type) {
		return fmt.Errorf("empty or invalid timing point type")
	}
	if tp.DistanceFromStart < 0 {
		return fmt.Errorf("distance from start must be equal or greater than 0")
	}

	// check box name for timing point
	if tp.BoxName == "" {
		return fmt.Errorf("empty box name")
	}
	unknownBoxName := true
	for _, l := range locs {
		if l.BoxName == tp.BoxName {
			unknownBoxName = false
		}
	}
	if unknownBoxName {
		return fmt.Errorf("unknown box name for timing point")
	}

	// check time restrictions
	if tp.MinTimeSec < 0 || tp.MaxTimeSec < 0 || tp.MinLapTimeSec < 0 {
		return fmt.Errorf("min, max and lap times must be equal or greater than 0")
	}
	return nil
}
