package config

import (
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

// Config holds application configuration from environment variables.
type Config struct {
	Token       string `env:"TOKEN,required"`
	Port        string `env:"PORT" envDefault:"8080"`
	Healthcheck bool   `env:"HEALTHCHECK" envDefault:"false"`
}

// Load reads the environment and returns the parsed config.
func Load() Config {
	_ = godotenv.Load()

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("config: %v", err)
	}
	return cfg
}
