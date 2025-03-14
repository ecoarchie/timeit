package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ecoarchie/timeit/internal/database"
	"github.com/ecoarchie/timeit/internal/entity"
	"github.com/google/uuid"
)

type ResultsManager interface {
	GetResults(ctx context.Context, raceID, eventID uuid.UUID) ([][]*entity.AthleteSplit, error)
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
func (rs ResultsService) ResultForAthlete(id uuid.UUID) *entity.AthleteSplit {
	return nil
}

// TODO
func (rs ResultsService) GetResults(ctx context.Context, raceID, eventID uuid.UUID) ([][]*entity.AthleteSplit, error) {
	recs, splits, err := rs.AthleteRepo.GetRecordsAndSplitsForEventAthlete(ctx, raceID, eventID)
	if err != nil {
		return nil, err
	}
	for _, s := range splits {
		fmt.Printf("splid ID: %v, name: %s, prevLap: %v, min_time: %v, max_time: %v, min_lap_time: %v\n", s.ID, s.Name, s.PreviousLapSplitID, s.MinTime, s.MaxTime, s.MinLapTime)
	}

	allRecords := [][]*entity.AthleteSplit{}
	splitMap := arrangeSplitsByReaderName(splits)
	for _, r := range recs {
		if len(r.Records) == 0 {
			continue
		}
		res, err := getResultForSingleAthlete(r, splits, splitMap)
		if err != nil {
			return nil, err
		}
		allRecords = append(allRecords, res)
	}
	return allRecords, nil
}

func arrangeSplitsByReaderName(ss []*entity.Split) map[SplitID]*entity.Split {
	ssMap := map[SplitID]*entity.Split{}
	for _, s := range ss {
		ssMap[s.ID] = s
	}
	return ssMap
}

func getResultForSingleAthlete(r database.GetEventAthleteRecordsRow, splits []*entity.Split, splitMap map[SplitID]*entity.Split) ([]*entity.AthleteSplit, error) {
	// create slice for athlete's splits which will be populated further
	singleAthleteRecords := make([]*entity.AthleteSplit, len(splits))
	athleteResultsMap := make(map[SplitID]*entity.AthleteSplit, len(splits))
	fmt.Printf("Check Athlete with ID: %v\n", r.AthleteID)

	for i, recTOD := range r.Records {
		recReader := r.ReaderIds[i]
		// iterate over splits for this event to find valid split for record's tod
		fmt.Printf("Checking TOD %v\n", recTOD)
		for j, s := range splits {
			// check if split reader id matches record's reader name
			if s.TimeReaderID != recReader {
				continue
			}
			// check min_time, max_time constraint
			prevLapSplitResult := athleteResultsMap[s.PreviousLapSplitID.UUID]
			if !isValidSplit(r.WaveStart.Time, recTOD.Time, s, prevLapSplitResult) {
				continue
			}

			// check if such athlete result for this split is already in results map
			_, exist := athleteResultsMap[s.ID]

			// for type 'start' existing results must be overwritten, for 'standard' and 'finish' existing must be kept unchanged
			if !exist || s.Type == entity.SplitTypeStart {
				// FIXME add checking if start type split is not configured at all
				var netTime time.Duration
				if s.Type != entity.SplitTypeStart {
					netTime = recTOD.Time.Sub(singleAthleteRecords[0].TOD)
				}
				res := &entity.AthleteSplit{
					RaceID:    s.RaceID,
					EventID:   s.EventID,
					AthleteID: r.AthleteID,
					SplitID:   s.ID,
					TOD:       recTOD.Time,
					GunTime:   recTOD.Time.Sub(r.WaveStart.Time),
					NetTime:   netTime,
				}
				athleteResultsMap[s.ID] = res
				singleAthleteRecords[j] = res
				continue
			}
		}
		fmt.Println()
		fmt.Println()
	}
	for _, s := range singleAthleteRecords {
		fmt.Println(s)
	}
	return singleAthleteRecords, nil
}

func isValidSplit(waveStart time.Time, tod time.Time, s *entity.Split, prev *entity.AthleteSplit) bool {
	fmt.Printf("Checking split %s, for record %v, with prevResult %v\n", s.Name, tod, prev)
	if tod.Before(waveStart) {
		return false
	}
	validMinTime := waveStart.Add(s.MinTime)
	var validMaxTime time.Time
	if s.MaxTime == 0 {
		validMaxTime = waveStart.Add(time.Hour * 240) // FIXME replace this magic with const max time
	} else {
		validMaxTime = waveStart.Add(s.MaxTime)
	}

	if !((tod.Equal(validMinTime) || tod.After(validMinTime)) && (tod.Equal(validMaxTime) || tod.Before(validMaxTime))) {
		return false
	}
	if prev != nil {
		if prev.TOD.Add(s.MinLapTime).After(tod) {
			return false
		}
	}
	fmt.Printf("Split %s is valid\n\n", s.Name)
	return true
}
