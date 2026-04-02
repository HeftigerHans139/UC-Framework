package api

import (
	"encoding/json"
	"net/http"
	"sort"
	"uc_framework/internal/bot"
)

var PluginRegistry *bot.PluginRegistry
var SyncPluginEnabledStateFunc func(name string, active bool) error
var GetPluginEnabledStateFunc func(name string) (bool, bool)
var ExtraPluginsFunc func() []PluginStatus
var ToggleExtraPluginFunc func(name string, active bool) error

// PluginStatus reprÃ¤sentiert den Status eines Plugins fÃ¼r das Webinterface
type PluginStatus struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
}

// GET /api/plugins
func PluginsHandler(w http.ResponseWriter, r *http.Request) {
	if PluginRegistry == nil {
		http.Error(w, "plugin registry unavailable", http.StatusServiceUnavailable)
		return
	}

	plugins := PluginRegistry.All()
	var result []PluginStatus
	names := make([]string, 0, len(plugins))
	for name := range plugins {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		state := plugins[name]
		active := state.Loaded
		if GetPluginEnabledStateFunc != nil {
			if persisted, ok := GetPluginEnabledStateFunc(name); ok {
				active = persisted
			}
		}
		result = append(result, PluginStatus{
			Name:        name,
			Description: state.Description,
			Active:      active,
		})
	}

	if ExtraPluginsFunc != nil {
		result = append(result, ExtraPluginsFunc()...)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// POST /api/plugins/toggle
func TogglePluginHandler(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		Name   string `json:"name"`
		Active bool   `json:"active"`
	}
	var req Req
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if PluginRegistry == nil {
		http.Error(w, "plugin registry unavailable", http.StatusServiceUnavailable)
		return
	}

	// Check if it's a real registered plugin
	registeredPlugins := PluginRegistry.All()
	if _, isRegistered := registeredPlugins[req.Name]; isRegistered {
		if req.Active {
			if err := PluginRegistry.Load(req.Name); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		} else {
			if err := PluginRegistry.Unload(req.Name); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		if SyncPluginEnabledStateFunc != nil {
			if err := SyncPluginEnabledStateFunc(req.Name, req.Active); err != nil {
				http.Error(w, "sync plugin state failed: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PluginStatus{
			Name:   req.Name,
			Active: PluginRegistry.IsLoaded(req.Name),
		})
		return
	}

	// Try virtual/extra plugin toggle
	if ToggleExtraPluginFunc != nil {
		if err := ToggleExtraPluginFunc(req.Name, req.Active); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PluginStatus{
			Name:   req.Name,
			Active: req.Active,
		})
		return
	}

	http.Error(w, "plugin not found", http.StatusNotFound)
}
