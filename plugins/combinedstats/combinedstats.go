package combinedstats

import (
	"uc_framework/internal/bot"
	"uc_framework/web/api"
)

type CombinedStatsPlugin struct {
	admins  int
	members int
}

func (p *CombinedStatsPlugin) Start(d *bot.Dispatcher) error {
	d.RegisterEventHandler(p.Name(), "user_update", p.onUserUpdate)
	// API-Getter registrieren
	api.GetStatsFunc = func() (int, int) {
		return p.admins, p.members
	}
	return nil
}

func (p *CombinedStatsPlugin) Name() string {
	return "CombinedStats"
}

func (p *CombinedStatsPlugin) Description() string {
	return "Combines and displays admin/member stats."
}

func (p *CombinedStatsPlugin) Stop(d *bot.Dispatcher) error {
	api.GetStatsFunc = nil
	return nil
}

func (p *CombinedStatsPlugin) onUserUpdate(event bot.Event) {
	admins, okA := event.Payload["admins_online"].(int)
	members, okM := event.Payload["members_online"].(int)
	if okA {
		p.admins = admins
	}
	if okM {
		p.members = members
	}
	// TODO: Channel-Name/Description aktualisieren oder Web-API bereitstellen
}

func New() *CombinedStatsPlugin {
	return &CombinedStatsPlugin{}
}
