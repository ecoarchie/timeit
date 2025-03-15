package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type AthleteResultsManager interface {
	RaceConfigurator
	AthleteManager
	ResultsManager
	RecalculateAthleteResult(ctx context.Context, raceID uuid.UUID) error
}

type AthleteResultsService struct {
	RaceConfigurator
	AthleteManager
	ResultsManager
	Cache *RaceCache
}

func NewAthleteResultsService(r RaceConfigurator, a AthleteManager, res ResultsManager) *AthleteResultsService {
	return &AthleteResultsService{
		RaceConfigurator: r,
		AthleteManager:   a,
		ResultsManager:   res,
	}
}

func (prs *AthleteResultsService) RecalculateAthleteResult(ctx context.Context, raceID uuid.UUID) error {
	IDs, err := prs.RaceConfigurator.GetEventIDsWithWaveStarted(ctx, raceID)
	if err != nil {
		return err
	}
	if len(IDs) == 0 {
		return nil
	}
	// FIXME
	for _, eventID := range IDs {
		eventResults, err := prs.ResultsManager.GetResultsForEvent(ctx, raceID, eventID)
		if err != nil {
			return err
		}
		for _, e := range eventResults {
			if e == nil {
				fmt.Println("WE HAVE NIL HERE")
			}
		}
		_, err = prs.ResultsManager.CalculateRanks(ctx, eventResults)
		if err != nil {
			return err
		}
		// fmt.Println("Res: ", res)
	}
	return nil
}
