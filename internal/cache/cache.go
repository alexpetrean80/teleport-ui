package cache

import (
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
)

const dirName = "teleport-ui"

func dir() (string, error) {
	cacheHome := os.Getenv("XDG_CACHE_HOME")
	if cacheHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("getting home directory: %w", err)
		}
		cacheHome = filepath.Join(home, ".cache")
	}
	return filepath.Join(cacheHome, dirName), nil
}

func filePath(name string) (string, error) {
	d, err := dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(d, name+".gob"), nil
}

func Load[T any](name string, dest *[]T) error {
	path, err := filePath(name)
	if err != nil {
		return err
	}

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return gob.NewDecoder(f).Decode(dest)
}

func Save[T any](name string, data []T) error {
	d, err := dir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(d, 0o755); err != nil {
		return fmt.Errorf("creating cache directory: %w", err)
	}

	path := filepath.Join(d, name+".gob")
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating cache file: %w", err)
	}
	defer f.Close()

	return gob.NewEncoder(f).Encode(data)
}

func Clear(name string) error {
	path, err := filePath(name)
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("removing cache file: %w", err)
	}
	return nil
}

func ClearAll() error {
	d, err := dir()
	if err != nil {
		return err
	}
	if err := os.RemoveAll(d); err != nil {
		return fmt.Errorf("removing cache directory: %w", err)
	}
	return nil
}
