package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SelectorModel represents the state of the commit message selector
type SelectorModel struct {
	choices  []string
	cursor   int
	selected int
	quitting bool
}

// SelectorMsg represents messages for the selector
type SelectorMsg struct {
	Type    string
	Content string
}

// Styles for the UI
var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	itemStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			Foreground(lipgloss.Color("#FAFAFA"))

	selectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(lipgloss.Color("#7D56F4")).
				Background(lipgloss.Color("#FAFAFA")).
				Bold(true)

	cursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Italic(true)
)

// NewSelectorModel creates a new selector model
func NewSelectorModel(choices []string) *SelectorModel {
	return &SelectorModel{
		choices:  choices,
		cursor:   0,
		selected: -1,
		quitting: false,
	}
}

// Init initializes the model
func (m SelectorModel) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m SelectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case "enter", " ":
			m.selected = m.cursor
			return m, tea.Quit

		case "e":
			// Edit mode - return the selected choice for editing
			m.selected = m.cursor
			return m, tea.Quit
		}
	}

	return m, nil
}

// View renders the UI
func (m SelectorModel) View() string {
	if m.quitting {
		return "No commit message selected.\n"
	}

	var s strings.Builder

	// Title
	s.WriteString(titleStyle.Render("🤖 Choose a commit message:"))
	s.WriteString("\n\n")

	// Choices
	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = cursorStyle.Render(">")
		}

		style := itemStyle
		if m.cursor == i {
			style = selectedItemStyle
		}

		s.WriteString(fmt.Sprintf("%s %s\n", cursor, style.Render(choice)))
	}

	// Help text
	s.WriteString("\n")
	s.WriteString(helpStyle.Render("↑/↓ or k/j: navigate • enter: select • e: edit • q: quit"))

	return s.String()
}

// GetSelected returns the selected choice
func (m SelectorModel) GetSelected() string {
	if m.selected >= 0 && m.selected < len(m.choices) {
		return m.choices[m.selected]
	}
	return ""
}

// WasEditRequested checks if the user requested to edit
func (m SelectorModel) WasEditRequested() bool {
	return m.selected >= 0 && m.quitting
}
