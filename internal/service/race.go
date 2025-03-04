package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ecoarchie/timeit/internal/entity"
	"github.com/ecoarchie/timeit/pkg/logger"
	"github.com/ecoarchie/timeit/pkg/validator"
	"github.com/google/uuid"
)

type ValidationErrors map[string]string

// TODO separate errors to ErrorToSave (preventing from saving) and Warning (can save, just pay attention)
// TODO add category boundaries validation

type RaceConfigurator interface {
	Validate(ctx context.Context, v *validator.Validator, rc *entity.RaceConfig)
	Save(ctx context.Context, rc *entity.RaceConfig) error
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
		return nil, fmt.Errorf("error parsing race UUID")
	}
	rconfig, err := rc.repo.GetRaceConfig(ctx, uuID)
	if err != nil {
		return nil, err
	}
	if rconfig == nil {
		return nil, nil
	}

	// FIXME update RaceCache

	return rconfig, nil
}

func (rs RaceService) Save(ctx context.Context, rc *entity.RaceConfig) error {
	err := rs.repo.SaveRaceConfig(ctx, rc)
	if err != nil {
		const msg = "error saving race to repo"
		rs.l.Error(msg, err)
		return err
	}
	rs.raceCache.StoreRaceConfig(rc)
	rs.l.Info("race cache updated")
	return nil
}

func (rs RaceService) Validate(ctx context.Context, v *validator.Validator, rc *entity.RaceConfig) {
	validateRace(v, rc.Race)

	v.Check(len(rc.TimeReaders) > 0, "time readers", "race must have at least one time reader")
	if len(rc.TimeReaders) > 0 {
		validateTimeReaders(v, rc.TimeReaders, rc.Race.ID)
	}

	v.Check(len(rc.Events) != 0, "events", "must be at least one")
	if len(rc.Events) > 0 {
		var eventNames []string
		for _, e := range rc.Events {
			eventNames = append(eventNames, e.Name)
		}
		v.Check(validator.Unique(eventNames), "event names", "must be unique")
		for _, ec := range rc.Events {
			validateEventConfig(v, rc.Race, rc.TimeReaders, ec)
		}
	}
}

func validateRace(v *validator.Validator, race *entity.Race) {
	v.Check(entity.IsIANATimezone(race.Timezone), "timezone", "must be valid IANA timezone")
	v.Check(race.Name != "", "race name", "must not be empty")
}

func validateTimeReaders(v *validator.Validator, readers []*entity.TimeReader, raceID uuid.UUID) {
	var timeReadersNames []string
	for _, r := range readers {
		timeReadersNames = append(timeReadersNames, r.ReaderName)
		v.Check(r.RaceID == raceID, "timers raceID", "must correspond to ID of configurated race")
	}
	v.Check(validator.Unique(timeReadersNames), "time readers names", "must be unique")
}

func validateEventConfig(v *validator.Validator, race *entity.Race, readers []*entity.TimeReader, ec *entity.EventConfig) {
	v.Check(race.ID == ec.RaceID, "race_id for event", "must correspond to ID of configurated race")
	v.Check(ec.ID != uuid.Nil, "event_id", "must not be empty")
	v.Check(ec.Name != "", "event_name", "must not be empty")
	v.Check(ec.DistanceInMeters > 0, "event distance_in_meters", "must be greater than 0")

	v.Check(len(ec.Splits) != 0, "splits", "event must have at least one split")
	if len(ec.Splits) > 0 {
		var splitsNames []string
		splitTypeQty := make(map[entity.SplitType]int)
		for _, split := range ec.Splits {
			splitsNames = append(splitsNames, split.Name)
			splitTypeQty[split.Type]++
			validateSplit(v, ec.RaceID, ec.ID, readers, split)
		}
		v.Check(validator.Unique(splitsNames), "splits", "must have unique names for event")
		v.Check(splitTypeQty[entity.SplitTypeStart] < 2, "split with type start", "must be 0 or 1")
		v.Check(splitTypeQty[entity.SplitTypeFinish] != 0, "split with type finish", "must be configured")
		v.Check(splitTypeQty[entity.SplitTypeFinish] < 2, "split with type finish", "must not be more than 1")
	}

	v.Check(len(ec.Waves) > 0, "waves", "must be at least one for event")
	if len(ec.Waves) > 0 {
		var wavesNames []string
		for _, w := range ec.Waves {
			wavesNames = append(wavesNames, w.Name)
			validateWave(v, race.ID, ec.ID, w)
		}
		v.Check(validator.Unique(wavesNames), "waves", "must have unique names for event")
	}

	if len(ec.Categories) > 0 {
		var categoryNames []string
		for _, c := range ec.Categories {
			categoryNames = append(categoryNames, c.Name)
			validateCategory(v, race.ID, ec.ID, ec.EventDate, c)
		}
		v.Check(validator.Unique(categoryNames), "categories", "must have unique names for event")
	}
}

