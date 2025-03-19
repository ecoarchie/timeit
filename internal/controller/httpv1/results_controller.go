package httpv1

import (
	"context"
	"net/http"

	"github.com/ecoarchie/timeit/internal/service"
	"github.com/ecoarchie/timeit/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type resultsRoutes struct {
	service service.ResultsManager
	logger  *logger.Logger
}

func newResultsRoutes(logger *logger.Logger, service service.ResultsManager) http.Handler {
	logger.Info("creating new race routes")
	rr := &resultsRoutes{
		service: service,
		logger:  logger,
	}
	r := chi.NewRouter()
	r.Get("/results", rr.getResults)
	return r
}

func (p resultsRoutes) getResults(w http.ResponseWriter, r *http.Request) {
	rID := chi.URLParam(r, "race_id")
	raceID, _ := uuid.Parse(rID)
	res, err := p.service.GetResults(context.Background(), raceID)
	if err != nil {
		p.logger.Error("Calculate results: ", "err", err)
		serverErrorResponse(w, err)
		return
	}
	err = writeJSON(w, http.StatusOK, res, nil)
	if err != nil {
		serverErrorResponse(w, err)
	}
}
