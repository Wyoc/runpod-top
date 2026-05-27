package tui

import "github.com/charmbracelet/lipgloss"

var (
	focusedBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("62"))

	unfocusedBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("240"))

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("62")).
			PaddingLeft(1)

	statusRunning = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	statusStopped = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	statusExited  = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	statusOther   = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	cursorStyle   = lipgloss.NewStyle().Background(lipgloss.Color("237")).Bold(true)
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))

	labelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	valueStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("255"))

	barFilledLow  = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	barFilledMid  = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	barFilledHigh = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	barEmpty      = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

	errorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
	statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))

	confirmBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("214")).
			Padding(1, 2)
)
