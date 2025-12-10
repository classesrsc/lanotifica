package config

import (
	"os"
	"path/filepath"
	"testing"
)

var (
	originalConfigDir  string
	originalConfigPath string
)

func TestMain(m *testing.M) {
	// Save original values.
	originalConfigDir = configDir
	originalConfigPath = configPath

	// Run tests.
	code := m.Run()

	// Restore original values.
	configDir = originalConfigDir
	configPath = originalConfigPath

	os.Exit(code)
}

func setupTestConfig(t *testing.T) {
	t.Helper()

	tempDir := t.TempDir()
	configDir = filepath.Join(tempDir, "la-notify")
	configPath = filepath.Join(configDir, "config.json")

	t.Cleanup(func() {
		configDir = originalConfigDir
		configPath = originalConfigPath
	})
}
