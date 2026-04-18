package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	
	"verybusy/internal/sim"
	"verybusy/internal/theme"
)

type GraphsView struct {
	engine *sim.Engine
	width  int
	height int
}

func NewGraphsView(engine *sim.Engine) *GraphsView {
	return &GraphsView{engine: engine}
}

func (v *GraphsView) SetSize(w, h int) {
	v.width = w
	v.height = h
}

func (v *GraphsView) Update(msg tea.Msg) {
}

func (v *GraphsView) View() string {
	if v.width <= 0 || v.height <= 0 {
		return ""
	}

	// simple graph implementation
	contentHeight := v.height - 4
	var lines []string

	metrics := v.engine.Metrics["Network Traffic (Anomaly Score)"]
	
	graphH := contentHeight - 2
	if graphH < 1 {
		return ""
	}
	
	// Max value for scaling
	maxVal := 100.0
	for _, m := range metrics {
		if m.Value > maxVal {
			maxVal = m.Value
		}
	}

	graphW := v.width - 4
	if graphW < 1 {
		return ""
	}
	
	start := 0
	if len(metrics) > graphW {
		start = len(metrics) - graphW
	}
	
	displayMetrics := metrics[start:]
	
	// Draw from top to bottom
	for row := graphH; row > 0; row-- {
		threshold := (float64(row) / float64(graphH)) * maxVal
		var line strings.Builder
		for _, m := range displayMetrics {
			if m.Value >= threshold {
				line.WriteString(theme.WarningStyle.Render("█"))
			} else {
				line.WriteString(" ")
			}
		}
		lines = append(lines, line.String())
	}
	
	// Base line
	lines = append(lines, theme.MutedStyle.Render(strings.Repeat("-", graphW)))
	lines = append(lines, fmt.Sprintf("Latest Anomaly Score: %.1f", displayMetrics[len(displayMetrics)-1].Value))
	
	for len(lines) < contentHeight {
		lines = append(lines, "")
	}

	content := strings.Join(lines, "\n")
	return theme.PanelStyle.
		Width(v.width - 2).
		Height(v.height - 2).
		Render(lipgloss.JoinVertical(lipgloss.Left, theme.TitleStyle.Render("TELEMETRY (NETWORK ANOMALY)"), content))
}
