package handler

import "net/http"

// GET /healthz
func HandleGetHealthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
