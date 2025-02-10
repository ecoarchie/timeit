package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Participant struct {
	ID          uuid.UUID      `json:"participant_id"`
	RaceID      uuid.UUID      `json:"race_id"`
	FirstName   string         `json:"first_name"`
	LastName    string         `json:"last_name"`
	Gender      CategoryGender `json:"gender"`
	DateOfBirth time.Time      `json:"date_of_birth"`
	Phone       string         `json:"phone"`
	Comments    string         `json:"comments"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type ParticipantRequest struct {
	RaceID      string         `json:"race_id"`
	FirstName   string         `json:"first_name"`
	LastName    string         `json:"last_name"`
	Gender      CategoryGender `json:"gender"`
	DateOfBirth time.Time      `json:"date_of_birth"`
	Phone       string         `json:"phone"`
	Comments    string         `json:"comments"`
}

func NewParticipant(req ParticipantRequest) (*Participant, error) {
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

	raceID, err := uuid.Parse(req.RaceID)
	if err != nil {
		return nil, fmt.Errorf("wrong race id for participant")
	}
	id := uuid.New()
	return &Participant{
		ID:          id,
		RaceID:      raceID,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Gender:      req.Gender,
		DateOfBirth: req.DateOfBirth,
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
