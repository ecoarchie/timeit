package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Athlete struct {
	ID          uuid.UUID      `json:"athlete_id"`
	RaceID      uuid.UUID      `json:"race_id"`
	EventID     uuid.UUID      `json:"event_id"`
	WaveID      uuid.UUID      `json:"wave_id"`
	Bib         int            `json:"bib"`
	Chip        int            `json:"chip"`
	FirstName   string         `json:"first_name"`
	LastName    string         `json:"last_name"`
	Gender      CategoryGender `json:"gender"`
	DateOfBirth time.Time      `json:"date_of_birth"`
	CategoryID  uuid.NullUUID  `json:"category_id"`
	Phone       string         `json:"phone"`
	Comments    string         `json:"comments"`
}

type AthleteCreateRequest struct {
	ID          uuid.UUID      `json:"athlete_id"`
	RaceID      uuid.UUID      `json:"race_id"`
	EventID     uuid.UUID      `json:"event_id"`
	WaveID      uuid.UUID      `json:"wave_id"`
	Bib         int            `json:"bib"`
	Chip        int            `json:"chip"`
	FirstName   string         `json:"first_name"`
	LastName    string         `json:"last_name"`
	Gender      CategoryGender `json:"gender"`
	DateOfBirth time.Time      `json:"date_of_birth"`
	CategoryID  uuid.NullUUID  `json:"category_id"`
	Phone       string         `json:"phone"`
	Comments    string         `json:"comments"`
}

type AthleteUpdateRequest struct {
	ID uuid.UUID
	AthleteCreateRequest
}

type Status string

const (
	NYS Status = "not yet started"
	RUN Status = "running"
	FIN Status = "finished"
	DSQ Status = "disqualified"
	QRT Status = "qurantine"
	DNS Status = "pre-race withdrawal"
	DNF Status = "withdrawn during race"
)

func NewAthlete(req AthleteCreateRequest) (*Athlete, error) {
	if req.RaceID == uuid.Nil {
		return nil, fmt.Errorf("athlete race must be assigned")
	}
	if req.EventID == uuid.Nil {
		return nil, fmt.Errorf("athlete event must be assigned")
	}
	if req.WaveID == uuid.Nil {
		return nil, fmt.Errorf("athlete wave be assigned")
	}
	if req.Bib <= 0 {
		return nil, fmt.Errorf("athlete bib must be greater than 0")
	}
	if req.Chip <= 0 {
		return nil, fmt.Errorf("athlete chip must be greater than 0")
	}
	if req.FirstName == "" {
		req.FirstName = "athlete"
	}
	if req.LastName == "" {
		req.LastName = "unknown"
	}

	// gender check
	if req.Gender == "" {
		req.Gender = "unknown"
	}
	if !IsValidGender(req.Gender) {
		return nil, fmt.Errorf("invalid gender for athlete")
	}

	// birth date check
	const zeroBirthDate = "1900-01-01"
	zbd, err := time.Parse(time.DateOnly, zeroBirthDate)
	if err != nil {
		return nil, fmt.Errorf("error parsing zero birth date")
	}
	if req.DateOfBirth.IsZero() {
		req.DateOfBirth = zbd
	}
	if req.DateOfBirth.Before(zbd) || req.DateOfBirth.After(time.Now()) {
		return nil, fmt.Errorf("athlete's birth date is incorrect")
	}

	var id uuid.UUID
	if req.ID == uuid.Nil {
		id = uuid.New()
	} else {
		id = req.ID
	}
	return &Athlete{
		ID:          id,
		RaceID:      req.RaceID,
		EventID:     req.EventID,
		WaveID:      req.WaveID,
		Bib:         req.Bib,
		Chip:        req.Chip,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Gender:      req.Gender,
		DateOfBirth: req.DateOfBirth,
		CategoryID:  req.CategoryID,
		Phone:       req.Phone,
		Comments:    req.Comments,
	}, nil
}

func IsValidGender(c CategoryGender) bool {
	switch c {
	case CategoryGenderFemale, CategoryGenderMale, CategoryGenderMixed, CategoryGenderUnknown:
		return true
	default:
		return false
	}
}

func RandomAthlete(name, surname string, gender CategoryGender, bib, chip int) *Athlete {
	return &Athlete{
		ID:          uuid.New(),
		RaceID:      uuid.New(),
		EventID:     uuid.New(),
		WaveID:      uuid.New(),
		Bib:         bib,
		Chip:        chip,
		FirstName:   name,
		LastName:    name,
		Gender:      gender,
		DateOfBirth: time.Time{},
		CategoryID:  uuid.NullUUID{},
		Phone:       "",
		Comments:    "",
	}
}
