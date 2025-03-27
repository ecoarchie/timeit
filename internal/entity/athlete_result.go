package entity

import (
	"time"

	"github.com/google/uuid"
)

type AthleteSplit struct {
	RaceID          uuid.UUID
	EventID         uuid.UUID
	AthleteID       uuid.UUID
	SplitID         uuid.UUID
	TOD             time.Time
	GunTime         time.Duration
	NetTime         time.Duration
	Gender          CategoryGender
	CategoryID      uuid.NullUUID
	GunRankOverall  int
	GunRankGender   int
	GunRankCategory int
	NetRankOverall  int
	NetRankGender   int
	NetRankCategory int
	Visited         bool
}

type SplitData struct {
	SplitID         uuid.UUID
	TOD             time.Time
	GunTime         time.Duration
	NetTime         time.Duration
	GunRankOverall  int
	GunRankGender   int
	GunRankCategory int
	NetRankOverall  int
	NetRankGender   int
	NetRankCategory int
}

type AthleteSplitResults struct {
	AthleteID  uuid.UUID
	Gender     CategoryGender
	CategoryID uuid.NullUUID
	Splits     map[string]SplitData
}

func NewAthleteSplitsTemlate(ss []*Split, athleteID uuid.UUID, categoryID uuid.NullUUID, gender CategoryGender) []*AthleteSplit {
	result := make([]*AthleteSplit, len(ss))
	for i, s := range ss {
		result[i] = &AthleteSplit{
			RaceID:     s.RaceID,
			EventID:    s.EventID,
			AthleteID:  athleteID,
			SplitID:    s.ID,
			Gender:     gender,
			CategoryID: categoryID,
		}
	}
	return result
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
