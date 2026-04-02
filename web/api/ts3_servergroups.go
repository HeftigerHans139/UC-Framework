package api

import (
	"encoding/json"
	"net/http"
)

type TS3ServerGroup struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var GetTS3ServerGroupsFunc func() ([]TS3ServerGroup, error)

func TS3ServerGroupsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if GetTS3ServerGroupsFunc == nil {
		http.Error(w, "ts3 server groups unavailable", http.StatusServiceUnavailable)
		return
	}

	groups, err := GetTS3ServerGroupsFunc()
	if err != nil {
		http.Error(w, "failed to load server groups: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"ok":     true,
		"groups": groups,
	})
}
