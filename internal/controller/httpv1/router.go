package httpv1

import (
	"net/http"

	"github.com/ecoarchie/timeit/internal/service"
	"github.com/ecoarchie/timeit/pkg/logger"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func NewRaceRouter(handler *chi.Mux, l logger.Interface, raceService service.RaceConfigurator) {
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
