package entity

import (
	"time"

	"github.com/google/uuid"
)

const zeroBirthDate = "1900-Jan-01"

type Participant struct {
	ID          uuid.UUID      `json:"id"`
	FirstName   string         `json:"first_name"`
	LastName    string         `json:"last_name"`
	Gender      CategoryGender `json:"gender"`
	DateOfBirth time.Time      `json:"date_of_birth"`
	Phone       string         `json:"phone"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type NewParticipantRequest struct {
	FirstName   string         `json:"first_name"`
	LastName    string         `json:"last_name"`
	Gender      CategoryGender `json:"gender"`
	DateOfBirth time.Time      `json:"date_of_birth"`
	Phone       string         `json:"phone"`
}

func NewParticipant(req NewParticipantRequest) (*Participant, error) {
	if req.FirstName == "" {
		req.FirstName = "athlete"
	}
	if req.LastName == "" {
		req.LastName = "unknown"
	}
	if req.Gender == "" {
		req.Gender = "unknown"
	}
	if req.DateOfBirth.IsZero() {
		req.DateOfBirth, _ = time.Parse("2006-Jan-02", zeroBirthDate)
	}
	id := uuid.New()
	return &Participant{
		ID:          id,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Gender:      req.Gender,
		DateOfBirth: req.DateOfBirth,
		Phone:       req.Phone,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}
