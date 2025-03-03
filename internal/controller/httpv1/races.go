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
		rr.l.Error("error parsing new race form", err)
		serverErrorResponse(w, err)
		return
	}
	race, err := rr.rcs.CreateRace(r.Context(), req)
	if err != nil {
		rr.l.Error("error creating race", err)
		serverErrorResponse(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(race)
	// w.WriteHeader(http.StatusOK)
}

func (rr *raceRoutes) saveRaceConfig(w http.ResponseWriter, r *http.Request) {
	var conf *entity.RaceConfig
	err := json.NewDecoder(r.Body).Decode(&conf)
	if err != nil {
		rr.l.Error("error parsing race config form data", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO add check for db save error for proper error reponse code
	errs := rr.rcs.Save(r.Context(), conf)
	if len(errs) != 0 {
		rr.l.Error("error saving race config", errs)
		var resp []byte
		for _, e := range errs {
			resp = append(resp, []byte(e.Error())...)
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		return
	}
	rr.l.Info("Config for race saved")
	w.Write([]byte("ok"))
}

// func (rr *raceRoutes) getResultsForAthlete(w http.ResponseWriter, r *http.Request) {
// 	rr.results.ResultsForAthlete()
// }
