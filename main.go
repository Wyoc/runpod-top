package main

import (
	"flag"
	"fmt"
	"os"
	"runpod-top/internal/api"
	"runpod-top/internal/config"
	"runpod-top/internal/tui"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	configPath := flag.String("config", config.DefaultPath(), "config file path")
	interval := flag.Duration("interval", 0, "polling interval (e.g. 3s, 5s)")
	initConfig := flag.Bool("init-config", false, "create default config file and exit")
	flag.Parse()

	if *initConfig {
		path := *configPath
		if err := config.WriteDefault(path); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating config: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Config file created at %s\n", path)
		return
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config %s: %v\n", *configPath, err)
		os.Exit(1)
	}

	// Precedence: CLI flags > env vars > config file > defaults
	apiKey := cfg.APIKey
	if v := os.Getenv("RUNPOD_API_KEY"); v != "" {
		apiKey = v
	}

	pollInterval := 3 * time.Second
	if cfg.Interval != 0 {
		pollInterval = cfg.Interval
	}
	if *interval != 0 {
		pollInterval = *interval
	}

	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "Error: API key not configured.")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Set it via one of:")
		fmt.Fprintf(os.Stderr, "  1. Config file: %s\n", *configPath)
		fmt.Fprintln(os.Stderr, "  2. Environment: export RUNPOD_API_KEY=<key>")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintf(os.Stderr, "Run `runpod-top --init-config` to create a default config file.\n")
		fmt.Fprintln(os.Stderr, "Get your API key from https://console.runpod.io/")
		os.Exit(1)
	}

	client := api.NewClient(apiKey)
	model := tui.NewModel(client, pollInterval)

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
