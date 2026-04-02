package api

import (
	"encoding/json"
	"net/http"
)

type FrameworkSettings struct {
	PlatformMode string `json:"platform_mode"`
}

var (
	GetFrameworkSettingsFunc  func() (FrameworkSettings, error)
	SaveFrameworkSettingsFunc func(settings FrameworkSettings) error
)

func FrameworkSettingsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if GetFrameworkSettingsFunc == nil {
			http.Error(w, "framework settings unavailable", http.StatusServiceUnavailable)
			return
		}
		settings, err := GetFrameworkSettingsFunc()
		if err != nil {
			http.Error(w, "failed to load framework settings: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(settings)
	case http.MethodPost:
		if SaveFrameworkSettingsFunc == nil {
			http.Error(w, "framework settings unavailable", http.StatusServiceUnavailable)
			return
		}
		var req FrameworkSettings
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		if err := SaveFrameworkSettingsFunc(req); err != nil {
			http.Error(w, "failed to save framework settings: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"ok": true})
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
