package api

import (
	"encoding/json"
	"net/http"
	"time"
)

type TS3ConnectionStatus struct {
	BotRunning     bool   `json:"bot_running"`
	Connected      bool   `json:"connected"`
	Host           string `json:"host"`
	Port           int    `json:"port"`
	LastCheckAt    string `json:"last_check_at"`
	LastError      string `json:"last_error,omitempty"`
	Implementation string `json:"implementation"`
}

var (
	GetTS3ConnectionStatusFunc func() TS3ConnectionStatus
	TestTS3ConnectionFunc      func() (TS3ConnectionStatus, error)
)

func TS3ConnectionStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if GetTS3ConnectionStatusFunc == nil {
		http.Error(w, "ts3 connection status unavailable", http.StatusServiceUnavailable)
		return
	}

	status := GetTS3ConnectionStatusFunc()
	if status.LastCheckAt == "" {
		status.LastCheckAt = time.Now().Format(time.RFC3339)
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"ok":     true,
		"status": status,
	})
}

func TS3ConnectionTestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if TestTS3ConnectionFunc == nil {
		http.Error(w, "ts3 connection test unavailable", http.StatusServiceUnavailable)
		return
	}

	status, err := TestTS3ConnectionFunc()
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"ok":     false,
			"error":  err.Error(),
			"status": status,
		})
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]any{
		"ok":     true,
		"status": status,
	})
}
