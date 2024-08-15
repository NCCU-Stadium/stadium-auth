package app

import (
	"auth-service/internal/config"
	"log"
	"net/http"
)

func StartServer() {
	c := config.NewConfig() // Initialize the config

	// Start REST server
	a := New(c)
	srv := &http.Server{
		Addr:    ":8085",
		Handler: a.Routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
