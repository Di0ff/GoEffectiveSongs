package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBDsn       string
	MusicAPIURL string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found, using environment variables")
	}

	config := &Config{
		DBDsn:       os.Getenv("DB_DSN"),
		MusicAPIURL: os.Getenv("MUSIC_API_URL"),
	}

	if config.DBDsn == "" || config.MusicAPIURL == "" {
		log.Fatal("Required environment variables are missing")
	}

	return config, nil
}
