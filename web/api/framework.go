package api

import (
	"encoding/json"
	"net/http"
)

type FrameworkInfo struct {
	Name          string `json:"name"`
	Version       string `json:"version"`
	UpdatedAt     string `json:"updated_at"`
	LatestVersion string `json:"latest_version"`
	IsLatest      bool   `json:"is_latest"`
}

var (
	GetFrameworkInfoFunc       func() (FrameworkInfo, error)
	RestartFrameworkFunc       func() error
	SendServerAnnouncementFunc func(message string) error
)

type serverAnnouncementRequest struct {
	Message string `json:"message"`
}

func FrameworkInfoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if GetFrameworkInfoFunc == nil {
		http.Error(w, "framework info unavailable", http.StatusServiceUnavailable)
		return
	}

	info, err := GetFrameworkInfoFunc()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"ok":   true,
		"info": info,
	})
}

func FrameworkRestartHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if RestartFrameworkFunc == nil {
		http.Error(w, "framework restart unavailable", http.StatusServiceUnavailable)
		return
	}

	if err := RestartFrameworkFunc(); err != nil {
		http.Error(w, "failed to restart framework: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"ok":      true,
		"message": "framework restart scheduled",
	})
}

func FrameworkAnnouncementHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if SendServerAnnouncementFunc == nil {
		http.Error(w, "announcement unavailable", http.StatusServiceUnavailable)
		return
	}

	var req serverAnnouncementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if err := SendServerAnnouncementFunc(req.Message); err != nil {
		http.Error(w, "failed to send announcement: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"ok":      true,
		"message": "announcement sent",
	})
}
