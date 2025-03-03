package httpv1

import (
	"encoding/json"
	"fmt"
	"maps"
	"net/http"

	"github.com/ecoarchie/timeit/internal/service"
	"github.com/ecoarchie/timeit/pkg/logger"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func NewRaceRouter(handler *chi.Mux, l logger.Interface, raceService service.RaceConfigurator) {
	handler.NotFound(notFoundResponse)
	handler.MethodNotAllowed(methodNotAllowedResponse)

	handler.Use(cors.AllowAll().Handler)
	handler.Use(middleware.Heartbeat("/ping"))
	handler.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello from index page"))
	})

	// Routers
	handler.Mount("/races", newRaceRoutes(l, raceService))
}

func NewAthleteResultsRouter(handler *chi.Mux, l logger.Interface, manager service.AthleteResultsManager) {
	handler.Mount("/races/{race_id}/athletes", newAthletesResultsRoutes(l, manager))
}

func writeJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}
	js = append(js, '\n')

	maps.Copy(w.Header(), headers)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func errorResponse(w http.ResponseWriter, status int, message string) {
	mes := map[string]string{"error": message}
	err := writeJSON(w, status, mes, nil)
	if err != nil {
		w.WriteHeader(500)
	}
}

func serverErrorResponse(w http.ResponseWriter, err error) {
	message := fmt.Errorf("server error: %w", err)
	errorResponse(w, http.StatusInternalServerError, message.Error())
}

func notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	errorResponse(w, http.StatusNotFound, message)
}

func methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	errorResponse(w, http.StatusMethodNotAllowed, message)
}
