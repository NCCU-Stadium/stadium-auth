package restapp

import "net/http"

func (s *RestServer) Refresh(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Refresh token"))
	return
}
