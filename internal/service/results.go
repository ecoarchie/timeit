package service

import (
	"time"

	"github.com/ecoarchie/timeit/internal/entity"
	"github.com/google/uuid"
)

type ResultsManager interface {
	ResultForAthlete(id uuid.UUID) *entity.AthleteResult
}

type ResultsService struct {
	AthleteRepo AthleteRepo
}

func NewResultsService(repo AthleteRepo) *ResultsService {
	return &ResultsService{
		AthleteRepo: repo,
	}
}

// TODO
func (rs ResultsService) ResultForAthlete(id uuid.UUID) *entity.AthleteResult {
	return nil
}

// TODO
func (rs ResultsService) GetResults(pr *entity.AthleteResult, recs []entity.ReaderRecord, waveStartTime time.Time, tps []*entity.Split) (*entity.AthleteResult, error) {
	// tps are splits for athletes' event from race cache
	// tps are sorted by distance from start
	// recs are only for particular athlete chip, sorted by TOD
	// resc must be valid, canUse = true
	// wave is athletes' wave. assuming that wave status is 'started'
	tpBoxes := make(map[TimeReaderID][]*entity.Split)
	var startingSplit *entity.Split
	// ps := entity.NewAthleteResults(athlete)
	for _, tp := range tps {
		tpBoxes[tp.TimeReaderID] = append(tpBoxes[tp.TimeReaderID], tp)
		if tp.Type == entity.SplitTypeStart {
			startingSplit = tp
		}
	}
	for _, r := range recs {
		// BUG r.RaceID should be r.ReaderName. But since tp doesn't have boxname field - I need to rewrite logic completly
		res := rs.ResultForRecord(r, waveStartTime, tpBoxes[r.RaceID], pr, startingSplit)
		if res == nil {
			continue
		}
		// BUG check if split of type 'start' doesn't exist at all
		_, exists := pr.Results[res.SplitID]
		if exists && res.SplitID != startingSplit.ID {
			// TODO skip creating result for finish and standard types of Splits. Later the rule for finish and standard type may change from 'first read' to 'last read'
			continue
		}
		// BUG check for prevlap rule. After that delete pr passed to ResultForRecord func
		pr.Results[res.SplitID] = res
	}
	return pr, nil
}

func (rs ResultsService) ResultForRecord(r entity.ReaderRecord, waveStartTime time.Time, tps []*entity.Split, pr *entity.AthleteResult, startingSplit *entity.Split) *entity.SplitResult {
	if r.TOD.Before(waveStartTime) {
		// record is before wave start, thus check next record
		return nil
	}
	res := &entity.SplitResult{}
	for i, tp := range tps {
		// if tp.Type == entity.SplitTypeFinish || tp.Type == entity.SplitTypeStandard {
		// 	// if there is already result for standard or finish Split then skip
		// 	if _, ok := ps.Results[tp.ID]; ok {
		// 		return nil
		// 	}
		// }
		// ps.Results[tp.ID] = &entity.SplitResult{
		// 	SplitID: tp.ID,
		// }
		validMinTime := waveStartTime.Add(time.Duration(tp.MinTimeSec))
		var validMaxTime time.Time
		if tp.MaxTimeSec == 0 {
			validMaxTime = waveStartTime.Add(time.Hour * 240)
		} else {
			validMaxTime = waveStartTime.Add(time.Duration(tp.MaxTimeSec) * time.Second)
		}
		if (r.TOD.Equal(validMinTime) || r.TOD.After(validMinTime)) && (r.TOD.Equal(validMaxTime) || r.TOD.Before(validMaxTime)) {
			if i > 0 {
				prevLapSplit := tps[i-1]
				if r.TOD.Sub(pr.Results[prevLapSplit.ID].TOD) < time.Duration(tp.MinLapTimeSec)*time.Second {
					return nil
				}
			}
			recTODfromStart := r.TOD.Sub(waveStartTime)
			var startCalculated bool
			if startingSplit != nil {
				_, startCalculated = pr.Results[startingSplit.ID]
			}

			res.SplitID = tp.ID
			// Calculate gun time
			res.GunTime = recTODfromStart

			// TODO add test case with absent start point
			// Calculate net time
			if tp.Type == entity.SplitTypeStart || (tp.Type != entity.SplitTypeStart && !startCalculated) || startingSplit == nil {
				res.NetTime = recTODfromStart
			} else {
				res.NetTime = r.TOD.Sub(pr.Results[startingSplit.ID].TOD)
			}

			// Calculate TOD
			res.TOD = r.TOD

			// skip checking the rest timming points for that reader_name
			return res
		}
	}
	return nil
}
