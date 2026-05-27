package tui

import (
	"context"
	"runpod-top/internal/api"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	client   *api.Client
	interval time.Duration

	pods   []api.Pod
	err    error
	width  int
	height int

	podList PodListModel
	detail  DetailModel
	helpBar help.Model

	focus      int
	confirming *ConfirmAction
	statusMsg  string
	statusTTL  int

	keys KeyMap
}

func NewModel(client *api.Client, interval time.Duration) Model {
	return Model{
		client:   client,
		interval: interval,
		podList:  NewPodListModel(),
		detail:   NewDetailModel(),
		helpBar:  help.New(),
		keys:     DefaultKeyMap(),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		fetchPodsCmd(m.client),
		tickCmd(m.interval),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.helpBar.Width = msg.Width
		return m, nil

	case TickMsg:
		if m.statusTTL > 0 {
			m.statusTTL--
			if m.statusTTL == 0 {
				m.statusMsg = ""
			}
		}
		return m, tea.Batch(fetchPodsCmd(m.client), tickCmd(m.interval))

	case PodsUpdatedMsg:
		m.pods = msg.Pods
		m.err = nil
		m.podList.SetPods(msg.Pods)
		m.detail.SetPods(m.podList.SelectedPods())
		return m, nil

	case APIErrorMsg:
		m.err = msg.Err
		return m, nil

	case ActionResultMsg:
		if msg.Err != nil {
			m.statusMsg = "Failed to " + msg.Action + ": " + msg.Err.Error()
		} else {
			m.statusMsg = "Pod " + msg.Action + " successful"
		}
		m.statusTTL = 5
		return m, fetchPodsCmd(m.client)

	case tea.KeyMsg:
		if m.confirming != nil {
			return m.handleConfirmKey(msg)
		}
		return m.handleKey(msg)
	}

	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Quit):
		return m, tea.Quit

	case key.Matches(msg, m.keys.Help):
		m.helpBar.ShowAll = !m.helpBar.ShowAll
		return m, nil

	case key.Matches(msg, m.keys.Tab):
		m.focus = (m.focus + 1) % 2
		return m, nil

	case key.Matches(msg, m.keys.Start):
		if pod := m.podList.CursorPod(); pod != nil {
			m.confirming = &ConfirmAction{PodID: pod.ID, PodName: pod.Name, Action: "start"}
		}
		return m, nil

	case key.Matches(msg, m.keys.Stop):
		if pod := m.podList.CursorPod(); pod != nil {
			m.confirming = &ConfirmAction{PodID: pod.ID, PodName: pod.Name, Action: "stop"}
		}
		return m, nil

	case key.Matches(msg, m.keys.Restart):
		if pod := m.podList.CursorPod(); pod != nil {
			m.confirming = &ConfirmAction{PodID: pod.ID, PodName: pod.Name, Action: "restart"}
		}
		return m, nil
	}

	if m.focus == 0 {
		var cmd tea.Cmd
		m.podList, cmd = m.podList.Update(msg, m.keys)
		m.detail.SetPods(m.podList.SelectedPods())
		return m, cmd
	}

	var cmd tea.Cmd
	m.detail, cmd = m.detail.Update(msg, m.keys)
	return m, cmd
}

func (m Model) handleConfirmKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Confirm):
		action := m.confirming
		m.confirming = nil
		return m, executePodAction(m.client, action.PodID, action.Action)
	case key.Matches(msg, m.keys.Cancel):
		m.confirming = nil
		return m, nil
	}
	return m, nil
}

func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	if m.width < 60 || m.height < 10 {
		return "Terminal too small. Resize to at least 60x10."
	}

	helpView := helpStyle.Render(m.helpBar.View(m.keys))
	helpHeight := lipgloss.Height(helpView)

	statusLine := ""
	if m.statusMsg != "" {
		statusLine = statusStyle.Render(m.statusMsg)
	}
	if m.err != nil {
		statusLine = errorStyle.Render("API: " + m.err.Error())
	}
	statusHeight := 0
	if statusLine != "" {
		statusHeight = 1
	}

	availableHeight := m.height - helpHeight - statusHeight - 2

	leftWidth := m.width * 35 / 100
	if leftWidth < 30 {
		leftWidth = 30
	}
	rightWidth := m.width - leftWidth

	panelHeight := availableHeight - 2

	leftBorder := unfocusedBorderStyle
	rightBorder := unfocusedBorderStyle
	if m.focus == 0 {
		leftBorder = focusedBorderStyle
	} else {
		rightBorder = focusedBorderStyle
	}

	innerLeftW := leftWidth - 2
	innerRightW := rightWidth - 2
	if innerLeftW < 1 {
		innerLeftW = 1
	}
	if innerRightW < 1 {
		innerRightW = 1
	}

	leftPanel := leftBorder.
		Width(innerLeftW).
		Height(panelHeight).
		Render(m.podList.View(innerLeftW, panelHeight))

	rightPanel := rightBorder.
		Width(innerRightW).
		Height(panelHeight).
		Render(m.detail.View(innerRightW, panelHeight))

	mainArea := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)

	parts := []string{mainArea}
	if statusLine != "" {
		parts = append(parts, statusLine)
	}
	parts = append(parts, helpView)

	screen := lipgloss.JoinVertical(lipgloss.Left, parts...)

	if m.confirming != nil {
		popup := ConfirmView(*m.confirming, m.width)
		screen = OverlayCenter(screen, popup, m.width, m.height)
	}

	return screen
}

func tickCmd(interval time.Duration) tea.Cmd {
	return tea.Tick(interval, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func fetchPodsCmd(client *api.Client) tea.Cmd {
	return func() tea.Msg {
		pods, err := client.FetchPods(context.Background())
		if err != nil {
			return APIErrorMsg{Err: err}
		}
		return PodsUpdatedMsg{Pods: pods}
	}
}

func executePodAction(client *api.Client, podID, action string) tea.Cmd {
	return func() tea.Msg {
		var err error
		switch action {
		case "start":
			err = client.StartPod(context.Background(), podID)
		case "stop":
			err = client.StopPod(context.Background(), podID)
		case "restart":
			err = client.StopPod(context.Background(), podID)
			if err == nil {
				err = client.StartPod(context.Background(), podID)
			}
		}
		return ActionResultMsg{PodID: podID, Action: action, Err: err}
	}
}
