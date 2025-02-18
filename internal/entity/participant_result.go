package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ParticipantResult struct {
	*Participant
	ResultsForTPs `json:"results"`
}
type TimingPointID = uuid.UUID

type TimingPointResult struct {
	TimingPointID uuid.UUID     `json:"timing_point_id"`
	TOD           time.Time     `json:"tod"`
	GunTime       time.Duration `json:"gun_time"`
	NetTime       time.Duration `json:"net_time"`
}

type ResultsForTPs map[TimingPointID]*TimingPointResult

func NewParticipantResults(p *Participant) *ParticipantResult {
	return &ParticipantResult{
		p,
		make(ResultsForTPs),
	}
}

func (pr *ParticipantResult) String() string {
	return fmt.Sprintf("{Participant:%v\nresults:%v}", pr.Participant, pr.ResultsForTPs)
}

func (tp *TimingPointResult) String() string {
	return fmt.Sprintf("\n{id: %v,\nTOD: %v,\nGunTime: %v,\nNetTime: %v}", tp.TimingPointID, tp.TOD, tp.GunTime, tp.NetTime)
}
