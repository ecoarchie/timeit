// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package database

import (
	"database/sql/driver"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type CategoryGender string

const (
	CategoryGenderMale    CategoryGender = "male"
	CategoryGenderFemale  CategoryGender = "female"
	CategoryGenderMixed   CategoryGender = "mixed"
	CategoryGenderUnknown CategoryGender = "unknown"
)

func (e *CategoryGender) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = CategoryGender(s)
	case string:
		*e = CategoryGender(s)
	default:
		return fmt.Errorf("unsupported scan type for CategoryGender: %T", src)
	}
	return nil
}

type NullCategoryGender struct {
	CategoryGender CategoryGender
	Valid          bool // Valid is true if CategoryGender is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullCategoryGender) Scan(value interface{}) error {
	if value == nil {
		ns.CategoryGender, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.CategoryGender.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullCategoryGender) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.CategoryGender), nil
}

type TpType string

const (
	TpTypeStart    TpType = "start"
	TpTypeStandard TpType = "standard"
	TpTypeFinish   TpType = "finish"
)

func (e *TpType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = TpType(s)
	case string:
		*e = TpType(s)
	default:
		return fmt.Errorf("unsupported scan type for TpType: %T", src)
	}
	return nil
}

type NullTpType struct {
	TpType TpType
	Valid  bool // Valid is true if TpType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullTpType) Scan(value interface{}) error {
	if value == nil {
		ns.TpType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.TpType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullTpType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.TpType), nil
}

type BoxRecord struct {
	ID      int32
	RaceID  uuid.UUID
	Chip    int32
	Tod     pgtype.Timestamptz
	BoxName string
	CanUse  bool
}

type Category struct {
	ID           uuid.UUID
	RaceID       uuid.UUID
	EventID      uuid.UUID
	Name         string
	Gender       CategoryGender
	FromAge      int32
	FromRaceDate bool
	ToAge        int32
	ToRaceDate   bool
}

type ChipBib struct {
	RaceID  uuid.UUID
	EventID uuid.UUID
	Chip    int32
	Bib     int32
}

type Event struct {
	ID               uuid.UUID
	RaceID           uuid.UUID
	Name             string
	DistanceInMeters int32
	EventDate        pgtype.Timestamptz
}

type EventLocation struct {
	EventID uuid.UUID
	RaceID  uuid.UUID
	BoxName string
}

type EventParticipant struct {
	RaceID        uuid.UUID
	EventID       uuid.UUID
	ParticipantID uuid.UUID
	WaveID        uuid.UUID
	CategoryID    uuid.UUID
	Bib           pgtype.Int4
}

type Participant struct {
	ID          uuid.UUID
	FirstName   pgtype.Text
	LastName    pgtype.Text
	Gender      CategoryGender
	DateOfBirth pgtype.Date
	Phone       pgtype.Text
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
}

type ParticipantTimingPoint struct {
	RaceID        uuid.UUID
	EventID       uuid.UUID
	TimingPointID uuid.UUID
	ParticipantID uuid.UUID
	Tod           pgtype.Timestamp
	GunTime       int64
	NetTime       int64
}

type PhysicalLocation struct {
	RaceID  uuid.UUID
	BoxName string
}

type Race struct {
	ID       uuid.UUID
	Name     string
	RaceDate pgtype.Date
	Timezone string
}

type TimingPoint struct {
	ID                uuid.UUID
	RaceID            uuid.UUID
	EventID           uuid.UUID
	Name              string
	Type              TpType
	DistanceFromStart int32
	BoxName           string
	MinTimeSec        pgtype.Int4
	MaxTimeSec        pgtype.Int4
	MinLapTimeSec     pgtype.Int4
}

type Wave struct {
	ID         uuid.UUID
	RaceID     uuid.UUID
	EventID    uuid.UUID
	Name       string
	StartTime  pgtype.Timestamptz
	IsLaunched bool
}
