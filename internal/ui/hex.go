package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/yuunaka1/VeryBusy/internal/sim"
	"github.com/yuunaka1/VeryBusy/internal/theme"
)

type HexView struct {
	engine *sim.Engine
	width  int
	height int
	ticks  int
}

func NewHexView(engine *sim.Engine) *HexView {
	return &HexView{engine: engine}
}

func (v *HexView) SetSize(w, h int) {
	v.width = w
	v.height = h
}

func (v *HexView) Update(msg tea.Msg) {
	if _, ok := msg.(tickMsg); ok {
		// Advance faster to look like it's analyzing quickly
		v.ticks += 3
	}
	if _, ok := msg.(logTickMsg); ok {
		// Also advance on log ticks to make it smoother and more random
		v.ticks += 1
	}
}

func (v *HexView) View() string {
	if v.width <= 0 || v.height <= 0 {
		return ""
	}

	contentWidth := v.width - 4
	contentHeight := v.height - 4

	if contentHeight <= 0 {
		return ""
	}

	var lines []string
	
	// Create scrolling effect
	startOffset := v.ticks % 0xFFFF
	
	for i := 0; i < contentHeight; i++ {
		offset := (startOffset + i) * 16 // 16 bytes per line
		
		var hexPart strings.Builder
		var asciiPart strings.Builder
		
		for j := 0; j < 16; j++ {
			b := byte((offset*7 + j*13) % 256)
			
			// Inject fake PE headers looping
			offsetInFile := offset % 0x0FFF
			if offsetInFile == 0 && j == 0 { b = 0x4D } // M
			if offsetInFile == 0 && j == 1 { b = 0x5A } // Z
			if offsetInFile == 0xE0 && j == 0 { b = 0x50 } // P
			if offsetInFile == 0xE0 && j == 1 { b = 0x45 } // E
			
			hexByte := fmt.Sprintf("%02X", b)
			if b == 0x00 {
				hexByte = theme.MutedStyle.Render(hexByte)
			} else if b == 0x90 {
				// NOP sled illusion
				hexByte = theme.WarningStyle.Render(hexByte)
			} else if b >= 0x40 && b <= 0x5A {
				hexByte = theme.GoodStyle.Render(hexByte)
			}

			hexPart.WriteString(hexByte)
			hexPart.WriteString(" ")
			if j == 7 {
				hexPart.WriteString(" ")
			}
			
			if b >= 32 && b <= 126 {
				asciiPart.WriteByte(b)
			} else {
				asciiPart.WriteByte('.')
			}
		}
		
		line := fmt.Sprintf("%08X  %s |%s|", offset, hexPart.String(), theme.InfoStyle.Render(asciiPart.String()))
		
		if lipgloss.Width(line) > contentWidth && contentWidth > 0 {
			line = line[:contentWidth-3] + "..."
		}
		lines = append(lines, line)
	}

	content := strings.Join(lines, "\n")
	
	headerText := fmt.Sprintf("MALWARE BINARY ANALYSIS (suspicious_file_0x%04X.exe)", v.ticks % 0xFFFF)
	return theme.PanelStyle.
		Width(v.width - 2).
		Height(v.height - 2).
		Render(lipgloss.JoinVertical(lipgloss.Left, theme.TitleStyle.Render(headerText), content))
}
