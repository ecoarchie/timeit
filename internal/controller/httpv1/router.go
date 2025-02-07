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
	// Options
	// handler.Use(cors.Handler(cors.Options{
	// 	// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
	// 	AllowedOrigins: []string{"https://*", "http://*"},
	// 	// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
	// 	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	// 	AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	// 	ExposedHeaders:   []string{"Link"},
	// 	AllowCredentials: false,
	// 	MaxAge:           300, // Maximum value not ignored by any of major browsers
	// }))
	handler.Use(cors.AllowAll().Handler)
	handler.Use(middleware.Heartbeat("/ping"))
	handler.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello from index page"))
	})

	// Routers
	handler.Mount("/races", newRaceRoutes(l, service))
}
