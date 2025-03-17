package httpv1

import (
	"context"
	"net/http"
	"time"

	"github.com/ecoarchie/timeit/internal/controller/httpv1/dto"
	"github.com/ecoarchie/timeit/internal/entity"
	"github.com/ecoarchie/timeit/internal/service"
	"github.com/ecoarchie/timeit/pkg/logger"
	"github.com/ecoarchie/timeit/pkg/validator"
	"github.com/go-chi/chi/v5"
)

type raceRoutes struct {
	conf service.RaceConfigurator
	log  *logger.Logger
}

func newRaceRoutes(log *logger.Logger, conf service.RaceConfigurator) http.Handler {
	rr := &raceRoutes{
		conf: conf,
		log:  log,
	}
	r := chi.NewRouter()
	r.Get("/", rr.getRaces)
	r.Post("/", rr.createRace)
	r.Get("/{race_id}", rr.getRaceConfig)
	r.Post("/{race_id}", rr.saveRaceConfig)
	r.Delete("/{race_id}", rr.deleteRace)
	r.Get("/{race_id}/waves", rr.getWavesForRace)
	r.Post("/{race_id}/waves/start", rr.startWave)
	return r
}

func (rr *raceRoutes) startWave(w http.ResponseWriter, r *http.Request) {
	v := validator.New()
	rID := chi.URLParam(r, "race_id")
	// wID := chi.URLParam(r, "wave_id")
	var waveStart entity.WaveStart
	err := readJSON(w, r, &waveStart)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	v.Check(rID != "" && validator.IsUUID(rID), "race_id", "must be provided and be valid uuid")
	v.Check(waveStart.WaveID.String() != "" && validator.IsUUID(waveStart.WaveID.String()), "wave_id", "must be provided and be valid uuid")
	if !v.Valid() {
		errorResponse(w, http.StatusBadRequest, v.Errors)
		return
	}
	startTime, waveFound, err := rr.conf.StartWave(context.Background(), rID, waveStart)
	if err != nil {
		serverErrorResponse(w, err)
		return
	}
	if !waveFound {
		errorResponse(w, http.StatusNotFound, "wave not found")
		return
	}
	res := map[string]string{
		"start_time": startTime.Format(time.RFC3339Nano),
	}
	writeJSON(w, http.StatusOK, res, nil)
}

func (rr *raceRoutes) getWavesForRace(w http.ResponseWriter, r *http.Request) {
	v := validator.New()
	rID := chi.URLParam(r, "race_id")
	v.Check(rID != "", "race_id", "must be provided ")
	v.Check(validator.IsUUID(rID), "race_id", "must be valid UUID")
	if !v.Valid() {
		errorResponse(w, http.StatusBadRequest, v.Errors)
		return
	}
	waves, err := rr.conf.GetWavesForRace(context.Background(), rID)
	if err != nil {
		serverErrorResponse(w, err)
		return
	}
	if waves == nil {
		errorResponse(w, http.StatusNotFound, "waves for race not found")
		return
	}
	writeJSON(w, http.StatusOK, waves, nil)
}

func (rr *raceRoutes) deleteRace(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "race_id")
	err := rr.conf.DeleteRace(context.Background(), id)
	if err != nil {
		serverErrorResponse(w, err)
		return
	}
	writeJSON(w, http.StatusNoContent, nil, nil)
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
		rr.log.Debug("race not found", "raceID", id)
		errorResponse(w, http.StatusNotFound, "race not found")
		return
	}
	writeJSON(w, http.StatusOK, cfg, nil)
}

func (rr *raceRoutes) createRace(w http.ResponseWriter, r *http.Request) {
	var dto *dto.RaceDTO
	err := readJSON(w, r, &dto)
	if err != nil {
		mes := "error parsing new race form"
		rr.log.Error(mes, err)
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	v := validator.New()
	race, err := rr.conf.CreateRace(r.Context(), dto, v)
	if err != nil {
		rr.log.Error("error creating race", err)
		serverErrorResponse(w, err)
		return
	}
	writeJSON(w, http.StatusOK, race, nil)
}

func (rr *raceRoutes) saveRaceConfig(w http.ResponseWriter, r *http.Request) {
	var raceConfig *dto.RaceConfig
	err := readJSON(w, r, &raceConfig)
	if err != nil {
		mes := "error parsing race config form data"
		rr.log.Error(mes, err)
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	v := validator.New()
	raceConfig.Validate(r.Context(), v)

	if !v.Valid() {
		failedValidationResponse(w, v.Errors)
		return
	}
	err = rr.conf.SaveRaceConfig(r.Context(), raceConfig, v)
	if err != nil {
		mes := "error saving race config"
		rr.log.Error(mes, err)
		serverErrorResponse(w, err)
		return
	}
	rr.log.Info("Config for race saved")
	writeJSON(w, http.StatusNoContent, "", nil)
}
