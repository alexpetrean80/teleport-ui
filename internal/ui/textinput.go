package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// RunTextInput shows an input box with an optional default value.
// Returns the entered string and whether the user confirmed (enter vs esc).
func RunTextInput(prompt, defaultValue string) (string, bool, error) {
	m := newTextInput(prompt, defaultValue)
	program := tea.NewProgram(m)
	finalModel, err := program.Run()
	if err != nil {
		return "", false, err
	}
	result := finalModel.(textInputModel)
	return result.value, result.confirmed, nil
}

type textInputModel struct {
	prompt      string
	value       string
	confirmed   bool
	promptStyle lipgloss.Style
	inputStyle  lipgloss.Style
	hintStyle   lipgloss.Style
}

func newTextInput(prompt, defaultValue string) textInputModel {
	return textInputModel{
		prompt: prompt,
		value:  defaultValue,
		promptStyle: lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#1A7F1A", Dark: "#00FF00"}).
			Bold(true),
		inputStyle: lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#111111", Dark: "#FFFFFF"}),
		hintStyle: lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#444444", Dark: "#666666"}),
	}
}

func (m textInputModel) Init() tea.Cmd {
	return nil
}

func (m textInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			m.confirmed = true
			return m, tea.Quit
		case "backspace":
			if len(m.value) > 0 {
				m.value = m.value[:len(m.value)-1]
			}
		default:
			if len(msg.String()) == 1 {
				m.value += msg.String()
			}
		}
	}
	return m, nil
}

func (m textInputModel) View() string {
	var b strings.Builder
	b.WriteString(m.promptStyle.Render(m.prompt+": "))
	b.WriteString(m.inputStyle.Render(m.value))
	b.WriteString("█\n\n")
	b.WriteString(m.hintStyle.Render("enter: confirm • esc: cancel"))
	return b.String()
}
