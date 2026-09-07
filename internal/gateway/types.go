package gateway

import (
	"encoding/json"

	"gostatus/internal/store"
)

// Payload is a gateway event frame.
type Payload struct {
	Op int             `json:"op"`
	D  json.RawMessage `json:"d"`
	S  *int            `json:"s"`
	T  *string         `json:"t"`
}

// PresenceUpdate is a PRESENCE_UPDATE gateway event.
type PresenceUpdate struct {
	User         *UserMin          `json:"user,omitempty"`
	Status       string            `json:"status"`
	ClientStatus map[string]string `json:"client_status"`
	Activities   []store.Activity  `json:"activities"`
}

// UserMin is the minimal user payload sent with presence updates.
type UserMin struct {
	ID string `json:"id"`
}

type guildCreate struct {
	Presences []PresenceUpdate `json:"presences"`
}
