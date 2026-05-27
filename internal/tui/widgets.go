package tui

import (
	"fmt"
	"strings"
)

func RenderBar(percent float64, width int) string {
	if width < 10 {
		width = 10
	}
	barWidth := width - 6
	if barWidth < 4 {
		barWidth = 4
	}

	filled := int(percent / 100.0 * float64(barWidth))
	if filled > barWidth {
		filled = barWidth
	}
	if filled < 0 {
		filled = 0
	}

	var style = barFilledLow
	if percent >= 85 {
		style = barFilledHigh
	} else if percent >= 60 {
		style = barFilledMid
	}

	bar := style.Render(strings.Repeat("█", filled)) +
		barEmpty.Render(strings.Repeat("░", barWidth-filled))

	return fmt.Sprintf("%s %3.0f%%", bar, percent)
}

func FormatUptime(seconds int) string {
	if seconds <= 0 {
		return "—"
	}
	d := seconds / 86400
	h := (seconds % 86400) / 3600
	m := (seconds % 3600) / 60
	s := seconds % 60

	if d > 0 {
		return fmt.Sprintf("%dd %dh %dm", d, h, m)
	}
	if h > 0 {
		return fmt.Sprintf("%dh %dm %ds", h, m, s)
	}
	if m > 0 {
		return fmt.Sprintf("%dm %ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}

func FormatCost(costPerHr float64, uptimeSeconds int) string {
	session := costPerHr * float64(uptimeSeconds) / 3600.0
	return fmt.Sprintf("$%.2f/hr | Session: $%.2f", costPerHr, session)
}

func StatusBadge(status string) string {
	switch status {
	case "RUNNING":
		return statusRunning.Render("● RUNNING")
	case "STOPPED":
		return statusStopped.Render("○ STOPPED")
	case "EXITED":
		return statusExited.Render("◌ EXITED")
	default:
		return statusOther.Render("◌ " + status)
	}
}

func Truncate(s string, maxWidth int) string {
	if len(s) <= maxWidth {
		return s
	}
	if maxWidth <= 3 {
		return s[:maxWidth]
	}
	return s[:maxWidth-3] + "..."
}
