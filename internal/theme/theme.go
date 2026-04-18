package theme

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	ColorBg        = lipgloss.Color("#0d1117")
	ColorFg        = lipgloss.Color("#e6edf3")
	ColorBorder    = lipgloss.Color("#30363d")
	ColorInfo      = lipgloss.Color("#58a6ff")
	ColorGood      = lipgloss.Color("#3fb950")
	ColorWarning   = lipgloss.Color("#d29922")
	ColorError     = lipgloss.Color("#f85149")
	ColorCritical  = lipgloss.Color("#ff7b72")
	ColorHighlight = lipgloss.Color("#1f6feb")
	ColorMuted     = lipgloss.Color("#8b949e")

	// Styles
	BaseStyle = lipgloss.NewStyle().Foreground(ColorFg)

	PanelStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorder).
		Padding(0, 1)

	TitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorHighlight).
		MarginBottom(1)

	StatusBarStyle = lipgloss.NewStyle().
		Background(ColorBorder).
		Foreground(ColorFg).
		Padding(0, 1)

	InfoStyle     = lipgloss.NewStyle().Foreground(ColorInfo)
	GoodStyle     = lipgloss.NewStyle().Foreground(ColorGood)
	WarningStyle  = lipgloss.NewStyle().Foreground(ColorWarning)
	ErrorStyle    = lipgloss.NewStyle().Foreground(ColorError)
	CriticalStyle = lipgloss.NewStyle().Foreground(ColorCritical).Bold(true)
	MutedStyle    = lipgloss.NewStyle().Foreground(ColorMuted)
)
