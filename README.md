# UC-Framework

UC-Framework is a modular TeamSpeak 3 and discord bot framework written in Go, with a web interface, plugin system, bot lifecycle control, and configurable authentication.

compatible with: Linux and Windows

## Table of Contents

1. [What is UC-Framework?](#what-is-uc-framework)
2. [Features](#features)
3. [Installation (Windows and Linux)](#installation-windows-and-linux)
4. [Starting and Running](#starting-and-running)
5. [Use Cases](#use-cases)
6. [Configuration](#configuration)
7. [Project Structure](#project-structure)
8. [Troubleshooting](#troubleshooting)
9. [Development](#development)
10. [Project Docs](#project-docs)

## What is UC-Framework?

UC-Framework is a centralized runtime and management layer for TeamSpeak 3 automation.

- Written in Go (Go 1.21)
- Web interface running on port `8080`
- Configuration via `config/config.json`
- Plugin-based architecture for common community tasks
- Bot lifecycle control (start/stop/restart + watchdog)

The goal is to manage TS3 server automation through a single API and UI instead of maintaining multiple disconnected scripts.

## Features

### Core

- TeamSpeak query connection and status handling
- Modular plugin system
- REST API for status, plugins, and configuration
- Web UI for operations and settings
- German/English language support

### Built-in Plugins

- `AdminCounter`: counts online admins based on configured server groups
- `MemberCounter`: counts online members with exclusion rules
- `CombinedStats`: combines statistics for dashboard/API
- `AfkMover`: moves inactive users to the AFK channel

### Web Interface

- Dashboard and plugin status overview
- Plugin enable/disable controls
- Live plugin configuration save/load
- TS3 connection test page
- TS3 settings page (including channel selection)
- Bot control page (including watchdog)

### Authentication

Supported login modes:

- `none` (no login)
- `local` (local login)
- `ranksystem` (TSN-Ranksystem login)
- `local_ranksystem` (local login + ranksystem fallback)

Important: `none` is insecure and should only be used in isolated test environments.

## Installation (Windows and Linux)

### Requirements

- Go `1.21` or newer
- Access to a TS3 server with query credentials
- Write access to:
  - `config/`
  - `runtime/`

### 1) Clone the repository

```bash
git clone <YOUR_REPO_URL>
cd UC-Framework
```

### 2) Download dependencies

```bash
go mod download
```

### 3) Configure the project

Edit `config/config.json`.

Most important fields:

- `ts3.host`, `ts3.query_port`, `ts3.username`, `ts3.password`
- `web_auth.enabled`
- `web_auth.provider`
- optional: `web_auth.ranksystem.login_url`

## Starting and Running

### Quick start

```bash
go run .
```

Then open:

- `http://localhost:8080`

### Build and run binary

```bash
go build -o uc-framework .
./uc-framework
```

On Windows, the binary will be `uc-framework.exe`.

### Windows notes

- Supervisor/watchdog scripts:
  - `scripts/bot-supervisor.ps1`
  - `scripts/bot-watchdog.ps1`
- If execution policies block scripts, run with an appropriate policy in controlled environments.

### Linux notes

- Supervisor/watchdog scripts:
  - `scripts/bot-supervisor.sh`
  - `scripts/bot-watchdog.sh`
- Make scripts executable:

```bash
chmod +x scripts/*.sh
```

## Use Cases

Typical scenarios:

- Community servers with automated counter/stat channels
- TeamSpeak servers with AFK automation
- Browser-based bot operation instead of shell-only workflows
- Plugin-based extension for custom behavior
- Operation with local auth or external ranksystem auth

## Configuration

The main configuration file is `config/config.json`.

### Key Sections

- `ts3`: TeamSpeak query connection settings
- `bot_control`: bot process control and watchdog intervals
- `web_auth`: authentication settings
- `plugin_configs`: plugin-specific settings

### Auth Modes (Summary)

- `none`: no login (insecure)
- `local`: local credentials only
- `ranksystem`: external TSN-Ranksystem login
- `local_ranksystem`: local first, then ranksystem fallback

## Project Structure

```text
cmd/                # alternative entry point
config/             # runtime configuration
internal/           # core, bot, TS3 client
plugins/            # plugin implementations
scripts/            # supervisor/watchdog (Windows + Linux)
web/                # API + static web interface
runtime/            # PIDs, state, runtime data
```

## Troubleshooting

### Web interface not reachable

- Check whether the process is running
- Check if port `8080` is already in use
- Check firewall/reverse proxy configuration

### TS3 connection failed

- Verify host/port/query user/password in `config/config.json`
- Verify query permissions on your TS3 server
- Use the TS3 connection test page in the web interface

### Login does not work

- Verify auth mode in auth settings
- For `ranksystem`/`local_ranksystem`: set `web_auth.ranksystem.login_url`
- For `local`: verify local credentials/password hash

## Development

### Local commands

```bash
go test ./...
go build ./...
```

### Recommended workflow

1. Create a feature branch
2. Commit focused changes with clear messages
3. Run build/tests locally
4. Open a pull request

## Project Docs

- Changelog: `CHANGELOG.md`
- Security policy: `SECURITY.md`
- Contributing guide: `CONTRIBUTING.md`
- License: `LICENSE`
