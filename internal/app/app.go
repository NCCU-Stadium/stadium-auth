package app

import (
	"auth-service/internal/config"
	"auth-service/internal/grpcapp"
	"log"
	"net"
)

func StartServer() {
	c := config.NewConfig() // Initialize the config

	grpca := grpcapp.New(c)
	lis, err := net.Listen("tcp", ":8080")

	log.Printf("Server started successfully!\n\tRunning gRPC on port 8080")

	if err = grpca.Server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
