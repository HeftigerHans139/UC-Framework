package api

import (
	"encoding/json"
	"net/http"
	"strings"
)

type DiscordSettings struct {
	Enabled               bool     `json:"enabled"`
	BotToken              string   `json:"bot_token"`
	ApplicationID         string   `json:"application_id"`
	GuildID               string   `json:"guild_id"`
	AFKKickEnabled        bool     `json:"afk_kick_enabled"`
	AFKInactivityMinutes  int      `json:"afk_inactivity_minutes"`
	BotDisplayName        string   `json:"bot_display_name"`
	StatusText            string   `json:"status_text"`
	CommandPrefix         string   `json:"command_prefix"`
	LogChannelID          string   `json:"log_channel_id"`
	AnnouncementChannelID string   `json:"announcement_channel_id"`
	SupportCategoryID     string   `json:"support_category_id"`
	SupportLogChannelID   string   `json:"support_log_channel_id"`
	AdminRoleIDs          []string `json:"admin_role_ids"`
	SupporterRoleIDs      []string `json:"supporter_role_ids"`
	BotRoleIDs            []string `json:"bot_role_ids"`
}

var (
	GetDiscordSettingsFunc  func() (DiscordSettings, error)
	SaveDiscordSettingsFunc func(settings DiscordSettings) error
)

func DiscordSettingsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if GetDiscordSettingsFunc == nil {
			http.Error(w, "discord settings unavailable", http.StatusServiceUnavailable)
			return
		}
		settings, err := GetDiscordSettingsFunc()
		if err != nil {
			http.Error(w, "failed to load discord settings: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(settings)
	case http.MethodPost:
		if SaveDiscordSettingsFunc == nil {
			http.Error(w, "discord settings unavailable", http.StatusServiceUnavailable)
			return
		}
		var req DiscordSettings
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		req.ApplicationID = strings.TrimSpace(req.ApplicationID)
		req.GuildID = strings.TrimSpace(req.GuildID)
		if req.AFKInactivityMinutes <= 0 {
			req.AFKInactivityMinutes = 30
		}
		req.BotDisplayName = strings.TrimSpace(req.BotDisplayName)
		req.CommandPrefix = strings.TrimSpace(req.CommandPrefix)
		if req.Enabled {
			if strings.TrimSpace(req.BotToken) == "" {
				http.Error(w, "bot_token is required when Discord is enabled", http.StatusBadRequest)
				return
			}
			if req.GuildID == "" {
				http.Error(w, "guild_id is required when Discord is enabled", http.StatusBadRequest)
				return
			}
		}
		if req.CommandPrefix == "" {
			req.CommandPrefix = "!"
		}
		if err := SaveDiscordSettingsFunc(req); err != nil {
			http.Error(w, "failed to save discord settings: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"ok":               true,
			"restart_required": true,
		})
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
