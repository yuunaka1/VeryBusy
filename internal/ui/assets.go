package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	
	"verybusy/internal/sim"
	"verybusy/internal/theme"
)

type AssetsView struct {
	engine *sim.Engine
	width  int
	height int
}

func NewAssetsView(engine *sim.Engine) *AssetsView {
	return &AssetsView{engine: engine}
}

func (v *AssetsView) SetSize(w, h int) {
	v.width = w
	v.height = h
}

func (v *AssetsView) Update(msg tea.Msg) {
}

func (v *AssetsView) View() string {
	if v.width <= 0 || v.height <= 0 {
		return ""
	}

	contentHeight := v.height - 4
	var lines []string

	// Header
	header := fmt.Sprintf("%-15s %-10s %-10s %-10s", "HOSTNAME", "RISK", "ALERTS", "EDR")
	lines = append(lines, theme.MutedStyle.Render(header))
	lines = append(lines, theme.MutedStyle.Render(strings.Repeat("-", lipgloss.Width(header))))

	assets := v.engine.Assets
	displayCount := len(assets)
	if displayCount > contentHeight-2 {
		displayCount = contentHeight - 2
	}

	for i := 0; i < displayCount; i++ {
		a := assets[i]
		
		riskStr := fmt.Sprintf("%d", a.RiskScore)
		riskStyle := theme.InfoStyle
		if a.RiskScore > 10 {
			riskStyle = theme.CriticalStyle
		} else if a.RiskScore > 5 {
			riskStyle = theme.WarningStyle
		} else {
			riskStyle = theme.GoodStyle
		}

		lines = append(lines, fmt.Sprintf("%-15s %s %-10d %-10s", 
			a.Hostname,
			riskStyle.Render(fmt.Sprintf("%-10s", riskStr)),
			a.ActiveAlerts,
			string(a.EDRState),
		))
	}

	for len(lines) < contentHeight {
		lines = append(lines, "")
	}

	content := strings.Join(lines, "\n")
	return theme.PanelStyle.
		Width(v.width - 2).
		Height(v.height - 2).
		Render(lipgloss.JoinVertical(lipgloss.Left, theme.TitleStyle.Render("ASSET STATUS"), content))
}
