package restapp

import (
	"auth-service/internal/config"
	"log"
	"net/http"
)

type RestServer struct {
	config *config.Config
}

func New(c *config.Config) *RestServer {
	log.Print("Rest server started")
	return &RestServer{config: c}
}

func (s *RestServer) Routes() http.Handler {
	// Start the server
	mux := http.NewServeMux()

	mux.HandleFunc("POST /login/", s.Login)
	mux.HandleFunc("POST /register", s.Register)
	mux.HandleFunc("POST /refresh", s.Refresh)

	return mux
}
