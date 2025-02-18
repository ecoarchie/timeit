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
func (rs ResultsService) GetResults(pr *entity.ParticipantResult, recs []entity.BoxRecord, waveStartTime time.Time, tps []*entity.TimingPoint) (*entity.ParticipantResult, error) {
	// tps are timing points for participants' event from race cache
	// tps are sorted by distance from start
	// recs are only for particular participant chip, sorted by TOD
	// resc must be valid, canUse = true
	// wave is participants' wave. assuming that wave status is 'started'
	tpBoxes := make(map[string][]*entity.TimingPoint)
	var startingTP *entity.TimingPoint
	// ps := entity.NewParticipantResults(participant)
	for _, tp := range tps {
		tpBoxes[tp.BoxName] = append(tpBoxes[tp.BoxName], tp)
		if tp.Type == entity.TPTypeStart {
			startingTP = tp
		}
	}
	for _, r := range recs {
		res := rs.ResultForRecord(r, waveStartTime, tpBoxes[r.BoxName], pr, startingTP)
		if res == nil {
			continue
		}
		_, exists := pr.ResultsForTPs[res.TimingPointID]
		if exists && res.TimingPointID != startingTP.ID {
			// TODO skip creating result for finish and standard types of TPs. Later check validation rule to accept the record
			continue
		}
		// BUG check for prevlap rule. After that delete pr passed to ResultForRecord func
		pr.ResultsForTPs[res.TimingPointID] = res
	}
	return pr, nil
}

func (rs ResultsService) ResultForRecord(r entity.BoxRecord, waveStartTime time.Time, tps []*entity.TimingPoint, pr *entity.ParticipantResult, startingTP *entity.TimingPoint) *entity.TimingPointResult {
	if r.TOD.Before(waveStartTime) {
		// record is before wave start, thus check next record
		return nil
	}
	res := &entity.TimingPointResult{}
	for i, tp := range tps {
		// if tp.Type == entity.TPTypeFinish || tp.Type == entity.TPTypeStandard {
		// 	// if there is already result for standard or finish TP then skip
		// 	if _, ok := ps.ResultsForTPs[tp.ID]; ok {
		// 		return nil
		// 	}
		// }
		// ps.ResultsForTPs[tp.ID] = &entity.TimingPointResult{
		// 	TimingPointID: tp.ID,
		// }
		validMinTime := waveStartTime.Add(time.Duration(tp.MinTimeSec))
		var validMaxTime time.Time
		if tp.MaxTimeSec == 0 {
			validMaxTime = waveStartTime.Add(time.Hour * 24)
		} else {
			validMaxTime = waveStartTime.Add(time.Duration(tp.MaxTimeSec) * time.Second)
		}
		if (r.TOD.Equal(validMinTime) || r.TOD.After(validMinTime)) && (r.TOD.Equal(validMaxTime) || r.TOD.Before(validMaxTime)) {
			if i > 0 {
				prevLapTP := tps[i-1]
				if r.TOD.Sub(pr.ResultsForTPs[prevLapTP.ID].TOD) < time.Duration(tp.MinLapTimeSec)*time.Second {
					return nil
				}
			}
			recTODfromStart := r.TOD.Sub(waveStartTime)
			_, startExists := pr.ResultsForTPs[startingTP.ID]

			res.TimingPointID = tp.ID
			// Calculate gun time
			res.GunTime = recTODfromStart

			// Calculate net time
			if tp.Type == entity.TPTypeStart || (tp.Type != entity.TPTypeStart && !startExists) {
				res.NetTime = recTODfromStart
			} else {
				res.NetTime = r.TOD.Sub(pr.ResultsForTPs[startingTP.ID].TOD)
			}

			// Calculate TOD
			res.TOD = r.TOD

			// skip checking the rest timming points for that box_name
			return res
		}
	}
	return nil
}
