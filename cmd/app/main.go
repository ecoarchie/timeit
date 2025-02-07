package main

import (
	"log"

	"github.com/ecoarchie/timeit/config"
	"github.com/ecoarchie/timeit/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
