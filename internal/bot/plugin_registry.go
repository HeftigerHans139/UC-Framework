package bot

import (
	"fmt"
	"sync"
)

type PluginFactory func() Plugin

type PluginState struct {
	Name        string
	Description string
	Loaded      bool
	factory     PluginFactory
	plugin      Plugin
}

type PluginRegistry struct {
	plugins    map[string]*PluginState
	dispatcher *Dispatcher
	lock       sync.RWMutex
}

func NewPluginRegistry(dispatcher *Dispatcher) *PluginRegistry {
	return &PluginRegistry{
		plugins:    make(map[string]*PluginState),
		dispatcher: dispatcher,
	}
}

func (r *PluginRegistry) RegisterFactory(name string, factory PluginFactory) {
	plugin := factory()

	r.lock.Lock()
	defer r.lock.Unlock()
	r.plugins[name] = &PluginState{
		Name:        name,
		Description: plugin.Description(),
		Loaded:      false,
		factory:     factory,
	}
}

func (r *PluginRegistry) Load(name string) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	state, ok := r.plugins[name]
	if !ok {
		return fmt.Errorf("plugin %s not registered", name)
	}
	if state.Loaded {
		return nil
	}

	plugin := state.factory()
	if err := plugin.Start(r.dispatcher); err != nil {
		return fmt.Errorf("start plugin %s: %w", name, err)
	}

	state.plugin = plugin
	state.Loaded = true
	return nil
}

func (r *PluginRegistry) Unload(name string) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	state, ok := r.plugins[name]
	if !ok {
		return fmt.Errorf("plugin %s not registered", name)
	}
	if !state.Loaded || state.plugin == nil {
		return nil
	}

	if err := state.plugin.Stop(r.dispatcher); err != nil {
		return fmt.Errorf("stop plugin %s: %w", name, err)
	}

	r.dispatcher.UnregisterPlugin(name)
	state.plugin = nil
	state.Loaded = false
	return nil
}

func (r *PluginRegistry) IsLoaded(name string) bool {
	r.lock.RLock()
	defer r.lock.RUnlock()
	if p, ok := r.plugins[name]; ok {
		return p.Loaded
	}
	return false
}

func (r *PluginRegistry) Get(name string) (Plugin, bool) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	p, ok := r.plugins[name]
	if !ok || !p.Loaded || p.plugin == nil {
		return nil, false
	}
	return p.plugin, true
}

func (r *PluginRegistry) All() map[string]PluginState {
	r.lock.RLock()
	defer r.lock.RUnlock()
	copy := make(map[string]PluginState, len(r.plugins))
	for k, v := range r.plugins {
		copy[k] = PluginState{
			Name:        v.Name,
			Description: v.Description,
			Loaded:      v.Loaded,
		}
	}
	return copy
}

// GetPluginConfig gibt die aktuelle JSON-Konfiguration eines geladenen Plugins zurück
func (r *PluginRegistry) GetPluginConfig(name string) ([]byte, error) {
	r.lock.RLock()
	state, ok := r.plugins[name]
	r.lock.RUnlock()
	if !ok {
		return nil, fmt.Errorf("plugin %s not registered", name)
	}
	if !state.Loaded || state.plugin == nil {
		return nil, fmt.Errorf("plugin %s not loaded", name)
	}
	c, ok := state.plugin.(Configurable)
	if !ok {
		return nil, fmt.Errorf("plugin %s is not configurable", name)
	}
	return c.GetConfig()
}

// UpdatePluginConfig aktualisiert die In-Memory-Konfiguration eines geladenen Plugins.
// Gibt (false, nil) zurück wenn das Plugin nicht geladen ist (kein Fehler).
func (r *PluginRegistry) UpdatePluginConfig(name string, raw []byte) (bool, error) {
	r.lock.RLock()
	state, ok := r.plugins[name]
	r.lock.RUnlock()
	if !ok {
		return false, fmt.Errorf("plugin %s not registered", name)
	}
	if !state.Loaded || state.plugin == nil {
		return false, nil
	}
	c, ok := state.plugin.(Configurable)
	if !ok {
		return false, fmt.Errorf("plugin %s is not configurable", name)
	}
	if err := c.SetConfig(raw); err != nil {
		return false, err
	}
	return true, nil
}
