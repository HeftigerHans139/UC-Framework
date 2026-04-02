package api

import (
	"encoding/json"
	"net/http"
)

// SavePluginConfigFunc wird von core.go gesetzt und speichert die Plugin-Konfiguration persistent.
var SavePluginConfigFunc func(name string, raw []byte) error
var GetSavedPluginConfigFunc func(name string) ([]byte, error)

// GET  /api/plugins/config?name=<PluginName>  – gibt die aktuelle JSON-Konfiguration zurück
// POST /api/plugins/config?name=<PluginName>  – aktualisiert die Konfiguration in-memory und auf Disk
func PluginConfigHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "missing name parameter", http.StatusBadRequest)
		return
	}
	if PluginRegistry == nil {
		http.Error(w, "plugin registry unavailable", http.StatusServiceUnavailable)
		return
	}

	switch r.Method {
	case http.MethodGet:
		if GetSavedPluginConfigFunc != nil {
			raw, err := GetSavedPluginConfigFunc(name)
			if err == nil && raw != nil {
				w.Header().Set("Content-Type", "application/json")
				w.Write(raw)
				return
			}
		}

		raw, err := PluginRegistry.GetPluginConfig(name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(raw)

	case http.MethodPost:
		var raw json.RawMessage
		if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		// In-Memory-Update versuchen (schlägt fehl, falls Plugin nicht geladen)
		_, err := PluginRegistry.UpdatePluginConfig(name, raw)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Konfiguration persistent speichern
		if SavePluginConfigFunc != nil {
			if err := SavePluginConfigFunc(name, raw); err != nil {
				http.Error(w, "disk write failed: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}

		if GetSavedPluginConfigFunc != nil {
			saved, err := GetSavedPluginConfigFunc(name)
			if err == nil && saved != nil {
				w.Header().Set("Content-Type", "application/json")
				w.Write(saved)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(raw)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
