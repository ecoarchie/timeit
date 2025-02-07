package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ecoarchie/timeit/config"
	"github.com/ecoarchie/timeit/pkg/httpserver"
	"github.com/ecoarchie/timeit/pkg/logger"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func Run(cfg *config.Config) {
	logger := logger.New(cfg.Log.Level)

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
