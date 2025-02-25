package httpv1

import (
	"fmt"
	"net/http"

	"github.com/ecoarchie/timeit/internal/service"
	"github.com/ecoarchie/timeit/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type athletesResultsRoutes struct {
	service service.AthleteResultsManager
	l       logger.Interface
}

func newAthletesResultsRoutes(l logger.Interface, service service.AthleteResultsManager) http.Handler {
	l.Info("creating new race routes")
	rr := &athletesResultsRoutes{
		service: service,
		l:       l,
	}
	r := chi.NewRouter()
	r.Get("/{athlete_id}", rr.athleteByID)
	return r
}

func (p athletesResultsRoutes) athleteByID(w http.ResponseWriter, r *http.Request) {
	raceID := chi.URLParam(r, "race_id")
	rUUID, _ := uuid.Parse(raceID)
	athleteID := chi.URLParam(r, "athlete_id")
	pUUID, _ := uuid.Parse(athleteID)
	fmt.Printf("raceID = %s, athleteID = %s", raceID, athleteID)
	p.service.GetAthleteByID(rUUID, pUUID)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("athlete"))
}
