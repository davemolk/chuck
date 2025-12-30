package handlers

import "net/http"

type HealthHandlers struct{}

func NewHealthHandlers() *HealthHandlers {
	return &HealthHandlers{}
}

func (h *HealthHandlers) HealthCheck(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("i'm healthy!"))
}
