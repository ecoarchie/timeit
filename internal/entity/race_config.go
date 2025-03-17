package entity

import (
	"context"
	"fmt"

	"github.com/ecoarchie/timeit/pkg/validator"
	"github.com/google/uuid"
)

type RaceConfig struct {
	*Race
	TimeReaders []*TimeReader  `json:"time_readers"`
	Events      []*EventConfig `json:"events"`
}

type EventConfig struct {
	*Event
	Splits     []*SplitConfig `json:"splits"`
	Waves      []*Wave        `json:"waves"`
	Categories []*Category    `json:"categories"`
}
type SplitConfig struct {
	ID                 uuid.UUID `json:"split_id"`
	RaceID             uuid.UUID `json:"race_id"`
	EventID            uuid.UUID `json:"event_id"`
	Name               string    `json:"split_name"`
	Type               SplitType `json:"split_type"`
	DistanceFromStart  int       `json:"distance_from_start"`
	TimeReaderID       uuid.UUID `json:"time_reader_id"`
	MinTime            Duration  `json:"min_time_sec"`
	MaxTime            Duration  `json:"max_time_sec"`
	MinLapTime         Duration  `json:"min_lap_time_sec"`
	PreviousLapSplitID uuid.NullUUID
}

type RaceFormData struct {
	Id       string `json:"race_id"`
	Name     string `json:"race_name"`
	Timezone string `json:"timezone"`
}

func (rc RaceConfig) String() string {
	return fmt.Sprintf(`{
	RaceID: %s,
	Name: %s,
	Timezone: %s,
	Events: 
		%+v
}`, rc.Name, rc.Name, rc.Timezone, rc.Events)
}

func (rc *RaceConfig) Validate(ctx context.Context, v *validator.Validator) {
	rc.validateRace(v)

	v.Check(len(rc.TimeReaders) > 0, "time readers", "race must have at least one time reader")
	if len(rc.TimeReaders) > 0 {
		rc.validateTimeReaders(v)
	}

	v.Check(len(rc.Events) != 0, "events", "must be at least one")
	if len(rc.Events) > 0 {
		var eventNames []string
		for _, e := range rc.Events {
			eventNames = append(eventNames, e.Name)
		}
		v.Check(validator.Unique(eventNames), "event names", "must be unique")
		for _, ec := range rc.Events {
			validateEventConfig(v, rc.Race.ID, rc.TimeReaders, ec)
		}
	}
}

func (rc *RaceConfig) validateRace(v *validator.Validator) {
	v.Check(IsIANATimezone(rc.Timezone), "timezone", "must be valid IANA timezone")
	v.Check(rc.Name != "", "race name", "must not be empty")
}

func (rc *RaceConfig) validateTimeReaders(v *validator.Validator) {
	var timeReadersNames []string
	for _, r := range rc.TimeReaders {
		timeReadersNames = append(timeReadersNames, r.ReaderName)
		v.Check(r.RaceID == rc.ID, "timers raceID", "must correspond to ID of configurated race")
	}
	v.Check(validator.Unique(timeReadersNames), "time readers names", "must be unique")
}

func validateEventConfig(v *validator.Validator, raceID uuid.UUID, readers []*TimeReader, ec *EventConfig) {
	v.Check(raceID == ec.RaceID, "race_id for event", "must correspond to ID of configurated race")
	v.Check(ec.ID != uuid.Nil, "event_id", "must not be empty")
	v.Check(ec.Name != "", "event_name", "must not be empty")
	v.Check(ec.DistanceInMeters > 0, "event distance_in_meters", "must be greater than 0")

	v.Check(len(ec.Splits) != 0, "splits", "event must have at least one split")
	if len(ec.Splits) > 0 {
		var splitsNames []string
		splitTypeQty := make(map[SplitType]int)
		for _, split := range ec.Splits {
			splitsNames = append(splitsNames, split.Name)
			splitTypeQty[split.Type]++
			validateSplit(v, ec.RaceID, ec.ID, readers, split)
		}
		v.Check(validator.Unique(splitsNames), "splits", "must have unique names for event")
		v.Check(splitTypeQty[SplitTypeStart] < 2, "split with type start", "must be 0 or 1")
		v.Check(splitTypeQty[SplitTypeFinish] == 1, "split with type finish", "must be only 1")
	}

	v.Check(len(ec.Waves) > 0, "waves", "must be at least one for event")
	if len(ec.Waves) > 0 {
		var wavesNames []string
		for _, w := range ec.Waves {
			wavesNames = append(wavesNames, w.Name)
			validateWave(v, raceID, ec.ID, w)
		}
		v.Check(validator.Unique(wavesNames), "waves", "must have unique names for event")
	}

	if len(ec.Categories) > 0 {
		var categoryNames []string
		for _, c := range ec.Categories {
			categoryNames = append(categoryNames, c.Name)
			validateCategory(v, raceID, ec.ID, c)
		}
		v.Check(validator.Unique(categoryNames), "categories", "must have unique names for event")
	}
}

func validateCategory(v *validator.Validator, raceID, eventID uuid.UUID, c *Category) {
	v.Check(raceID != uuid.Nil, "category race_id", "must not be empty")
	v.Check(raceID == c.RaceID, "category race_id", "must correspond to ID of configurated race")
	v.Check(eventID != uuid.Nil, "category event_id", "must not be empty")
	v.Check(eventID == c.EventID, "category event_id", "invalid event ID for category")
	v.Check(c.Name != "", "category name", "must not be empty")
	v.Check(IsValidGender(c.Gender), "gender", "must be male, female or mixed")
	v.Check(c.AgeFrom >= 0, "category age from", "must be greater or equal to 0")
	v.Check(c.AgeTo > 0, "category age to", "must be greater than 0")
	v.Check(c.AgeFrom < c.AgeTo, "category age", "upper age limit must be greater than lower age limit")
}

func validateWave(v *validator.Validator, raceID, eventID uuid.UUID, w *Wave) {
	v.Check(raceID != uuid.Nil, "wave race ID", "must not be null")
	v.Check(raceID == w.RaceID, "wave's race ID", "must correspond to ID of configurated race")
	v.Check(eventID != uuid.Nil, "wave's event ID", "must not be null")
	v.Check(eventID == w.EventID, "wave's evengt ID", "must correspond to ID of configurated event")
	v.Check(w.Name != "", "wave name", "must not be empty")
}

func validateSplit(v *validator.Validator, raceID, eventID uuid.UUID, readers []*TimeReader, split *SplitConfig) {
	v.Check(raceID != uuid.Nil, "split's race ID", "must not be nil")
	v.Check(raceID == split.RaceID, "split's race ID", "must correspond to ID of configurated race")
	v.Check(eventID == split.EventID, "split's event ID", "must correspond to ID of configurated event")
	v.Check(eventID != uuid.Nil, "split's event ID", "must not be null")
	v.Check(split.Name != "", "split name", "must not be empty")
	v.Check(split.Type != "", "split type", "must not be empty")
	v.Check(IsValidSplitType(split.Type), "split type", "must be start, standard or finish")
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
