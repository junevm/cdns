package ui

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

// IsTTY checks if the output is a terminal
func IsTTY() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

// GetTerminalSize returns the current terminal dimensions
func GetTerminalSize() (width int, height int, err error) {
	if !IsTTY() {
		return 80, 24, nil // Default for non-TTY
	}
	return term.GetSize(int(os.Stdout.Fd()))
}

// Styles contains reusable UI styles
type Styles struct {
	Success lipgloss.Style
	Error   lipgloss.Style
	Info    lipgloss.Style
	Warning lipgloss.Style
	Header  lipgloss.Style
	Box     lipgloss.Style
	Brand   lipgloss.Style
}

// NewStyles creates a new Styles instance with default styling
func NewStyles() *Styles {
	return &Styles{
		Success: lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")).
			Bold(true),
		Error: lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true),
		Info: lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Bold(true),
		Warning: lipgloss.NewStyle().
			Foreground(lipgloss.Color("214")).
			Bold(true),
		Header: lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true),
		Box: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(1, 2),
		Brand: lipgloss.NewStyle().
			Foreground(lipgloss.Color("39")).
			Bold(true),
	}
}

// RenderInfo renders an info message
func (s *Styles) RenderInfo(msg string) string {
	if !IsTTY() {
		return msg
	}
	return s.Info.Render(msg)
}

// RenderSuccess renders a success message with icon
func (s *Styles) RenderSuccess(msg string) string {
	if !IsTTY() {
		return "✓ " + msg
	}
	return s.Success.Render("✓ " + msg)
}

// RenderError renders an error message with icon
func (s *Styles) RenderError(msg string) string {
	if !IsTTY() {
		return "✗ " + msg
	}
	return s.Error.Render("✗ " + msg)
}

// RenderWarning renders a warning message with icon
func (s *Styles) RenderWarning(msg string) string {
	if !IsTTY() {
		return "⚠ " + msg
	}
	return s.Warning.Render("⚠ " + msg)
}

// RenderBox renders content in a styled box
func (s *Styles) RenderBox(content string) string {
	if !IsTTY() {
		return content
	}
	return s.Box.Render(content)
}

// RenderBold renders text in bold with brand color
func (s *Styles) RenderBold(text string) string {
	if !IsTTY() {
		return text
	}
	return s.Brand.Render(text)
}

// RenderDim renders text dimmed/faded
func (s *Styles) RenderDim(text string) string {
	if !IsTTY() {
		return text
	}
	return lipgloss.NewStyle().Faint(true).Render(text)
}
