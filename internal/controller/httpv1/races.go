package httpv1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ecoarchie/timeit/internal/service"
	"github.com/ecoarchie/timeit/pkg/logger"
	"github.com/go-chi/chi/v5"
)

type raceRoutes struct {
	s service.RaceConfigurator
	l logger.Interface
}

func newRaceRoutes(l logger.Interface, service service.RaceConfigurator) http.Handler {
	l.Info("creating new race routes")
	rr := &raceRoutes{
		s: service,
		l: l,
	}
	r := chi.NewRouter()
	r.Post("/", rr.create)
	return r
}

func (rr *raceRoutes) create(w http.ResponseWriter, r *http.Request) {
	rr.l.Info("create new race")
	var req service.CreateRaceRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		rr.l.Error("error parsing request for race creation", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rr.l.Info("create race req: ", req)
	id, err := rr.s.Create(r.Context(), req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte(fmt.Sprintf("created race with id: %s", id)))
}
