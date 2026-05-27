package tui

import (
	"fmt"
	"runpod-top/internal/api"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DetailModel struct {
	pods     []api.Pod
	viewport viewport.Model
	width    int
	height   int
}

func NewDetailModel() DetailModel {
	return DetailModel{
		viewport: viewport.New(0, 0),
	}
}

func (d *DetailModel) SetPods(pods []api.Pod) {
	d.pods = pods
}

func (d DetailModel) Update(msg tea.KeyMsg, keys KeyMap) (DetailModel, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.ScrollUp):
		d.viewport.LineUp(5)
	case key.Matches(msg, keys.ScrollDown):
		d.viewport.LineDown(5)
	case key.Matches(msg, keys.Up):
		d.viewport.LineUp(1)
	case key.Matches(msg, keys.Down):
		d.viewport.LineDown(1)
	}
	return d, nil
}

func (d DetailModel) View(width, height int) string {
	if len(d.pods) == 0 {
		return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center,
			labelStyle.Render("Select a pod to view details"))
	}

	var content string
	if len(d.pods) == 1 {
		content = d.renderSinglePod(d.pods[0], width)
	} else {
		content = d.renderMultiPod(width)
	}

	d.viewport.Width = width
	d.viewport.Height = height
	d.viewport.SetContent(content)

	return d.viewport.View()
}

func (d DetailModel) renderSinglePod(pod api.Pod, width int) string {
	var sections []string
	barWidth := width - 16
	if barWidth < 15 {
		barWidth = 15
	}

	header := titleStyle.Render(pod.Name) + "  " + StatusBadge(pod.DesiredStatus)
	sections = append(sections, header)

	info := fmt.Sprintf("  %s  %s  %s",
		labelStyle.Render("GPU:")+valueStyle.Render(" "+pod.Machine.GpuDisplayName),
		labelStyle.Render("Loc:")+valueStyle.Render(" "+pod.Machine.Location),
		labelStyle.Render("vCPU:")+valueStyle.Render(fmt.Sprintf(" %d", pod.VcpuCount)),
	)
	sections = append(sections, info)

	uptime := pod.UptimeSeconds
	if pod.Runtime != nil {
		uptime = pod.Runtime.UptimeInSeconds
	}
	costLine := fmt.Sprintf("  %s  %s",
		labelStyle.Render("Uptime:")+valueStyle.Render(" "+FormatUptime(uptime)),
		labelStyle.Render("Cost:")+valueStyle.Render(" "+FormatCost(pod.CostPerHr, uptime)),
	)
	sections = append(sections, costLine)
	sections = append(sections, "")

	if pod.Runtime != nil {
		if len(pod.Runtime.Gpus) > 0 {
			sections = append(sections, "  "+labelStyle.Render("GPU Utilization"))
			for i, gpu := range pod.Runtime.Gpus {
				label := fmt.Sprintf("  GPU %d ", i)
				sections = append(sections, label+RenderBar(gpu.GpuUtilPercent, barWidth))
			}
			sections = append(sections, "")

			sections = append(sections, "  "+labelStyle.Render("GPU VRAM"))
			for i, gpu := range pod.Runtime.Gpus {
				label := fmt.Sprintf("  GPU %d ", i)
				sections = append(sections, label+RenderBar(gpu.MemoryUtilPercent, barWidth))
			}
			sections = append(sections, "")
		}

		sections = append(sections, "  "+labelStyle.Render("Container"))
		sections = append(sections, "  CPU    "+RenderBar(pod.Runtime.Container.CpuPercent, barWidth))
		sections = append(sections, "  Memory "+RenderBar(pod.Runtime.Container.MemoryPercent, barWidth))
		sections = append(sections, "")

		if len(pod.Runtime.Ports) > 0 {
			sections = append(sections, "  "+labelStyle.Render("Ports"))
			for _, port := range pod.Runtime.Ports {
				visibility := "private"
				if port.IsIpPublic {
					visibility = "public"
				}
				line := fmt.Sprintf("    %d → %d  %s  %s (%s)",
					port.PrivatePort, port.PublicPort,
					port.IP, port.Type, visibility)
				sections = append(sections, line)
			}
		}
	} else {
		sections = append(sections, "  "+labelStyle.Render("No runtime data (pod not running)"))
	}

	return strings.Join(sections, "\n")
}

func (d DetailModel) renderMultiPod(width int) string {
	var sections []string
	barWidth := width - 20
	if barWidth < 12 {
		barWidth = 12
	}

	for i, pod := range d.pods {
		if i > 0 {
			sections = append(sections, strings.Repeat("─", width))
		}

		header := fmt.Sprintf(" %s  %s  %s",
			titleStyle.Render(pod.Name),
			StatusBadge(pod.DesiredStatus),
			labelStyle.Render(pod.Machine.GpuDisplayName),
		)
		sections = append(sections, header)

		if pod.Runtime != nil {
			for j, gpu := range pod.Runtime.Gpus {
				line := fmt.Sprintf("  GPU %d  %s  VRAM %s",
					j,
					RenderBar(gpu.GpuUtilPercent, barWidth/2),
					RenderBar(gpu.MemoryUtilPercent, barWidth/2),
				)
				sections = append(sections, line)
			}
			line := fmt.Sprintf("  CPU %s  Mem %s",
				RenderBar(pod.Runtime.Container.CpuPercent, barWidth/2),
				RenderBar(pod.Runtime.Container.MemoryPercent, barWidth/2),
			)
			sections = append(sections, line)

			uptime := pod.Runtime.UptimeInSeconds
			costLine := fmt.Sprintf("  %s",
				labelStyle.Render("Cost:")+valueStyle.Render(" "+FormatCost(pod.CostPerHr, uptime)),
			)
			sections = append(sections, costLine)
		} else {
			sections = append(sections, "  "+labelStyle.Render("Not running"))
		}
	}

	return strings.Join(sections, "\n")
}
