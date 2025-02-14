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

func (pr *ParticipantResult) String() string {
	return fmt.Sprintf("{Participant:%v\nresults:%v}", pr.Participant, pr.ResultsForTPs)
}

func (tp *TimingPointResult) String() string {
	return fmt.Sprintf("\n{id: %v,\nTOD: %v,\nGunTime: %v,\nNetTime: %v}", tp.TimingPointID, tp.TOD, tp.GunTime, tp.NetTime)
}

func NewParticipantResults(p *Participant, l int) *ParticipantResult {
	return &ParticipantResult{
		p,
		make(ResultsForTPs, l),
	}
}

type TimingPointResult struct {
	TimingPointID uuid.UUID `json:"timing_point_id"`
	TOD           time.Time `json:"tod"`
	GunTime       int64     `json:"gun_time"`
	NetTime       int64     `json:"net_time"`
}

type ResultsForTPs map[string]*TimingPointResult

func (rf ResultsForTPs) SetGun(tpName string, t int64) {
	rf[tpName].GunTime = t
}

func (rf ResultsForTPs) SetNet(tpName string, t int64) {
	rf[tpName].NetTime = t
}

func (rf ResultsForTPs) SetTOD(tpName string, t time.Time) {
	rf[tpName].TOD = t
}
