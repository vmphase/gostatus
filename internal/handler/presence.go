package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// Presence serves the full cached Discord presence payload for a user,
// falling back to an offline/empty response if not found.
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
		writeJSON(w, map[string]string{"error": "missing user id"})
		return
	}

	if p, ok := h.store.Get(id); ok {
		writeJSON(w, p)
	} else {
		writeJSON(w, map[string]any{
			"Status":       "offline",
			"ClientStatus": map[string]any{},
			"Activities":   []any{},
		})
	}
}

func writeJSON(w http.ResponseWriter, v any) {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("encode JSON response: %v", err)
	}
}
