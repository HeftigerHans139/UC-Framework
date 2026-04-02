package api

import (
	"encoding/json"
	"net/http"
)

type SupportSettings struct {
	Enabled               bool   `json:"enabled"`
	SupportChannelIDs     []int  `json:"support_channel_ids"`
	WaitingAreaChannel    int    `json:"waiting_area_channel_id"`
	OpenPokeMessage       string `json:"open_poke_message"`
	ClosedPokeMessage     string `json:"closed_poke_message"`
	JoinOpenPokeMessage   string `json:"join_open_poke_message"`
	JoinClosedPokeMessage string `json:"join_closed_poke_message"`
	SupporterPokeMessage  string `json:"supporter_poke_message"`
	SupporterGroupIDs     []int  `json:"supporter_group_ids"`
	AutoScheduleEnabled   bool   `json:"auto_schedule_enabled"`
	AutoOpenTime          string `json:"auto_open_time"`
	AutoCloseTime         string `json:"auto_close_time"`
}

type SupportStatus struct {
	Open                bool   `json:"open"`
	LastAction          string `json:"last_action"`
	LastError           string `json:"last_error,omitempty"`
	AutoScheduleEnabled bool   `json:"auto_schedule_enabled"`
	AutoOpenTime        string `json:"auto_open_time"`
	AutoCloseTime       string `json:"auto_close_time"`
}

type supportActionRequest struct {
	Action string `json:"action"`
}

var (
	GetSupportSettingsFunc   func() (SupportSettings, error)
	SaveSupportSettingsFunc  func(settings SupportSettings) error
	GetSupportStatusFunc     func() (SupportStatus, error)
	ExecuteSupportActionFunc func(action string) (SupportStatus, error)
)

func SupportSettingsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if GetSupportSettingsFunc == nil {
			http.Error(w, "support settings unavailable", http.StatusServiceUnavailable)
			return
		}
		settings, err := GetSupportSettingsFunc()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(settings)
	case http.MethodPost:
		if SaveSupportSettingsFunc == nil {
			http.Error(w, "support settings unavailable", http.StatusServiceUnavailable)
			return
		}
		var req SupportSettings
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		if err := SaveSupportSettingsFunc(req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"ok": true})
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func SupportStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if GetSupportStatusFunc == nil {
		http.Error(w, "support status unavailable", http.StatusServiceUnavailable)
		return
	}
	status, err := GetSupportStatusFunc()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"ok": true, "status": status})
}

func SupportActionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if ExecuteSupportActionFunc == nil {
		http.Error(w, "support action unavailable", http.StatusServiceUnavailable)
		return
	}

	var req supportActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	status, err := ExecuteSupportActionFunc(req.Action)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"ok": true, "status": status})
}
