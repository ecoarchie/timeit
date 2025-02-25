package service

type AthleteResultsManager interface {
	AthleteManager
	ResultsManager
	RecalculateAthleteResult() error
}

type AthleteResultsService struct {
	AthleteManager
	ResultsManager
}

func NewAthleteResultsService(p AthleteManager, res ResultsManager) *AthleteResultsService {
	return &AthleteResultsService{
		p,
		res,
	}
}

func (prs *AthleteResultsService) RecalculateAthleteResult() error {
	return nil
}
