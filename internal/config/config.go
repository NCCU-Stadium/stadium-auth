package config

import (
	"auth-service/dotenv"
	"log"
	"os"
)

type Config struct {
	DatabaseURI string
	Secret      string
}

func NewConfig() *Config {
	err := dotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	return &Config{
		DatabaseURI: os.Getenv("DATABASE_URI"),
		Secret:      os.Getenv("SECRET"),
	}
}
