# Security Policy

## Supported Versions

The following versions currently receive security updates:

| Version | Supported |
| --- | --- |
| 1.0.x | :white_check_mark: |
| < 1.0.0 | :x: |

Note: Security fixes are published to `main` first and then included in a release.

## Security Scope

This policy applies to:

- Web interface and API (`web/`)
- Authentication and session handling
- TS3 connection and bot control
- Plugin runtime and plugin configuration
- Supervisor/watchdog scripts (`scripts/`)

Out of scope:

- Security issues in third-party services operated separately (for example external TSN Ranksystem), unless UC-Framework code has a direct impact.

## Reporting a Vulnerability

Please do **not** report vulnerabilities publicly via GitHub issues.

Instead, please report responsibly and include:

- Title and short description
- Affected version/commit
- Reproduction steps (PoC)
- Expected vs. actual behavior
- Risk assessment (for example auth bypass, RCE, data leak)
- Suggested mitigation (if possible)

Contact:

- Security E-Mail: `zitroneemu139@proton.me`
- Alternatively: private message to the repository owner

If you want to change the contact path, update this file accordingly.

## Response Process

Target timelines (best effort):

- Initial acknowledgment: within 72 hours
- First technical assessment: within 7 days
- Status updates: at least every 14 days until resolution

By severity:

- Critical: prioritized handling
- High/Medium/Low: based on risk and impact

## Disclosure Policy

- Coordinated disclosure after a fix is available
- No public disclosure of sensitive details before a fix
- After release: changelog entry and short security note

## Security Best Practices for Operators

- Use auth mode `none` only in isolated test environments
- Use strong passwords and, if possible, password hashes instead of plaintext
- Protect `config/config.json` and runtime files (`runtime/`) from unauthorized access
- Expose the web interface only internally or behind reverse proxy/TLS
- Update regularly to new releases and security fixes
