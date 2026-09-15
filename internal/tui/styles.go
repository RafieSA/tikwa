package tui

import "github.com/charmbracelet/lipgloss"

// Palet TIKWA — dipol-polin, premium
var (
	// Warna brand
	ColorTikTokMagenta = lipgloss.Color("#FF0050")
	ColorTikTokCyan    = lipgloss.Color("#00F2EA")
	ColorWAGreen       = lipgloss.Color("#25D366")
	ColorGold          = lipgloss.Color("#FFD700")
	ColorDim           = lipgloss.Color("#6C6C6C")
	ColorError         = lipgloss.Color("#FF3B30")
	ColorSuccess       = lipgloss.Color("#34C759")

	StyleLogo = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorGold).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorGold).
			Padding(0, 2).
			Align(lipgloss.Center)

	StyleSubtitle = lipgloss.NewStyle().
			Foreground(ColorDim).
			Italic(true).
			Align(lipgloss.Center)

	StyleMenuBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorTikTokCyan).
			Padding(1, 2).
			MarginTop(1)

	StyleMenuItem = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF"))

	StyleMenuKey = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorTikTokMagenta).
			Padding(0, 1).
			MarginRight(1)

	StyleSelected = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorWAGreen).
			Background(lipgloss.Color("#1A1A1A"))

	StyleLoadingBar = lipgloss.NewStyle().
			Foreground(ColorTikTokCyan)

	StyleErrorBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorError).
			Foreground(ColorError).
			Padding(0, 1)

	StyleSuccessBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorSuccess).
			Foreground(ColorSuccess).
			Padding(0, 1)

	StyleEmptyBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorDim).
			Foreground(ColorDim).
			Padding(1, 2).
			Align(lipgloss.Center)

	StyleHelp = lipgloss.NewStyle().
			Foreground(ColorDim).
			MarginTop(1)
)
