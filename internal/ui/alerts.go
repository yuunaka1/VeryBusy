package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	
	"verybusy/internal/sim"
	"verybusy/internal/theme"
)

type AlertsView struct {
	engine *sim.Engine
	width  int
	height int
}

func NewAlertsView(engine *sim.Engine) *AlertsView {
	return &AlertsView{engine: engine}
}

func (v *AlertsView) SetSize(w, h int) {
	v.width = w
	v.height = h
}

func (v *AlertsView) Update(msg tea.Msg) {
}

func (v *AlertsView) View() string {
	if v.width <= 0 || v.height <= 0 {
		return ""
	}

	contentHeight := v.height - 4
	var lines []string

	alerts := v.engine.Alerts
	start := 0
	if len(alerts) > contentHeight {
		start = len(alerts) - contentHeight
	}

	for i := start; i < len(alerts); i++ {
		a := alerts[i]
		
		sevStyle := theme.CriticalStyle
		if a.Severity == sim.Medium {
			sevStyle = theme.WarningStyle
		}

		lines = append(lines, fmt.Sprintf("%s %s - %s", 
			theme.MutedStyle.Render(a.Timestamp.Format("15:04:05")),
			sevStyle.Render(a.ID),
			a.RuleName,
		))
	}

	for len(lines) < contentHeight {
		lines = append(lines, "")
	}

	content := strings.Join(lines, "\n")
	return theme.PanelStyle.
		Width(v.width - 2).
		Height(v.height - 2).
		Render(lipgloss.JoinVertical(lipgloss.Left, theme.TitleStyle.Render("ACTIVE ALERTS"), content))
}