func validateCategory(v *validator.Validator, raceID, eventID uuid.UUID, eventDate time.Time, c *entity.Category) {
	v.Check(raceID != uuid.Nil, "category race_id", "must not be empty")
	v.Check(raceID == c.RaceID, "category race_id", "must correspond to ID of configurated race")
	v.Check(eventID != uuid.Nil, "category event_id", "must not be empty")
	v.Check(eventID == c.EventID, "category event_id", "invalid event ID for category")
	v.Check(c.Name != "", "category name", "must not be empty")
	v.Check(entity.IsValidGender(c.Gender), "gender", "must be male, female or mixed")
	v.Check(c.AgeFrom >= 0, "category age from", "must be greater or equal to 0")
	v.Check(c.AgeTo > 0, "category age to", "must be greater than 0")
	v.Check(c.AgeFrom < c.AgeTo, "category age", "upper age limit must be greater than lower age limit")
}

func validateWave(v *validator.Validator, raceID, eventID uuid.UUID, w *entity.Wave) {
	v.Check(raceID != uuid.Nil, "wave race ID", "must not be null")
	v.Check(raceID == w.RaceID, "wave's race ID", "must correspond to ID of configurated race")
	v.Check(eventID != uuid.Nil, "wave's event ID", "must not be null")
	v.Check(eventID == w.EventID, "wave's evengt ID", "must correspond to ID of configurated event")
	v.Check(w.Name != "", "wave name", "must not be empty")
}

func validateSplit(v *validator.Validator, raceID, eventID uuid.UUID, readers []*entity.TimeReader, split *entity.Split) {
	v.Check(raceID != uuid.Nil, "split's race ID", "must not be nil")
	v.Check(raceID == split.RaceID, "split's race ID", "must correspond to ID of configurated race")
	v.Check(eventID == split.EventID, "split's event ID", "must correspond to ID of configurated event")
	v.Check(eventID != uuid.Nil, "split's event ID", "must not be null")
	v.Check(split.Name != "", "split name", "must not be empty")
	v.Check(split.Type != "", "split type", "must not be empty")
	v.Check(entity.IsValidSplitType(split.Type), "split type", "must be start, standard or finish")
	v.Check(split.DistanceFromStart >= 0, "split distance from start", "must be greater or equal to 0")

	v.Check(split.TimeReaderID.String() != "", "split ID", "must not be empty")
	var tpIDsForLocs []uuid.UUID
	for _, l := range readers {
		tpIDsForLocs = append(tpIDsForLocs, l.ID)
	}

	v.Check(validator.PermittedValue(split.TimeReaderID, tpIDsForLocs...), "split ID", "must have valid corresponded time reader")

	v.Check(split.MinTime >= 0, "split min time", "must be greater or equal to 0")
	v.Check(split.MaxTime >= 0, "split max time", "must be greater or equal to 0")
	v.Check(split.MinLapTime >= 0, "split min lap time", "must be greater or equal to 0")
}
