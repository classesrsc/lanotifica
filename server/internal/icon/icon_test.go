package icon

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewCache(t *testing.T) {
	t.Parallel()

	cache := NewCache()
	if cache == nil {
		t.Fatal("NewCache returned nil")
	}
	if cache.dir == "" {
		t.Error("Cache dir should not be empty")
	}
}

func TestCache_GetIconPath_EmptyPackage(t *testing.T) {
	t.Parallel()

	cache := NewCache()
	result := cache.GetIconPath("")
	if result != "" {
		t.Errorf("Expected empty string for empty package, got %s", result)
	}
}

func TestCache_GetIconPath_CachedIcon(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	cache := &Cache{dir: tempDir}

	// Create a fake cached icon.
	testPath := filepath.Join(tempDir, "com.test.app.png")
	if err := os.WriteFile(testPath, []byte("fake icon"), 0o600); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	result := cache.GetIconPath("com.test.app")
	if result != testPath {
		t.Errorf("Expected %s, got %s", testPath, result)
	}
}

func TestCache_GetIconPath_InvalidPackage(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	cache := &Cache{dir: tempDir}

	result := cache.GetIconPath("com.invalid.nonexistent.package.xyz123")
	if result != "" {
		t.Errorf("Expected empty string for invalid package, got %s", result)
	}
}

func TestCache_existsInCache(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	cache := &Cache{dir: tempDir}

	// Non-existent file.
	if cache.existsInCache(filepath.Join(tempDir, "nonexistent.png")) {
		t.Error("existsInCache should return false for non-existent file")
	}

	// Existing file.
	existingPath := filepath.Join(tempDir, "existing.png")
	if err := os.WriteFile(existingPath, []byte("test"), 0o600); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	if !cache.existsInCache(existingPath) {
		t.Error("existsInCache should return true for existing file")
	}
}

func TestCache_buildIconPath(t *testing.T) {
	t.Parallel()

	cache := &Cache{dir: "/tmp/test-icons"}
	result := cache.buildIconPath("com.example.app")
	expected := "/tmp/test-icons/com.example.app.png"
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}
