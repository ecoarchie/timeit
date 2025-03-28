package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type (
	// Config -.
	Config struct {
		App
		HTTP
		Log
		PG
	}

	// App -.
	App struct {
		Name    string `env-required:"true" env:"APP_NAME"`
		Version string `env-required:"true" env:"APP_VERSION"`
	}

	// HTTP -.
	HTTP struct {
		Port string `env-required:"true" env:"HTTP_PORT"`
	}

	// Log -.
	Log struct {
		Level string `env-required:"true" env:"LOG_LEVEL"`
	}

	// PG -.
	PG struct {
		URL     string `env-required:"true" env:"PG_URL"`
		PoolMax int    `env:"PG_POOL_MAX"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		fmt.Println("error reading env")
		return nil, err
	}

	return cfg, nil
}
