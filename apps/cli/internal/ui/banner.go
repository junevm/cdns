package ui

import "github.com/charmbracelet/lipgloss"

// GetBanner returns the styled ASCII art banner
func GetBanner() string {
	banner := `
  ___ ___  _  _ ___ 
 / __|   \| \| / __|
| (__| |) | .  \__ \
 \___|___/|_|\_|___/
`
	// Use a cyan color to match the theme (Color 39 is used elsewhere)
	style := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39")).
		MarginBottom(1)

	return style.Render(banner)
}
