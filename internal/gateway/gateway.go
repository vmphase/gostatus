package gateway

import (
	"encoding/json"
	"gostatus/internal/store"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	gatewayURL = "wss://gateway.discord.gg/?v=10&encoding=json"

	opHeartbeat = 1
	opIdentify  = 2

	intentGuilds         = 1 << 0
	intentGuildPresences = 1 << 8
)

func Connect(token string, s *store.Store) {
	for {
		if err := run(token, s); err != nil {
			log.Printf("Gateway disconnected: %v - reconnecting in 5s", err)
		}
		time.Sleep(5 * time.Second)
	}
}

func run(token string, s *store.Store) error {
	conn, _, err := websocket.DefaultDialer.Dial(gatewayURL, nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	var seq *int
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
			seq = p.S
		}

		switch p.Op {
		case 10: // Hello
			var hello struct {
				HeartbeatInterval int `json:"heartbeat_interval"`
			}
			json.Unmarshal(p.D, &hello)
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

func heartbeat(conn *websocket.Conn, intervalMs int, seq **int, stop <-chan struct{}) {
	ticker := time.NewTicker(time.Duration(intervalMs) * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			payload, _ := json.Marshal(Payload{Op: opHeartbeat, D: seqJSON(*seq)})
			conn.WriteMessage(websocket.TextMessage, payload)
		}
	}
}

func sendIdentify(conn *websocket.Conn, token string) {
	payload, _ := json.Marshal(map[string]any{
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
	conn.WriteMessage(websocket.TextMessage, payload)
}

func dispatch(event string, d json.RawMessage, s *store.Store) {
	switch event {
	case "GUILD_CREATE":
		var gc guildCreate
		json.Unmarshal(d, &gc)
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
		json.Unmarshal(d, &pu)
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
