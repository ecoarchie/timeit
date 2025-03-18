package dto

import (
	"context"
	"fmt"
	"time"

	"github.com/ecoarchie/timeit/pkg/validator"
	"github.com/google/uuid"
)

// FIXME extract all business validation logic to correspondent layer
type RaceDTO struct {
	ID       uuid.UUID `json:"race_id"`
	Name     string    `json:"race_name"`
	Timezone string    `json:"timezone"`
}

type TimeReaderDTO struct {
	ID         uuid.UUID `json:"time_reader_id"`
	RaceID     uuid.UUID `json:"race_id"`
	ReaderName string    `json:"reader_name"`
}

type EventDTO struct {
	ID               uuid.UUID `json:"event_id"`
	RaceID           uuid.UUID `json:"race_id"`
	Name             string    `json:"event_name"`
	DistanceInMeters int       `json:"distance_in_meters"`
	EventDate        string    `json:"event_date"`
}

type SplitDTO struct {
	ID                 uuid.UUID `json:"split_id"`
	RaceID             uuid.UUID `json:"race_id"`
	EventID            uuid.UUID `json:"event_id"`
	Name               string    `json:"split_name"`
	Type               string    `json:"split_type"`
	DistanceFromStart  int       `json:"distance_from_start"`
	TimeReaderID       uuid.UUID `json:"time_reader_id"`
	MinTime            string    `json:"min_time_sec"`
	MaxTime            string    `json:"max_time_sec"`
	MinLapTime         string    `json:"min_lap_time_sec"`
	PreviousLapSplitID uuid.NullUUID
}

type WaveDTO struct {
	ID         uuid.UUID `json:"wave_id"`
	RaceID     uuid.UUID `json:"race_id"`
	EventID    uuid.UUID `json:"event_id"`
	Name       string    `json:"wave_name"`
	StartTime  string    `json:"wave_start_time"`
	IsLaunched bool      `json:"is_launched"`
}

type CategoryDTO struct {
	ID           uuid.UUID `json:"category_id"`
	RaceID       uuid.UUID `json:"race_id"`
	EventID      uuid.UUID `json:"event_id"`
	Name         string    `json:"category_name"`
	Gender       string    `json:"category_gender"`
	AgeFrom      int       `json:"age_from"`
	FromRaceDate bool      `json:"from_race_date"`
	AgeTo        int       `json:"age_to"`
	ToRaceDate   bool      `json:"to_race_date"`
}

type RaceConfig struct {
	*RaceDTO
	TimeReaders []*TimeReaderDTO `json:"time_readers"`
	Events      []*EventConfig   `json:"events"`
}

type EventConfig struct {
	*EventDTO
	Splits     []*SplitDTO    `json:"splits"`
	Waves      []*WaveDTO     `json:"waves"`
	Categories []*CategoryDTO `json:"categories"`
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
			validateEventConfig(v, rc.RaceDTO.ID, rc.TimeReaders, ec)
		}
	}
}

func (rc *RaceConfig) validateRace(v *validator.Validator) {
	v.Check(rc.ID != uuid.Nil, "race_id", "must be valid UUID")
	v.Check(rc.Name != "", "race name", "must not be empty")
}

func (rc *RaceConfig) validateTimeReaders(v *validator.Validator) {
	for _, r := range rc.TimeReaders {
		v.Check(r.RaceID == rc.ID, "timers raceID", "must correspond to ID of configurated race")
	}
}

func validateEventConfig(v *validator.Validator, raceID uuid.UUID, readers []*TimeReaderDTO, ec *EventConfig) {
	v.Check(raceID == ec.RaceID, "race_id for event", "must correspond to ID of configurated race")
	v.Check(ec.ID != uuid.Nil, "event_id", "must not be empty")
	v.Check(ec.Name != "", "event_name", "must not be empty")
	v.Check(ec.DistanceInMeters > 0, "event distance_in_meters", "must be greater than 0")
	v.Check(validator.IsValidTime(time.RFC3339, ec.EventDate), "event_date", "must be date in RFC3339 format")

	if len(ec.Splits) > 0 {
		for _, split := range ec.Splits {
			validateSplit(v, ec.RaceID, ec.ID, readers, split)
		}
	}

	if len(ec.Waves) > 0 {
		for _, w := range ec.Waves {
			validateWave(v, raceID, ec.ID, w)
		}
	}

	if len(ec.Categories) > 0 {
		for _, c := range ec.Categories {
			validateCategory(v, raceID, ec.ID, c)
		}
	}
}

func validateCategory(v *validator.Validator, raceID, eventID uuid.UUID, c *CategoryDTO) {
	v.Check(raceID != uuid.Nil, "category race_id", "must not be empty")
	v.Check(raceID == c.RaceID, "category race_id", "must correspond to ID of configurated race")
	v.Check(eventID != uuid.Nil, "category event_id", "must not be empty")
	v.Check(eventID == c.EventID, "category event_id", "invalid event ID for category")
	v.Check(c.Name != "", "category name", "must not be empty")
}

func validateWave(v *validator.Validator, raceID, eventID uuid.UUID, w *WaveDTO) {
	v.Check(raceID != uuid.Nil, "wave race ID", "must not be null")
	v.Check(raceID == w.RaceID, "wave's race ID", "must correspond to ID of configurated race")
	v.Check(eventID != uuid.Nil, "wave's event ID", "must not be null")
	v.Check(eventID == w.EventID, "wave's event ID", "must correspond to ID of configurated event")
	v.Check(w.Name != "", "wave name", "must not be empty")
	v.Check(validator.IsValidTime(time.RFC3339, w.StartTime), "start_time", "must be date in RFC3339 format")
}

func validateSplit(v *validator.Validator, raceID, eventID uuid.UUID, readers []*TimeReaderDTO, split *SplitDTO) {
	v.Check(raceID != uuid.Nil, "split's race ID", "must not be nil")
	v.Check(raceID == split.RaceID, "split's race ID", "must correspond to ID of configurated race")
	v.Check(eventID == split.EventID, "split's event ID", "must correspond to ID of configurated event")
	v.Check(eventID != uuid.Nil, "split's event ID", "must not be null")
	v.Check(split.Name != "", "split name", "must not be empty")
	v.Check(split.Type != "", "split type", "must not be empty")
	v.Check(
		validator.IsValidDuration(split.MinTime) && validator.IsValidDuration(split.MaxTime) && validator.IsValidDuration(split.MinLapTime),
		"min_max_min_lap_time",
		"must be duration string",
	)

	v.Check(split.TimeReaderID.String() != "", "split ID", "must not be empty")
}
