package gateway

import "encoding/json"

const (
	ActivityTypePlaying   = iota
	ActivityTypeStreaming
	ActivityTypeListening
	ActivityTypeWatching
)

type Payload struct {
	Op int             `json:"op"`
	D  json.RawMessage `json:"d"`
	S  *int            `json:"s"`
	T  *string         `json:"t"`
}

type PresenceUpdate struct {
	User         *UserMin   `json:"user,omitempty"`
	Status       string     `json:"status"`
	ClientStatus any        `json:"client_status"`
	Activities   []Activity `json:"activities"`
}

type UserMin struct {
	ID string `json:"id"`
}

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
