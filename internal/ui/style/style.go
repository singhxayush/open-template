package style

import "github.com/charmbracelet/lipgloss"

// ----- Styling -----
// Using styling references from your provided code.
var (
	DocStyle = lipgloss.NewStyle().
			MarginLeft(2).
			MarginTop(1).
			Align(lipgloss.Left)

	ErrorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))

	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Padding(1).
			PaddingRight(15).
			Border(lipgloss.NormalBorder(), true).
			BorderForeground()

	ListItemStyle = lipgloss.NewStyle().PaddingLeft(2)

	PromptStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#BD93F9"))

	LeftPanelStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			MarginRight(1).
			PaddingLeft(1).
			PaddingRight(1).
			Foreground(lipgloss.Color("#cbcbcb")).
			Height(16).Width(35)

	RightPanelStyle = lipgloss.NewStyle().
		// Border(lipgloss.RoundedBorder()).
		// BorderRight(false).
		PaddingRight(2).
		MarginTop(1).
		Foreground(lipgloss.Color("#cbcbcb")).
		Height(12)

	InstructionStyle = lipgloss.
				NewStyle().
				Foreground(lipgloss.Color("#666666")).
				MarginTop(1)

	// A simple cursor style.
	CursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true)
)
