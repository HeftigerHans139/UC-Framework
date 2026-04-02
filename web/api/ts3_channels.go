package api

import (
	"encoding/json"
	"net/http"
)

type TS3Channel struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var GetTS3ChannelsFunc func() ([]TS3Channel, error)

func TS3ChannelsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if GetTS3ChannelsFunc == nil {
		http.Error(w, "ts3 channels unavailable", http.StatusServiceUnavailable)
		return
	}

	channels, err := GetTS3ChannelsFunc()
	if err != nil {
		http.Error(w, "failed to load channels: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"ok":       true,
		"channels": channels,
	})
}
