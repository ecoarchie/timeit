package httpv1

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ecoarchie/timeit/internal/entity"
	"github.com/ecoarchie/timeit/internal/service"
	"github.com/ecoarchie/timeit/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type athletesResultsRoutes struct {
	service service.AthleteResultsManager
	logger  *logger.Logger
}

func newAthletesResultsRoutes(logger *logger.Logger, service service.AthleteResultsManager) http.Handler {
	logger.Info("creating new race routes")
	rr := &athletesResultsRoutes{
		service: service,
		logger:  logger,
	}
	r := chi.NewRouter()
	r.Get("/{athlete_id}", rr.athleteByID)
	r.Post("/", rr.createSingleAthlete)
	r.Post("/csvheaders", rr.checkHeadersCSV)
	r.Post("/csv/{file_token}", rr.createBulkFromCSV)
	r.Delete("/{athlete_id}", rr.deleteAthleteByID)
	r.Delete("/", rr.deleteAthletesForRace)
	r.Get("/results", rr.getResults)
	return r
}

func (p athletesResultsRoutes) getResults(w http.ResponseWriter, r *http.Request) {
	rID := chi.URLParam(r, "race_id")
	raceID, _ := uuid.Parse(rID)
	err := p.service.RecalculateAthleteResult(context.Background(), raceID)
	if err != nil {
		fmt.Println(err)
	}
}

func (p athletesResultsRoutes) athleteByID(w http.ResponseWriter, r *http.Request) {
	athleteID := chi.URLParam(r, "athlete_id")
	pUUID, _ := uuid.Parse(athleteID)
	a := p.service.GetAthleteByID(r.Context(), pUUID)
	if a == nil {
		notFoundResponse(w, r)
		return
	}
	err := writeJSON(w, http.StatusOK, a, nil)
	if err != nil {
		serverErrorResponse(w, err)
	}
}

func (p athletesResultsRoutes) createSingleAthlete(w http.ResponseWriter, r *http.Request) {
	var req entity.AthleteCreateRequest
	err := readJSON(w, r, &req)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	a, err := p.service.CreateAthlete(r.Context(), req)
	if err != nil {
		mes := "error creating athlete"
		p.logger.Error(mes, err)
		errorResponse(w, http.StatusBadRequest, mes)
		return
	}
	err = writeJSON(w, http.StatusOK, a, nil)
	if err != nil {
		http.Error(w, "server error", http.StatusBadRequest)
		return
	}
}

func (p athletesResultsRoutes) deleteAthleteByID(w http.ResponseWriter, r *http.Request) {
	athleteID := chi.URLParam(r, "athlete_id")
	aUUID, _ := uuid.Parse(athleteID)
	err := p.service.DeleteAthlete(r.Context(), aUUID)
	if err != nil {
		p.logger.Error("error deleting athlete", err)
		serverErrorResponse(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (p athletesResultsRoutes) deleteAthletesForRace(w http.ResponseWriter, r *http.Request) {
	rID := chi.URLParam(r, "race_id")
	raceID, _ := uuid.Parse(rID)
	eID := r.URL.Query().Get("event_id")
	eventID, _ := uuid.Parse(eID)
	err := p.service.DeleteAthletesForRace(r.Context(), raceID, eventID)
	if err != nil {
		serverErrorResponse(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (p athletesResultsRoutes) checkHeadersCSV(w http.ResponseWriter, r *http.Request) {
	token, err := service.StoreTmpFile(r)
	if err != nil {
		serverErrorResponse(w, err)
		return
	}
	par := service.NewAthleteImporterCSV(token, ";")
	userHeaders, matchingHeaders, err := par.CompareHeaders()
	if err != nil {
		serverErrorResponse(w, err)
		return
	}
	resp := map[string]any{
		"file_token":       token,
		"user_headers":     userHeaders,
		"matching_headers": matchingHeaders,
	}
	err = writeJSON(w, http.StatusOK, resp, nil)
	if err != nil {
		serverErrorResponse(w, err)
	}
}

func (p athletesResultsRoutes) createBulkFromCSV(w http.ResponseWriter, r *http.Request) {
	rID := chi.URLParam(r, "race_id")
	raceID, _ := uuid.Parse(rID)
	fileToken := chi.URLParam(r, "file_token")
	var headers struct {
		Headers []string `json:"headers"`
	}
	err := readJSON(w, r, &headers)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	par := service.NewAthleteImporterCSV(fileToken, ";")
	athletes, err := par.ReadCSV(headers.Headers)
	if err != nil {
		serverErrorResponse(w, err)
		return
	}
	// FIXME
	athletReqs := p.service.FromCSVtoRequestAthlete(raceID, athletes)
	for _, a := range athletReqs {
		_, err := p.service.CreateAthlete(r.Context(), a)
		if err != nil {
			p.logger.Error("error create athlete with bib from csv: ", strconv.Itoa(a.Bib), err)
		}
	}
}
