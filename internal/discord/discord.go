package discord

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Channel struct {
	ID       string
	Name     string
	Type     discordgo.ChannelType
	ParentID string
}

type Role struct {
	ID   string
	Name string
}

type Client struct {
	token   string
	guildID string

	mu        sync.RWMutex
	session   *discordgo.Session
	connected bool
	afkStop   chan struct{}

	afkKickEnabled       bool
	afkInactivityMinutes int
	voiceActivity        map[string]time.Time // userID -> letzter Aktivitätszeitpunkt
	logf                 func(format string, args ...any)
}

func NewClient(token, guildID string) *Client {
	return &Client{
		token:         strings.TrimSpace(token),
		guildID:       strings.TrimSpace(guildID),
		voiceActivity: make(map[string]time.Time),
	}
}

func (c *Client) Connect() error {
	if c.token == "" {
		return fmt.Errorf("discord token is empty")
	}
	if c.guildID == "" {
		return fmt.Errorf("discord guild id is empty")
	}

	session, err := discordgo.New("Bot " + c.token)
	if err != nil {
		return err
	}
	session.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates

	// Verfolge Voice-Aktivität aller User im Guild.
	session.AddHandler(func(s *discordgo.Session, vs *discordgo.VoiceStateUpdate) {
		if vs == nil || strings.TrimSpace(vs.GuildID) != c.guildID {
			return
		}
		// Bot selbst ignorieren.
		if s.State != nil && s.State.User != nil && strings.TrimSpace(vs.UserID) == strings.TrimSpace(s.State.User.ID) {
			return
		}
		userID := vs.UserID
		if strings.TrimSpace(vs.ChannelID) == "" {
			// User hat Voice verlassen → aus Tracking entfernen.
			c.mu.Lock()
			delete(c.voiceActivity, userID)
			c.mu.Unlock()
		} else {
			// User hat einen Voice-Channel betreten oder seinen Status geändert → Aktivität aktualisieren.
			c.mu.Lock()
			c.voiceActivity[userID] = time.Now()
			c.mu.Unlock()
		}
	})

	if err := session.Open(); err != nil {
		return err
	}

	if _, err := session.Guild(c.guildID); err != nil {
		session.Close()
		return err
	}

	stopChan := make(chan struct{})
	c.mu.Lock()
	c.session = session
	c.connected = true
	c.afkStop = stopChan
	c.mu.Unlock()

	go c.runAfkChecker(session, stopChan)
	return nil
}

// runAfkChecker prüft jede Minute ob Voice-User zu lange inaktiv waren und trennt sie.
func (c *Client) runAfkChecker(session *discordgo.Session, stopChan <-chan struct{}) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-stopChan:
			return
		case <-ticker.C:
			c.mu.RLock()
			enabled := c.afkKickEnabled
			inactivityMinutes := c.afkInactivityMinutes
			logf := c.logf
			guildID := c.guildID
			var toKick []string
			if enabled && inactivityMinutes > 0 {
				cutoff := time.Now().Add(-time.Duration(inactivityMinutes) * time.Minute)
				for userID, lastActive := range c.voiceActivity {
					if lastActive.Before(cutoff) {
						toKick = append(toKick, userID)
					}
				}
			}
			c.mu.RUnlock()

			for _, userID := range toKick {
				if err := session.GuildMemberMove(guildID, userID, nil); err != nil {
					if logf != nil {
						logf("[WARN] [discord] AFK inactivity kick failed for user=%s: %v", userID, err)
					}
				} else {
					if logf != nil {
						logf("[INFO] [discord] AFK inactivity kick: user=%s was inactive for %d min", userID, inactivityMinutes)
					}
					c.mu.Lock()
					delete(c.voiceActivity, userID)
					c.mu.Unlock()
				}
			}
		}
	}
}

func (c *Client) Close() error {
	c.mu.Lock()
	session := c.session
	c.session = nil
	c.connected = false
	stopChan := c.afkStop
	c.afkStop = nil
	c.mu.Unlock()

	if stopChan != nil {
		close(stopChan)
	}
	if session == nil {
		return nil
	}
	return session.Close()
}

func (c *Client) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.connected && c.session != nil
}

func (c *Client) SetAFKKickConfig(enabled bool, inactivityMinutes int) {
	c.mu.Lock()
	c.afkKickEnabled = enabled
	if inactivityMinutes <= 0 {
		inactivityMinutes = 30
	}
	c.afkInactivityMinutes = inactivityMinutes
	c.mu.Unlock()
}

func (c *Client) SetLogger(logf func(format string, args ...any)) {
	c.mu.Lock()
	c.logf = logf
	c.mu.Unlock()
}

func (c *Client) ListChannels() ([]Channel, error) {
	c.mu.RLock()
	session := c.session
	guildID := c.guildID
	c.mu.RUnlock()
	if session == nil {
		return nil, fmt.Errorf("discord not connected")
	}

	channels, err := session.GuildChannels(guildID)
	if err != nil {
		return nil, err
	}
	out := make([]Channel, 0, len(channels))
	for _, ch := range channels {
		if ch == nil {
			continue
		}
		out = append(out, Channel{
			ID:       ch.ID,
			Name:     ch.Name,
			Type:     ch.Type,
			ParentID: ch.ParentID,
		})
	}
	return out, nil
}

func (c *Client) ListRoles() ([]Role, error) {
	c.mu.RLock()
	session := c.session
	guildID := c.guildID
	c.mu.RUnlock()
	if session == nil {
		return nil, fmt.Errorf("discord not connected")
	}

	guild, err := session.Guild(guildID)
	if err != nil {
		return nil, err
	}
	out := make([]Role, 0, len(guild.Roles))
	for _, role := range guild.Roles {
		if role == nil {
			continue
		}
		out = append(out, Role{ID: role.ID, Name: role.Name})
	}
	return out, nil
}

func (c *Client) SendMessage(channelID, message string) error {
	trimmedChannelID := strings.TrimSpace(channelID)
	trimmedMessage := strings.TrimSpace(message)
	if trimmedChannelID == "" {
		return fmt.Errorf("discord channel id is empty")
	}
	if trimmedMessage == "" {
		return fmt.Errorf("discord message is empty")
	}

	c.mu.RLock()
	session := c.session
	c.mu.RUnlock()
	if session == nil {
		return fmt.Errorf("discord not connected")
	}

	_, err := session.ChannelMessageSend(trimmedChannelID, trimmedMessage)
	return err
}
