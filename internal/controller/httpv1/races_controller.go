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
	"github.com/google/uuid"
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
	rID := chi.URLParam(r, "race_id")
	var waveStart entity.WaveStart
	err := readJSON(w, r, &waveStart)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	v := validator.New()
	v.Check(validator.IsUUID(rID), "race_id", "must be provided and be valid uuid")
	v.Check(validator.IsUUID(waveStart.WaveID.String()), "wave_id", "must be provided and be valid uuid")
	if !v.Valid() {
		failedValidationResponse(w, v.Errors)
		return
	}
	startTime, waveFound, err := rr.conf.StartWave(context.Background(), uuid.MustParse(rID), waveStart)
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
	rID := chi.URLParam(r, "race_id")

	v := validator.New()
	v.Check(rID != "", "race_id", "must be provided ")
	v.Check(validator.IsUUID(rID), "race_id", "must be valid UUID")
	if !v.Valid() {
		failedValidationResponse(w, v.Errors)
		return
	}
	waves, err := rr.conf.GetWavesForRace(context.Background(), uuid.MustParse(rID))
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
	v := *validator.New()
	v.Check(validator.IsUUID(id), "race_id", "must be valid uuid")
	if !v.Valid() {
		failedValidationResponse(w, v.Errors)
		return
	}
	err := rr.conf.DeleteRace(context.Background(), uuid.MustParse(id))
	if err != nil {
		serverErrorResponse(w, err)
		return
	}
	writeJSON(w, http.StatusNoContent, nil, nil)
}

func (rr *raceRoutes) getRaces(w http.ResponseWriter, r *http.Request) {
	races, err := rr.conf.GetRaces(context.Background())
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
	v := *validator.New()
	v.Check(validator.IsUUID(id), "race_id", "must be valid uuid")
	if !v.Valid() {
		failedValidationResponse(w, v.Errors)
		return
	}
	cfg, err := rr.conf.GetRaceConfig(context.Background(), uuid.MustParse(id))
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
	race, err := rr.conf.CreateRace(context.Background(), dto, v)
	if err != nil {
		rr.log.Error("error creating race", err)
		serverErrorResponse(w, err)
		return
	}
	if !v.Valid() {
		rr.log.Error("error validating race config to save")
		failedValidationResponse(w, v.Errors)
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
	ctx := context.Background()

	v := validator.New()
	raceConfig.Validate(ctx, v)

	if !v.Valid() {
		failedValidationResponse(w, v.Errors)
		return
	}
	err = rr.conf.SaveRaceConfig(ctx, raceConfig, v)
	if err != nil {
		mes := "error saving race config"
		rr.log.Error(mes, err)
		serverErrorResponse(w, err)
		return
	}
	rr.log.Info("Config for race saved")
	writeJSON(w, http.StatusNoContent, "", nil)
}
