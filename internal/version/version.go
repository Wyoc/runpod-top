package version

import "fmt"

// Populated at build time via -ldflags "-X runpod-top/internal/version.Version=..."
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

func String() string {
	return fmt.Sprintf("runpod-top %s (commit %s, built %s)", Version, Commit, Date)
}
