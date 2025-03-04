package httpv1

import (
	"net/http"

	"github.com/ecoarchie/timeit/internal/entity"
	"github.com/ecoarchie/timeit/internal/service"
	"github.com/ecoarchie/timeit/pkg/logger"
	"github.com/ecoarchie/timeit/pkg/validator"
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
	r.Get("/{race_id}", rr.getRaceConfig)
	r.Post("/save", rr.saveRaceConfig)
	return r
}

func (rr *raceRoutes) getRaceConfig(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "race_id")
	rr.rcs.GetRaceConfig(r.Context(), id)
}

func (rr *raceRoutes) createRace(w http.ResponseWriter, r *http.Request) {
	var req *entity.RaceFormData
	err := readJSON(w, r, &req)
	if err != nil {
		mes := "error parsing new race form"
		rr.l.Error(mes, err)
		errorResponse(w, http.StatusBadRequest, err.Error())
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
	err := readJSON(w, r, &conf)
	if err != nil {
		mes := "error parsing race config form data"
		rr.l.Error(mes, err)
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	v := validator.New()
	rr.rcs.Validate(r.Context(), v, conf)

	if !v.Valid() {
		errorResponse(w, http.StatusBadRequest, v.Errors)
		return
	}
	err = rr.rcs.Save(r.Context(), conf)
	if err != nil {
		mes := "error saving race config"
		rr.l.Error(mes, err)
		serverErrorResponse(w, err)
		return
	}
	rr.l.Info("Config for race saved")
	writeJSON(w, http.StatusNoContent, "", nil)
}
