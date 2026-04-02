package api

import (
	"encoding/json"
	"net/http"
)

const FrameworkVersion = "1.0.0"

type StatusResponse struct {
	Status  string `json:"status"`
	Uptime  string `json:"uptime"`
	Version string `json:"version"`
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	resp := StatusResponse{
		Status:  "ok",
		Uptime:  "TODO",
		Version: FrameworkVersion,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
