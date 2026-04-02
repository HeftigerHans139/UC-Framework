package api

import (
	"encoding/json"
	"net/http"
)

type StatsResponse struct {
	Admins  int `json:"admins_online"`
	Members int `json:"members_online"`
}

var (
	GetStatsFunc func() (int, int)
)

func StatsHandler(w http.ResponseWriter, r *http.Request) {
	admins, members := 0, 0
	if GetStatsFunc != nil {
		admins, members = GetStatsFunc()
	}
	resp := StatsResponse{
		Admins:  admins,
		Members: members,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
