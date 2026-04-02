package api

import (
	"encoding/json"
	"net/http"
	"time"
)

type AnnouncementSettings struct {
	Message               string `json:"message"`
	RepeatEnabled         bool   `json:"repeat_enabled"`
	ScheduleMode          string `json:"schedule_mode"`
	RepeatIntervalMinutes int    `json:"repeat_interval_minutes"`
	RepeatIntervalCount   int    `json:"repeat_interval_count"`
	RepeatTime            string `json:"repeat_time"`
}

type AnnouncementStatus struct {
	Message             string     `json:"message"`
	LastSentAt          *time.Time `json:"last_sent_at,omitempty"`
	RepeatEnabled       bool       `json:"repeat_enabled"`
	ScheduleMode        string     `json:"schedule_mode"`
	RepeatIntervalCount int        `json:"repeat_interval_count"`
}

type announcementSendRequest struct {
	Message string `json:"message"`
}

var (
	GetAnnouncementSettingsFunc  func() (AnnouncementSettings, error)
	SaveAnnouncementSettingsFunc func(settings AnnouncementSettings) error
	GetAnnouncementStatusFunc    func() (AnnouncementStatus, error)
	SendAnnouncementFunc         func(message string) error
)

func AnnouncementSettingsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if GetAnnouncementSettingsFunc == nil {
			http.Error(w, "announcement settings unavailable", http.StatusServiceUnavailable)
			return
		}
		settings, err := GetAnnouncementSettingsFunc()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(settings)
	case http.MethodPost:
		if SaveAnnouncementSettingsFunc == nil {
			http.Error(w, "announcement settings unavailable", http.StatusServiceUnavailable)
			return
		}
		var req AnnouncementSettings
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		if err := SaveAnnouncementSettingsFunc(req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"ok": true})
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func AnnouncementStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if GetAnnouncementStatusFunc == nil {
		http.Error(w, "announcement status unavailable", http.StatusServiceUnavailable)
		return
	}
	status, err := GetAnnouncementStatusFunc()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"ok": true, "status": status})
}

func AnnouncementSendHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if SendAnnouncementFunc == nil {
		http.Error(w, "announcement send unavailable", http.StatusServiceUnavailable)
		return
	}

	var req announcementSendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if err := SendAnnouncementFunc(req.Message); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"ok": true})
}
