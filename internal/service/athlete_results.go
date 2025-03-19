package service

import (
	"context"
	"fmt"

	"github.com/ecoarchie/timeit/internal/entity"
	"github.com/google/uuid"
)

type AthleteResultsManager interface {
	RaceConfigurator
	AthleteManager
	ResultsManager
	RecalculateAthleteResult(ctx context.Context, raceID uuid.UUID) error
	GetResults(ctx context.Context, raceID uuid.UUID) (map[EventID][]*entity.AthleteSplit, error)
}

type AthleteResultsService struct {
	RaceConfigurator
	AthleteManager
	ResultsManager
}

func NewAthleteResultsService(r RaceConfigurator, a AthleteManager, res ResultsManager) *AthleteResultsService {
	return &AthleteResultsService{
		RaceConfigurator: r,
		AthleteManager:   a,
		ResultsManager:   res,
	}
}

func (prs *AthleteResultsService) RecalculateAthleteResult(ctx context.Context, raceID uuid.UUID) error {
	return nil
}

func (prs *AthleteResultsService) GetResults(ctx context.Context, raceID uuid.UUID) (map[EventID][]*entity.AthleteSplit, error) {
	IDs, err := prs.RaceConfigurator.GetEventIDsWithWaveStarted(ctx, raceID)
	if err != nil {
		return nil, err
	}
	if len(IDs) == 0 {
		return nil, nil
	}
	// FIXME
	m := make(map[EventID][]*entity.AthleteSplit)
	for _, eventID := range IDs {
		eventResults, err := prs.ResultsManager.GetResultsForEvent(ctx, raceID, eventID)
		if err != nil {
			return nil, err
		}
		for _, e := range eventResults {
			if e == nil {
				fmt.Println("WE HAVE NIL HERE")
			}
		}
		res, err := prs.ResultsManager.CalculateRanks(ctx, eventResults)
		if err != nil {
			return nil, err
		}
		m[eventID] = res
		// fmt.Println("Res: ", res)
	}
	return m, nil
}
