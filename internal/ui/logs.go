package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	
	"github.com/yuunaka1/VeryBusy/internal/sim"
	"github.com/yuunaka1/VeryBusy/internal/theme"
)

type LogsView struct {
	engine *sim.Engine
	width  int
	height int
}

func NewLogsView(engine *sim.Engine) *LogsView {
	return &LogsView{engine: engine}
}

func (v *LogsView) SetSize(w, h int) {
	v.width = w
	v.height = h
}

func (v *LogsView) Update(msg tea.Msg) {
	// just stubbing out for ticks
}

func (v *LogsView) View() string {
	if v.width <= 0 || v.height <= 0 {
		return ""
	}

	contentWidth := v.width - 4 // border padding
	contentHeight := v.height - 4 // border + title

	var lines []string
	
	// Print right aligned log limit
	logs := v.engine.Logs
	start := 0
	if len(logs) > contentHeight {
		start = len(logs) - contentHeight
	}

	for i := start; i < len(logs); i++ {
		log := logs[i]
		
		sevStyle := theme.InfoStyle
		switch log.Severity {
		case sim.Low:
			sevStyle = theme.GoodStyle
		case sim.Medium:
			sevStyle = theme.WarningStyle
		case sim.High, sim.Critical:
			sevStyle = theme.CriticalStyle
		}
		
		sevStr := sevStyle.Render(fmt.Sprintf("[%-4d]", log.Severity)) // fake severity number or level
		ts := theme.MutedStyle.Render(log.Timestamp.Format("15:04:05"))
		
		line := fmt.Sprintf("%s %s %-12s | %s", ts, sevStr, log.Hostname, log.Message)
		// Truncate to contentWidth
		if lipgloss.Width(line) > contentWidth && contentWidth > 0 {
			line = line[:contentWidth-3] + "..."
		}
		lines = append(lines, line)
	}

	// Pad rest with blanks if needed
	for len(lines) < contentHeight {
		lines = append(lines, "")
	}

	content := strings.Join(lines, "\n")
	
	return theme.PanelStyle.
		Width(v.width - 2).
		Height(v.height - 2).
		Render(lipgloss.JoinVertical(lipgloss.Left, theme.TitleStyle.Render("SYSTEM LOGS"), content))
}
