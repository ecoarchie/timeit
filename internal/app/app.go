package app

import (
	"context"
	"fmt"
	"log"
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
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Run(cfg *config.Config) {
	logger := logger.New(cfg.Log.Level)

	// Postgres pool
	logger.Info("Initializint postgres pool")
	pool, err := pgxpool.New(context.Background(), cfg.PG.URL)
	if err != nil {
		log.Fatal("Cannot connect to database")
		os.Exit(1)
	}
	defer pool.Close()

	db := database.New(pool)

	// Services
	logger.Info("Creating services")
	raceService := service.NewRaceService(logger, repo.NewRaceRepoPG(db))

	// Routers
	logger.Info("Creating routers")
	router := chi.NewRouter()
	httpv1.NewRouter(router, logger, raceService)
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
