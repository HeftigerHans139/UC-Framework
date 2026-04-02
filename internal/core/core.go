package core

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
	"uc_framework/internal/bot"
	"uc_framework/internal/discord"
	"uc_framework/internal/ts3"
	"uc_framework/plugins/admincounter"
	"uc_framework/plugins/afkmover"
	"uc_framework/plugins/combinedstats"
	"uc_framework/plugins/membercounter"
	"uc_framework/web"
	"uc_framework/web/api"
)

type Config struct {
	DefaultLanguage    string                     `json:"default_language"`
	SupportedLanguages []string                   `json:"supported_languages"`
	Framework          FrameworkConfig            `json:"framework"`
	TS3                TS3Config                  `json:"ts3"`
	Discord            DiscordConfig              `json:"discord"`
	Support            SupportConfig              `json:"support"`
	Announcement       AnnouncementConfig         `json:"announcement"`
	BotControl         BotControlConfig           `json:"bot_control"`
	WebAuth            WebAuthConfig              `json:"web_auth"`
	PluginConfigs      map[string]json.RawMessage `json:"plugin_configs"`
}

type FrameworkConfig struct {
	PlatformMode string `json:"platform_mode"`
}

type DiscordConfig struct {
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

type SupportConfig struct {
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

type AnnouncementConfig struct {
	Message               string `json:"message"`
	RepeatEnabled         bool   `json:"repeat_enabled"`
	ScheduleMode          string `json:"schedule_mode"`
	RepeatIntervalMinutes int    `json:"repeat_interval_minutes"`
	RepeatIntervalCount   int    `json:"repeat_interval_count"`
	RepeatTime            string `json:"repeat_time"`
}

type BotControlConfig struct {
	Enabled                bool     `json:"enabled"`
	BotExecutable          string   `json:"bot_executable"`
	WorkingDir             string   `json:"working_dir"`
	BotArgs                []string `json:"bot_args"`
	WatchdogEnabled        bool     `json:"watchdog_enabled"`
	WatchdogMinIntervalSec int      `json:"watchdog_min_interval_sec"`
	WatchdogMaxIntervalSec int      `json:"watchdog_max_interval_sec"`
}

type WebAuthConfig struct {
	Enabled                 bool                 `json:"enabled"`
	Provider                string               `json:"provider"`
	AllowRanksystemFallback bool                 `json:"allow_ranksystem_fallback"`
	Username                string               `json:"username"`
	Password                string               `json:"password"`
	PasswordHash            string               `json:"password_hash"`
	ForcePasswordChange     bool                 `json:"force_password_change"`
	SessionTTLMinutes       int                  `json:"session_ttl_minutes"`
	Ranksystem              WebAuthRanksystemCfg `json:"ranksystem"`
}

type WebAuthRanksystemCfg struct {
	LoginURL      string `json:"login_url"`
	UsernameField string `json:"username_field"`
	PasswordField string `json:"password_field"`
	APIKeyHeader  string `json:"api_key_header"`
	APIKeyValue   string `json:"api_key_value"`
	BearerToken   string `json:"bearer_token"`
}

type TS3Config struct {
	Host           string `json:"host"`
	QueryPort      int    `json:"query_port"`
	VoicePort      int    `json:"voice_port"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	BotNickname    string `json:"bot_nickname"`
	DefaultChannel string `json:"default_channel"`
	QuerySlowmode  int    `json:"query_slowmode_ms"`
	AdminGroups    []int  `json:"admin_groups"`
	BotGroups      []int  `json:"bot_groups"`
	Port           int    `json:"port,omitempty"` // legacy: wird als Fallback auf QueryPort gemappt
}

var (
	config         Config
	language       map[string]string
	configFilePath string
	liveConfigs    map[string]json.RawMessage
	liveConfigsMu  sync.RWMutex

	webAuthPasswordFromEnv     bool
	webAuthPasswordHashFromEnv bool
	ts3PasswordFromEnv         bool
	discordTokenFromEnv        bool

	pluginEnabledStates   map[string]bool
	pluginEnabledStatesMu sync.RWMutex
)

func parseBoolEnv(name string) bool {
	v := strings.ToLower(strings.TrimSpace(os.Getenv(name)))
	return v == "1" || v == "true" || v == "yes" || v == "on"
}

func applySecretEnvOverrides(cfg *Config) {
	ts3PasswordFromEnv = false
	discordTokenFromEnv = false
	webAuthPasswordFromEnv = false
	webAuthPasswordHashFromEnv = false

	if v := os.Getenv("UC_FRAMEWORK_TS3_PASSWORD"); v != "" {
		cfg.TS3.Password = v
		ts3PasswordFromEnv = true
	}
	if v := strings.TrimSpace(os.Getenv("UC_FRAMEWORK_DISCORD_BOT_TOKEN")); v != "" {
		cfg.Discord.BotToken = v
		discordTokenFromEnv = true
	}
	if v := os.Getenv("UC_FRAMEWORK_WEB_AUTH_PASSWORD"); v != "" {
		cfg.WebAuth.Password = v
		webAuthPasswordFromEnv = true
	}
	if v := strings.TrimSpace(os.Getenv("UC_FRAMEWORK_WEB_AUTH_PASSWORD_HASH")); v != "" {
		cfg.WebAuth.PasswordHash = v
		cfg.WebAuth.Password = ""
		webAuthPasswordHashFromEnv = true
	}
}

func enforceSecretPolicy(cfg *Config) {
	strict := parseBoolEnv("UC_FRAMEWORK_REQUIRE_ENV_SECRETS")

	if strings.TrimSpace(cfg.TS3.Password) != "" && !ts3PasswordFromEnv {
		msg := "[SECURITY] TS3 password comes from config file. Prefer UC_FRAMEWORK_TS3_PASSWORD env var"
		if strict {
			log.Fatal(msg)
		}
		log.Println(msg)
	}
	if strings.TrimSpace(cfg.Discord.BotToken) != "" && !discordTokenFromEnv {
		msg := "[SECURITY] Discord bot token comes from config file. Prefer UC_FRAMEWORK_DISCORD_BOT_TOKEN env var"
		if strict {
			log.Fatal(msg)
		}
		log.Println(msg)
	}
	if strings.TrimSpace(cfg.WebAuth.Password) != "" && !webAuthPasswordFromEnv {
		msg := "[SECURITY] Web auth plaintext password comes from config file. Prefer hash/env and avoid plaintext"
		if strict {
			log.Fatal(msg)
		}
		log.Println(msg)
	}
}

func logInternetSecurityStartupChecks(cfg *Config) {
	if !parseBoolEnv("UC_FRAMEWORK_INTERNET_MODE") {
		return
	}

	log.Println("[SECURITY] [STARTUP] Internet mode enabled - validating hardening prerequisites")

	enforceHTTPS := parseBoolEnv("UC_FRAMEWORK_ENFORCE_HTTPS")
	if !enforceHTTPS {
		log.Println("[SECURITY] [CRITICAL] UC_FRAMEWORK_ENFORCE_HTTPS is disabled while UC_FRAMEWORK_INTERNET_MODE=true")
	}

	requireEnvSecrets := parseBoolEnv("UC_FRAMEWORK_REQUIRE_ENV_SECRETS")
	if !requireEnvSecrets {
		log.Println("[SECURITY] [WARN] UC_FRAMEWORK_REQUIRE_ENV_SECRETS is disabled in internet mode")
	}

	if strings.TrimSpace(cfg.TS3.Password) != "" && !ts3PasswordFromEnv {
		log.Println("[SECURITY] [CRITICAL] TS3 password is loaded from config file in internet mode")
	}
	if strings.TrimSpace(cfg.Discord.BotToken) != "" && !discordTokenFromEnv {
		log.Println("[SECURITY] [CRITICAL] Discord bot token is loaded from config file in internet mode")
	}
	if strings.TrimSpace(cfg.WebAuth.Password) != "" && !webAuthPasswordFromEnv {
		log.Println("[SECURITY] [CRITICAL] Web auth plaintext password is loaded from config file in internet mode")
	}

	if !cfg.WebAuth.Enabled {
		log.Println("[SECURITY] [CRITICAL] Web authentication is disabled in internet mode")
	}
	if strings.TrimSpace(cfg.WebAuth.PasswordHash) == "" && !webAuthPasswordHashFromEnv {
		log.Println("[SECURITY] [WARN] Web auth password hash is empty in internet mode")
	}

	if parseBoolEnv("UC_FRAMEWORK_TRUST_PROXY") {
		log.Println("[SECURITY] [WARN] UC_FRAMEWORK_TRUST_PROXY is enabled - ensure requests come only through a trusted reverse proxy")
	}
}

func resolveConfigFilePath() string {
	candidates := make([]string, 0, 3)

	if cwd, err := os.Getwd(); err == nil {
		candidates = append(candidates, filepath.Join(cwd, "config", "config.json"))
	}

	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		candidates = append(candidates, filepath.Join(exeDir, "config", "config.json"))
	}

	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}

	if len(candidates) > 0 {
		return candidates[0]
	}

	return filepath.Join("config", "config.json")
}

func normalizeLanguageCode(language string) string {
	switch strings.ToLower(strings.TrimSpace(language)) {
	case "de":
		return "de"
	default:
		return "en"
	}
}

func normalizeSupportedLanguages(languages []string) []string {
	seen := map[string]struct{}{}
	result := make([]string, 0, len(languages))
	for _, language := range languages {
		normalized := normalizeLanguageCode(language)
		if _, ok := seen[normalized]; ok {
			continue
		}
		seen[normalized] = struct{}{}
		result = append(result, normalized)
	}
	if len(result) == 0 {
		return []string{"en", "de"}
	}
	return result
}

func isSupportedLanguage(language string, supported []string) bool {
	for _, candidate := range supported {
		if candidate == language {
			return true
		}
	}
	return false
}

func canonicalPluginConfigKey(name string) string {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "admincounter":
		return "admincounter"
	case "membercounter":
		return "membercounter"
	case "combinedstats":
		return "combinedstats"
	case "afkmover":
		return "afkmover"
	default:
		return strings.ToLower(strings.TrimSpace(name))
	}
}

func normalizeSupportConfig(in SupportConfig) SupportConfig {
	seen := map[int]struct{}{}
	channels := make([]int, 0, len(in.SupportChannelIDs))
	groupSeen := map[int]struct{}{}
	supporterGroups := make([]int, 0, len(in.SupporterGroupIDs))
	for _, id := range in.SupportChannelIDs {
		if id <= 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		channels = append(channels, id)
	}
	for _, id := range in.SupporterGroupIDs {
		if id <= 0 {
			continue
		}
		if _, ok := groupSeen[id]; ok {
			continue
		}
		groupSeen[id] = struct{}{}
		supporterGroups = append(supporterGroups, id)
	}
	out := in
	out.Enabled = in.Enabled
	out.SupportChannelIDs = channels
	out.SupporterGroupIDs = supporterGroups
	out.OpenPokeMessage = strings.TrimSpace(in.OpenPokeMessage)
	out.ClosedPokeMessage = strings.TrimSpace(in.ClosedPokeMessage)
	out.JoinOpenPokeMessage = strings.TrimSpace(in.JoinOpenPokeMessage)
	out.JoinClosedPokeMessage = strings.TrimSpace(in.JoinClosedPokeMessage)
	out.SupporterPokeMessage = strings.TrimSpace(in.SupporterPokeMessage)
	out.AutoOpenTime = strings.TrimSpace(in.AutoOpenTime)
	out.AutoCloseTime = strings.TrimSpace(in.AutoCloseTime)
	if out.AutoOpenTime == "" {
		out.AutoOpenTime = "08:00"
	}
	if out.AutoCloseTime == "" {
		out.AutoCloseTime = "22:00"
	}
	if out.WaitingAreaChannel < 0 {
		out.WaitingAreaChannel = 0
	}
	return out
}

func normalizeAnnouncementConfig(in AnnouncementConfig) AnnouncementConfig {
	out := in
	out.Message = strings.TrimSpace(in.Message)
	out.ScheduleMode = strings.TrimSpace(in.ScheduleMode)
	out.RepeatTime = strings.TrimSpace(in.RepeatTime)
	if out.ScheduleMode != "time" && out.ScheduleMode != "once" {
		out.ScheduleMode = "interval"
	}
	if out.RepeatIntervalMinutes < 10 {
		out.RepeatIntervalMinutes = 10
	}
	if out.RepeatIntervalCount < 1 {
		out.RepeatIntervalCount = 1
	}
	if out.RepeatTime == "" {
		out.RepeatTime = "08:00"
	}
	return out
}

func normalizePlatformMode(mode string) string {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case "discord":
		return "discord"
	case "both", "teamspeak_discord", "discord_teamspeak":
		return "both"
	default:
		return "teamspeak"
	}
}

func normalizeFrameworkConfig(in FrameworkConfig) FrameworkConfig {
	out := in
	out.PlatformMode = normalizePlatformMode(in.PlatformMode)
	return out
}

func normalizeStringIDs(values []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		if _, ok := seen[trimmed]; ok {
			continue
		}
		seen[trimmed] = struct{}{}
		out = append(out, trimmed)
	}
	return out
}

func normalizeDiscordConfig(in DiscordConfig) DiscordConfig {
	out := in
	out.BotToken = strings.TrimSpace(in.BotToken)
	out.ApplicationID = strings.TrimSpace(in.ApplicationID)
	out.GuildID = strings.TrimSpace(in.GuildID)
	if out.AFKInactivityMinutes <= 0 {
		out.AFKInactivityMinutes = 30
	}
	out.BotDisplayName = strings.TrimSpace(in.BotDisplayName)
	out.StatusText = strings.TrimSpace(in.StatusText)
	out.CommandPrefix = strings.TrimSpace(in.CommandPrefix)
	out.LogChannelID = strings.TrimSpace(in.LogChannelID)
	out.AnnouncementChannelID = strings.TrimSpace(in.AnnouncementChannelID)
	out.SupportCategoryID = strings.TrimSpace(in.SupportCategoryID)
	out.SupportLogChannelID = strings.TrimSpace(in.SupportLogChannelID)
	out.AdminRoleIDs = normalizeStringIDs(in.AdminRoleIDs)
	out.SupporterRoleIDs = normalizeStringIDs(in.SupporterRoleIDs)
	out.BotRoleIDs = normalizeStringIDs(in.BotRoleIDs)
	if out.CommandPrefix == "" {
		out.CommandPrefix = "!"
	}
	if out.BotDisplayName == "" {
		out.BotDisplayName = "UC-Framework"
	}
	return out
}

func isValidClockTime(v string) bool {
	_, err := time.Parse("15:04", v)
	return err == nil
}

func getLivePluginConfigRaw(name string) (json.RawMessage, bool) {
	canon := canonicalPluginConfigKey(name)
	if raw, ok := liveConfigs[canon]; ok {
		return raw, true
	}
	// Legacy fallback for historic mixed-case keys.
	for k, raw := range liveConfigs {
		if strings.EqualFold(k, name) {
			return raw, true
		}
	}
	return nil, false
}

func normalizePluginConfigMap(in map[string]json.RawMessage) map[string]json.RawMessage {
	out := make(map[string]json.RawMessage, len(in))
	for k, v := range in {
		canon := canonicalPluginConfigKey(k)
		// Canonicalize all plugin config keys to avoid duplicate legacy variants.
		out[canon] = v
	}

	return out
}

func Start() {
	setupProcessLogging()
	loadConfig()
	logInternetSecurityStartupChecks(&config)
	loadLanguage(config.DefaultLanguage)
	api.ConfigureAuth(api.AuthConfig{
		Enabled:                 config.WebAuth.Enabled,
		Provider:                config.WebAuth.Provider,
		AllowRanksystemFallback: config.WebAuth.AllowRanksystemFallback,
		Username:                config.WebAuth.Username,
		Password:                config.WebAuth.Password,
		PasswordHash:            config.WebAuth.PasswordHash,
		ForcePasswordChange:     config.WebAuth.ForcePasswordChange,
		SessionTTLMinutes:       config.WebAuth.SessionTTLMinutes,
		Ranksystem: api.RanksystemAuthConfig{
			LoginURL:      config.WebAuth.Ranksystem.LoginURL,
			UsernameField: config.WebAuth.Ranksystem.UsernameField,
			PasswordField: config.WebAuth.Ranksystem.PasswordField,
			APIKeyHeader:  config.WebAuth.Ranksystem.APIKeyHeader,
			APIKeyValue:   config.WebAuth.Ranksystem.APIKeyValue,
			BearerToken:   config.WebAuth.Ranksystem.BearerToken,
		},
	})

	supervisorScript := filepath.Join("scripts", "bot-supervisor.sh")
	watchdogScript := filepath.Join("scripts", "bot-watchdog.sh")
	if runtime.GOOS == "windows" {
		supervisorScript = filepath.Join("scripts", "bot-supervisor.ps1")
		watchdogScript = filepath.Join("scripts", "bot-watchdog.ps1")
	}

	api.ConfigureBotControl(api.BotControlConfig{
		Enabled:                config.BotControl.Enabled,
		SupervisorScript:       supervisorScript,
		WatchdogScript:         watchdogScript,
		BotExecutable:          config.BotControl.BotExecutable,
		WorkingDir:             config.BotControl.WorkingDir,
		BotArgs:                config.BotControl.BotArgs,
		StateFile:              filepath.Join("runtime", "bot", "state.json"),
		PidFile:                filepath.Join("runtime", "bot", "bot.pid"),
		WatchdogPidFile:        filepath.Join("runtime", "bot", "watchdog.pid"),
		LogFile:                filepath.Join("runtime", "bot", "watchdog.log"),
		WatchdogMinIntervalSec: config.BotControl.WatchdogMinIntervalSec,
		WatchdogMaxIntervalSec: config.BotControl.WatchdogMaxIntervalSec,
	})

	// Live-Config-Store aus config.json initialisieren
	liveConfigsMu.Lock()
	liveConfigs = make(map[string]json.RawMessage)
	for k, v := range config.PluginConfigs {
		liveConfigs[k] = v
	}
	liveConfigsMu.Unlock()

	// Start Webserver parallel
	go web.StartWebServer()

	// Initialize Dispatcher and plugin registry (runtime starts on demand via bot control actions).
	dispatcher := bot.NewDispatcher()
	pluginRegistry := bot.NewPluginRegistry(dispatcher)

	pluginRegistry.RegisterFactory("AdminCounter", func() bot.Plugin {
		return admincounter.New(loadAdminCounterConfig())
	})
	pluginRegistry.RegisterFactory("MemberCounter", func() bot.Plugin {
		return membercounter.New(loadMemberCounterConfig())
	})
	pluginRegistry.RegisterFactory("CombinedStats", func() bot.Plugin {
		return combinedstats.New()
	})
	pluginRegistry.RegisterFactory("AfkMover", func() bot.Plugin {
		return afkmover.New(loadAfkMoverConfig())
	})

	// Plugin-Registry für API verfügbar machen
	api.PluginRegistry = pluginRegistry

	// Gespeicherte Plugin-Aktivierungszustände laden
	pluginEnabledStatesMu.Lock()
	pluginEnabledStates = loadPluginStates()
	pluginEnabledStatesMu.Unlock()

	// Save-Funktion fÃ¼r API bereitstellen
	api.SavePluginConfigFunc = func(name string, raw []byte) error {
		key := canonicalPluginConfigKey(name)
		liveConfigsMu.Lock()
		for existing := range liveConfigs {
			if strings.EqualFold(existing, name) && canonicalPluginConfigKey(existing) == key && existing != key {
				delete(liveConfigs, existing)
			}
		}
		liveConfigs[key] = json.RawMessage(raw)
		liveConfigsMu.Unlock()
		return savePluginConfigs()
	}

	api.GetSavedPluginConfigFunc = func(name string) ([]byte, error) {
		liveConfigsMu.RLock()
		defer liveConfigsMu.RUnlock()
		raw, ok := getLivePluginConfigRaw(name)
		if !ok || raw == nil {
			return nil, fmt.Errorf("plugin %s config not found", name)
		}
		return raw, nil
	}

	runtimeMu := &sync.Mutex{}
	var runtimeTS3Client *ts3.TS3Client
	var runtimeDiscordClient *discord.Client
	var runtimeMonitorStop chan struct{}
	var runtimeWatchdogStop chan struct{}
	previousSupportChannelsByClient := map[int]int{}
	supportJoinObserverInitialized := false
	var supportLastAutoOpenKey string
	var supportLastAutoCloseKey string
	supportOpen := false
	supportLastAction := ""
	supportLastError := ""
	var announcementLastSentTime *time.Time
	announcementLastIntervalKey := ""
	announcementIntervalSentCount := 0
	announcementLastSettingsSig := ""
	announcementOncePending := false
	var stopRuntime func() error
	botRunning := false
	desiredRunning := false
	watchdogRunning := false
	lastAction := ""
	lastError := ""

	applySupportStateLocked := func(open bool) error {
		cfg := normalizeSupportConfig(config.Support)
		config.Support = cfg
		if len(cfg.SupportChannelIDs) == 0 {
			return fmt.Errorf("no support channels configured")
		}

		const (
			openJoinPower        = 0
			openSubscribePower   = 0
			closedJoinPower      = 75
			closedSubscribePower = 75
		)

		client := runtimeTS3Client
		if !botRunning || client == nil || !client.IsConnected() {
			return fmt.Errorf("TS3 bot is not running or connected")
		}

		neededJoinPower := closedJoinPower
		neededSubscribePower := closedSubscribePower
		if open {
			neededJoinPower = openJoinPower
			neededSubscribePower = openSubscribePower
		}

		for _, channelID := range cfg.SupportChannelIDs {
			if err := client.SetChannelAccess(channelID, neededJoinPower, neededSubscribePower); err != nil {
				return fmt.Errorf("set support channel access %d failed: %w", channelID, err)
			}
		}

		supportSet := map[int]struct{}{}
		for _, channelID := range cfg.SupportChannelIDs {
			supportSet[channelID] = struct{}{}
		}

		clients, err := client.ListClients()
		if err != nil {
			return err
		}

		statusMessage := cfg.ClosedPokeMessage
		if open {
			statusMessage = cfg.OpenPokeMessage
		}
		statusMessage = strings.TrimSpace(statusMessage)
		if statusMessage != "" {
			if err := client.SendServerMessage(statusMessage); err != nil {
				log.Printf("[WARN] [support] server message for action=%s failed: %v", map[bool]string{true: "open", false: "close"}[open], err)
			}
		}

		if !open && cfg.WaitingAreaChannel > 0 {
			for _, entry := range clients {
				if entry.IsQuery {
					continue
				}
				if _, ok := supportSet[entry.ChannelID]; !ok {
					continue
				}
				if err := client.MoveClient(entry.ID, cfg.WaitingAreaChannel); err != nil {
					log.Printf("[WARN] [support] move to waiting area failed for client=%d target_channel=%d: %v", entry.ID, cfg.WaitingAreaChannel, err)
				}
			}
		}

		supportOpen = open
		if open {
			supportLastAction = "open"
		} else {
			supportLastAction = "close"
		}
		supportLastError = ""
		return nil
	}

	supportStatus := func() api.SupportStatus {
		cfg := normalizeSupportConfig(config.Support)
		config.Support = cfg
		return api.SupportStatus{
			Open:                supportOpen,
			LastAction:          supportLastAction,
			LastError:           supportLastError,
			AutoScheduleEnabled: cfg.AutoScheduleEnabled,
			AutoOpenTime:        cfg.AutoOpenTime,
			AutoCloseTime:       cfg.AutoCloseTime,
		}
	}

	botStatus := func() map[string]interface{} {
		pid := 0
		if botRunning {
			pid = os.Getpid()
		}
		return map[string]interface{}{
			"ok":              true,
			"running":         botRunning,
			"pid":             pid,
			"desired_running": desiredRunning,
			"last_action":     lastAction,
			"last_error":      lastError,
		}
	}

	startRuntime := func() error {
		if botRunning {
			lastError = ""
			return nil
		}

		// Reload configuration on each start so config/script changes become active.
		loadConfig()
		loadLanguage(config.DefaultLanguage)

		liveConfigsMu.Lock()
		liveConfigs = make(map[string]json.RawMessage)
		for k, v := range config.PluginConfigs {
			liveConfigs[k] = v
		}
		liveConfigsMu.Unlock()

		api.ConfigureAuth(api.AuthConfig{
			Enabled:                 config.WebAuth.Enabled,
			Provider:                config.WebAuth.Provider,
			AllowRanksystemFallback: config.WebAuth.AllowRanksystemFallback,
			Username:                config.WebAuth.Username,
			Password:                config.WebAuth.Password,
			PasswordHash:            config.WebAuth.PasswordHash,
			ForcePasswordChange:     config.WebAuth.ForcePasswordChange,
			SessionTTLMinutes:       config.WebAuth.SessionTTLMinutes,
			Ranksystem: api.RanksystemAuthConfig{
				LoginURL:      config.WebAuth.Ranksystem.LoginURL,
				UsernameField: config.WebAuth.Ranksystem.UsernameField,
				PasswordField: config.WebAuth.Ranksystem.PasswordField,
				APIKeyHeader:  config.WebAuth.Ranksystem.APIKeyHeader,
				APIKeyValue:   config.WebAuth.Ranksystem.APIKeyValue,
				BearerToken:   config.WebAuth.Ranksystem.BearerToken,
			},
		})

		platformMode := normalizePlatformMode(config.Framework.PlatformMode)
		startedAny := false

		if platformMode == "teamspeak" || platformMode == "both" {
			client := ts3.NewTS3Client(
				config.TS3.Host,
				config.TS3.QueryPort,
				config.TS3.VoicePort,
				config.TS3.Username,
				config.TS3.Password,
			)
			if err := client.Connect(); err != nil {
				lastError = err.Error()
				return err
			}
			if strings.TrimSpace(config.TS3.BotNickname) != "" {
				if err := client.SetBotNickname(config.TS3.BotNickname); err != nil {
					if strings.Contains(err.Error(), "ts3 error id=513") {
						log.Printf("[WARN] [ts3] configured bot nickname already in use; starting anyway with current query nickname: %v", err)
					} else {
						client.MarkDisconnected(err)
						lastError = err.Error()
						return fmt.Errorf("failed to set bot nickname: %w", err)
					}
				}
			}
			client.SetDispatcher(dispatcher)

			for _, name := range []string{"AdminCounter", "MemberCounter", "CombinedStats", "AfkMover"} {
				pluginEnabledStatesMu.RLock()
				enabled, exists := pluginEnabledStates[name]
				pluginEnabledStatesMu.RUnlock()
				if exists && !enabled {
					continue
				}
				if err := pluginRegistry.Load(name); err != nil {
					client.MarkDisconnected(err)
					lastError = err.Error()
					return fmt.Errorf("failed to load plugin %s: %w", name, err)
				}
			}

			monitorStop := make(chan struct{})
			runtimeTS3Client = client
			runtimeMonitorStop = monitorStop
			startedAny = true

			log.Printf("[INFO] [bot] teamspeak started: host=%s query_port=%d voice_port=%d slowmode_ms=%d nickname=%s default_channel=%s",
				config.TS3.Host,
				config.TS3.QueryPort,
				config.TS3.VoicePort,
				config.TS3.QuerySlowmode,
				config.TS3.BotNickname,
				config.TS3.DefaultChannel,
			)

			go client.ListenEvents()
			go monitorTS3Connection(client, monitorStop, func(err error, attempts int) {
				runtimeMu.Lock()
				defer runtimeMu.Unlock()

				failureMessage := fmt.Sprintf("TS3 reconnect failed %d times: %v", attempts, err)
				if stopErr := stopRuntime(); stopErr != nil {
					log.Printf("[ERROR] [ts3] forced stop after reconnect failures failed: %v", stopErr)
				}
				lastError = failureMessage
				log.Printf("[ERROR] [ts3] stopping bot after %d failed reconnect attempts", attempts)
			})
		}

		if platformMode == "discord" || platformMode == "both" {
			discordClient := discord.NewClient(config.Discord.BotToken, config.Discord.GuildID)
			discordClient.SetAFKKickConfig(config.Discord.AFKKickEnabled, config.Discord.AFKInactivityMinutes)
			discordClient.SetLogger(log.Printf)
			if err := discordClient.Connect(); err != nil {
				if runtimeTS3Client != nil {
					runtimeTS3Client.MarkDisconnected(nil)
					runtimeTS3Client = nil
				}
				if runtimeMonitorStop != nil {
					close(runtimeMonitorStop)
					runtimeMonitorStop = nil
				}
				lastError = err.Error()
				return fmt.Errorf("failed to connect discord: %w", err)
			}
			runtimeDiscordClient = discordClient
			startedAny = true
			log.Printf("[INFO] [bot] discord started: guild_id=%s", config.Discord.GuildID)
			if config.Discord.AFKKickEnabled {
				log.Printf("[INFO] [discord] AFK inactivity kick active: timeout=%d min", config.Discord.AFKInactivityMinutes)
			}
		}

		if !startedAny {
			lastError = "no platform runtime configured"
			return fmt.Errorf("no platform runtime configured")
		}

		botRunning = true
		lastError = ""
		return nil
	}

	stopRuntime = func() error {
		if !botRunning {
			lastError = ""
			return nil
		}

		states := pluginRegistry.All()
		for _, name := range []string{"AfkMover", "CombinedStats", "MemberCounter", "AdminCounter"} {
			st, ok := states[name]
			if !ok || !st.Loaded {
				continue
			}
			if err := pluginRegistry.Unload(name); err != nil {
				log.Printf("[ERROR] [plugins] unload failed %s: %v", name, err)
			}
		}

		if runtimeMonitorStop != nil {
			close(runtimeMonitorStop)
			runtimeMonitorStop = nil
		}
		if runtimeTS3Client != nil {
			runtimeTS3Client.MarkDisconnected(nil)
			runtimeTS3Client = nil
		}
		if runtimeDiscordClient != nil {
			if err := runtimeDiscordClient.Close(); err != nil {
				log.Printf("[ERROR] [discord] close failed: %v", err)
			}
			runtimeDiscordClient = nil
		}

		botRunning = false
		lastError = ""
		return nil
	}

	watchdogStatus := func() map[string]interface{} {
		pid := 0
		if watchdogRunning {
			pid = os.Getpid()
		}
		return map[string]interface{}{
			"ok":      true,
			"running": watchdogRunning,
			"pid":     pid,
			"mode":    "runtime",
		}
	}

	startRuntimeWatchdog := func() error {
		if watchdogRunning {
			return nil
		}
		stopCh := make(chan struct{})
		runtimeWatchdogStop = stopCh
		watchdogRunning = true

		go func() {
			ticker := time.NewTicker(10 * time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-stopCh:
					return
				case <-ticker.C:
				}

				runtimeMu.Lock()
				if watchdogRunning && desiredRunning && !botRunning {
					if err := startRuntime(); err != nil {
						lastError = err.Error()
						log.Printf("[ERROR] [watchdog] runtime restart failed: %v", err)
					}
				}
				runtimeMu.Unlock()
			}
		}()

		log.Printf("[INFO] [watchdog] started: mode=runtime interval=10s")
		return nil
	}

	stopRuntimeWatchdog := func() error {
		if !watchdogRunning {
			return nil
		}
		watchdogRunning = false
		if runtimeWatchdogStop != nil {
			close(runtimeWatchdogStop)
			runtimeWatchdogStop = nil
		}
		return nil
	}

	api.BotStatusProvider = func() (map[string]interface{}, error) {
		runtimeMu.Lock()
		defer runtimeMu.Unlock()
		return botStatus(), nil
	}

	api.BotActionExecutor = func(action string) (map[string]interface{}, error) {
		runtimeMu.Lock()
		defer runtimeMu.Unlock()

		switch action {
		case "start":
			desiredRunning = true
			lastAction = "start"
			if err := startRuntime(); err != nil {
				return botStatus(), err
			}
		case "stop":
			desiredRunning = false
			lastAction = "stop"
			if err := stopRuntime(); err != nil {
				return botStatus(), err
			}
		case "restart":
			desiredRunning = true
			lastAction = "restart"
			if err := stopRuntime(); err != nil {
				return botStatus(), err
			}
			if err := startRuntime(); err != nil {
				return botStatus(), err
			}
		default:
			return botStatus(), fmt.Errorf("invalid action")
		}

		return botStatus(), nil
	}

	api.WatchdogStatusProvider = func() (map[string]interface{}, error) {
		runtimeMu.Lock()
		defer runtimeMu.Unlock()
		return watchdogStatus(), nil
	}

	api.WatchdogActionExecutor = func(action string) (map[string]interface{}, error) {
		runtimeMu.Lock()
		defer runtimeMu.Unlock()

		switch action {
		case "start":
			if err := startRuntimeWatchdog(); err != nil {
				return watchdogStatus(), err
			}
		case "stop":
			if err := stopRuntimeWatchdog(); err != nil {
				return watchdogStatus(), err
			}
		default:
			return watchdogStatus(), fmt.Errorf("invalid action")
		}
		return watchdogStatus(), nil
	}

	api.SyncPluginEnabledStateFunc = func(name string, active bool) error {
		// Plugin-State in plugin_states.json persistieren
		pluginEnabledStatesMu.Lock()
		pluginEnabledStates[name] = active
		statesCopy := make(map[string]bool, len(pluginEnabledStates))
		for k, v := range pluginEnabledStates {
			statesCopy[k] = v
		}
		pluginEnabledStatesMu.Unlock()

		if err := savePluginStates(statesCopy); err != nil {
			log.Printf("[ERROR] [plugins] failed to save plugin states: %v", err)
		}

		// AfkMover: zusätzlich Enabled-Feld im Plugin-Config aktualisieren (Runtime-Verhalten)
		if strings.ToLower(strings.TrimSpace(name)) == "afkmover" {
			liveConfigsMu.Lock()
			defer liveConfigsMu.Unlock()

			key := canonicalPluginConfigKey("afkmover")
			cfg := afkmover.AfkMoverConfig{TimeoutSeconds: 600, Enabled: true, ExcludedChannels: []int{}}
			if raw, ok := getLivePluginConfigRaw("afkmover"); ok {
				_ = json.Unmarshal(raw, &cfg)
			}
			if cfg.ExcludedChannels == nil {
				cfg.ExcludedChannels = []int{}
			}
			cfg.Enabled = active

			raw, err := json.Marshal(cfg)
			if err != nil {
				return err
			}
			liveConfigs[key] = json.RawMessage(raw)

			config.PluginConfigs = make(map[string]json.RawMessage, len(liveConfigs))
			for k, v := range liveConfigs {
				config.PluginConfigs[k] = v
			}
			return saveConfig()
		}

		return nil
	}

	api.GetPluginEnabledStateFunc = func(name string) (bool, bool) {
		pluginEnabledStatesMu.RLock()
		defer pluginEnabledStatesMu.RUnlock()
		active, ok := pluginEnabledStates[name]
		return active, ok
	}

	api.SaveWebAuthPasswordHashFunc = func(passwordHash string) error {
		config.WebAuth.PasswordHash = passwordHash
		config.WebAuth.Password = ""
		return saveConfig()
	}

	api.SaveWebAuthForcePasswordChangeFunc = func(force bool) error {
		config.WebAuth.ForcePasswordChange = force
		return saveConfig()
	}

	api.SaveWebAuthUsernameFunc = func(username string) error {
		config.WebAuth.Username = username
		return saveConfig()
	}

	api.SaveWebAuthEnabledFunc = func(enabled bool) error {
		config.WebAuth.Enabled = enabled
		return saveConfig()
	}

	api.SaveWebAuthProviderFunc = func(provider string) error {
		config.WebAuth.Provider = provider
		return saveConfig()
	}

	api.SaveWebAuthAllowRanksystemFallbackFunc = func(enabled bool) error {
		config.WebAuth.AllowRanksystemFallback = enabled
		return saveConfig()
	}

	api.GetFrameworkInfoFunc = func() (api.FrameworkInfo, error) {
		updatedAt := ""
		if exePath, err := os.Executable(); err == nil {
			if st, err := os.Stat(exePath); err == nil {
				updatedAt = st.ModTime().Format(time.RFC3339)
			}
		}

		return api.FrameworkInfo{
			Name:          "UC-Framework",
			Version:       api.FrameworkVersion,
			UpdatedAt:     updatedAt,
			LatestVersion: api.FrameworkVersion,
			IsLatest:      true,
		}, nil
	}

	api.RestartFrameworkFunc = func() error {
		exePath, err := os.Executable()
		if err != nil {
			return err
		}

		workDir, err := os.Getwd()
		if err != nil {
			workDir = filepath.Dir(exePath)
		}

		args := append([]string(nil), os.Args[1:]...)
		if runtime.GOOS == "windows" {
			psQuote := func(s string) string {
				return strings.ReplaceAll(s, "'", "''")
			}

			argList := make([]string, 0, len(args))
			for _, arg := range args {
				argList = append(argList, fmt.Sprintf("'%s'", psQuote(arg)))
			}

			script := fmt.Sprintf(
				"$ErrorActionPreference='Stop'; $exe='%s'; $wd='%s'; $args=@(%s); Start-Sleep -Milliseconds 700; Start-Process -FilePath $exe -ArgumentList $args -WorkingDirectory $wd -WindowStyle Hidden",
				psQuote(exePath),
				psQuote(workDir),
				strings.Join(argList, ", "),
			)

			cmd := exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", script)
			if err := cmd.Start(); err != nil {
				return err
			}
		} else {
			shellQuote := func(s string) string {
				return "'" + strings.ReplaceAll(s, "'", "'\"'\"'") + "'"
			}

			parts := make([]string, 0, len(args)+1)
			parts = append(parts, shellQuote(exePath))
			for _, arg := range args {
				parts = append(parts, shellQuote(arg))
			}

			script := fmt.Sprintf("sleep 0.7; cd %s; %s >/dev/null 2>&1 &", shellQuote(workDir), strings.Join(parts, " "))
			cmd := exec.Command("sh", "-c", script)
			if err := cmd.Start(); err != nil {
				return err
			}
		}

		go func() {
			time.Sleep(250 * time.Millisecond)
			os.Exit(0)
		}()

		return nil
	}

	api.SendServerAnnouncementFunc = func(message string) error {
		runtimeMu.Lock()
		client := runtimeTS3Client
		running := botRunning
		runtimeMu.Unlock()

		if !running || client == nil || !client.IsConnected() {
			return fmt.Errorf("TS3 bot is not running or connected")
		}

		return client.SendServerMessage(message)
	}

	api.GetLanguageSettingsFunc = func() (string, []string, error) {
		return config.DefaultLanguage, append([]string(nil), config.SupportedLanguages...), nil
	}

	api.SaveLanguageSettingsFunc = func(language string) error {
		normalized := normalizeLanguageCode(language)
		if !isSupportedLanguage(normalized, config.SupportedLanguages) {
			return fmt.Errorf("unsupported language: %s", normalized)
		}
		config.DefaultLanguage = normalized
		loadLanguage(normalized)
		return saveConfig()
	}

	api.GetFrameworkSettingsFunc = func() (api.FrameworkSettings, error) {
		cfg := normalizeFrameworkConfig(config.Framework)
		config.Framework = cfg
		return api.FrameworkSettings{
			PlatformMode: cfg.PlatformMode,
		}, nil
	}

	api.SaveFrameworkSettingsFunc = func(s api.FrameworkSettings) error {
		config.Framework = normalizeFrameworkConfig(FrameworkConfig{
			PlatformMode: s.PlatformMode,
		})
		return saveConfig()
	}

	api.GetDiscordSettingsFunc = func() (api.DiscordSettings, error) {
		cfg := normalizeDiscordConfig(config.Discord)
		config.Discord = cfg
		return api.DiscordSettings{
			Enabled:               cfg.Enabled,
			BotToken:              cfg.BotToken,
			ApplicationID:         cfg.ApplicationID,
			GuildID:               cfg.GuildID,
			AFKKickEnabled:        cfg.AFKKickEnabled,
			AFKInactivityMinutes:  cfg.AFKInactivityMinutes,
			BotDisplayName:        cfg.BotDisplayName,
			StatusText:            cfg.StatusText,
			CommandPrefix:         cfg.CommandPrefix,
			LogChannelID:          cfg.LogChannelID,
			AnnouncementChannelID: cfg.AnnouncementChannelID,
			SupportCategoryID:     cfg.SupportCategoryID,
			SupportLogChannelID:   cfg.SupportLogChannelID,
			AdminRoleIDs:          append([]string(nil), cfg.AdminRoleIDs...),
			SupporterRoleIDs:      append([]string(nil), cfg.SupporterRoleIDs...),
			BotRoleIDs:            append([]string(nil), cfg.BotRoleIDs...),
		}, nil
	}

	api.SaveDiscordSettingsFunc = func(s api.DiscordSettings) error {
		config.Discord = normalizeDiscordConfig(DiscordConfig{
			Enabled:               s.Enabled,
			BotToken:              s.BotToken,
			ApplicationID:         s.ApplicationID,
			GuildID:               s.GuildID,
			AFKKickEnabled:        s.AFKKickEnabled,
			AFKInactivityMinutes:  s.AFKInactivityMinutes,
			BotDisplayName:        s.BotDisplayName,
			StatusText:            s.StatusText,
			CommandPrefix:         s.CommandPrefix,
			LogChannelID:          s.LogChannelID,
			AnnouncementChannelID: s.AnnouncementChannelID,
			SupportCategoryID:     s.SupportCategoryID,
			SupportLogChannelID:   s.SupportLogChannelID,
			AdminRoleIDs:          append([]string(nil), s.AdminRoleIDs...),
			SupporterRoleIDs:      append([]string(nil), s.SupporterRoleIDs...),
			BotRoleIDs:            append([]string(nil), s.BotRoleIDs...),
		})
		return saveConfig()
	}

	getDiscordMetaClient := func() (*discord.Client, func(), error) {
		runtimeMu.Lock()
		client := runtimeDiscordClient
		runtimeMu.Unlock()
		if client != nil && client.IsConnected() {
			return client, func() {}, nil
		}

		cfg := normalizeDiscordConfig(config.Discord)
		temp := discord.NewClient(cfg.BotToken, cfg.GuildID)
		if err := temp.Connect(); err != nil {
			return nil, nil, err
		}
		return temp, func() {
			_ = temp.Close()
		}, nil
	}

	api.GetDiscordChannelsFunc = func() ([]api.DiscordChannel, error) {
		client, cleanup, err := getDiscordMetaClient()
		if err != nil {
			return nil, err
		}
		defer cleanup()
		channels, err := client.ListChannels()
		if err != nil {
			return nil, err
		}
		out := make([]api.DiscordChannel, 0, len(channels))
		for _, ch := range channels {
			channelType := "unknown"
			switch ch.Type {
			case 0:
				channelType = "text"
			case 2:
				channelType = "voice"
			case 4:
				channelType = "category"
			case 5:
				channelType = "announcement"
			case 13:
				channelType = "stage"
			case 15:
				channelType = "forum"
			}
			out = append(out, api.DiscordChannel{ID: ch.ID, Name: ch.Name, Type: channelType, ParentID: ch.ParentID})
		}
		return out, nil
	}

	api.GetDiscordRolesFunc = func() ([]api.DiscordRole, error) {
		client, cleanup, err := getDiscordMetaClient()
		if err != nil {
			return nil, err
		}
		defer cleanup()
		roles, err := client.ListRoles()
		if err != nil {
			return nil, err
		}
		out := make([]api.DiscordRole, 0, len(roles))
		for _, role := range roles {
			out = append(out, api.DiscordRole{ID: role.ID, Name: role.Name})
		}
		return out, nil
	}

	api.GetTS3SettingsFunc = func() api.TS3Settings {
		return api.TS3Settings{
			Host:           config.TS3.Host,
			QueryPort:      config.TS3.QueryPort,
			VoicePort:      config.TS3.VoicePort,
			QueryUsername:  config.TS3.Username,
			QueryPassword:  config.TS3.Password,
			BotNickname:    config.TS3.BotNickname,
			DefaultChannel: config.TS3.DefaultChannel,
			QuerySlowmode:  config.TS3.QuerySlowmode,
		}
	}

	api.GetTS3ChannelsFunc = func() ([]api.TS3Channel, error) {
		runtimeMu.Lock()
		client := runtimeTS3Client
		runtimeMu.Unlock()
		if client == nil || !client.IsConnected() {
			return nil, fmt.Errorf("TS3 bot is not running")
		}

		channels, err := client.ListChannels()
		if err != nil {
			return nil, err
		}

		out := make([]api.TS3Channel, 0, len(channels))
		for _, ch := range channels {
			out = append(out, api.TS3Channel{ID: ch.ID, Name: ch.Name})
		}
		return out, nil
	}

	api.GetTS3ServerGroupsFunc = func() ([]api.TS3ServerGroup, error) {
		runtimeMu.Lock()
		client := runtimeTS3Client
		runtimeMu.Unlock()
		if client == nil || !client.IsConnected() {
			return nil, fmt.Errorf("TS3 bot is not running")
		}

		groups, err := client.ListServerGroups()
		if err != nil {
			return nil, err
		}

		out := make([]api.TS3ServerGroup, 0, len(groups))
		for _, g := range groups {
			out = append(out, api.TS3ServerGroup{ID: g.ID, Name: g.Name})
		}
		return out, nil
	}

	api.GetTS3ConnectionStatusFunc = func() api.TS3ConnectionStatus {
		runtimeMu.Lock()
		client := runtimeTS3Client
		running := botRunning
		runtimeMu.Unlock()
		connected := false
		if client != nil {
			connected = client.IsConnected()
		}
		status := api.TS3ConnectionStatus{
			BotRunning:     running,
			Connected:      connected,
			Host:           config.TS3.Host,
			Port:           config.TS3.QueryPort,
			LastCheckAt:    time.Now().Format(time.RFC3339),
			Implementation: "serverquery",
		}
		if !running {
			status.LastError = "bot disabled"
		}
		return status
	}

	api.GetSupportSettingsFunc = func() (api.SupportSettings, error) {
		cfg := normalizeSupportConfig(config.Support)
		config.Support = cfg
		return api.SupportSettings{
			Enabled:               cfg.Enabled,
			SupportChannelIDs:     append([]int(nil), cfg.SupportChannelIDs...),
			WaitingAreaChannel:    cfg.WaitingAreaChannel,
			OpenPokeMessage:       cfg.OpenPokeMessage,
			ClosedPokeMessage:     cfg.ClosedPokeMessage,
			JoinOpenPokeMessage:   cfg.JoinOpenPokeMessage,
			JoinClosedPokeMessage: cfg.JoinClosedPokeMessage,
			SupporterPokeMessage:  cfg.SupporterPokeMessage,
			SupporterGroupIDs:     append([]int(nil), cfg.SupporterGroupIDs...),
			AutoScheduleEnabled:   cfg.AutoScheduleEnabled,
			AutoOpenTime:          cfg.AutoOpenTime,
			AutoCloseTime:         cfg.AutoCloseTime,
		}, nil
	}

	api.SaveSupportSettingsFunc = func(s api.SupportSettings) error {
		next := normalizeSupportConfig(SupportConfig{
			Enabled:               s.Enabled,
			SupportChannelIDs:     append([]int(nil), s.SupportChannelIDs...),
			WaitingAreaChannel:    s.WaitingAreaChannel,
			OpenPokeMessage:       s.OpenPokeMessage,
			ClosedPokeMessage:     s.ClosedPokeMessage,
			JoinOpenPokeMessage:   s.JoinOpenPokeMessage,
			JoinClosedPokeMessage: s.JoinClosedPokeMessage,
			SupporterPokeMessage:  s.SupporterPokeMessage,
			SupporterGroupIDs:     append([]int(nil), s.SupporterGroupIDs...),
			AutoScheduleEnabled:   s.AutoScheduleEnabled,
			AutoOpenTime:          s.AutoOpenTime,
			AutoCloseTime:         s.AutoCloseTime,
		})

		if len(next.SupportChannelIDs) == 0 {
			return fmt.Errorf("at least one support channel is required")
		}
		if next.AutoScheduleEnabled {
			if !isValidClockTime(next.AutoOpenTime) || !isValidClockTime(next.AutoCloseTime) {
				return fmt.Errorf("auto open and close time must use HH:MM format")
			}
		}

		config.Support = next
		return saveConfig()
	}

	api.GetSupportStatusFunc = func() (api.SupportStatus, error) {
		runtimeMu.Lock()
		defer runtimeMu.Unlock()
		return supportStatus(), nil
	}

	api.ExecuteSupportActionFunc = func(action string) (api.SupportStatus, error) {
		runtimeMu.Lock()
		defer runtimeMu.Unlock()

		switch strings.ToLower(strings.TrimSpace(action)) {
		case "open":
			if err := applySupportStateLocked(true); err != nil {
				supportLastError = err.Error()
				return supportStatus(), err
			}
		case "close":
			if err := applySupportStateLocked(false); err != nil {
				supportLastError = err.Error()
				return supportStatus(), err
			}
		default:
			return supportStatus(), fmt.Errorf("invalid action")
		}

		return supportStatus(), nil
	}

	api.ExtraPluginsFunc = func() []api.PluginStatus {
		cfg := normalizeSupportConfig(config.Support)
		return []api.PluginStatus{{
			Name:   "SupportControl",
			Active: cfg.Enabled,
		}}
	}

	api.ToggleExtraPluginFunc = func(name string, active bool) error {
		if !strings.EqualFold(name, "SupportControl") {
			return fmt.Errorf("unknown virtual plugin: %s", name)
		}
		runtimeMu.Lock()
		defer runtimeMu.Unlock()
		cfg := normalizeSupportConfig(config.Support)
		cfg.Enabled = active
		config.Support = cfg
		return saveConfig()
	}

	api.GetAnnouncementSettingsFunc = func() (api.AnnouncementSettings, error) {
		cfg := normalizeAnnouncementConfig(config.Announcement)
		config.Announcement = cfg
		return api.AnnouncementSettings{
			Message:               cfg.Message,
			RepeatEnabled:         cfg.RepeatEnabled,
			ScheduleMode:          cfg.ScheduleMode,
			RepeatIntervalMinutes: cfg.RepeatIntervalMinutes,
			RepeatIntervalCount:   cfg.RepeatIntervalCount,
			RepeatTime:            cfg.RepeatTime,
		}, nil
	}

	api.SaveAnnouncementSettingsFunc = func(s api.AnnouncementSettings) error {
		next := normalizeAnnouncementConfig(AnnouncementConfig{
			Message:               s.Message,
			RepeatEnabled:         s.RepeatEnabled,
			ScheduleMode:          s.ScheduleMode,
			RepeatIntervalMinutes: s.RepeatIntervalMinutes,
			RepeatIntervalCount:   s.RepeatIntervalCount,
			RepeatTime:            s.RepeatTime,
		})

		if len(strings.TrimSpace(next.Message)) == 0 {
			return fmt.Errorf("announcement message cannot be empty")
		}
		if next.RepeatEnabled {
			if next.ScheduleMode == "time" && !isValidClockTime(next.RepeatTime) {
				return fmt.Errorf("repeat time must use HH:MM format")
			}
			if next.ScheduleMode == "interval" && next.RepeatIntervalMinutes < 10 {
				return fmt.Errorf("repeat interval must be at least 10 minutes")
			}
			if next.ScheduleMode == "interval" && next.RepeatIntervalCount < 1 {
				return fmt.Errorf("repeat count must be at least 1")
			}
		}

		config.Announcement = next
		announcementLastSentTime = nil
		announcementIntervalSentCount = 0
		announcementLastSettingsSig = ""
		announcementLastIntervalKey = ""
		announcementOncePending = next.RepeatEnabled && next.ScheduleMode == "once"
		return saveConfig()
	}

	api.GetAnnouncementStatusFunc = func() (api.AnnouncementStatus, error) {
		runtimeMu.Lock()
		defer runtimeMu.Unlock()

		cfg := normalizeAnnouncementConfig(config.Announcement)
		status := api.AnnouncementStatus{
			Message:             cfg.Message,
			RepeatEnabled:       cfg.RepeatEnabled,
			ScheduleMode:        cfg.ScheduleMode,
			LastSentAt:          announcementLastSentTime,
			RepeatIntervalCount: cfg.RepeatIntervalCount,
		}
		return status, nil
	}

	api.SendAnnouncementFunc = func(message string) error {
		runtimeMu.Lock()
		client := runtimeTS3Client
		discordClient := runtimeDiscordClient
		running := botRunning
		runtimeMu.Unlock()

		if !running {
			return fmt.Errorf("bot is not running or connected")
		}

		platformMode := normalizePlatformMode(config.Framework.PlatformMode)
		sent := false
		if (platformMode == "teamspeak" || platformMode == "both") && client != nil && client.IsConnected() {
			if err := client.SendServerMessage(message); err != nil {
				return err
			}
			sent = true
		}
		if (platformMode == "discord" || platformMode == "both") && discordClient != nil && discordClient.IsConnected() && strings.TrimSpace(config.Discord.AnnouncementChannelID) != "" {
			if err := discordClient.SendMessage(config.Discord.AnnouncementChannelID, message); err != nil {
				return err
			}
			sent = true
		}
		if !sent {
			return fmt.Errorf("no active platform available for announcements")
		}

		runtimeMu.Lock()
		now := time.Now()
		announcementLastSentTime = &now
		runtimeMu.Unlock()

		return nil
	}

	api.TestTS3ConnectionFunc = func() (api.TS3ConnectionStatus, error) {
		runtimeMu.Lock()
		client := runtimeTS3Client
		running := botRunning
		runtimeMu.Unlock()
		connected := false
		if client != nil {
			connected = client.IsConnected()
		}
		status := api.TS3ConnectionStatus{
			BotRunning:     running,
			Connected:      connected,
			Host:           config.TS3.Host,
			Port:           config.TS3.QueryPort,
			LastCheckAt:    time.Now().Format(time.RFC3339),
			Implementation: "serverquery",
		}

		if !running {
			status.LastError = "bot disabled"
			return status, fmt.Errorf("bot is disabled")
		}

		if client == nil || !status.Connected {
			status.LastError = "not connected"
			return status, fmt.Errorf("TS3 is not connected")
		}

		if _, err := client.ListClients(); err != nil {
			status.Connected = false
			status.LastError = err.Error()
			client.MarkDisconnected(err)
			return status, err
		}

		return status, nil
	}

	api.SaveTS3SettingsFunc = func(s api.TS3Settings) error {
		config.TS3.Host = s.Host
		config.TS3.QueryPort = s.QueryPort
		config.TS3.VoicePort = s.VoicePort
		config.TS3.Username = s.QueryUsername
		config.TS3.Password = s.QueryPassword
		config.TS3.BotNickname = s.BotNickname
		config.TS3.DefaultChannel = s.DefaultChannel
		config.TS3.QuerySlowmode = s.QuerySlowmode
		return saveConfig()
	}

	desiredRunning = false
	lastAction = "startup"
	log.Printf("[INFO] [bot] runtime ready")

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			runtimeMu.Lock()
			cfg := normalizeSupportConfig(config.Support)
			config.Support = cfg
			if !cfg.Enabled || !cfg.AutoScheduleEnabled {
				runtimeMu.Unlock()
				continue
			}

			now := time.Now()
			currentTime := now.Format("15:04")
			day := now.Format("2006-01-02")

			if currentTime == cfg.AutoOpenTime {
				key := day + "-" + cfg.AutoOpenTime
				if supportLastAutoOpenKey != key {
					if err := applySupportStateLocked(true); err != nil {
						supportLastError = err.Error()
					} else {
						supportLastAutoOpenKey = key
					}
				}
			}

			if currentTime == cfg.AutoCloseTime {
				key := day + "-" + cfg.AutoCloseTime
				if supportLastAutoCloseKey != key {
					if err := applySupportStateLocked(false); err != nil {
						supportLastError = err.Error()
					} else {
						supportLastAutoCloseKey = key
					}
				}
			}

			runtimeMu.Unlock()
		}
	}()

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			runtimeMu.Lock()
			cfg := normalizeSupportConfig(config.Support)
			config.Support = cfg
			client := runtimeTS3Client
			if !cfg.Enabled || !botRunning || client == nil || !client.IsConnected() {
				runtimeMu.Unlock()
				continue
			}

			clients, err := client.ListClients()
			if err != nil {
				runtimeMu.Unlock()
				continue
			}

			supportSet := map[int]struct{}{}
			for _, id := range cfg.SupportChannelIDs {
				supportSet[id] = struct{}{}
			}
			if cfg.WaitingAreaChannel > 0 {
				supportSet[cfg.WaitingAreaChannel] = struct{}{}
			}

			nextByClient := map[int]int{}
			entrantClientIDs := make([]int, 0)
			entrantNames := map[int]string{}
			for _, c := range clients {
				if c.IsQuery {
					continue
				}
				nextByClient[c.ID] = c.ChannelID
				_, inSupportNow := supportSet[c.ChannelID]
				if !inSupportNow {
					continue
				}

				prevChannel, hadPrev := previousSupportChannelsByClient[c.ID]
				if !hadPrev || prevChannel != c.ChannelID {
					entrantClientIDs = append(entrantClientIDs, c.ID)
					entrantNames[c.ID] = c.Nickname
				}
			}

			if !supportJoinObserverInitialized {
				previousSupportChannelsByClient = nextByClient
				supportJoinObserverInitialized = true
				runtimeMu.Unlock()
				continue
			}

			previousSupportChannelsByClient = nextByClient

			userJoinMessage := cfg.JoinClosedPokeMessage
			if supportOpen {
				userJoinMessage = cfg.JoinOpenPokeMessage
			}
			userJoinMessage = strings.TrimSpace(userJoinMessage)

			supporterMessage := strings.TrimSpace(cfg.SupporterPokeMessage)
			supporterGroups := map[int]struct{}{}
			for _, groupID := range cfg.SupporterGroupIDs {
				supporterGroups[groupID] = struct{}{}
			}

			for _, entrantID := range entrantClientIDs {
				entrantName := entrantNames[entrantID]
				if userJoinMessage != "" {
					if err := client.PokeClient(entrantID, userJoinMessage); err != nil {
						log.Printf("[WARN] [support] join poke failed for client=%d: %v", entrantID, err)
					}
				}

				if !supportOpen || supporterMessage == "" || len(supporterGroups) == 0 {
					continue
				}

				notice := strings.ReplaceAll(supporterMessage, "{user}", entrantName)
				for _, candidate := range clients {
					if candidate.IsQuery || candidate.ID == entrantID {
						continue
					}
					isSupporter := false
					for _, g := range candidate.ServerGroups {
						if _, ok := supporterGroups[g]; ok {
							isSupporter = true
							break
						}
					}
					if isSupporter {
						if err := client.PokeClient(candidate.ID, notice); err != nil {
							log.Printf("[WARN] [support] supporter poke failed for client=%d entrant=%d: %v", candidate.ID, entrantID, err)
						}
					}
				}
			}

			runtimeMu.Unlock()
		}
	}()

	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			now := time.Now()
			shouldSend := false
			sendMessage := ""
			timeModeKey := ""
			intervalSend := false
			var client *ts3.TS3Client

			runtimeMu.Lock()
			cfg := normalizeAnnouncementConfig(config.Announcement)
			config.Announcement = cfg
			if !cfg.RepeatEnabled || len(strings.TrimSpace(cfg.Message)) == 0 {
				runtimeMu.Unlock()
				continue
			}

			client = runtimeTS3Client
			if !botRunning || client == nil || !client.IsConnected() {
				runtimeMu.Unlock()
				continue
			}

			sig := fmt.Sprintf("%t|%s|%s|%d|%d", cfg.RepeatEnabled, cfg.ScheduleMode, cfg.Message, cfg.RepeatIntervalMinutes, cfg.RepeatIntervalCount)
			if announcementLastSettingsSig != sig {
				announcementLastSettingsSig = sig
				announcementIntervalSentCount = 0
				announcementLastSentTime = nil
				announcementLastIntervalKey = ""
				announcementOncePending = cfg.RepeatEnabled && cfg.ScheduleMode == "once"
			}

			switch cfg.ScheduleMode {
			case "once":
				if announcementOncePending {
					shouldSend = true
					sendMessage = cfg.Message
				}
			case "interval":
				if announcementIntervalSentCount >= cfg.RepeatIntervalCount {
					break
				}
				if announcementLastSentTime == nil {
					baseline := now
					announcementLastSentTime = &baseline
				} else {
					elapsed := now.Sub(*announcementLastSentTime)
					intervalDuration := time.Duration(cfg.RepeatIntervalMinutes) * time.Minute
					if elapsed >= intervalDuration {
						shouldSend = true
						sendMessage = cfg.Message
						intervalSend = true
					}
				}
			case "time":
				currentTime := now.Format("15:04")
				day := now.Format("2006-01-02")
				if currentTime == cfg.RepeatTime {
					key := day + "-" + cfg.RepeatTime
					if announcementLastIntervalKey != key {
						shouldSend = true
						sendMessage = cfg.Message
						timeModeKey = key
					}
				}
			}
			runtimeMu.Unlock()

			if !shouldSend {
				continue
			}

			if err := client.SendServerMessage(sendMessage); err != nil {
				continue
			}

			runtimeMu.Lock()
			sentAt := now
			announcementLastSentTime = &sentAt
			if intervalSend {
				announcementIntervalSentCount++
			}
			if timeModeKey != "" {
				announcementLastIntervalKey = timeModeKey
			}
			if cfg := normalizeAnnouncementConfig(config.Announcement); cfg.ScheduleMode == "once" {
				announcementOncePending = false
			}
			runtimeMu.Unlock()
		}
	}()

	// TODO: Event/command loop, API-Anbindung
}

func pluginStatesFilePath() string {
	return filepath.Join("runtime", "bot", "plugin_states.json")
}

func loadPluginStates() map[string]bool {
	defaults := map[string]bool{
		"AdminCounter":  true,
		"MemberCounter": true,
		"CombinedStats": true,
		"AfkMover":      true,
	}
	data, err := os.ReadFile(pluginStatesFilePath())
	if err != nil {
		return defaults
	}
	loaded := map[string]bool{}
	if err := json.Unmarshal(data, &loaded); err != nil {
		return defaults
	}
	for k, v := range loaded {
		defaults[k] = v
	}
	return defaults
}

func savePluginStates(states map[string]bool) error {
	data, err := json.MarshalIndent(states, "", "  ")
	if err != nil {
		return fmt.Errorf("serialize plugin states: %w", err)
	}
	path := pluginStatesFilePath()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("ensure plugin states dir: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

func setupProcessLogging() {
	logFilePath := strings.TrimSpace(os.Getenv("UC_FRAMEWORK_LOG_FILE"))
	if logFilePath == "" {
		return
	}
	if err := os.MkdirAll(filepath.Dir(logFilePath), 0755); err != nil {
		log.Printf("[ERROR] [logging] ensure log dir failed: %v", err)
		return
	}
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("[ERROR] [logging] open log file failed: %v", err)
		return
	}
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
}

func loadAdminCounterConfig() admincounter.AdminCounterConfig {
	liveConfigsMu.RLock()
	raw, ok := getLivePluginConfigRaw("admincounter")
	liveConfigsMu.RUnlock()
	if ok {
		var cfg admincounter.AdminCounterConfig
		if err := json.Unmarshal(raw, &cfg); err == nil {
			return cfg
		}
	}
	return admincounter.AdminCounterConfig{AdminGroups: config.TS3.AdminGroups}
}

func loadMemberCounterConfig() membercounter.MemberCounterConfig {
	liveConfigsMu.RLock()
	raw, ok := getLivePluginConfigRaw("membercounter")
	liveConfigsMu.RUnlock()
	if ok {
		var cfg membercounter.MemberCounterConfig
		if err := json.Unmarshal(raw, &cfg); err == nil {
			return cfg
		}
	}
	excluded := append(append([]int(nil), config.TS3.AdminGroups...), config.TS3.BotGroups...)
	return membercounter.MemberCounterConfig{ExcludedGroups: excluded, ExcludedNicknames: []string{}}
}

func loadAfkMoverConfig() afkmover.AfkMoverConfig {
	liveConfigsMu.RLock()
	raw, ok := getLivePluginConfigRaw("afkmover")
	liveConfigsMu.RUnlock()
	if ok {
		var cfg afkmover.AfkMoverConfig
		if err := json.Unmarshal(raw, &cfg); err == nil {
			return cfg
		}
	}
	return afkmover.AfkMoverConfig{TimeoutSeconds: 600, Enabled: true, ExcludedChannels: []int{}}
}

func savePluginConfigs() error {
	liveConfigsMu.RLock()
	config.PluginConfigs = make(map[string]json.RawMessage, len(liveConfigs))
	for k, v := range liveConfigs {
		config.PluginConfigs[canonicalPluginConfigKey(k)] = v
	}
	liveConfigsMu.RUnlock()

	return saveConfig()
}

func saveConfig() error {
	persisted := config
	if ts3PasswordFromEnv {
		persisted.TS3.Password = ""
	}
	if discordTokenFromEnv {
		persisted.Discord.BotToken = ""
	}
	if webAuthPasswordFromEnv {
		persisted.WebAuth.Password = ""
	}
	if webAuthPasswordHashFromEnv {
		persisted.WebAuth.PasswordHash = ""
		persisted.WebAuth.Password = ""
	}

	data, err := json.MarshalIndent(persisted, "", "  ")
	if err != nil {
		return fmt.Errorf("serialize config: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(configFilePath), 0755); err != nil {
		return fmt.Errorf("ensure config dir: %w", err)
	}
	return os.WriteFile(configFilePath, data, 0644)
}

func loadConfig() {
	configFilePath = resolveConfigFilePath()
	bytes, err := os.ReadFile(configFilePath)
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}
	if err := json.Unmarshal(bytes, &config); err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}
	applySecretEnvOverrides(&config)
	if config.PluginConfigs == nil {
		config.PluginConfigs = make(map[string]json.RawMessage)
	}
	config.PluginConfigs = normalizePluginConfigMap(config.PluginConfigs)
	config.DefaultLanguage = normalizeLanguageCode(config.DefaultLanguage)
	config.SupportedLanguages = normalizeSupportedLanguages(config.SupportedLanguages)
	if !isSupportedLanguage(config.DefaultLanguage, config.SupportedLanguages) {
		config.DefaultLanguage = "en"
	}
	if config.WebAuth.Username == "" {
		config.WebAuth.Username = "admin"
	}
	if config.WebAuth.Provider == "" {
		config.WebAuth.Provider = "local"
	}
	if config.WebAuth.SessionTTLMinutes <= 0 {
		config.WebAuth.SessionTTLMinutes = 1440
	}
	if config.WebAuth.Password == "" {
		if config.WebAuth.PasswordHash == "" {
			config.WebAuth.Password = "change-me"
		}
	}

	if config.WebAuth.Provider == "local" {
		var raw map[string]json.RawMessage
		var rawWebAuth map[string]json.RawMessage
		if err := json.Unmarshal(bytes, &raw); err == nil {
			if webAuthRaw, ok := raw["web_auth"]; ok {
				_ = json.Unmarshal(webAuthRaw, &rawWebAuth)
			}
		}
		_, forceFieldPresent := rawWebAuth["force_password_change"]
		if !forceFieldPresent {
			if config.WebAuth.PasswordHash == "" && config.WebAuth.Password == "change-me" {
				config.WebAuth.ForcePasswordChange = true
			}
		}
	}

	config.Framework = normalizeFrameworkConfig(config.Framework)
	config.Discord = normalizeDiscordConfig(config.Discord)

	if config.BotControl.BotExecutable == "" {
		if runtime.GOOS == "windows" {
			config.BotControl.BotExecutable = "uc-framework_bot.exe"
		} else {
			config.BotControl.BotExecutable = "uc-framework_bot"
		}
	}
	if config.BotControl.BotExecutable == "uz_bot_bot" || config.BotControl.BotExecutable == "uz_bot_bot.exe" {
		if runtime.GOOS == "windows" {
			config.BotControl.BotExecutable = "uc-framework_bot.exe"
		} else {
			config.BotControl.BotExecutable = "uc-framework_bot"
		}
	}
	if config.BotControl.WorkingDir == "" {
		config.BotControl.WorkingDir = "."
	}
	if config.BotControl.WatchdogMinIntervalSec < 60 {
		config.BotControl.WatchdogMinIntervalSec = 60
	}
	if config.BotControl.WatchdogMaxIntervalSec < config.BotControl.WatchdogMinIntervalSec {
		config.BotControl.WatchdogMaxIntervalSec = 120
	}

	if config.TS3.QueryPort == 0 {
		if config.TS3.Port != 0 {
			config.TS3.QueryPort = config.TS3.Port
		} else {
			config.TS3.QueryPort = 10011
		}
	}
	if config.TS3.VoicePort == 0 {
		config.TS3.VoicePort = 9987
	}
	if config.TS3.QuerySlowmode <= 0 {
		config.TS3.QuerySlowmode = 250
	}
	if config.TS3.BotNickname == "" {
		config.TS3.BotNickname = "UC-Framework"
	}

	config.Support = normalizeSupportConfig(config.Support)
	if !config.WebAuth.Enabled {
		// Standard: Auth ist aktiv, auÃŸer sie wird explizit auf false gesetzt.
		var raw map[string]json.RawMessage
		if err := json.Unmarshal(bytes, &raw); err == nil {
			if _, ok := raw["web_auth"]; !ok {
				config.WebAuth.Enabled = true
			}
		}
	}

	enforceSecretPolicy(&config)
}

func loadLanguage(lang string) {
	langPath := filepath.Join("locales", lang+".json")
	bytes, err := os.ReadFile(langPath)
	if err != nil {
		log.Fatalf("Failed to read language file: %v", err)
	}
	if err := json.Unmarshal(bytes, &language); err != nil {
		log.Fatalf("Failed to parse language file: %v", err)
	}
}

func monitorTS3Connection(client *ts3.TS3Client, stopCh <-chan struct{}, onMaxFailures func(error, int)) {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	disconnectLogged := false
	failedReconnects := 0
	for {
		select {
		case <-stopCh:
			return
		case <-ticker.C:
		}

		if client.IsConnected() {
			if disconnectLogged {
				disconnectLogged = false
			}
			failedReconnects = 0
			continue
		}

		if !disconnectLogged {
			log.Printf("[ERROR] [ts3] connection lost: host=%s query_port=%d", client.Host, client.Port)
			disconnectLogged = true
		}

		if err := client.Connect(); err != nil {
			failedReconnects++
			log.Printf("[ERROR] [ts3] reconnect attempt failed: %v", err)
			if failedReconnects >= 5 {
				if onMaxFailures != nil {
					onMaxFailures(err, failedReconnects)
				}
				return
			}
			continue
		}

		failedReconnects = 0
	}
}
