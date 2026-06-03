package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestTextInputDefaultValue(t *testing.T) {
	t.Parallel()

	m := newTextInput("Database name", "appdb")
	if m.value != "appdb" {
		t.Fatalf("value = %q, want %q", m.value, "appdb")
	}
	if m.confirmed {
		t.Fatal("confirmed = true before any key, want false")
	}
}

func TestTextInputTyping(t *testing.T) {
	t.Parallel()

	m := newTextInput("Database name", "")
	for _, r := range "prod" {
		model, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		m = model.(textInputModel)
	}
	if m.value != "prod" {
		t.Fatalf("value = %q, want %q", m.value, "prod")
	}
}

func TestTextInputBackspace(t *testing.T) {
	t.Parallel()

	m := newTextInput("Database name", "abc")

	model, _ := m.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	m = model.(textInputModel)
	if m.value != "ab" {
		t.Fatalf("value = %q, want %q", m.value, "ab")
	}
}

func TestTextInputBackspaceEmpty(t *testing.T) {
	t.Parallel()

	m := newTextInput("Database name", "")

	model, _ := m.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	m = model.(textInputModel)
	if m.value != "" {
		t.Fatalf("value = %q, want empty", m.value)
	}
}

func TestTextInputConfirm(t *testing.T) {
	t.Parallel()

	m := newTextInput("Database name", "appdb")

	model, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = model.(textInputModel)
	if !m.confirmed {
		t.Fatal("confirmed = false after enter, want true")
	}
	if m.value != "appdb" {
		t.Fatalf("value = %q, want %q", m.value, "appdb")
	}
}

func TestTextInputCancel(t *testing.T) {
	t.Parallel()

	m := newTextInput("Database name", "appdb")

	model, _ := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = model.(textInputModel)
	if m.confirmed {
		t.Fatal("confirmed = true after esc, want false")
	}
}
