package config

import (
	"os"
	"log"
	"errors"

	"github.com/joho/godotenv"
	"github.com/caarlos0/env/v11"
)

type Config struct{
	Port   string `env:"PORT" envDefault:"8080"`
	DBUrl  string `env:"DB_URL,required,notEmpty"`
	IsProd bool   `env:"IS_PROD" envDefault:false`
}

func Load() (*Config) {
	err := godotenv.Load()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Println(".env file not found, reading from system environment variables")
		} else {
			log.Fatalf("Failed to load .env: %v", err)
		}
	}
	var cfg Config
	err = env.Parse(&cfg)
	if err != nil {
		log.Fatalf("Failed to load config: %v")
	}

	if cfg.Port != "" && cfg.Port[0] != ':' {
		cfg.Port = ":" + cfg.Port
	}

	return &cfg
}
