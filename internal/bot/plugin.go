package bot

// Plugin interface for modular extensions
// Plugins can register events and commands

type Plugin interface {
	Name() string
	Description() string
	Start(d *Dispatcher) error
	Stop(d *Dispatcher) error
}

// Configurable ist ein optionales Interface für Plugins, die Live-Konfiguration unterstützen
type Configurable interface {
	GetConfig() ([]byte, error)
	SetConfig(raw []byte) error
}
