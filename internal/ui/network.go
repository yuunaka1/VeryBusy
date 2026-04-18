package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	
	"github.com/yuunaka1/VeryBusy/internal/sim"
	"github.com/yuunaka1/VeryBusy/internal/theme"
)

type NetworkView struct {
	engine *sim.Engine
	width  int
	height int
	ticks  int
}

func NewNetworkView(engine *sim.Engine) *NetworkView {
	return &NetworkView{engine: engine}
}

func (v *NetworkView) SetSize(w, h int) {
	v.width = w
	v.height = h
}

func (v *NetworkView) Update(msg tea.Msg) {
	// advance scrolling on generic tick messages
	if _, ok := msg.(tickMsg); ok {
		v.ticks++
	}
}

func (v *NetworkView) View() string {
	if v.width <= 0 || v.height <= 0 {
		return ""
	}

	contentWidth := v.width - 4
	contentHeight := v.height - 4

	if contentHeight <= 0 {
		return ""
	}

	var lines []string
	states := []string{"ESTABLISHED", "SYN_SENT", "TIME_WAIT", "CLOSE_WAIT"}
	
	// Create scrolling effect using ticks
	startOffset := v.ticks % 1000
	
	for i := 0; i < contentHeight; i++ {
		idx := startOffset + i
		
		srcIP := fmt.Sprintf("10.%d.%d.%d", (idx*13)%254+1, (idx*19)%254+1, (idx*7)%254+1)
		dstIP := fmt.Sprintf("192.168.%d.%d", (idx*23)%254+1, (idx*29)%254+1)
		state := states[(idx*3)%len(states)]
		sport := 1024 + (idx*11)%60000
		dport := []int{80, 443, 22, 3389, 53, 445}[idx%6]
		
		stateStr := state
		if state == "ESTABLISHED" {
			stateStr = theme.GoodStyle.Render(fmt.Sprintf("%-11s", state))
		} else if state == "SYN_SENT" {
			stateStr = theme.WarningStyle.Render(fmt.Sprintf("%-11s", state))
		} else {
			stateStr = theme.MutedStyle.Render(fmt.Sprintf("%-11s", state))
		}

		line := fmt.Sprintf("%-15s:%-5d -> %-15s:%-4d %s", srcIP, sport, dstIP, dport, stateStr)
		
		if lipgloss.Width(line) > contentWidth && contentWidth > 0 {
			line = line[:contentWidth-3] + "..."
		}
		lines = append(lines, line)
	}

	content := strings.Join(lines, "\n")
	
	return theme.PanelStyle.
		Width(v.width - 2).
		Height(v.height - 2).
		Render(lipgloss.JoinVertical(lipgloss.Left, theme.TitleStyle.Render("NETWORK TRAFFIC"), content))
}
