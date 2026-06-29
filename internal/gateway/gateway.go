package gateway

import (
	"encoding/json"
	"gostatus/internal/store"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	gatewayURL = "wss://gateway.discord.gg/?v=10&encoding=json"

	opHeartbeat = 1
	opIdentify  = 2

	intentGuilds         = 1 << 0
	intentGuildPresences = 1 << 8

	defaultHeartbeatInterval = 41250
	maxBackoff              = 60
)

type seqHolder struct {
	mu  sync.Mutex
	seq *int
}

func (s *seqHolder) set(seq *int) {
	s.mu.Lock()
	s.seq = seq
	s.mu.Unlock()
}

func (s *seqHolder) get() *int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.seq
}

func Connect(token string, s *store.Store) {
	backoff := 1
	for {
		if err := run(token, s); err != nil {
			log.Printf("Gateway disconnected: %v - reconnecting in %ds", err, backoff)
		}
		time.Sleep(time.Duration(backoff) * time.Second)
		if backoff < maxBackoff {
			backoff *= 2
		}
	}
}

func run(token string, s *store.Store) error {
	conn, _, err := websocket.DefaultDialer.Dial(gatewayURL, nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	var seq seqHolder
	heartbeatStop := make(chan struct{})
	defer close(heartbeatStop)

	for {
		_, raw, err := conn.ReadMessage()
		if err != nil {
			return err
		}

		var p Payload
		if err := json.Unmarshal(raw, &p); err != nil {
			continue
		}
		if p.S != nil {
			seq.set(p.S)
		}

		switch p.Op {
		case 10: // Hello
			var hello struct {
				HeartbeatInterval int `json:"heartbeat_interval"`
			}
			if err := json.Unmarshal(p.D, &hello); err != nil {
				return err
			}
			if hello.HeartbeatInterval <= 0 {
				hello.HeartbeatInterval = defaultHeartbeatInterval
			}
			go heartbeat(conn, hello.HeartbeatInterval, &seq, heartbeatStop)
			sendIdentify(conn, token)
		case 11: // HeartbeatAck - no op
		case 0: // Dispatch
			if p.T != nil {
				dispatch(*p.T, p.D, s)
			}
		}
	}
}

func heartbeat(conn *websocket.Conn, intervalMs int, seq *seqHolder, stop <-chan struct{}) {
	ticker := time.NewTicker(time.Duration(intervalMs) * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			payload, err := json.Marshal(Payload{Op: opHeartbeat, D: seqJSON(seq.get())})
			if err != nil {
				log.Printf("Heartbeat marshal error: %v", err)
				continue
			}
			if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
				log.Printf("Heartbeat write error: %v", err)
			}
		}
	}
}

func sendIdentify(conn *websocket.Conn, token string) {
	payload, err := json.Marshal(map[string]any{
		"op": opIdentify,
		"d": map[string]any{
			"token":   token,
			"intents": intentGuilds | intentGuildPresences,
			"properties": map[string]string{
				"os":      "linux",
				"browser": "gostatus",
				"device":  "gostatus",
			},
		},
	})
	if err != nil {
		log.Printf("Identify marshal error: %v", err)
		return
	}
	if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
		log.Printf("Identify write error: %v", err)
	}
}

func dispatch(event string, d json.RawMessage, s *store.Store) {
	switch event {
	case "GUILD_CREATE":
		var gc guildCreate
		if err := json.Unmarshal(d, &gc); err != nil {
			log.Printf("Failed to unmarshal GUILD_CREATE: %v", err)
			return
		}
		for _, pr := range gc.Presences {
			if pr.User != nil && pr.User.ID != "" {
				s.Set(pr.User.ID, store.Presence{
					Status:       pr.Status,
					ClientStatus: pr.ClientStatus,
					Activities:   toStoreActivities(pr.Activities),
				})
			}
		}
		log.Printf("Cached %d presences from GUILD_CREATE", len(gc.Presences))

	case "PRESENCE_UPDATE":
		var pu PresenceUpdate
		if err := json.Unmarshal(d, &pu); err != nil {
			log.Printf("Failed to unmarshal PRESENCE_UPDATE: %v", err)
			return
		}
		if pu.User != nil && pu.User.ID != "" {
			s.Set(pu.User.ID, store.Presence{
				Status:       pu.Status,
				ClientStatus: pu.ClientStatus,
				Activities:   toStoreActivities(pu.Activities),
			})
		}
	}
}

func toStoreActivities(in []Activity) []store.Activity {
	out := make([]store.Activity, len(in))
	for i, a := range in {
		out[i] = store.Activity{
			Name:    a.Name,
			Type:    a.Type,
			Details: a.Details,
			State:   a.State,
			SyncID:  a.SyncID,
		}
	}
	return out
}

func seqJSON(seq *int) json.RawMessage {
	if seq == nil {
		return json.RawMessage("null")
	}
	b, _ := json.Marshal(*seq)
	return b
}
