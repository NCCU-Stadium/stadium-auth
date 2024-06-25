package restapp

import "net/http"

func (s *RestServer) Register(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Register"))
	return
}
