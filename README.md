Here's a professional English README without emojis:

---

# cmdbook - CLI Command Manager

[![Go Report Card](https://goreportcard.com/badge/github.com/pHo9UBenaA/cmdbook)](https://goreportcard.com/report/github.com/pHo9UBenaA/cmdbook)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

A terminal productivity tool for managing frequently used commands with shortcut aliases.

![demo](https://raw.githubusercontent.com/pHo9UBenaA/cmdbook/main/assets/demo.gif)

## Key Features
- **Command Shortcut Management**
  - Store commands with custom aliases
  - Automatic prefix detection (e.g., extracts `git` from `git add .`)
  - Interactive scrollable list view
- **Quick Execution**
  - Execute stored commands with simple shortcuts
- **Cross-Platform Support**
  - Works on macOS, Linux, and Windows WSL
- **Persistent Configuration**
  - Automatically saves to TOML file (`~/.cmdbook.toml`)

## Installation

### Install latest release
```bash
go install github.com/pHo9UBenaA/cmdbook@latest
```

### Local build (development)
```bash
git clone git@github.com:pHo9UBenaA/cmdbook.git
cd cmdbook
task install # or `go install cb.go`
```

## Basic Usage

### Add Command
```bash
# Basic (auto-generate prefix/shortcut)
cb add "docker compose up --build"

# Custom options
cb add "git push origin main" --prefix git --short push-main
```

### Execute Command
```bash
cb exec git push-main
```

### List Commands
```bash
# Interactive view (arrow keys to scroll)
cb list
```

### Remove Command
```bash
cb remove git push-main
```

## Configuration File
Commands are stored in `~/.cmdbook.toml`:

## Dependencies
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [go-toml](https://github.com/pelletier/go-toml) - TOML configuration handling
- [keyboard](https://github.com/eiannone/keyboard) - Keyboard input handling
- [term](https://pkg.go.dev/golang.org/x/term) - Terminal handling

## License
MIT License - See [LICENSE](LICENSE) for details.
