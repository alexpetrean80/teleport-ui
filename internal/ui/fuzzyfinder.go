package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func RunFuzzyFinder[T Stringable](items []T) (*T, error) {
	finder := NewFuzzyFinder(items)

	program := tea.NewProgram(finder)
	finalModel, err := program.Run()
	if err != nil {
		return nil, err
	}

	finderModel := finalModel.(FuzzyFinder[T])
	return finderModel.Selected(), nil
}

type Stringable interface {
	String() string
}

type FuzzyFinder[T Stringable] struct {
	items         []T
	filteredItems []T
	cursor        int
	query         string
	selected      *T
	width         int
	height        int
	selectedStyle lipgloss.Style
	normalStyle   lipgloss.Style
	promptStyle   lipgloss.Style
	matchStyle    lipgloss.Style
}

func NewFuzzyFinder[T Stringable](items []T) FuzzyFinder[T] {
	return FuzzyFinder[T]{
		items:         items,
		filteredItems: items,
		cursor:        0,
		query:         "",
		selectedStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Background(lipgloss.Color("#3C3C3C")).
			Bold(true),
		normalStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")),
		promptStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true),
		matchStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF00FF")).
			Bold(true),
	}
}

func (f FuzzyFinder[T]) Init() tea.Cmd {
	return nil
}

func (f FuzzyFinder[T]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		f.width = msg.Width
		f.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return f, tea.Quit

		case "enter":
			if len(f.filteredItems) > 0 && f.cursor < len(f.filteredItems) {
				f.selected = &f.filteredItems[f.cursor]
			}
			return f, tea.Quit

		case "up", "ctrl+k":
			if f.cursor > 0 {
				f.cursor--
			}

		case "down", "ctrl+j":
			if f.cursor < len(f.filteredItems)-1 {
				f.cursor++
			}

		case "backspace":
			if len(f.query) > 0 {
				f.query = f.query[:len(f.query)-1]
				f.filterItems()
				if f.cursor >= len(f.filteredItems) && len(f.filteredItems) > 0 {
					f.cursor = len(f.filteredItems) - 1
				}
				if len(f.filteredItems) == 0 {
					f.cursor = 0
				}
			}

		default:
			if len(msg.String()) == 1 {
				f.query += msg.String()
				f.filterItems()
				f.cursor = 0
			}
		}
	}

	return f, nil
}

func (f FuzzyFinder[T]) View() string {
	var b strings.Builder

	prompt := f.promptStyle.Render("> ") + f.query
	b.WriteString(prompt)
	b.WriteString("\n\n")

	maxItems := f.height - 5
	if maxItems < 1 {
		maxItems = 10
	}

	start := 0
	end := len(f.filteredItems)

	if len(f.filteredItems) > maxItems {
		if f.cursor >= maxItems/2 {
			start = f.cursor - maxItems/2
			if start+maxItems > len(f.filteredItems) {
				start = len(f.filteredItems) - maxItems
			}
		}
		end = start + maxItems
		if end > len(f.filteredItems) {
			end = len(f.filteredItems)
		}
	}

	for i := start; i < end; i++ {
		item := f.filteredItems[i]
		itemStr := item.String()

		if i == f.cursor {
			b.WriteString(f.selectedStyle.Render("▶ " + f.highlightMatches(itemStr)))
		} else {
			b.WriteString(f.normalStyle.Render("  " + f.highlightMatches(itemStr)))
		}
		b.WriteString("\n")
	}

	if len(f.filteredItems) == 0 {
		b.WriteString(
			lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).Render("  No matches"),
		)
		b.WriteString("\n")
	}

	b.WriteString("\n")
	info := lipgloss.NewStyle().Foreground(lipgloss.Color("#666666")).Render(
		"↑/↓: navigate • enter: select • esc: cancel",
	)
	b.WriteString(info)

	return b.String()
}

func (f *FuzzyFinder[T]) filterItems() {
	if f.query == "" {
		f.filteredItems = f.items
		return
	}

	filtered := make([]T, 0)
	queryLower := strings.ToLower(f.query)

	for _, item := range f.items {
		if fuzzyMatch(strings.ToLower(item.String()), queryLower) {
			filtered = append(filtered, item)
		}
	}

	f.filteredItems = filtered
}

func (f FuzzyFinder[T]) highlightMatches(text string) string {
	if f.query == "" {
		return text
	}

	queryLower := strings.ToLower(f.query)
	textLower := strings.ToLower(text)

	var result strings.Builder
	queryIdx := 0

	for i, char := range text {
		if queryIdx < len(queryLower) && textLower[i] == queryLower[queryIdx] {
			result.WriteString(f.matchStyle.Render(string(char)))
			queryIdx++
		} else {
			result.WriteRune(char)
		}
	}

	return result.String()
}

func fuzzyMatch(text, query string) bool {
	if query == "" {
		return true
	}

	queryIdx := 0
	for _, char := range text {
		if queryIdx < len(query) && char == rune(query[queryIdx]) {
			queryIdx++
		}
		if queryIdx == len(query) {
			return true
		}
	}

	return queryIdx == len(query)
}

func (f FuzzyFinder[T]) Selected() *T {
	return f.selected
}
