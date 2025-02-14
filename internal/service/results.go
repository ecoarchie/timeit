package service

import (
	"time"

	"github.com/ecoarchie/timeit/internal/entity"
	"github.com/google/uuid"
)

type Results interface {
	ResultForParticipant(id uuid.UUID) *entity.ParticipantResult
}

type ResultsService struct {
	ParticipantRepo ParticipantRepo
}

// TODO
func (rs ResultsService) ResultForParticipant(id uuid.UUID) *entity.ParticipantResult {
	return nil
}

// TODO
func (rs ResultsService) GetResult(p *entity.Participant, recs []entity.BoxRecord, w entity.Wave, tps []*entity.TimingPoint) (*entity.ParticipantResult, error) { // tps are possible timing points for participants' event
	// recs are sorted by TOD
	// tps are sorted by distance from start
	// resc must be valid, canUse = true
	// w is participants' wave. assuming that wave status is started
	tpBoxes := make(map[string][]*entity.TimingPoint)
	var startingTP *entity.TimingPoint
	ps := entity.NewParticipantResults(p, len(tps))
	for _, tp := range tps {
		tp.SetValidMinMaxTimes(w.StartTime)
		tpBoxes[tp.BoxName] = append(tpBoxes[tp.BoxName], tp)
		if tp.Type == entity.TPTypeStart {
			startingTP = tp
		}
	}
outer:
	for _, r := range recs {
		if r.TOD.Before(w.StartTime) {
			// record is before wave start, thus check next record
			continue
		}
		for i, tp := range tpBoxes[r.BoxName] {
			if tp.Type == entity.TPTypeFinish || tp.Type == entity.TPTypeStandard {
				if _, ok := ps.ResultsForTPs[tp.Name]; ok {
					continue outer
				}
			}
			ps.ResultsForTPs[tp.Name] = &entity.TimingPointResult{
				TimingPointID: tp.ID,
			}
			if (r.TOD.Equal(tp.ValidMinTime) || r.TOD.After(tp.ValidMinTime)) && (r.TOD.Equal(tp.ValidMaxTime) || r.TOD.Before(tp.ValidMaxTime)) {
				if i > 0 {
					prevLapTP := tpBoxes[r.BoxName][i-1]
					if r.TOD.Sub(ps.ResultsForTPs[prevLapTP.Name].TOD) < time.Duration(tp.MinLapTimeSec)*time.Second {
						continue outer
					}
				}
				ps.ResultsForTPs[tp.Name].GunTime = int64(r.TOD.Sub(w.StartTime).Microseconds())
				if tp.Name == startingTP.Name {
					ps.ResultsForTPs[tp.Name].NetTime = ps.ResultsForTPs[tp.Name].GunTime
				} else {
					ps.ResultsForTPs[tp.Name].NetTime = int64(r.TOD.Sub(ps.ResultsForTPs[startingTP.Name].TOD).Microseconds())
				}
				ps.ResultsForTPs[tp.Name].TOD = r.TOD
			}
		}
	}
	return ps, nil
}
