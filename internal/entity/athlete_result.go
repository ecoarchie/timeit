package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type AthleteResult struct {
	*Athlete
	Results `json:"results"`
}
type (
	ReaderName = string
	SplitID    = uuid.UUID
)

type SplitResult struct {
	SplitID uuid.UUID     `json:"split_id"`
	TOD     time.Time     `json:"tod"`
	GunTime time.Duration `json:"gun_time"`
	NetTime time.Duration `json:"net_time"`
}

type Results map[SplitID]*SplitResult

func NewAthleteResults(p *Athlete) *AthleteResult {
	return &AthleteResult{
		p,
		make(Results),
	}
}

func (pr *AthleteResult) String() string {
	return fmt.Sprintf("{Athlete:%v\nresults:%v}", pr.Athlete, pr.Results)
}

func (tp *SplitResult) String() string {
	return fmt.Sprintf("\n{id: %v,\nTOD: %v,\nGunTime: %v,\nNetTime: %v}", tp.SplitID, tp.TOD, tp.GunTime, tp.NetTime)
}
