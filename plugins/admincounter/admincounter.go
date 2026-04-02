package admincounter

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"sync"
	"uc_framework/internal/bot"
	"uc_framework/internal/ts3"
)

type AdminCounterConfig struct {
	AdminGroups        []int  `json:"admin_groups"`
	RenameChannelID    int    `json:"rename_channel_id,omitempty"`
	RenameNameTemplate string `json:"rename_name_template,omitempty"`
	RenameCountToken   string `json:"rename_count_token,omitempty"`
}

type AdminCounterPlugin struct {
	mu               sync.RWMutex
	adminGroups      []int
	renameChannelID  int
	renameTemplate   string
	renameCountToken string
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

func (p *AdminCounterPlugin) Start(d *bot.Dispatcher) error {
	d.RegisterEventHandler(p.Name(), "user_update", p.onUserUpdate)
	d.RegisterCommand(p.Name(), bot.Command{
		Name:        "admins",
		Description: "Show online admin count",
		Execute: func(args []string, ctx bot.CommandContext) {
			log.Printf("[INFO] [AdminCounter] admins command executed")
		},
	})
	return nil
}

func (p *AdminCounterPlugin) Stop(d *bot.Dispatcher) error { return nil }
func (p *AdminCounterPlugin) Name() string                 { return "AdminCounter" }
func (p *AdminCounterPlugin) Description() string {
	return "ZÃ¤hlt Online-Admins anhand konfigurierbarer Servergruppen."
}

func (p *AdminCounterPlugin) GetConfig() ([]byte, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return json.Marshal(AdminCounterConfig{
		AdminGroups:        p.adminGroups,
		RenameChannelID:    p.renameChannelID,
		RenameNameTemplate: p.renameTemplate,
		RenameCountToken:   p.renameCountToken,
	})
}

func (p *AdminCounterPlugin) SetConfig(raw []byte) error {
	var cfg AdminCounterConfig
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return err
	}
	if cfg.AdminGroups == nil {
		cfg.AdminGroups = []int{}
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.adminGroups = cfg.AdminGroups
	p.renameChannelID = cfg.RenameChannelID
	p.renameTemplate = strings.TrimSpace(cfg.RenameNameTemplate)
	p.renameCountToken = resolveRenameCountToken(strings.TrimSpace(cfg.RenameCountToken), p.renameTemplate)
	p.lastRenamedCount = -1
	return nil
}

func (p *AdminCounterPlugin) onUserUpdate(event bot.Event) {
	ts3Client, ok := event.Payload["ts3_client"].(*ts3.TS3Client)
	if !ok || ts3Client == nil {
		log.Printf("[ERROR] [AdminCounter] missing TS3 client in event payload")
		return
	}
	clients, err := ts3Client.ListClients()
	if err != nil {
		log.Printf("[ERROR] [AdminCounter] listing clients failed: %v", err)
		return
	}
	p.mu.RLock()
	adminGroups := append([]int(nil), p.adminGroups...)
	renameChannelID := p.renameChannelID
	renameTemplate := p.renameTemplate
	renameCountToken := p.renameCountToken
	lastRenamedCount := p.lastRenamedCount
	p.mu.RUnlock()

	adminCount := 0
	for _, c := range clients {
		if c.IsQuery {
			continue
		}
		for _, group := range c.ServerGroups {
			if containsInt(adminGroups, group) {
				adminCount++
				break
			}
		}
	}
	if event.Payload != nil {
		event.Payload["admins_online"] = adminCount
	}

	if renameChannelID > 0 && strings.Contains(renameTemplate, renameCountToken) && adminCount != lastRenamedCount {
		newName := strings.ReplaceAll(renameTemplate, renameCountToken, strconv.Itoa(adminCount))
		currentName, found, err := channelNameByID(ts3Client, renameChannelID)
		if err != nil {
			log.Printf("[ERROR] [AdminCounter] channel lookup failed: %v", err)
		} else if found && currentName == newName {
			p.mu.Lock()
			p.lastRenamedCount = adminCount
			p.mu.Unlock()
			return
		}
		if err := ts3Client.RenameChannel(renameChannelID, newName); err != nil {
			log.Printf("[ERROR] [AdminCounter] channel rename failed: %v", err)
		} else {
			p.mu.Lock()
			p.lastRenamedCount = adminCount
			p.mu.Unlock()
		}
	}

	// TODO: Channel-Name/Description via TS3-Query setzen
}

func New(cfg AdminCounterConfig) *AdminCounterPlugin {
	groups := cfg.AdminGroups
	if groups == nil {
		groups = []int{}
	}
	return &AdminCounterPlugin{
		adminGroups:      groups,
		renameChannelID:  cfg.RenameChannelID,
		renameTemplate:   strings.TrimSpace(cfg.RenameNameTemplate),
		renameCountToken: resolveRenameCountToken(strings.TrimSpace(cfg.RenameCountToken), strings.TrimSpace(cfg.RenameNameTemplate)),
		lastRenamedCount: -1,
	}
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
