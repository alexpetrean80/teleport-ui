package ui

import (
	"slices"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

type testStr string

func (s testStr) String() string { return string(s) }

func TestFuzzyMatch(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		text  string
		query string
		want  bool
	}{
		{name: "exact match", text: "hello", query: "hello", want: true},
		{name: "subsequence", text: "abcde", query: "ace", want: true},
		{name: "no match", text: "abc", query: "xyz", want: false},
		{name: "empty query", text: "anything", query: "", want: true},
		{name: "empty text non-empty query", text: "", query: "x", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := fuzzyMatch(tt.text, tt.query); got != tt.want {
				t.Fatalf("fuzzyMatch(%q, %q) = %v, want %v", tt.text, tt.query, got, tt.want)
			}
		})
	}
}

func TestFilterByTyping(t *testing.T) {
	t.Parallel()

	finder := NewFuzzyFinder([]testStr{"apple", "banana", "apricot"})

	model, _ := finder.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	finder = model.(FuzzyFinder[testStr])

	model, _ = finder.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
	finder = model.(FuzzyFinder[testStr])

	want := []testStr{"apple", "apricot"}
	if !slices.Equal(finder.filteredItems, want) {
		t.Fatalf("filteredItems = %v, want %v", finder.filteredItems, want)
	}
}

func TestCursorNavigation(t *testing.T) {
	t.Parallel()

	finder := NewFuzzyFinder([]testStr{"a", "b", "c"})
	if finder.cursor != 0 {
		t.Fatalf("initial cursor = %d, want 0", finder.cursor)
	}

	model, _ := finder.Update(tea.KeyMsg{Type: tea.KeyDown})
	finder = model.(FuzzyFinder[testStr])
	if finder.cursor != 1 {
		t.Fatalf("after down: cursor = %d, want 1", finder.cursor)
	}

	model, _ = finder.Update(tea.KeyMsg{Type: tea.KeyDown})
	finder = model.(FuzzyFinder[testStr])
	if finder.cursor != 2 {
		t.Fatalf("after second down: cursor = %d, want 2", finder.cursor)
	}

	model, _ = finder.Update(tea.KeyMsg{Type: tea.KeyDown})
	finder = model.(FuzzyFinder[testStr])
	if finder.cursor != 2 {
		t.Fatalf("after third down (clamp): cursor = %d, want 2", finder.cursor)
	}

	model, _ = finder.Update(tea.KeyMsg{Type: tea.KeyUp})
	finder = model.(FuzzyFinder[testStr])
	if finder.cursor != 1 {
		t.Fatalf("after up: cursor = %d, want 1", finder.cursor)
	}
}

func TestSelectItem(t *testing.T) {
	t.Parallel()

	finder := NewFuzzyFinder([]testStr{"alpha", "beta"})

	model, _ := finder.Update(tea.KeyMsg{Type: tea.KeyDown})
	finder = model.(FuzzyFinder[testStr])

	model, _ = finder.Update(tea.KeyMsg{Type: tea.KeyEnter})
	finder = model.(FuzzyFinder[testStr])

	sel := finder.Selected()
	if sel == nil {
		t.Fatal("Selected() = nil, want pointer to beta")
	}
	if *sel != "beta" {
		t.Fatalf("Selected() = %q, want beta", *sel)
	}
}

func TestCancelSelection(t *testing.T) {
	t.Parallel()

	finder := NewFuzzyFinder([]testStr{"alpha", "beta"})

	model, _ := finder.Update(tea.KeyMsg{Type: tea.KeyEsc})
	finder = model.(FuzzyFinder[testStr])

	if sel := finder.Selected(); sel != nil {
		t.Fatalf("Selected() = %v, want nil after esc", sel)
	}
}

func TestBackspace(t *testing.T) {
	t.Parallel()

	finder := NewFuzzyFinder([]testStr{"apple", "banana"})

	model, _ := finder.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})
	finder = model.(FuzzyFinder[testStr])

	if len(finder.filteredItems) != 1 {
		t.Fatalf("after 'b': len(filteredItems) = %d, want 1", len(finder.filteredItems))
	}

	model, _ = finder.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	finder = model.(FuzzyFinder[testStr])

	if len(finder.filteredItems) != 2 {
		t.Fatalf("after backspace: len(filteredItems) = %d, want 2", len(finder.filteredItems))
	}
}
