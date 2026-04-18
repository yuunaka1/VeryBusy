package ui

import (
	"fmt"
	"math/rand"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	
	"github.com/yuunaka1/VeryBusy/internal/sim"
	"github.com/yuunaka1/VeryBusy/internal/theme"
)

// MainModel wraps the currently active view
type MainModel struct {
	engine *sim.Engine
	mode   string // "soc", "logs", "alerts", ...
	width  int
	height int

	// TUI components (views)
	logsView    *LogsView
	alertsView  *AlertsView
	assetsView  *AssetsView
	graphsView  *GraphsView
	networkView *NetworkView
	
	// Layout properties
	splitLayout bool
}

type tickMsg time.Time
type logTickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func logTickCmd() tea.Cmd {
	// Random interval between 100ms and 800ms
	d := time.Duration(rand.Intn(700)+100) * time.Millisecond
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return logTickMsg(t)
	})
}

func NewMainModel(engine *sim.Engine, mode string) *MainModel {
	vLogs := NewLogsView(engine)
	vAlerts := NewAlertsView(engine)
	vAssets := NewAssetsView(engine)
	vGraphs := NewGraphsView(engine)
	vNetwork := NewNetworkView(engine)
	
	return &MainModel{
		engine:      engine,
		mode:        mode,
		splitLayout: mode == "soc",
		logsView:    vLogs,
		alertsView:  vAlerts,
		assetsView:  vAssets,
		graphsView:  vGraphs,
		networkView: vNetwork,
	}
}

func (m *MainModel) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen, tickCmd(), logTickCmd())
}

func (m *MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		
		statusHeight := 1
		contentHeight := m.height - statusHeight - 1 // -1 for general padding

		// Update sub-views based on mode
		if m.splitLayout {
			// Split into left 2 screens and right 3 screens
			leftW := m.width / 2
			rightW := m.width - leftW
			
			leftH1 := contentHeight / 2
			leftH2 := contentHeight - leftH1
			
			rightH1 := contentHeight / 3
			rightH2 := contentHeight / 3
			rightH3 := contentHeight - rightH1 - rightH2
			
			m.logsView.SetSize(leftW, leftH1)
			m.alertsView.SetSize(leftW, leftH2)
			
			m.graphsView.SetSize(rightW, rightH1)
			m.networkView.SetSize(rightW, rightH2)
			m.assetsView.SetSize(rightW, rightH3)
		} else {
			switch m.mode {
			case "logs":
				m.logsView.SetSize(m.width, contentHeight)
			case "alerts":
				m.alertsView.SetSize(m.width, contentHeight)
			case "graphs":
				m.graphsView.SetSize(m.width, contentHeight)
			case "assets":
				m.assetsView.SetSize(m.width, contentHeight)
			case "network":
				m.networkView.SetSize(m.width, contentHeight)
			}
		}
		
	case tickMsg:
		m.engine.Tick()
		
		// Tick views
		m.logsView.Update(msg)
		m.alertsView.Update(msg)
		m.assetsView.Update(msg)
		m.graphsView.Update(msg)
		m.networkView.Update(msg)
		
		cmds = append(cmds, tickCmd())
		
	case logTickMsg:
		m.engine.GenerateLogs(time.Time(msg))
		m.logsView.Update(msg)
		cmds = append(cmds, logTickCmd())
	}

	return m, tea.Batch(cmds...)
}

func (m *MainModel) statusBar() string {
	info := fmt.Sprintf(" Mode: %s | Theme: %s | Time: %s", m.mode, m.engine.Theme(), time.Now().Format("15:04:05"))
	w := m.width - lipgloss.Width(info)
	if w < 0 {
		w = 0
	}
	return theme.StatusBarStyle.Width(m.width).Render(info)
}

func (m *MainModel) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}

	var content string

	if m.splitLayout {
		leftPanel := lipgloss.JoinVertical(lipgloss.Left, m.logsView.View(), m.alertsView.View())
		rightPanel := lipgloss.JoinVertical(lipgloss.Left, m.graphsView.View(), m.networkView.View(), m.assetsView.View())
		content = lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)
	} else {
		switch m.mode {
		case "logs":
			content = m.logsView.View()
		case "alerts":
			content = m.alertsView.View()
		case "graphs":
			content = m.graphsView.View()
		case "assets":
			content = m.assetsView.View()
		case "network":
			content = m.networkView.View()
		default:
			content = "Unknown mode"
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, content, m.statusBar())
}
