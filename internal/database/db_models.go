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

type Athlete struct {
	ID              uuid.UUID
	RaceID          uuid.UUID
	FirstName       pgtype.Text
	LastName        pgtype.Text
	Gender          CategoryGender
	DateOfBirth     pgtype.Date
	Phone           pgtype.Text
	AthleteComments pgtype.Text
	CreatedAt       pgtype.Timestamp
	UpdatedAt       pgtype.Timestamp
}

type AthleteSplit struct {
	RaceID    uuid.UUID
	EventID   uuid.UUID
	SplitID   uuid.UUID
	AthleteID uuid.UUID
	Tod       pgtype.Timestamp
	GunTime   int64
	NetTime   int64
}

type Category struct {
	ID           uuid.UUID
	RaceID       uuid.UUID
	EventID      uuid.UUID
	CategoryName string
	Gender       CategoryGender
	AgeFrom      int32
	DateFrom     pgtype.Timestamp
	AgeTo        int32
	DateTo       pgtype.Timestamp
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
	EventName        string
	DistanceInMeters int32
	EventDate        pgtype.Timestamp
}

type EventAthlete struct {
	RaceID     uuid.UUID
	EventID    uuid.UUID
	AthleteID  uuid.UUID
	WaveID     uuid.UUID
	CategoryID uuid.NullUUID
	Bib        pgtype.Int4
}

type Race struct {
	ID       uuid.UUID
	RaceName string
	Timezone string
}

type ReaderRecord struct {
	ID         int32
	RaceID     uuid.UUID
	Chip       int32
	Tod        pgtype.Timestamp
	ReaderName string
	CanUse     bool
}

type Split struct {
	ID                uuid.UUID
	RaceID            uuid.UUID
	EventID           uuid.UUID
	SplitName         string
	SplitType         TpType
	DistanceFromStart int32
	TimeReaderID      uuid.UUID
	MinTime           pgtype.Int8
	MaxTime           pgtype.Int8
	MinLapTime        pgtype.Int8
}

type TimeReader struct {
	ID         uuid.UUID
	RaceID     uuid.UUID
	ReaderName string
}

type Wave struct {
	ID         uuid.UUID
	RaceID     uuid.UUID
	EventID    uuid.UUID
	WaveName   string
	StartTime  pgtype.Timestamp
	IsLaunched bool
}
