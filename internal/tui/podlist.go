package tui

import (
	"fmt"
	"runpod-top/internal/api"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type PodListModel struct {
	pods     []api.Pod
	cursor   int
	selected map[string]bool
	offset   int
}

func NewPodListModel() PodListModel {
	return PodListModel{
		selected: make(map[string]bool),
	}
}

func (p *PodListModel) SetPods(pods []api.Pod) {
	p.pods = pods
	if p.cursor >= len(pods) && len(pods) > 0 {
		p.cursor = len(pods) - 1
	}
	// Prune selections for pods that no longer exist
	ids := make(map[string]bool, len(pods))
	for _, pod := range pods {
		ids[pod.ID] = true
	}
	for id := range p.selected {
		if !ids[id] {
			delete(p.selected, id)
		}
	}
}

func (p PodListModel) Update(msg tea.KeyMsg, keys KeyMap) (PodListModel, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.Up):
		if p.cursor > 0 {
			p.cursor--
		}
	case key.Matches(msg, keys.Down):
		if p.cursor < len(p.pods)-1 {
			p.cursor++
		}
	case key.Matches(msg, keys.Select):
		if len(p.pods) > 0 {
			id := p.pods[p.cursor].ID
			if p.selected[id] {
				delete(p.selected, id)
			} else {
				p.selected[id] = true
			}
		}
	}
	return p, nil
}

func (p PodListModel) CursorPod() *api.Pod {
	if len(p.pods) == 0 || p.cursor >= len(p.pods) {
		return nil
	}
	pod := p.pods[p.cursor]
	return &pod
}

func (p PodListModel) SelectedPods() []api.Pod {
	if len(p.selected) > 0 {
		var result []api.Pod
		for _, pod := range p.pods {
			if p.selected[pod.ID] {
				result = append(result, pod)
			}
		}
		if len(result) > 0 {
			return result
		}
	}
	if pod := p.CursorPod(); pod != nil {
		return []api.Pod{*pod}
	}
	return nil
}

func (p PodListModel) View(width, height int) string {
	if len(p.pods) == 0 {
		return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center,
			labelStyle.Render("No pods found"))
	}

	title := titleStyle.Render("Pods")
	headerLine := labelStyle.Render(fmt.Sprintf(" %d pod(s)", len(p.pods)))

	visibleRows := height - 2
	if visibleRows < 1 {
		visibleRows = 1
	}

	if p.cursor < p.offset {
		p.offset = p.cursor
	}
	if p.cursor >= p.offset+visibleRows {
		p.offset = p.cursor - visibleRows + 1
	}

	var rows []string
	end := p.offset + visibleRows
	if end > len(p.pods) {
		end = len(p.pods)
	}

	for i := p.offset; i < end; i++ {
		pod := p.pods[i]
		row := p.renderRow(pod, i, width)
		rows = append(rows, row)
	}

	content := strings.Join(rows, "\n")
	return lipgloss.JoinVertical(lipgloss.Left, title+headerLine, content)
}

func (p PodListModel) renderRow(pod api.Pod, index, width int) string {
	isCursor := index == p.cursor
	isSelected := p.selected[pod.ID]

	var statusIcon string
	switch pod.DesiredStatus {
	case "RUNNING":
		statusIcon = statusRunning.Render("●")
	case "STOPPED":
		statusIcon = statusStopped.Render("○")
	case "EXITED":
		statusIcon = statusExited.Render("◌")
	default:
		statusIcon = statusOther.Render("◌")
	}

	selectMark := " "
	if isSelected {
		selectMark = selectedStyle.Render("✓")
	}

	gpuInfo := pod.Machine.GpuDisplayName
	if gpuInfo == "" {
		gpuInfo = fmt.Sprintf("%dx GPU", pod.GpuCount)
	}

	gpuSummary := ""
	if pod.Runtime != nil && len(pod.Runtime.Gpus) > 0 {
		total := 0.0
		for _, g := range pod.Runtime.Gpus {
			total += g.GpuUtilPercent
		}
		avg := total / float64(len(pod.Runtime.Gpus))
		gpuSummary = fmt.Sprintf(" %3.0f%%", avg)
	}

	nameWidth := width - 15
	if nameWidth < 8 {
		nameWidth = 8
	}
	name := Truncate(pod.Name, nameWidth)

	line := fmt.Sprintf(" %s %s %-*s %s%s",
		selectMark, statusIcon, nameWidth, name,
		labelStyle.Render(Truncate(gpuInfo, 12)), gpuSummary)

	if isCursor {
		line = cursorStyle.Width(width).Render(line)
	}

	return line
}
