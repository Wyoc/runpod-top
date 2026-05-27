package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type ConfirmAction struct {
	PodID   string
	PodName string
	Action  string
}

var (
	confirmTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("214")).
				MarginBottom(1)

	confirmTextStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("255"))

	confirmHintStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("241")).
				MarginTop(1)

	confirmKeyYes = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("42"))

	confirmKeyNo = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("196"))
)

func actionIcon(action string) string {
	switch action {
	case "start":
		return "▶"
	case "stop":
		return "■"
	case "restart":
		return "↻"
	default:
		return "?"
	}
}

func ConfirmView(action ConfirmAction, width int) string {
	icon := actionIcon(action.Action)
	title := confirmTitleStyle.Render(fmt.Sprintf("%s  Confirm %s", icon, strings.Title(action.Action)))

	body := confirmTextStyle.Render(
		fmt.Sprintf("Are you sure you want to %s pod\n%s?",
			lipgloss.NewStyle().Bold(true).Render(action.Action),
			lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("62")).Render(action.PodName),
		),
	)

	hint := confirmHintStyle.Render(
		confirmKeyYes.Render("[y]") + " confirm    " + confirmKeyNo.Render("[n/esc]") + " cancel",
	)

	content := lipgloss.JoinVertical(lipgloss.Center, title, body, hint)

	boxWidth := 44
	if len(action.PodName) > 28 {
		boxWidth = len(action.PodName) + 16
	}

	return confirmBoxStyle.
		Width(boxWidth).
		Render(content)
}

func OverlayCenter(bg, fg string, bgWidth, bgHeight int) string {
	fgLines := strings.Split(fg, "\n")
	fgHeight := len(fgLines)
	fgWidth := 0
	for _, line := range fgLines {
		if w := lipgloss.Width(line); w > fgWidth {
			fgWidth = w
		}
	}

	bgLines := strings.Split(bg, "\n")
	for len(bgLines) < bgHeight {
		bgLines = append(bgLines, strings.Repeat(" ", bgWidth))
	}

	startRow := (bgHeight - fgHeight) / 2
	startCol := (bgWidth - fgWidth) / 2
	if startRow < 0 {
		startRow = 0
	}
	if startCol < 0 {
		startCol = 0
	}

	for i, fgLine := range fgLines {
		row := startRow + i
		if row >= len(bgLines) {
			break
		}

		bgLine := bgLines[row]
		bgRunes := []rune(bgLine)

		before := string(bgRunes[:min(startCol, len(bgRunes))])
		after := ""
		endCol := startCol + lipgloss.Width(fgLine)
		if endCol < len(bgRunes) {
			after = string(bgRunes[endCol:])
		}

		padding := ""
		if startCol > len(bgRunes) {
			padding = strings.Repeat(" ", startCol-len(bgRunes))
		}

		bgLines[row] = before + padding + fgLine + after
	}

	return strings.Join(bgLines, "\n")
}
