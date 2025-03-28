package service

import (
	"cmp"
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/ecoarchie/timeit/internal/database"
	"github.com/ecoarchie/timeit/internal/entity"
	"github.com/google/uuid"
)

type ResultsManager interface {
	CalculateSplitResultsForEvent(ctx context.Context, raceID, eventID uuid.UUID) ([]*entity.AthleteSplit, error)
	CalculateRanks(ctx context.Context, eventResults []*entity.AthleteSplit)
	CalculateSplitResults(ctx context.Context, raceID uuid.UUID) error
	GetSplitResults(ctx context.Context, raceID uuid.UUID) (map[EventID][]entity.AthleteSplitResults, error)
}

type ResultsService struct {
	AthleteRepo AthleteRepo
}

func NewResultsService(repo AthleteRepo) *ResultsService {
	return &ResultsService{
		AthleteRepo: repo,
	}
}

func (rs ResultsService) GetSplitResults(ctx context.Context, raceID uuid.UUID) (map[EventID][]entity.AthleteSplitResults, error) {
	panic("Not implemented")
}

func (rs *ResultsService) CalculateSplitResults(ctx context.Context, raceID uuid.UUID) error {
	IDs, err := rs.AthleteRepo.GetEventIDsWithWavesStarted(ctx, raceID)
	if err != nil {
		return err
	}
	if len(IDs) == 0 {
		return nil
	}

	var allRecords []*entity.AthleteSplit
	for _, eventID := range IDs {
		eventResults, err := rs.CalculateSplitResultsForEvent(ctx, raceID, eventID)
		if err != nil {
			return err
		}
		allRecords = append(allRecords, eventResults...)
		for _, e := range eventResults {
			if e == nil {
				fmt.Println("WE HAVE NIL HERE")
			}
		}
	}
	start := time.Now()
	err = rs.AthleteRepo.SaveBulkAthleteSplits(ctx, raceID, allRecords)
	if err != nil {
		fmt.Println("error saving bulk athlete split for the whole race")
		return err
	}
	fmt.Printf("saving whole bulk athlete splits took %v\n", time.Since(start))
	return nil
}

func (rs ResultsService) CalculateSplitResultsForEvent(ctx context.Context, raceID, eventID uuid.UUID) ([]*entity.AthleteSplit, error) {
	start := time.Now()
	recs, splits, err := rs.AthleteRepo.GetRecordsAndSplitsForEventAthlete(ctx, raceID, eventID)
	if err != nil {
		fmt.Println("error getting records and splits for event athlete: ", err)
		return nil, err
	}
	fmt.Printf("Time for getting records for eventID = %v is %v\n", eventID, time.Since(start))
	var startSplit *entity.Split
	for _, s := range splits {
		if s.Type == entity.SplitTypeStart {
			startSplit = s
			break
		}
	}

	manualAthleteSplits, err := rs.AthleteRepo.GetManualAthleteSplits(ctx, raceID, eventID)
	if err != nil {
		fmt.Println("Error getting manual")
		return nil, err
	}
	fmt.Println("manual: ", manualAthleteSplits)

	var allRecords []*entity.AthleteSplit
	// start = time.Now()
	for _, r := range recs {
		athleteSplits, potentialStatus, err := calculateSplitResultForSingleAthlete(r, splits, startSplit)
		if err != nil {
			fmt.Println("Error getting result for single athlete: ", err)
			return nil, err
		}
		if mans, ok := manualAthleteSplits[r.AthleteID]; ok {
			replaceWithManual(athleteSplits, mans)
		}
		if entity.ValidStatusTransition(entity.Status(r.StatusFull), potentialStatus) {
			err = rs.AthleteRepo.UpdateStatus(ctx, potentialStatus, raceID, eventID, r.AthleteID)
			if err != nil {
				fmt.Println("error updating status after split calculation")
				return nil, err
			}
		}
		allRecords = append(allRecords, athleteSplits...)
	}
	return allRecords, nil
}

func replaceWithManual(original, manual []*entity.AthleteSplit) []*entity.AthleteSplit {
	fmt.Println("WE HAVE MANUAL")
	for _, m := range manual {
		for i, o := range original {
			if m.SplitID == o.SplitID {
				original[i] = m
			}
		}
	}
	return original
}

func calculateSplitResultForSingleAthlete(r database.GetEventAthleteRecordsCRow, splits []*entity.Split, startSplit *entity.Split) ([]*entity.AthleteSplit, entity.Status, error) {
	// create slice for athlete's splits with zero times values and visited is false
	singleAthleteRecords := entity.NewAthleteSplitsTemlate(splits, r.AthleteID, r.CategoryID, entity.CategoryGender(r.Gender))
	if len(r.RrTod) == 0 {
		return singleAthleteRecords, entity.NYS, nil
	}

	athleteResultsMap := make(map[entity.SplitID]*entity.AthleteSplit, len(splits))
	currentStatus := entity.Status(r.StatusFull)

	for _, rec := range r.RrTod {
		// iterate over splits for this event to find valid split for record's tod
		for j, s := range splits {
			// check if split reader id matches record's reader name
			if s.TimeReaderID != rec.ReaderID {
				continue
			}
			// check min_time, max_time constraint
			prevLapSplitResult := athleteResultsMap[s.PreviousLapSplitID.UUID]
			if !s.IsValidForRecord(r.WaveStart.Time, rec.TOD, prevLapSplitResult) {
				continue
			}

			// check if such athlete result for this split is already in results map
			_, exist := athleteResultsMap[s.ID]

			// for type 'start' existing results must be overwritten, for 'standard' and 'finish' existing must be kept unchanged
			if !exist || s.Type == entity.SplitTypeStart {
				var netTime time.Duration
				if s.Type != entity.SplitTypeStart {
					if startSplit != nil && singleAthleteRecords[0].IsVisited() {
						netTime = rec.TOD.Sub(singleAthleteRecords[0].TOD)
					} else {
						netTime = rec.TOD.Sub(r.WaveStart.Time)
					}
				}
				singleAthleteRecords[j].TOD = rec.TOD
				singleAthleteRecords[j].GunTime = rec.TOD.Sub(r.WaveStart.Time)
				singleAthleteRecords[j].NetTime = netTime
				athleteResultsMap[s.ID] = singleAthleteRecords[j]

				if s.Type == entity.SplitTypeFinish {
					currentStatus = entity.FIN
					// no need to calculate records further since we found the finish record
					break
				}
				currentStatus = entity.RUN
			}
		}
	}
	return singleAthleteRecords, currentStatus, nil
}

func (rs ResultsService) CalculateRanks(ctx context.Context, eventResults []*entity.AthleteSplit) {
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
	// for i, r := range eventResults {
	// 	fmt.Printf("%d. res = %+v\n", i, r)
	// }
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
