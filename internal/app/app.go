package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ecoarchie/timeit/config"
	"github.com/ecoarchie/timeit/internal/database"
	"github.com/ecoarchie/timeit/internal/service"
	"github.com/ecoarchie/timeit/internal/service/repo"
	"github.com/ecoarchie/timeit/pkg/httpserver"
	"github.com/ecoarchie/timeit/pkg/logger"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Run(cfg *config.Config) {
	logger := logger.New(cfg.Log.Level)

	// Repository
	pool, err := pgxpool.New(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal("Cannot connect to database")
	}
	defer pool.Close()

	db := database.New(pool)

	// Services

	raceService := service.NewRaceService(repo.NewRaceRepoPG(db))
	// TODO inject service into router

	router := chi.NewRouter()
	router.Use(middleware.Heartbeat("/ping"))
	// router.Get("/test", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("Hello World!"))
	// })
	httpServer := httpserver.New(router, httpserver.Port(cfg.HTTP.Port))

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
	err := httpServer.Shutdown()
	if err != nil {
		logger.Error(fmt.Sprintf("app - Run - httpServer.Shutdown: %s", err.Error()))
	}
}
