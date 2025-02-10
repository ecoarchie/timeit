package httpv1

import (
	"net/http"

	"github.com/ecoarchie/timeit/internal/service"
	"github.com/ecoarchie/timeit/pkg/logger"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func NewRouter(handler *chi.Mux, l logger.Interface, service service.RaceConfigurator) {
	handler.Use(cors.AllowAll().Handler)
	handler.Use(middleware.Heartbeat("/ping"))
	handler.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello from index page"))
	})

	// Routers
	handler.Mount("/races", newRaceRoutes(l, service))
}
