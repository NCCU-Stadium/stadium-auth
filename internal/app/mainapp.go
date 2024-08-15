package app

import (
	"auth-service/internal/config"
	"auth-service/internal/helper"
	"log"
	"net/http"
)

type Server struct {
	config        *config.Config
	refreshHelper *helper.RefreshHelper
	helper        *helper.Helper
}

func New(c *config.Config) *Server {
	log.Print("Rest server started")
	refreshHelper := helper.NewRefreshDB(c)
	helper := helper.NewRestHelper(c)

	return &Server{config: c, refreshHelper: refreshHelper, helper: helper}
}

func (s Server) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		w.Header().Set("Access-Control-Allow-Origin", origin)
		// Allow CORS here By http://localhost:3000
		// w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func use(r *http.ServeMux, middlewares ...func(next http.Handler) http.Handler) http.Handler {
	var s http.Handler
	s = r

	for _, middleware := range middlewares {
		s = middleware(s)
	}

	return s
}

func hello(w http.ResponseWriter, r *http.Request) {
	// Get cookie test
	cookie, err := r.Cookie("accessToken")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte("accessToken value: " + cookie.Value))
	return
	// w.Write([]byte("HELLO from auth-service"))
	// return
}

func (s *Server) Routes() http.Handler {
	// Start the server
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", hello)
	mux.HandleFunc("POST /login", s.Login)
	mux.HandleFunc("POST /register", s.Register)
	mux.HandleFunc("POST /refresh", s.Refresh)

	wrapped := use(mux, s.middleware)

	return wrapped
}
