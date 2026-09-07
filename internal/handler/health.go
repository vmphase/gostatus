package handler

import (
	"net/http"

	"gostatus/internal/gateway"
)

// Healthz reports liveness of the server and gateway connection.
func (h *Handler) Healthz(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")

	if !gateway.Connected() {
		w.WriteHeader(http.StatusServiceUnavailable)
		writeJSON(w, map[string]string{"status": "gateway disconnected"})
		return
	}
	writeJSON(w, map[string]string{"status": "ok"})
}
