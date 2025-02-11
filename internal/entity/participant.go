package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Participant struct {
	ID          uuid.UUID      `json:"participant_id"`
	RaceID      uuid.UUID      `json:"race_id"`
	EventID     uuid.UUID      `json:"event_id"`
	WaveID      uuid.UUID      `json:"wave_id"`
	Bib         int            `json:"bib"`
	Tag         int            `json:"tag"`
	FirstName   string         `json:"first_name"`
	LastName    string         `json:"last_name"`
	Gender      CategoryGender `json:"gender"`
	DateOfBirth time.Time      `json:"date_of_birth"`
	CategoryID  uuid.NullUUID  `json:"category_id"`
	Phone       string         `json:"phone"`
	Comments    string         `json:"comments"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type ParticipantCreateRequest struct {
	RaceID      uuid.UUID      `json:"race_id"`
	EventID     uuid.UUID      `json:"event_id"`
	WaveID      uuid.UUID      `json:"wave_id"`
	Bib         int            `json:"bib"`
	Tag         int            `json:"tag"`
	FirstName   string         `json:"first_name"`
	LastName    string         `json:"last_name"`
	Gender      CategoryGender `json:"gender"`
	DateOfBirth time.Time      `json:"date_of_birth"`
	CategoryID  uuid.NullUUID  `json:"category_id"`
	Phone       string         `json:"phone"`
	Comments    string         `json:"comments"`
}

func NewParticipant(req ParticipantCreateRequest) (*Participant, error) {
	if req.RaceID == uuid.Nil {
		return nil, fmt.Errorf("participant must have race assigned")
	}
	if req.EventID == uuid.Nil {
		return nil, fmt.Errorf("participant must have event assigned")
	}
	if req.WaveID == uuid.Nil {
		return nil, fmt.Errorf("participant must have wave assigned")
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
	if !isValidGender(req.Gender) {
		return nil, fmt.Errorf("invalid gender for participant")
	}

	// birth date check
	const zeroBirthDate = "1900-Jan-01"
	zbd, err := time.Parse(time.DateOnly, zeroBirthDate)
	if err != nil {
		return nil, fmt.Errorf("error parsing zero birth date")
	}
	if req.DateOfBirth.Before(zbd) {
		return nil, fmt.Errorf("participant's birth year is less than 1900")
	}
	if req.DateOfBirth.IsZero() {
		req.DateOfBirth = zbd
	}

	id := uuid.New()
	return &Participant{
		ID:          id,
		RaceID:      req.RaceID,
		EventID:     req.EventID,
		WaveID:      req.WaveID,
		Bib:         req.Bib,
		Tag:         req.Tag,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Gender:      req.Gender,
		DateOfBirth: req.DateOfBirth,
		CategoryID:  req.CategoryID,
		Phone:       req.Phone,
		Comments:    req.Comments,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

func isValidGender(c CategoryGender) bool {
	switch c {
	case CategoryGenderFemale, CategoryGenderMale, CategoryGenderMixed, CategoryGenderUnknown:
		return true
	default:
		return false
	}
}
