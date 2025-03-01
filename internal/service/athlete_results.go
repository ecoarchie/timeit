package service

type AthleteResultsManager interface {
	AthleteManager
	ResultsManager
	RecalculateAthleteResult() error
}

type AthleteResultsService struct {
	AthleteManager
	ResultsManager
	Cache *RaceCache
}

func NewAthleteResultsService(a AthleteManager, res ResultsManager) *AthleteResultsService {
	return &AthleteResultsService{
		AthleteManager: a,
		ResultsManager: res,
	}
}

func (prs *AthleteResultsService) RecalculateAthleteResult() error {
	return nil
}
