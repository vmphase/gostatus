package handler

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (h *Handler) Presence(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/presence/")

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if p, ok := h.store.Get(id); ok {
		json.NewEncoder(w).Encode(p)
	} else {
		json.NewEncoder(w).Encode(map[string]any{
			"status":        "offline",
			"client_status": map[string]any{},
			"activities":    []any{},
		})
	}
}
