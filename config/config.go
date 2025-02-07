package config

import "github.com/ilyakaznacheev/cleanenv"

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
		Name    string `env-required:"true"     env:"APP_NAME"`
		Version string `env-required:"true"  env:"APP_VERSION"`
	}

	// HTTP -.
	HTTP struct {
		Port string `env-required:"true"  env:"HTTP_PORT"`
	}

	// Log -.
	Log struct {
		Level string `env-required:"true"    env:"LOG_LEVEL"`
	}

	// PG -.
	PG struct {
		PoolMax int    `env-required:"true"  env:"PG_POOL_MAX"`
		URL     string `env-required:"true"                 env:"PG_URL"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
