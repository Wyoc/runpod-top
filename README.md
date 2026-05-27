# runpod-top

A terminal UI for monitoring RunPod GPU pods in real-time. Think `htop`, but for your RunPod instances.

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go&logoColor=white)

## Features

- **Live telemetry** — GPU utilization, VRAM, container CPU/memory updated every few seconds
- **Split-pane dashboard** — pod list on the left, detailed metrics with visual bars on the right
- **Multi-select** — select multiple pods to compare metrics side-by-side
- **Pod management** — start, stop, and restart pods directly from the TUI
- **Confirmation dialogs** — destructive actions require explicit confirmation

## Install

```bash
go install runpod-top@latest
```

Or build from source:

```bash
git clone https://github.com/jdupuy/runpod-top.git
cd runpod-top
go build -o runpod-top .
```

## Usage

```bash
export RUNPOD_API_KEY="your-api-key"
./runpod-top
```

Get your API key from [console.runpod.io](https://console.runpod.io/).

### Configuration

runpod-top looks for a config file at `~/.config/runpod-top/config.toml`. Generate a default one with:

```bash
runpod-top --init-config
```

This creates:

```toml
# runpod-top configuration
# Get your API key from https://console.runpod.io/

# api_key = ""

# Polling interval (e.g. "3s", "5s", "10s")
# interval = "3s"
```

**Precedence:** CLI flags > environment variables > config file > defaults.

| Source | API key | Interval |
|--------|---------|----------|
| Config file | `api_key` | `interval` |
| Environment | `RUNPOD_API_KEY` | — |
| CLI flag | — | `-interval` |

### Options

| Flag | Default | Description |
|------|---------|-------------|
| `-config` | `~/.config/runpod-top/config.toml` | Config file path |
| `-interval` | `3s` | Polling interval (e.g. `5s`, `10s`, `1m`) |
| `-init-config` | | Create default config file and exit |

```bash
./runpod-top -interval 5s
./runpod-top -config /path/to/config.toml
```

## Key Bindings

| Key | Action |
|-----|--------|
| `j` / `k` / `arrows` | Navigate pod list |
| `Space` | Toggle multi-select |
| `Tab` | Switch focus between panels |
| `s` | Start selected pod |
| `x` | Stop selected pod |
| `r` | Restart selected pod |
| `Ctrl+u` / `Ctrl+d` | Scroll detail panel |
| `?` | Toggle full help |
| `q` / `Ctrl+c` | Quit |

## Metrics Displayed

**Per pod:**

- GPU utilization % (per GPU, with color-coded bar)
- GPU VRAM usage % (per GPU)
- Container CPU %
- Container memory %
- Uptime
- Cost per hour and estimated session cost
- Status (running/stopped/exited)
- GPU type and datacenter location
- Port mappings

## Architecture

Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) (Elm Architecture) and [Lip Gloss](https://github.com/charmbracelet/lipgloss) for styling. Polls the RunPod GraphQL API — no websockets or agents required.

```
runpod-top/
  main.go                  # Entry point
  internal/
    api/                   # RunPod GraphQL client
      client.go            # Query pods
      mutations.go         # Start/stop pods
      types.go             # Response types
    config/
      config.go            # TOML config loading
    tui/                   # Terminal UI
      model.go             # Root model (Init/Update/View)
      podlist.go           # Left panel — pod list
      detail.go            # Right panel — metrics
      confirm.go           # Confirmation popup
      keys.go              # Key bindings
      styles.go            # Colors and layout
      widgets.go           # Progress bars, formatting
      messages.go          # Message types
```

## License

MIT
