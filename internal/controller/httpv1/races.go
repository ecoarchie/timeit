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
	conf service.RaceConfigurator
	log  logger.Interface
}

func newRaceRoutes(log logger.Interface, conf service.RaceConfigurator) http.Handler {
	log.Info("creating new race routes")
	rr := &raceRoutes{
		conf: conf,
		log:  log,
	}
	r := chi.NewRouter()
	r.Get("/", rr.getRaces)
	r.Post("/", rr.createRace)
	r.Get("/{race_id}", rr.getRaceConfig)
	r.Post("/save", rr.saveRaceConfig)
	return r
}

func (rr *raceRoutes) getRaces(w http.ResponseWriter, r *http.Request) {
	races, err := rr.conf.GetRaces(r.Context())
	if err != nil {
		serverErrorResponse(w, err)
		return
	}
	if races == nil {
		notFoundResponse(w, r)
		return
	}
	writeJSON(w, http.StatusOK, races, nil)
}

func (rr *raceRoutes) getRaceConfig(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "race_id")
	cfg, err := rr.conf.GetRaceConfig(r.Context(), id)
	if err != nil {
		serverErrorResponse(w, err)
		return
	}
	if cfg == nil {
		errorResponse(w, http.StatusNotFound, "race not found")
		return
	}
	writeJSON(w, http.StatusOK, cfg, nil)
}

func (rr *raceRoutes) createRace(w http.ResponseWriter, r *http.Request) {
	var req *entity.RaceFormData
	err := readJSON(w, r, &req)
	if err != nil {
		mes := "error parsing new race form"
		rr.log.Error(mes, err)
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	race, err := rr.conf.CreateRace(r.Context(), req)
	if err != nil {
		rr.log.Error("error creating race", err)
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
		rr.log.Error(mes, err)
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	v := validator.New()
	conf.Validate(r.Context(), v)

	if !v.Valid() {
		failedValidationResponse(w, v.Errors)
		return
	}
	err = rr.conf.SaveRaceConfig(r.Context(), conf)
	if err != nil {
		mes := "error saving race config"
		rr.log.Error(mes, err)
		serverErrorResponse(w, err)
		return
	}
	rr.log.Info("Config for race saved")
	writeJSON(w, http.StatusNoContent, "", nil)
}
