package restapp

import (
	"auth-service/internal/config"
	"auth-service/internal/restapp/helper"
	"log"
	"net/http"
)

type RestServer struct {
	config        *config.Config
	refreshHelper *restapp_helper.RefreshHelper
	restHelper    *restapp_helper.RestHelper
}

func New(c *config.Config) *RestServer {
	log.Print("Rest server started")
	refreshHelper := restapp_helper.NewRefreshDB(c)
	restHelper := restapp_helper.NewRestHelper(c)

	return &RestServer{config: c, refreshHelper: refreshHelper, restHelper: restHelper}
}

func (s RestServer) middleware(next http.Handler) http.Handler {
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
	w.Write([]byte("HELLO from auth-service"))
	return
}

func (s *RestServer) Routes() http.Handler {
	// Start the server
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", hello)
	mux.HandleFunc("POST /login/", s.Login)
	mux.HandleFunc("POST /register", s.Register)
	mux.HandleFunc("POST /refresh", s.Refresh)

	wrapped := use(mux, s.middleware)

	return wrapped
}
