package bot

import "sync"

type EventType string

type Event struct {
	Type    EventType
	Payload map[string]interface{}
}

type EventHandler func(Event)

type Command struct {
	Name        string
	Description string
	Execute     func(args []string, ctx CommandContext)
}

type CommandContext struct {
	UserID   string
	Channel  string
	Language string
	// ... weitere Felder
}

type registeredEventHandler struct {
	pluginName string
	handler    EventHandler
}

type registeredCommand struct {
	pluginName string
	command    Command
}

// Dispatcher verwaltet Events und Commands
type Dispatcher struct {
	eventHandlers map[EventType][]registeredEventHandler
	commands      map[string]registeredCommand
	lock          sync.RWMutex
}

// HandleCommandInput parses and executes a command string (e.g. from chat)
func (d *Dispatcher) HandleCommandInput(input string, ctx CommandContext) {
	// Simple: assume "/command arg1 arg2 ..."
	if len(input) == 0 || input[0] != '/' {
		return
	}
	parts := splitArgs(input[1:])
	if len(parts) == 0 {
		return
	}
	name := parts[0]
	args := []string{}
	if len(parts) > 1 {
		args = parts[1:]
	}
	d.ExecuteCommand(name, args, ctx)
}

// splitArgs splits a string by spaces, respecting quoted arguments
func splitArgs(s string) []string {
	var args []string
	var current string
	inQuotes := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '"' {
			inQuotes = !inQuotes
			continue
		}
		if c == ' ' && !inQuotes {
			if current != "" {
				args = append(args, current)
				current = ""
			}
			continue
		}
		current += string(c)
	}
	if current != "" {
		args = append(args, current)
	}
	return args
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		eventHandlers: make(map[EventType][]registeredEventHandler),
		commands:      make(map[string]registeredCommand),
	}
}

func (d *Dispatcher) RegisterEventHandler(pluginName string, eventType EventType, handler EventHandler) {
	d.lock.Lock()
	defer d.lock.Unlock()

	d.eventHandlers[eventType] = append(d.eventHandlers[eventType], registeredEventHandler{
		pluginName: pluginName,
		handler:    handler,
	})
}

func (d *Dispatcher) DispatchEvent(event Event) {
	d.lock.RLock()
	handlers := append([]registeredEventHandler(nil), d.eventHandlers[event.Type]...)
	d.lock.RUnlock()

	for _, handler := range handlers {
		handler.handler(event)
	}
}

func (d *Dispatcher) RegisterCommand(pluginName string, cmd Command) {
	d.lock.Lock()
	defer d.lock.Unlock()

	d.commands[cmd.Name] = registeredCommand{
		pluginName: pluginName,
		command:    cmd,
	}
}

func (d *Dispatcher) ExecuteCommand(name string, args []string, ctx CommandContext) {
	d.lock.RLock()
	cmd, ok := d.commands[name]
	d.lock.RUnlock()
	if ok {
		cmd.command.Execute(args, ctx)
	}
}

func (d *Dispatcher) UnregisterPlugin(pluginName string) {
	d.lock.Lock()
	defer d.lock.Unlock()

	for eventType, handlers := range d.eventHandlers {
		filtered := handlers[:0]
		for _, handler := range handlers {
			if handler.pluginName != pluginName {
				filtered = append(filtered, handler)
			}
		}
		if len(filtered) == 0 {
			delete(d.eventHandlers, eventType)
			continue
		}
		d.eventHandlers[eventType] = filtered
	}

	for name, cmd := range d.commands {
		if cmd.pluginName == pluginName {
			delete(d.commands, name)
		}
	}
}
