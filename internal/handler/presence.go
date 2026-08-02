package handler

import (
	"encoding/json"
	"net/http"
	"strings"
)

// Full cached Discord presence payload for a user
// Fallbacks to an offline/empty presence response if not found
func (h *Handler) Presence(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/presence/")

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "missing user id"})
		return
	}

	if p, ok := h.store.Get(id); ok {
		json.NewEncoder(w).Encode(p)
	} else {
		json.NewEncoder(w).Encode(map[string]any{
			"Status":       "offline",
			"ClientStatus": map[string]any{},
			"Activities":   []any{},
		})
	}
}
