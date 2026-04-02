# Changelog

All notable changes to this project will be documented in this file.

The format is based on Keep a Changelog,
and this project follows Semantic Versioning.

## [Unreleased]

### Added
- README and LICENSE added for repository documentation and licensing.
- AFK Mover option to return users to their previous channel when they become active again.
- AFK Mover picker dialog and manual ID entry for excluded channels.
- Framework platform mode settings (TeamSpeak, Discord, Both).
- Dedicated Discord settings page in the web UI with bot access, behavior, channels, and role mappings.
- Discord REST endpoints for settings and Discord metadata:
  - `/api/settings/discord`
  - `/api/settings/discord/channels`
  - `/api/settings/discord/roles`
- Discord channel and role picker dialogs in the web UI (including multi-select role assignment).
- Discord runtime client integration using `discordgo` with connect/disconnect and guild validation.
- Discord AFK inactivity timeout setting (minutes) to disconnect inactive voice users.
- Framework page server announcement feature in the web UI (message input + API endpoint).
- Announcement scheduling support for one-time sends and configurable interval send counts.
- Startup security self-checks for internet mode with explicit `SECURITY CRITICAL/WARN` logs.
- Support management page in the web UI with:
  - Support channel selection
  - Waiting area selection
  - Configurable poke message for support open/close
  - Configurable join poke messages for open and closed support state
  - Configurable supporter poke message and supporter server group selection
  - Manual open/close actions
  - Automatic open/close scheduling by configured times
- New support API endpoints for settings, status, and actions.
- New TeamSpeak server group settings endpoint for loading selectable supporter groups in the web UI.
- TeamSpeak helper methods for support workflows:
  - Client poke message sending
  - Channel permission-based open/close switching
  - Server group listing for web configuration

### Changed
- AFK Mover return-on-activity option now uses a slider toggle in the web UI.
- AFK Mover excluded channel configuration now uses the same picker-plus-ID workflow as other TeamSpeak channel selections.
- Support channel selection, waiting area selection, and supporter group selection now use a consistent picker-plus-ID input workflow in the web UI.
- Support channel open/close now uses real channel permissions (join/subscribe power) instead of channel name prefixes.
- Runtime startup now respects configured platform mode and starts TeamSpeak, Discord, or both accordingly.
- Announcement delivery now supports Discord announcement channels when Discord runtime is active.
- Discord AFK behavior now disconnects users after configurable voice inactivity instead of disconnecting on AFK-channel join.
- Web server can run HTTP and HTTPS in parallel and supports optional HTTP->HTTPS redirect mode via environment flags.
- Security hardening supports environment-based secret overrides with optional strict startup enforcement.

### Fixed
- Staticcheck warning QF1003 in the announcement scheduler logic by replacing the `if / else if` chain with a `switch`.
- Missing translation binding for the closed-support join poke field on the support settings page.

## [1.0.0] - 2026-03-30

### Added
- Modular TeamSpeak 3 bot framework in Go with web UI and REST API.
- Plugin system with core plugins:
  - AdminCounter
  - MemberCounter
  - CombinedStats
  - AfkMover
- Web pages for dashboard, bot control, plugins, counter settings, AFK mover, TeamSpeak settings, and TS3 connection test.
- Bot lifecycle control via API and scripts (start, stop, restart, status).
- Watchdog support for automatic restart on crash/start failure.
- Localization support (German and English).

### Changed
- Authentication expanded to support four modes:
  - none
  - local
  - ranksystem
  - local_ranksystem
- Auth settings UI extended with mode switch and security warning for disabled login.
- Rank-System health handling adjusted to return mode-aware and configuration-aware responses.

### Fixed
- Rank-System login compatibility improved by supporting both JSON and form-based external login requests.
- Rank-System success detection improved for redirect-based login flows.

[Unreleased]: https://github.com/<owner>/<repo>/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/<owner>/<repo>/releases/tag/v1.0.0
