package service

import (
	"cmp"
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/ecoarchie/timeit/internal/database"
	"github.com/ecoarchie/timeit/internal/entity"
	"github.com/ecoarchie/timeit/pkg/pgxmapper"
	"github.com/google/uuid"
)

type ResultsManager interface {
	GetResultsForEvent(ctx context.Context, raceID, eventID uuid.UUID) ([]*entity.AthleteSplit, error)
	CalculateRanks(ctx context.Context, eventResults []*entity.AthleteSplit) ([]*entity.AthleteSplit, error)
	GetResults(ctx context.Context, raceID uuid.UUID) (map[EventID][]*entity.AthleteSplit, error)
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

func (rs *ResultsService) RecalculateAthleteResult(ctx context.Context, raceID uuid.UUID) error {
	return nil
}

func (rs *ResultsService) GetResults(ctx context.Context, raceID uuid.UUID) (map[EventID][]*entity.AthleteSplit, error) {
	IDs, err := rs.AthleteRepo.GetEventIDsWithWavesStarted(ctx, raceID)
	if err != nil {
		return nil, err
	}
	if len(IDs) == 0 {
		return nil, nil
	}
	// FIXME
	m := make(map[EventID][]*entity.AthleteSplit)
	for _, eventID := range IDs {
		eventResults, err := rs.GetResultsForEvent(ctx, raceID, eventID)
		if err != nil {
			return nil, err
		}
		for _, e := range eventResults {
			if e == nil {
				fmt.Println("WE HAVE NIL HERE")
			}
		}
		res, err := rs.CalculateRanks(ctx, eventResults)
		if err != nil {
			return nil, err
		}
		m[eventID] = res
		// fmt.Println("Res: ", res)
	}
	return m, nil
}

func (rs ResultsService) CalculateRanks(ctx context.Context, eventResults []*entity.AthleteSplit) ([]*entity.AthleteSplit, error) {
	gunCmp := func(a, b *entity.AthleteSplit) int {
		if a == nil && b == nil {
			return 0 // Consider nils equal
		}
		if a == nil {
			return -1 // Nil comes first
		}
		if b == nil {
			return 1 // Non-nil comes after nil
		}
		n := cmp.Compare(a.GunTime, b.GunTime)
		return n
	}
	netCmp := func(a, b *entity.AthleteSplit) int {
		if a == nil && b == nil {
			return 0 // Consider nils equal
		}
		if a == nil {
			return -1 // Nil comes first
		}
		if b == nil {
			return 1 // Non-nil comes after nil
		}
		return cmp.Compare(a.NetTime, b.NetTime)
	}

	start := time.Now()
	calculateRanks(eventResults, gunCmp, GUN)
	calculateRanks(eventResults, netCmp, NET)
	fmt.Printf("RANKS calculation for event took - %v\n", time.Since(start))
	for i, r := range eventResults {
		fmt.Printf("%d. res = %+v\n", i, r)
	}

	return eventResults, nil
}

const (
	GUN = "gun"
	NET = "net"
)

func calculateRanks(aSplits []*entity.AthleteSplit, compf func(a, b *entity.AthleteSplit) int, time string) {
	resGroupRank := make(map[string]int, len(aSplits))
	slices.SortFunc(aSplits, compf)
	for _, s := range aSplits {
		if s == nil {
			continue
		}
		oGroup := fmt.Sprintf("%s-%s", "overall", s.SplitID.String())
		cGroup := fmt.Sprintf("%s-%s", s.CategoryID.UUID.String(), s.SplitID.String())
		gGroup := fmt.Sprintf("%s-%s", s.Gender, s.SplitID.String())

		resGroupRank[oGroup]++
		if s.CategoryID.Valid {
			resGroupRank[cGroup]++
		}
		resGroupRank[gGroup]++

		switch time {
		case GUN:
			s.GunRankOverall = resGroupRank[oGroup]
			if s.CategoryID.Valid {
				s.GunRankCategory = resGroupRank[cGroup]
			}
			s.GunRankGender = resGroupRank[gGroup]
		case NET:
			s.NetRankOverall = resGroupRank[oGroup]
			if s.CategoryID.Valid {
				s.NetRankCategory = resGroupRank[cGroup]
			}
			s.NetRankGender = resGroupRank[gGroup]
		}
	}
}

func calculateNetRanks(splits []*entity.AthleteSplit, compf func(a, b *entity.AthleteSplit) int) {
	resGroupRank := make(map[string]int)
	slices.SortFunc(splits, compf)
	for _, s := range splits {
		if s == nil {
			continue
		}
		oGroup := fmt.Sprintf("%s-%s", "overall", s.SplitID.String())
		cGroup := fmt.Sprintf("%s-%s", s.CategoryID.UUID.String(), s.SplitID.String())
		gGroup := fmt.Sprintf("%s-%s", s.Gender, s.SplitID.String())

		// calculate rank for overall
		resGroupRank[oGroup]++
		s.NetRankOverall = resGroupRank[oGroup]

		// calculate rank for category
		if s.CategoryID.Valid {
			resGroupRank[cGroup]++
			s.NetRankCategory = resGroupRank[cGroup]
		}

		// calculate rank for gender
		resGroupRank[gGroup]++
		s.NetRankGender = resGroupRank[gGroup]
	}
}

// TODO
func (rs ResultsService) GetResultsForEvent(ctx context.Context, raceID, eventID uuid.UUID) ([]*entity.AthleteSplit, error) {
	start := time.Now()
	recs, splits, err := rs.AthleteRepo.GetRecordsAndSplitsForEventAthlete(ctx, raceID, eventID)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Time for getting records for eventID = %v is %v\n", eventID, time.Since(start))
	var startSplit *entity.Split
	for _, s := range splits {
		fmt.Printf("splid ID: %v, name: %s, prevLap: %v, min_time: %v, max_time: %v, min_lap_time: %v\n", s.ID, s.Name, s.PreviousLapSplitID, s.MinTime, s.MaxTime, s.MinLapTime)
		if s.Type == entity.SplitTypeStart {
			startSplit = s
			break
		}
	}

	allRecords := []*entity.AthleteSplit{}
	start = time.Now()
	for _, r := range recs {
		if len(r.RrTod) == 0 {
			continue
		}
		// start = time.Now()
		res, err := getResultForSingleAthlete(r, splits, startSplit)
		if err != nil {
			fmt.Println("Error getting result for single athlete: ", err)
			return nil, err
		}
		// fmt.Printf("Time for calculating results for single athleteID = %v is %v\n", r.AthleteID, time.Since(start))
		allRecords = append(allRecords, res...)
	}
	fmt.Printf("Time for calculating ALL records for event = %v is %v\n", eventID, time.Since(start))

	// saving calculated athlete splits results into db
	start = time.Now()
	saveAthleteSplitsParams := []database.CreateAthleteSplitsParams{}
	for _, ar := range allRecords {
		if ar == nil {
			continue
		}
		saveAthleteSplitsParams = append(saveAthleteSplitsParams, database.CreateAthleteSplitsParams{
			RaceID:    ar.RaceID,
			EventID:   ar.EventID,
			SplitID:   ar.SplitID,
			AthleteID: ar.AthleteID,
			Tod:       pgxmapper.TimeToPgxTimestamp(ar.TOD),
			GunTime:   pgxmapper.DurationToPgxInterval(ar.GunTime),
			NetTime:   pgxmapper.DurationToPgxInterval(ar.NetTime),
		})
	}
	err = rs.AthleteRepo.SaveAthleteSplits(ctx, saveAthleteSplitsParams)
	fmt.Printf("Time for inserting ALL athlete splits results is %v\n", time.Since(start))
	if err != nil {
		fmt.Println("Error inserting athlete splits", err)
		return nil, err
	}
	return allRecords, nil
}

func getResultForSingleAthlete(r database.GetEventAthleteRecordsCRow, splits []*entity.Split, startSplit *entity.Split) ([]*entity.AthleteSplit, error) {
	// create slice for athlete's splits which will be populated further
	singleAthleteRecords := make([]*entity.AthleteSplit, len(splits))
	athleteResultsMap := make(map[entity.SplitID]*entity.AthleteSplit, len(splits))
	// fmt.Printf("Check Athlete with ID: %v\n", r.AthleteID)

	for _, recTOD := range r.RrTod {
		// recReader := recTOD.ReaderID
		// iterate over splits for this event to find valid split for record's tod
		// fmt.Printf("Checking TOD %v\n", recTOD)
		for j, s := range splits {
			// check if split reader id matches record's reader name
			if s.TimeReaderID != recTOD.ReaderID {
				continue
			}
			// check min_time, max_time constraint
			prevLapSplitResult := athleteResultsMap[s.PreviousLapSplitID.UUID]
			if !s.IsValidForRecord(r.WaveStart.Time, recTOD.TOD, prevLapSplitResult) {
				continue
			}

			// check if such athlete result for this split is already in results map
			_, exist := athleteResultsMap[s.ID]

			// for type 'start' existing results must be overwritten, for 'standard' and 'finish' existing must be kept unchanged
			if !exist || s.Type == entity.SplitTypeStart {
				// FIXME add checking if start type split is not configured or missed at all
				var netTime time.Duration
				if s.Type != entity.SplitTypeStart {
					if startSplit != nil && singleAthleteRecords[0] != nil {
						netTime = recTOD.TOD.Sub(singleAthleteRecords[0].TOD)
					} else {
						netTime = recTOD.TOD.Sub(r.WaveStart.Time)
					}
				}
				res := &entity.AthleteSplit{
					RaceID:     s.RaceID,
					EventID:    s.EventID,
					AthleteID:  r.AthleteID,
					SplitID:    s.ID,
					TOD:        recTOD.TOD,
					GunTime:    recTOD.TOD.Sub(r.WaveStart.Time),
					NetTime:    netTime,
					Gender:     entity.CategoryGender(r.Gender),
					CategoryID: r.CategoryID,
				}
				athleteResultsMap[s.ID] = res
				singleAthleteRecords[j] = res
				continue
			}
		}
		// fmt.Println()
		// fmt.Println()
	}
	// for _, s := range singleAthleteRecords {
	// 	fmt.Println(s)
	// }
	return singleAthleteRecords, nil
}
