package ts3

import (
	"strings"
	"time"
)

type Client struct {
	ID           int
	Nickname     string
	ServerGroups []int
	IsQuery      bool
	ChannelID    int   // Channel, in dem sich der User befindet
	LastActive   int64 // Unix-Timestamp der letzten Aktivität
}

func (c *TS3Client) ListClients() ([]Client, error) {
	rows, err := c.query("clientlist -groups -times")
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()
	clients := make([]Client, 0, len(rows))
	for _, row := range rows {
		groupsRaw := row["client_servergroups"]
		groups := []int{}
		if groupsRaw != "" {
			for _, s := range splitComma(groupsRaw) {
				groups = append(groups, atoiDefault(s, 0))
			}
		}

		idleMS := atoiDefault(row["client_idle_time"], 0)
		lastActive := now - int64(idleMS/1000)

		clients = append(clients, Client{
			ID:           atoiDefault(row["clid"], 0),
			Nickname:     row["client_nickname"],
			ServerGroups: groups,
			IsQuery:      atoiDefault(row["client_type"], 0) == 1,
			ChannelID:    atoiDefault(row["cid"], 0),
			LastActive:   lastActive,
		})
	}

	return clients, nil
}

func (c *TS3Client) MoveClient(clientID int, channelID int) error {
	_, err := c.query("clientmove clid=" + intToString(clientID) + " cid=" + intToString(channelID))
	return err
}

func (c *TS3Client) PokeClient(clientID int, message string) error {
	if clientID <= 0 {
		return nil
	}

	msg := strings.TrimSpace(message)
	if msg == "" {
		return nil
	}

	// TS3 begrenzt die Groesse des msg-Parameters fuer clientpoke.
	msg = truncateRunes(msg, 100)
	_, err := c.query("clientpoke clid=" + intToString(clientID) + " msg=" + escapeTS3(msg))
	if err == nil {
		return nil
	}

	if strings.Contains(err.Error(), "id=1541") {
		fallback := truncateRunes(msg, 60)
		if fallback != "" && fallback != msg {
			if _, retryErr := c.query("clientpoke clid=" + intToString(clientID) + " msg=" + escapeTS3(fallback)); retryErr == nil {
				return nil
			}
		}
	}

	return err
}

func truncateRunes(v string, max int) string {
	if max <= 0 {
		return ""
	}
	r := []rune(v)
	if len(r) <= max {
		return v
	}
	if max <= 3 {
		return string(r[:max])
	}
	return string(r[:max-3]) + "..."
}

func splitComma(v string) []string {
	if v == "" {
		return nil
	}
	parts := make([]string, 0)
	start := 0
	for i := 0; i < len(v); i++ {
		if v[i] == ',' {
			parts = append(parts, v[start:i])
			start = i + 1
		}
	}
	parts = append(parts, v[start:])
	return parts
}
