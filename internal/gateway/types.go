package gateway

import "encoding/json"

// Discord activity types.
const (
	ActivityTypePlaying = iota
	ActivityTypeStreaming
	ActivityTypeListening
	ActivityTypeWatching
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
	User         *UserMin   `json:"user,omitempty"`
	Status       string     `json:"status"`
	ClientStatus any        `json:"client_status"`
	Activities   []Activity `json:"activities"`
}

// UserMin is the minimal user payload sent with presence updates.
type UserMin struct {
	ID string `json:"id"`
}

// Activity is a Discord activity as reported by the gateway.
type Activity struct {
	Name    string `json:"name"`
	Type    int    `json:"type"`
	Details string `json:"details"`
	State   string `json:"state"`
	SyncID  string `json:"sync_id"`
}

type guildCreate struct {
	Presences []PresenceUpdate `json:"presences"`
}
