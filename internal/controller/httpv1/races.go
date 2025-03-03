package httpv1

import (
	"encoding/json"
	"net/http"

	"github.com/ecoarchie/timeit/internal/entity"
	"github.com/ecoarchie/timeit/internal/service"
	"github.com/ecoarchie/timeit/pkg/logger"
	"github.com/go-chi/chi/v5"
)

type raceRoutes struct {
	rcs service.RaceConfigurator
	l   logger.Interface
}

func newRaceRoutes(l logger.Interface, service service.RaceConfigurator) http.Handler {
	l.Info("creating new race routes")
	rr := &raceRoutes{
		rcs: service,
		l:   l,
	}
	r := chi.NewRouter()
	r.Post("/", rr.createRace)
	r.Get("/{race_id}", rr.getRace)
	r.Post("/save", rr.saveRaceConfig)
	return r
}

func (rr *raceRoutes) getRace(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "race_id")
	rr.rcs.GetRaceConfig(r.Context(), id)
}

func (rr *raceRoutes) createRace(w http.ResponseWriter, r *http.Request) {
	var req *entity.RaceFormData
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		mes := "error parsing new race form"
		rr.l.Error(mes, err)
		errorResponse(w, http.StatusBadRequest, mes)
		return
	}
	race, err := rr.rcs.CreateRace(r.Context(), req)
	if err != nil {
		rr.l.Error("error creating race", err)
		serverErrorResponse(w, err)
		return
	}
	writeJSON(w, http.StatusOK, race, nil)
}

func (rr *raceRoutes) saveRaceConfig(w http.ResponseWriter, r *http.Request) {
	var conf *entity.RaceConfig
	err := json.NewDecoder(r.Body).Decode(&conf)
	if err != nil {
		mes := "error parsing race config form data"
		rr.l.Error(mes, err)
		errorResponse(w, http.StatusBadRequest, mes)
		return
	}

	// TODO add check for db save error for proper error reponse code
	errs := rr.rcs.Save(r.Context(), conf)
	if len(errs) != 0 {
		mes := "error saving race config"
		rr.l.Error(mes, errs)
		errorResponse(w, http.StatusBadRequest, mes)
		return
	}
	rr.l.Info("Config for race saved")
	w.Write([]byte("ok"))
}

// func (rr *raceRoutes) getResultsForAthlete(w http.ResponseWriter, r *http.Request) {
// 	rr.results.ResultsForAthlete()
// }
