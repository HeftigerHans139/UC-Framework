package ts3

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
	"uc_framework/internal/bot"
)

// TS3Client represents a basic TeamSpeak 3 Query client
// (spÃ¤ter kann hier ein externes Paket wie github.com/multiplay/go-ts3 eingebunden werden)
type TS3Client struct {
	Host      string
	Port      int
	VoicePort int
	Username  string
	Password  string

	conn      net.Conn
	reader    *bufio.Reader
	writer    *bufio.Writer
	lock      sync.RWMutex
	queryMu   sync.Mutex
	connected bool

	dispatcher *bot.Dispatcher
}

// NewTS3Client creates a new TS3Client from config
var defaultReconnectDelay = 5 * time.Second

func NewTS3Client(host string, port int, voicePort int, user, pass string) *TS3Client {
	return &TS3Client{
		Host:      host,
		Port:      port,
		VoicePort: voicePort,
		Username:  user,
		Password:  pass,
	}
}

// SetDispatcher sets the event dispatcher for event forwarding
func (c *TS3Client) SetDispatcher(d *bot.Dispatcher) {
	c.dispatcher = d
}

// Connect establishes a real TeamSpeak ServerQuery TCP connection.
func (c *TS3Client) Connect() error {
	addr := net.JoinHostPort(c.Host, intToString(c.Port))

	conn, err := net.DialTimeout("tcp", addr, 8*time.Second)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	c.lock.Lock()
	c.conn = conn
	c.reader = reader
	c.writer = writer
	c.connected = true
	c.lock.Unlock()

	// TS3 server sends a welcome banner immediately.
	_ = c.consumeWelcomeBanner()

	if _, err := c.query("login client_login_name=" + escapeTS3(c.Username) + " client_login_password=" + escapeTS3(c.Password)); err != nil {
		c.MarkDisconnected(err)
		return err
	}

	if c.VoicePort > 0 {
		if _, err := c.query("use port=" + intToString(c.VoicePort)); err != nil {
			c.MarkDisconnected(err)
			return err
		}
	}

	return nil
}

// SetBotNickname updates the TS3 display nickname of the connected query client.
func (c *TS3Client) SetBotNickname(nickname string) error {
	if strings.TrimSpace(nickname) == "" {
		return fmt.Errorf("bot nickname must not be empty")
	}
	_, err := c.query("clientupdate client_nickname=" + escapeTS3(nickname))
	return err
}

// SendServerMessage sends a server-wide text message via ServerQuery.
func (c *TS3Client) SendServerMessage(message string) error {
	trimmed := strings.TrimSpace(message)
	if trimmed == "" {
		return fmt.Errorf("message must not be empty")
	}
	_, err := c.query("sendtextmessage targetmode=3 target=1 msg=" + escapeTS3(trimmed))
	return err
}

// IsConnected returns connection status
func (c *TS3Client) IsConnected() bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.connected
}

// MarkDisconnected marks the current TS3 connection as lost and logs the reason.
func (c *TS3Client) MarkDisconnected(err error) {
	c.lock.Lock()
	c.connected = false
	if c.conn != nil {
		_ = c.conn.Close()
	}
	c.conn = nil
	c.reader = nil
	c.writer = nil
	c.lock.Unlock()

	if err != nil {
		log.Printf("[ERROR] [ts3] connection marked disconnected: %v", err)
		return
	}
}

// Reconnect tries to reconnect with delay
func (c *TS3Client) Reconnect() {
	for !c.IsConnected() {
		err := c.Connect()
		if err != nil {
			log.Printf("[ERROR] [ts3] reconnect failed: %v. retrying in %v", err, defaultReconnectDelay)
			time.Sleep(defaultReconnectDelay)
		}
	}
}

// ListenEvents polls TS3 periodically and dispatches update events.
func (c *TS3Client) ListenEvents() {
	for {
		if !c.IsConnected() {
			time.Sleep(3 * time.Second)
			continue
		}

		if _, err := c.ListClients(); err != nil {
			c.MarkDisconnected(err)
			time.Sleep(3 * time.Second)
			continue
		}

		if c.dispatcher != nil {
			event := bot.Event{
				Type: "user_update",
				Payload: map[string]interface{}{
					"ts3_client": c,
				},
			}
			c.dispatcher.DispatchEvent(event)
		}
		time.Sleep(10 * time.Second)
	}
}

func (c *TS3Client) consumeWelcomeBanner() error {
	c.lock.RLock()
	conn := c.conn
	reader := c.reader
	c.lock.RUnlock()
	if conn == nil || reader == nil {
		return nil
	}

	_ = conn.SetReadDeadline(time.Now().Add(250 * time.Millisecond))
	defer conn.SetReadDeadline(time.Time{})

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Timeout() {
				return nil
			}
			if strings.Contains(strings.ToLower(err.Error()), "timeout") {
				return nil
			}
			return err
		}
		_ = strings.TrimSpace(line)
	}
}
