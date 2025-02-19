package service

type ParticipantResultsManager interface {
	RecalculateParticipantResult() error
}

type ParticipantResultsService struct {
	p   ParticipantManager
	res ResultsManager
}

func NewParticipantResultsService(p ParticipantManager, res ResultsManager) *ParticipantResultsService {
	return &ParticipantResultsService{
		p:   p,
		res: res,
	}
}

func (prs *ParticipantResultsService) RecalculateParticipantResult() error {
	return nil
}
