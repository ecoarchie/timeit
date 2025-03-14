package entity

import (
	"time"

	"github.com/google/uuid"
)

type AthleteSplit struct {
	RaceID    uuid.UUID
	EventID   uuid.UUID
	AthleteID uuid.UUID
	SplitID   uuid.UUID     `json:"split_id"`
	TOD       time.Time     `json:"tod"`
	GunTime   time.Duration `json:"gun_time"`
	NetTime   time.Duration `json:"net_time"`
}

// func NewAthleteResults(p *Athlete) *AthleteResult {
// 	return &AthleteResult{
// 		p,
// 		make(Results),
// 	}
// }

// func (pr *AthleteResult) String() string {
// 	return fmt.Sprintf("{Athlete:%v\nresults:%v}", pr.Athlete, pr.Results)
// }

// func (tp *SplitResult) String() string {
// 	return fmt.Sprintf("\n{id: %v,\nTOD: %v,\nGunTime: %v,\nNetTime: %v}", tp.SplitID, tp.TOD, tp.GunTime, tp.NetTime)
// }
