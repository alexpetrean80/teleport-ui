package cache

import (
	"os"
	"path/filepath"
	"testing"
)

type testItem struct {
	Name  string
	Value int
}

func setTempCacheDir(t *testing.T) {
	t.Helper()
	t.Setenv("XDG_CACHE_HOME", t.TempDir())
}

func TestSaveAndLoad(t *testing.T) {
	setTempCacheDir(t)

	items := []testItem{
		{Name: "alpha", Value: 1},
		{Name: "beta", Value: 2},
	}

	if err := Save("test", items); err != nil {
		t.Fatalf("Save: %v", err)
	}

	var loaded []testItem
	if err := Load("test", &loaded); err != nil {
		t.Fatalf("Load: %v", err)
	}

	if len(loaded) != len(items) {
		t.Fatalf("got %d items, want %d", len(loaded), len(items))
	}
	for i := range items {
		if loaded[i] != items[i] {
			t.Errorf("item %d: got %+v, want %+v", i, loaded[i], items[i])
		}
	}
}

func TestSaveAndLoadEmpty(t *testing.T) {
	setTempCacheDir(t)

	var items []testItem
	if err := Save("empty", items); err != nil {
		t.Fatalf("Save: %v", err)
	}

	var loaded []testItem
	if err := Load("empty", &loaded); err != nil {
		t.Fatalf("Load: %v", err)
	}

	if len(loaded) != 0 {
		t.Fatalf("got %d items, want 0", len(loaded))
	}
}

func TestLoadMissing(t *testing.T) {
	setTempCacheDir(t)

	var loaded []testItem
	if err := Load("nonexistent", &loaded); err == nil {
		t.Fatal("expected error loading nonexistent cache, got nil")
	}
}

func TestClear(t *testing.T) {
	setTempCacheDir(t)

	items := []testItem{{Name: "x", Value: 42}}
	if err := Save("clearme", items); err != nil {
		t.Fatalf("Save: %v", err)
	}

	if err := Clear("clearme"); err != nil {
		t.Fatalf("Clear: %v", err)
	}

	var loaded []testItem
	if err := Load("clearme", &loaded); err == nil {
		t.Fatal("expected error after Clear, got nil")
	}
}

func TestClearNonexistent(t *testing.T) {
	setTempCacheDir(t)

	if err := Clear("neverexisted"); err != nil {
		t.Fatalf("Clear nonexistent: %v", err)
	}
}

func TestClearAll(t *testing.T) {
	setTempCacheDir(t)

	if err := Save("a", []testItem{{Name: "a", Value: 1}}); err != nil {
		t.Fatalf("Save a: %v", err)
	}
	if err := Save("b", []testItem{{Name: "b", Value: 2}}); err != nil {
		t.Fatalf("Save b: %v", err)
	}

	if err := ClearAll(); err != nil {
		t.Fatalf("ClearAll: %v", err)
	}

	d, err := dir()
	if err != nil {
		t.Fatalf("dir: %v", err)
	}
	if _, err := os.Stat(d); !os.IsNotExist(err) {
		t.Fatalf("cache directory still exists after ClearAll")
	}
}

func TestCacheFileLocation(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("XDG_CACHE_HOME", tmp)

	items := []testItem{{Name: "loc", Value: 99}}
	if err := Save("location", items); err != nil {
		t.Fatalf("Save: %v", err)
	}

	expected := filepath.Join(tmp, dirName, "location.gob")
	if _, err := os.Stat(expected); err != nil {
		t.Fatalf("cache file not at expected path %s: %v", expected, err)
	}
}
