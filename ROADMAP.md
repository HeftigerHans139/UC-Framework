# Roadmap: TS3 Bot Framework in Go

## 1. Project Setup & Core Structure
- [x] Create project structure (cmd, internal, pkg, plugins, web, config, etc.)
- [x] Initialize Go module
- [x] Basic configuration (config file, logging)

## 2. TS3 Query Integration
- [x] Integrate TS3 Query client (for example go-ts3)
- [x] Connection and reconnect logic
- [x] Event dispatcher for Query events

## 3. Event & Command System
- [x] Event handler architecture
- [x] Command handler with prefix/pattern recognition
- [x] Command registration for plugins
- [~] Server announcement command (web API/UI announcement implemented; bot chat command pending)
- [~] Support channel open/close command (web API/UI support control implemented; bot chat command pending)

## 4. Plugin System
- [x] Plugin interface (Init, Start, Stop, RegisterEvents, RegisterCommands)
- [x] Dynamic loading/unloading of plugins <!-- Runtime load/unload implemented via registry and API -->
- [x] Interfaces for events, commands, and web API

## 5. Core Plugins
- [x] Admin Counter (count admins, update channel) <!-- Channel update TODO -->
- [x] Member Counter (count members, update channel) <!-- Channel update TODO -->
- [x] Combined Stats (combine statistics, channel/web)

## 6. Web Interface (API)
- [x] REST API (for example Gin, Echo, Fiber)
- [x] Authentication: custom login system via API (token auth)
- [x] Switchable between local login and rank system login (provider-based)
- [x] API for status, statistics, plugin management, and configuration
- [x] Plugin configuration can be changed live and stored persistently (GET/POST /api/plugins/config)
- [x] TeamSpeak settings page including API (host, query/voice port, query user/password, bot nickname, default channel, slowmode)
- [x] Bot control page (start/stop/restart) including watchdog control
- [x] Framework server announcement page and API endpoint
- [x] Support control page and API (manual open/close, waiting area, poke texts, auto schedule)
- [x] Support open/close via channel permission switching (join/subscribe power)

## 7. Extensibility & Modularity
- [x] Plugins can register events, commands, and API routes
- [~] Hot reload/reload mechanism for plugins
- [x] Configuration management for plugins
- [x] Dedicated settings pages: AFK Mover (/afkmover.html), Counter (/counter.html), Support (/support.html)

## 8. Extensions (Later)
- [x] Rank system integration
- [ ] Discord sync
- [ ] Plugin store

## 9. Quality & Operations
- [x] Crash watchdog script (auto-restart only on crash/start failure, retry every 60-120 seconds)
- [~] Logging, monitoring, error handling <!-- Logging/error handling partially done, monitoring TODO -->
- [ ] Tests (unit, integration)
- [ ] Documentation (code, API, plugins)

---

This roadmap will be extended and adjusted as needed.