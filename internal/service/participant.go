package service

import "github.com/ecoarchie/timeit/pkg/logger"

type ParticipantEditor interface{}

type ParticipantRepo interface{}

type ParticipantService struct {
	l    logger.Interface
	repo ParticipantRepo
}

func NewParticipantService(logger logger.Interface, repo ParticipantRepo) *ParticipantService {
	return &ParticipantService{
		l:    logger,
		repo: repo,
	}
}
