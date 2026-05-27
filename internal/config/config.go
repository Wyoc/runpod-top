package config

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
)

type Config struct {
	APIKey   string        `toml:"api_key"`
	Interval time.Duration `toml:"interval"`
}

func DefaultPath() string {
	if dir, err := os.UserConfigDir(); err == nil {
		return filepath.Join(dir, "runpod-top", "config.toml")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "runpod-top", "config.toml")
}

func Load(path string) (Config, error) {
	var cfg Config
	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return cfg, nil
	}
	if err != nil {
		return cfg, err
	}
	_, err = toml.DecodeFile(path, &cfg)
	return cfg, err
}

func WriteDefault(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	content := `# runpod-top configuration
# Get your API key from https://console.runpod.io/

# api_key = ""

# Polling interval (e.g. "3s", "5s", "10s")
# interval = "3s"
`
	return os.WriteFile(path, []byte(content), 0o644)
}
