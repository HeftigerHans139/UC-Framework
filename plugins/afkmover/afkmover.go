package afkmover

import (
	"encoding/json"
	"log"
	"sync"
	"time"
	"uc_framework/internal/bot"
	"uc_framework/internal/ts3"
)

type AfkMoverConfig struct {
	AfkChannelID     int    `json:"afk_channel_id"`
	AfkChannelName   string `json:"afk_channel_name,omitempty"`
	TimeoutSeconds   int64  `json:"timeout_seconds"`
	ReturnOnActivity bool   `json:"return_on_activity"`
	ExcludedChannels []int  `json:"excluded_channels"`
	Enabled          bool   `json:"enabled"`
}

type AfkMoverPlugin struct {
	mu           sync.RWMutex
	config       AfkMoverConfig
	movedFromMap map[int]int
}

func (p *AfkMoverPlugin) Start(d *bot.Dispatcher) error {
	d.RegisterEventHandler(p.Name(), "user_update", p.onUserUpdate)
	return nil
}

func (p *AfkMoverPlugin) Stop(d *bot.Dispatcher) error { return nil }
func (p *AfkMoverPlugin) Name() string                 { return "AfkMover" }
func (p *AfkMoverPlugin) Description() string {
	return "Verschiebt inaktive User nach X Minuten in den AFK-Channel."
}

func (p *AfkMoverPlugin) GetConfig() ([]byte, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return json.Marshal(p.config)
}

func (p *AfkMoverPlugin) SetConfig(raw []byte) error {
	var cfg AfkMoverConfig
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return err
	}
	if cfg.ExcludedChannels == nil {
		cfg.ExcludedChannels = []int{}
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.config = cfg
	return nil
}

func (p *AfkMoverPlugin) onUserUpdate(event bot.Event) {
	p.mu.RLock()
	cfg := p.config
	p.mu.RUnlock()

	if !cfg.Enabled {
		return
	}
	ts3Client, ok := event.Payload["ts3_client"].(*ts3.TS3Client)
	if !ok || ts3Client == nil {
		log.Printf("[ERROR] [AfkMover] missing TS3 client in event payload")
		return
	}
	clients, err := ts3Client.ListClients()
	if err != nil {
		log.Printf("[ERROR] [AfkMover] listing clients failed: %v", err)
		return
	}
	now := time.Now().Unix()
	for _, c := range clients {
		if c.IsQuery || c.ID <= 0 {
			continue
		}

		if c.ChannelID == cfg.AfkChannelID {
			if !cfg.ReturnOnActivity || cfg.TimeoutSeconds <= 0 {
				continue
			}
			if now-c.LastActive > cfg.TimeoutSeconds {
				continue
			}

			fromChannel, tracked := p.getTrackedFromChannel(c.ID)
			if !tracked || fromChannel <= 0 || fromChannel == cfg.AfkChannelID {
				continue
			}

			if err := ts3Client.MoveClient(c.ID, fromChannel); err != nil {
				log.Printf("[ERROR] [AfkMover] return move failed for %s: %v", c.Nickname, err)
				continue
			}
			p.clearTrackedClient(c.ID)
			continue
		}

		p.clearTrackedClient(c.ID)

		excluded := false
		for _, excCh := range cfg.ExcludedChannels {
			if c.ChannelID == excCh {
				excluded = true
				break
			}
		}
		if excluded {
			continue
		}
		if cfg.TimeoutSeconds > 0 && now-c.LastActive > cfg.TimeoutSeconds {
			fromChannel := c.ChannelID
			if err := ts3Client.MoveClient(c.ID, cfg.AfkChannelID); err != nil {
				log.Printf("[ERROR] [AfkMover] move failed for %s: %v", c.Nickname, err)
				continue
			}
			p.trackMovedClient(c.ID, fromChannel)
		}
	}
}

func New(config AfkMoverConfig) *AfkMoverPlugin {
	if config.ExcludedChannels == nil {
		config.ExcludedChannels = []int{}
	}
	return &AfkMoverPlugin{config: config, movedFromMap: make(map[int]int)}
}

func (p *AfkMoverPlugin) trackMovedClient(clientID int, fromChannel int) {
	p.mu.Lock()
	p.movedFromMap[clientID] = fromChannel
	p.mu.Unlock()
}

func (p *AfkMoverPlugin) clearTrackedClient(clientID int) {
	p.mu.Lock()
	delete(p.movedFromMap, clientID)
	p.mu.Unlock()
}

func (p *AfkMoverPlugin) getTrackedFromChannel(clientID int) (int, bool) {
	p.mu.RLock()
	fromChannel, ok := p.movedFromMap[clientID]
	p.mu.RUnlock()
	return fromChannel, ok
}
