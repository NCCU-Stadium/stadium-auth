package config

import (
	"auth-service/dotenv"
	"log"
	"os"
)

type Config struct {
	DatabaseURI         string
	Secret              string
	RefreshDBURI        string
	RefreshDBName       string
	RefreshDBCollection string
}

func NewConfig() *Config {
	err := dotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	return &Config{
		DatabaseURI:         os.Getenv("DATABASE_URI"),
		Secret:              os.Getenv("SECRET"),
		RefreshDBURI:        os.Getenv("REFRESH_DB_URI"),
		RefreshDBName:       os.Getenv("REFRESH_DB_NAME"),
		RefreshDBCollection: os.Getenv("REFRESH_DB_COLLECTION"),
	}
}
