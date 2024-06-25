package restapp

import (
	"net/http"
	"strings"
)

func (s *RestServer) Login(w http.ResponseWriter, r *http.Request) {
	// Login logic here
	if strings.HasPrefix("/login/", r.URL.Path) {
		http.Error(w, "Invalid login option", http.StatusBadRequest)
		return
	}
	option := r.URL.Path[len("/login/"):]

	switch option {
	case "credentials":
		s.credentials(w, r)
	case "google":
		s.google(w, r)
	default:
		http.Error(w, "Invalid login option", http.StatusBadRequest)
	}
}

func (s *RestServer) google(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Google login"))
	return
}

func (s *RestServer) credentials(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Credentials login"))
	return
}
