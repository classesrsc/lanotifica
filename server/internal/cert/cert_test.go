package cert

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadOrCreate_CreatesNew(t *testing.T) {
	t.Parallel()

	// Create temp directory
	tempDir := t.TempDir()

	cert, err := LoadOrCreate(tempDir)
	if err != nil {
		t.Fatalf("LoadOrCreate failed: %v", err)
	}

	// Check certificate was created
	certPath := filepath.Join(tempDir, "cert.pem")
	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		t.Error("Expected cert.pem to be created")
	}

	// Check key was created
	keyPath := filepath.Join(tempDir, "key.pem")
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		t.Error("Expected key.pem to be created")
	}

	// Check fingerprint format (64 hex chars uppercase)
	if len(cert.Fingerprint) != 64 {
		t.Errorf("Expected fingerprint length 64, got %d", len(cert.Fingerprint))
	}

	if cert.Fingerprint != strings.ToUpper(cert.Fingerprint) {
		t.Error("Expected fingerprint to be uppercase")
	}

	// Check TLS certificate is valid
	if len(cert.TLSCert.Certificate) == 0 {
		t.Error("Expected TLS certificate to have at least one certificate")
	}
}

func TestLoadOrCreate_LoadsExisting(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()

	// Create first certificate
	cert1, err := LoadOrCreate(tempDir)
	if err != nil {
		t.Fatalf("First LoadOrCreate failed: %v", err)
	}

	// Load existing certificate
	cert2, err := LoadOrCreate(tempDir)
	if err != nil {
		t.Fatalf("Second LoadOrCreate failed: %v", err)
	}

	// Fingerprints should be identical
	if cert1.Fingerprint != cert2.Fingerprint {
		t.Errorf("Expected same fingerprint, got %s vs %s", cert1.Fingerprint, cert2.Fingerprint)
	}
}

func TestCalculateFingerprint_Valid(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()

	// Create a certificate first
	_, err := LoadOrCreate(tempDir)
	if err != nil {
		t.Fatalf("LoadOrCreate failed: %v", err)
	}

	certPath := filepath.Join(tempDir, "cert.pem")
	fingerprint, err := calculateFingerprint(certPath)
	if err != nil {
		t.Fatalf("calculateFingerprint failed: %v", err)
	}

	// Check fingerprint format
	if len(fingerprint) != 64 {
		t.Errorf("Expected fingerprint length 64, got %d", len(fingerprint))
	}

	// Check uppercase hex
	for _, c := range fingerprint {
		isDigit := c >= '0' && c <= '9'
		isUpperHex := c >= 'A' && c <= 'F'
		if !isDigit && !isUpperHex {
			t.Errorf("Expected uppercase hex, got character: %c", c)
		}
	}
}

func TestCalculateFingerprint_InvalidPEM(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	certPath := filepath.Join(tempDir, "invalid.pem")

	// Write invalid PEM content
	err := os.WriteFile(certPath, []byte("not a valid PEM"), 0o644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	_, err = calculateFingerprint(certPath)
	if err == nil {
		t.Error("Expected error for invalid PEM")
	}
}

func TestCalculateFingerprint_FileNotFound(t *testing.T) {
	t.Parallel()

	_, err := calculateFingerprint("/nonexistent/path/cert.pem")
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}
}

func TestGetLocalIPs_Contains127(t *testing.T) {
	t.Parallel()

	ips := getLocalIPs()

	found := false
	for _, ip := range ips {
		if ip.String() == "127.0.0.1" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected getLocalIPs to contain 127.0.0.1")
	}
}

func TestGetLocalIPs_NotEmpty(t *testing.T) {
	t.Parallel()

	ips := getLocalIPs()

	if len(ips) == 0 {
		t.Error("Expected at least one IP address")
	}
}

func TestConfigDir_NotEmpty(t *testing.T) {
	t.Parallel()

	dir := ConfigDir()

	if dir == "" {
		t.Error("Expected ConfigDir to return non-empty string")
	}

	if !strings.Contains(dir, "lanotifica") {
		t.Errorf("Expected ConfigDir to contain 'lanotifica', got: %s", dir)
	}
}

func TestConfigDir_WithXDGConfigHome(t *testing.T) {
	customConfig := "/custom/config"
	t.Setenv("XDG_CONFIG_HOME", customConfig)

	dir := ConfigDir()

	expected := filepath.Join(customConfig, "lanotifica")
	if dir != expected {
		t.Errorf("Expected %s, got %s", expected, dir)
	}
}

func TestFileExists(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()

	// Test non-existent file
	if fileExists(filepath.Join(tempDir, "nonexistent.txt")) {
		t.Error("Expected fileExists to return false for non-existent file")
	}

	// Create a file
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0o644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test existing file
	if !fileExists(testFile) {
		t.Error("Expected fileExists to return true for existing file")
	}
}
