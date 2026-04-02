package api

import (
	"encoding/json"
	"net/http"
	"strings"
)

type TS3Settings struct {
	Host           string `json:"host"`
	QueryPort      int    `json:"query_port"`
	VoicePort      int    `json:"voice_port"`
	QueryUsername  string `json:"query_username"`
	QueryPassword  string `json:"query_password"`
	BotNickname    string `json:"bot_nickname"`
	DefaultChannel string `json:"default_channel"`
	QuerySlowmode  int    `json:"query_slowmode_ms"`
}

var (
	GetTS3SettingsFunc  func() TS3Settings
	SaveTS3SettingsFunc func(settings TS3Settings) error
)

func TS3SettingsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if GetTS3SettingsFunc == nil {
			http.Error(w, "ts3 settings unavailable", http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(GetTS3SettingsFunc())

	case http.MethodPost:
		if SaveTS3SettingsFunc == nil {
			http.Error(w, "ts3 settings unavailable", http.StatusServiceUnavailable)
			return
		}
		var req TS3Settings
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		if strings.TrimSpace(req.Host) == "" {
			http.Error(w, "host is required", http.StatusBadRequest)
			return
		}
		if req.QueryPort <= 0 || req.QueryPort > 65535 {
			http.Error(w, "query_port must be between 1 and 65535", http.StatusBadRequest)
			return
		}
		if req.VoicePort <= 0 || req.VoicePort > 65535 {
			http.Error(w, "voice_port must be between 1 and 65535", http.StatusBadRequest)
			return
		}
		if strings.TrimSpace(req.QueryUsername) == "" {
			http.Error(w, "query_username is required", http.StatusBadRequest)
			return
		}
		if req.QuerySlowmode <= 0 {
			http.Error(w, "query_slowmode_ms must be > 0", http.StatusBadRequest)
			return
		}

		req.Host = strings.TrimSpace(req.Host)
		req.QueryUsername = strings.TrimSpace(req.QueryUsername)
		req.BotNickname = strings.TrimSpace(req.BotNickname)
		req.DefaultChannel = strings.TrimSpace(req.DefaultChannel)

		if err := SaveTS3SettingsFunc(req); err != nil {
			http.Error(w, "failed to save ts3 settings: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"ok":               true,
			"restart_required": true,
		})

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
