# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Git CZ Go is a CLI tool for creating Conventional Commit messages with an interactive TUI interface using Bubble Tea. The project follows Clean Architecture principles and is built with Go 1.24.6.

## Key Architecture

- **Entry Point**: `main.go` delegates to `cmd.Execute()`
- **CLI Layer**: `cmd/` uses Cobra for command definitions
- **Configuration**: `config/` handles YAML-based configuration loading
- **Business Logic**: `internal/app/` contains the main application logic
- **UI Models**: `internal/entity/` defines Bubble Tea models for TUI interaction

## Essential Commands

### Development Setup
```bash
# Install development tools
mise install

# Check tool versions  
mise current
```

### Build and Run
```bash
# Build the project
go build -o git-cz .

# Run directly
go run main.go

# Install to GOPATH/bin
go install
```

### Code Quality (Run Before Commits)
```bash
# Format code
go fmt ./...

# Run linter
golangci-lint run

# Clean up modules
go mod tidy

# Run tests
go test ./...
```

## Development Patterns

### Configuration Loading Pattern
```go
cfg, err := config.LoadConfig(configPath)
if err != nil {
    fmt.Fprintf(os.Stderr, "Error loading config: %s\n", err)
    os.Exit(1)
}
```

### Bubble Tea Model Structure
All TUI models in `internal/entity/` follow the standard Bubble Tea interface:
- `Init() tea.Cmd`
- `Update(tea.Msg) (tea.Model, tea.Cmd)`
- `View() string`

### Error Handling Convention
```go
if err != nil {
    fmt.Fprintf(os.Stderr, "Error: %s\n", err)
    os.Exit(1)
}
```

## Code Style

- **Go files**: Use tabs for indentation (defined in .editorconfig)
- **Other files**: 2 spaces indentation
- **Line endings**: LF (Unix style)
- **Encoding**: UTF-8
- **Package structure**: Follow `cmd/`, `config/`, `internal/` organization

## Dependencies

- **CLI Framework**: Cobra (`github.com/spf13/cobra`)
- **TUI Framework**: Bubble Tea (`github.com/charmbracelet/bubbletea`)
- **Configuration**: YAML v3 (`gopkg.in/yaml.v3`)
- **Development Tools**: Managed via mise.toml (golangci-lint, delve debugger)

## Testing and Debugging

```bash
# Run specific package tests
go test ./internal/...

# Debug with delve
dlv debug

# Test coverage
go test -cover ./...
```

## Important Files

- `config/config.yaml`: Default configuration template
- `mise.toml`: Development tool version management
- `.editorconfig`: Code formatting rules (tabs for Go, spaces for others)