package tui

import "github.com/charmbracelet/lipgloss"

// RenderPrefix applies the prompt prefix style used inside the TUI.
// Call this when printing a prompt line outside of a bubbletea session.
func RenderPrefix(s string) string {
	return prefixStyle.Render(s)
}

const completionNameWidth = 24

var (
	prefixStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("33")) // yellow

	completionNormalStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("252"))

	completionSelectedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("33")).
				Bold(true)

	completionDescStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("244"))

	completionSelectedDescStyle = lipgloss.NewStyle().
					Foreground(lipgloss.Color("240"))

	completionBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("240")).
				PaddingLeft(1).
				PaddingRight(1)
)
