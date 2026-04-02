package api

import (
	"encoding/json"
	"net/http"
)

type DiscordChannel struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	ParentID string `json:"parent_id"`
}

type DiscordRole struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var (
	GetDiscordChannelsFunc func() ([]DiscordChannel, error)
	GetDiscordRolesFunc    func() ([]DiscordRole, error)
)

func DiscordChannelsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if GetDiscordChannelsFunc == nil {
		http.Error(w, "discord channels unavailable", http.StatusServiceUnavailable)
		return
	}
	channels, err := GetDiscordChannelsFunc()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"channels": channels})
}

func DiscordRolesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if GetDiscordRolesFunc == nil {
		http.Error(w, "discord roles unavailable", http.StatusServiceUnavailable)
		return
	}
	roles, err := GetDiscordRolesFunc()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"roles": roles})
}
