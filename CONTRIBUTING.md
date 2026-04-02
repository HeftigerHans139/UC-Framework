# Contributing to UC-Framework

Thanks for your interest in contributing.
This document explains how to propose changes, submit pull requests, and keep contributions consistent.

## Code of Conduct

Please be respectful and constructive in discussions and reviews.

## Before You Start

- Check open issues and pull requests to avoid duplicate work.
- For security issues, do not open a public issue. Follow the policy in SECURITY.md.
- For larger changes, open an issue first and discuss the approach.

## Development Setup

### Requirements

- Go 1.21+
- Git
- Access to a TS3 test environment (recommended for integration testing)

### Run locally

```bash
go mod download
go run .
```

Web interface default URL:

- http://localhost:8080

### Build and test

```bash
go test ./...
go build ./...
```

## Branch and Commit Guidelines

### Branch naming

Use short, descriptive branch names, for example:

- feature/auth-mode-switch
- fix/ranksystem-health-response
- docs/readme-cleanup

### Commit messages

Prefer clear, action-oriented commit messages, for example:

- feat(auth): add local+ranksystem mode
- fix(api): return mode-aware health response
- docs: add security and changelog files

## Pull Request Process

1. Fork the repository (or create a branch if you have direct access).
2. Keep changes focused and small when possible.
3. Update documentation for user-facing behavior changes.
4. Run tests and build locally.
5. Open a pull request with a clear description.

### PR checklist

- [ ] Code builds successfully (`go build ./...`)
- [ ] Tests pass (`go test ./...`)
- [ ] No unrelated files were changed
- [ ] Docs updated (README/CHANGELOG/SECURITY) if applicable
- [ ] Config changes are explained
- [ ] UI changes include screenshots (if relevant)

## Coding Guidelines

### Go

- Use `gofmt` for formatting.
- Keep functions focused and readable.
- Avoid introducing global side effects unless required.
- Prefer explicit error handling and clear error messages.

### Frontend (web/static)

- Keep JavaScript simple and framework-free unless discussed.
- Reuse existing translation keys and patterns in i18n.js.
- Ensure new UI text is added in both English and German where needed.

### Configuration

- Do not commit secrets.
- If new config fields are added:
  - wire them through core and API consistently,
  - update documentation,
  - provide safe defaults.

## Issue Reporting

When opening an issue, include:

- What happened
- What you expected
- Steps to reproduce
- Environment (OS, Go version, config context)
- Logs or screenshots where helpful

## Security Contributions

- Follow SECURITY.md for disclosure.
- Do not publish exploit details before a fix is available.

## License

By contributing, you agree that your contributions are provided under the repository's proprietary license (see LICENSE).
