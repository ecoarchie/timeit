package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ecoarchie/timeit/config"
	"github.com/ecoarchie/timeit/internal/controller/httpv1"
	"github.com/ecoarchie/timeit/internal/database"
	"github.com/ecoarchie/timeit/internal/repo"
	"github.com/ecoarchie/timeit/internal/service"
	"github.com/ecoarchie/timeit/pkg/httpserver"
	"github.com/ecoarchie/timeit/pkg/logger"
	"github.com/ecoarchie/timeit/pkg/postgres"
	"github.com/go-chi/chi/v5"
)

func Run(cfg *config.Config) {
	logger := logger.New(cfg.Log.Level)

	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		logger.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err).Error())
	}
	defer pg.Close()

	queries := database.New(pg.Pool)

	// Race Cache
	raceCache := service.NewRaceCache()

	// Services
	logger.Info("Creating services")
	raceService := service.NewRaceService(logger, raceCache, repo.NewRaceRepoPG(queries, pg))

	athleteRepo := repo.NewAthleteRepoPG(queries, pg)
	athleteService := service.NewAthleteService(logger, athleteRepo, raceCache)
	resultsService := service.NewResultsService(athleteRepo)

	pRes := service.NewAthleteResultsService(athleteService, resultsService)

	// Routers
	logger.Info("Creating routers")
	router := chi.NewRouter()
	httpv1.NewRaceRouter(router, logger, raceService)
	httpv1.NewAthleteResultsRouter(router, logger, pRes)
	httpServer := httpserver.New(router, httpserver.Port(cfg.HTTP.Port))

	logger.Info("Starting server at", cfg.HTTP.Port)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logger.Info("app - Run - signal: " + s.String())
	case err := <-httpServer.Notify():
		logger.Error(fmt.Sprintf("app - Run - httpServer.Notify: %s", err.Error()))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		logger.Error(fmt.Sprintf("app - Run - httpServer.Shutdown: %s", err.Error()))
	}
}
