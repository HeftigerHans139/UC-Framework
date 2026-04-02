package membercounter

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"sync"
	"uc_framework/internal/bot"
	"uc_framework/internal/ts3"
)

type MemberCounterConfig struct {
	ExcludedGroups     []int    `json:"excluded_groups"`
	ExcludedNicknames  []string `json:"excluded_nicknames"`
	RenameChannelID    int      `json:"rename_channel_id,omitempty"`
	RenameNameTemplate string   `json:"rename_name_template,omitempty"`
	RenameCountToken   string   `json:"rename_count_token,omitempty"`
}

type MemberCounterPlugin struct {
	mu               sync.RWMutex
	config           MemberCounterConfig
	lastRenamedCount int
}

func containsInt(list []int, v int) bool {
	for _, i := range list {
		if i == v {
			return true
		}
	}
	return false
}

func containsString(list []string, v string) bool {
	for _, s := range list {
		if s == v {
			return true
		}
	}
	return false
}

func channelNameByID(ts3Client *ts3.TS3Client, channelID int) (string, bool, error) {
	channels, err := ts3Client.ListChannels()
	if err != nil {
		return "", false, err
	}
	for _, channel := range channels {
		if channel.ID == channelID {
			return channel.Name, true, nil
		}
	}
	return "", false, nil
}

func (p *MemberCounterPlugin) Start(d *bot.Dispatcher) error {
	d.RegisterEventHandler(p.Name(), "user_update", p.onUserUpdate)
	return nil
}

func (p *MemberCounterPlugin) Stop(d *bot.Dispatcher) error { return nil }
func (p *MemberCounterPlugin) Name() string                 { return "MemberCounter" }
func (p *MemberCounterPlugin) Description() string {
	return "ZÃ¤hlt alle Online-User auÃŸer konfigurierten Ausnahmen (Gruppen/Nicknames)."
}

func (p *MemberCounterPlugin) GetConfig() ([]byte, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return json.Marshal(p.config)
}

func (p *MemberCounterPlugin) SetConfig(raw []byte) error {
	var cfg MemberCounterConfig
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return err
	}
	if cfg.ExcludedGroups == nil {
		cfg.ExcludedGroups = []int{}
	}
	if cfg.ExcludedNicknames == nil {
		cfg.ExcludedNicknames = []string{}
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	cfg.RenameNameTemplate = strings.TrimSpace(cfg.RenameNameTemplate)
	p.config = cfg
	p.lastRenamedCount = -1
	return nil
}

func (p *MemberCounterPlugin) onUserUpdate(event bot.Event) {
	ts3Client, ok := event.Payload["ts3_client"].(*ts3.TS3Client)
	if !ok || ts3Client == nil {
		log.Printf("[ERROR] [MemberCounter] missing TS3 client in event payload")
		return
	}
	clients, err := ts3Client.ListClients()
	if err != nil {
		log.Printf("[ERROR] [MemberCounter] listing clients failed: %v", err)
		return
	}
	p.mu.RLock()
	cfg := p.config
	lastRenamedCount := p.lastRenamedCount
	p.mu.RUnlock()

	memberCount := 0
	for _, c := range clients {
		if c.IsQuery {
			continue
		}
		excluded := false
		for _, group := range c.ServerGroups {
			if containsInt(cfg.ExcludedGroups, group) {
				excluded = true
				break
			}
		}
		if !excluded && containsString(cfg.ExcludedNicknames, c.Nickname) {
			excluded = true
		}
		if !excluded {
			memberCount++
		}
	}
	if event.Payload != nil {
		event.Payload["members_online"] = memberCount
	}

	renameCountToken := resolveRenameCountToken(strings.TrimSpace(cfg.RenameCountToken), cfg.RenameNameTemplate)
	if cfg.RenameChannelID > 0 && strings.Contains(cfg.RenameNameTemplate, renameCountToken) && memberCount != lastRenamedCount {
		newName := strings.ReplaceAll(cfg.RenameNameTemplate, renameCountToken, strconv.Itoa(memberCount))
		currentName, found, err := channelNameByID(ts3Client, cfg.RenameChannelID)
		if err != nil {
			log.Printf("[ERROR] [MemberCounter] channel lookup failed: %v", err)
		} else if found && currentName == newName {
			p.mu.Lock()
			p.lastRenamedCount = memberCount
			p.mu.Unlock()
			return
		}
		if err := ts3Client.RenameChannel(cfg.RenameChannelID, newName); err != nil {
			log.Printf("[ERROR] [MemberCounter] channel rename failed: %v", err)
		} else {
			p.mu.Lock()
			p.lastRenamedCount = memberCount
			p.mu.Unlock()
		}
	}

	// TODO: Channel-Name/Description via TS3-Query setzen
}

func New(cfg MemberCounterConfig) *MemberCounterPlugin {
	if cfg.ExcludedGroups == nil {
		cfg.ExcludedGroups = []int{}
	}
	if cfg.ExcludedNicknames == nil {
		cfg.ExcludedNicknames = []string{}
	}
	cfg.RenameNameTemplate = strings.TrimSpace(cfg.RenameNameTemplate)
	cfg.RenameCountToken = resolveRenameCountToken(strings.TrimSpace(cfg.RenameCountToken), cfg.RenameNameTemplate)
	return &MemberCounterPlugin{config: cfg, lastRenamedCount: -1}
}

func resolveRenameCountToken(token, template string) string {
	t := strings.TrimSpace(token)
	if t != "" {
		return t
	}
	if strings.Contains(template, "{count}") {
		return "{count}"
	}
	if strings.Contains(template, "%%") {
		return "%%"
	}
	return "{count}"
}
