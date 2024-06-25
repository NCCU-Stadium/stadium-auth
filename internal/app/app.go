package app

import (
	"auth-service/internal/config"
	"auth-service/internal/grpcapp"
	"auth-service/internal/restapp"
	"log"
	"net"
	"net/http"
)

func StartServer() {
	c := config.NewConfig() // Initialize the config

	// Start gRPC server
	go func() {
		grpca := grpcapp.New(c)
		lis, err := net.Listen("tcp", ":50050")
		if err = grpca.Server.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Start REST server
	resta := restapp.New(c)
	srv := &http.Server{
		Addr:    ":8085",
		Handler: resta.Routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
