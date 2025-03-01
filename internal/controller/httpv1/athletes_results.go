package httpv1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ecoarchie/timeit/internal/entity"
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
	r.Post("/", rr.createSingleAthlete)
	r.Post("/csvheaders", rr.checkHeadersCSV)
	r.Post("/csv/{file_token}", rr.createBulkFromCSV)
	r.Delete("/{athlete_id}", rr.deleteAthleteByID)
	return r
}

func (p athletesResultsRoutes) athleteByID(w http.ResponseWriter, r *http.Request) {
	athleteID := chi.URLParam(r, "athlete_id")
	pUUID, _ := uuid.Parse(athleteID)
	a := p.service.GetAthleteByID(r.Context(), pUUID)
	if a == nil {
		p.l.Error("athlete not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(a)
}

func (p athletesResultsRoutes) createSingleAthlete(w http.ResponseWriter, r *http.Request) {
	var req entity.AthleteCreateRequest
	json.NewDecoder(r.Body).Decode(&req)
	fmt.Println("req HERE", req)
	// req.RaceID = rUUID
	a, err := p.service.CreateAthlete(r.Context(), req)
	if err != nil {
		p.l.Error("error creating athlete", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(a)
}

func (p athletesResultsRoutes) deleteAthleteByID(w http.ResponseWriter, r *http.Request) {
	athleteID := chi.URLParam(r, "athlete_id")
	aUUID, _ := uuid.Parse(athleteID)
	err := p.service.DeleteAthlete(r.Context(), aUUID)
	if err != nil {
		p.l.Error("error deleting athlete", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (p athletesResultsRoutes) checkHeadersCSV(w http.ResponseWriter, r *http.Request) {
	token, err := service.StoreTmpFile(r)
	if err != nil {
		http.Error(w, fmt.Errorf("error saving file: %w", err).Error(), http.StatusBadRequest)
		return
	}
	par := service.NewAthleteCSVParser(token+".csv", ";")
	userHeaders, matchingHeaders, err := par.CompareHeaders()
	if err != nil {
		http.Error(w, "error comparing headers in csv", http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string]any{"file_token": token, "user_headers": userHeaders, "matching_headers": matchingHeaders})
}

func (p athletesResultsRoutes) createBulkFromCSV(w http.ResponseWriter, r *http.Request) {
	fileToken := chi.URLParam(r, "file_token")
	var headers struct {
		Headers []string `json:"headers"`
	}
	json.NewDecoder(r.Body).Decode(&headers)
	par := service.NewAthleteCSVParser(fileToken+".csv", ";")
	athletes, err := par.ReadCSV(headers.Headers)
	if err != nil {
		http.Error(w, fmt.Errorf("error reading csv: %w", err).Error(), http.StatusBadRequest)
		return
	}
	// FIXME
	fmt.Printf("Athletes: %+v", athletes)
}
