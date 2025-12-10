// Package icon provides functionality to fetch and cache Android app icons from the Play Store.
package icon

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

var (
	// ErrPlayStoreNotFound is returned when the app is not found on Play Store.
	ErrPlayStoreNotFound = errors.New("app not found on Play Store")
	// ErrNoIconFound is returned when no icon URL is found in the Play Store page.
	ErrNoIconFound = errors.New("no icon found in Play Store page")
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

// Cache manages the icon cache directory.
type Cache struct {
	dir string
}

// NewCache creates a new icon cache using XDG_CACHE_HOME/la-notify/icons.
func NewCache() *Cache {
	xdgCache := os.Getenv("XDG_CACHE_HOME")
	if xdgCache == "" {
		home, _ := os.UserHomeDir()
		xdgCache = filepath.Join(home, ".cache")
	}
	return &Cache{
		dir: filepath.Join(xdgCache, "la-notify", "icons"),
	}
}

// GetIconPath returns the cached icon path for the given package name.
// If the icon is not in cache, it downloads it from Play Store first.
// Returns empty string if the icon cannot be obtained.
func (c *Cache) GetIconPath(packageName string) string {
	if packageName == "" {
		return ""
	}

	iconPath := c.buildIconPath(packageName)

	if c.existsInCache(iconPath) {
		return iconPath
	}

	if err := c.downloadAndCache(packageName, iconPath); err != nil {
		return ""
	}

	return iconPath
}

func (c *Cache) buildIconPath(packageName string) string {
	return filepath.Join(c.dir, packageName+".png")
}

func (c *Cache) existsInCache(iconPath string) bool {
	_, err := os.Stat(iconPath)
	return err == nil
}

func (c *Cache) downloadAndCache(packageName, iconPath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	iconURL, err := c.fetchIconURLFromPlayStore(ctx, packageName)
	if err != nil {
		return err
	}

	return c.saveIconToCache(ctx, iconURL, iconPath)
}

func (c *Cache) fetchIconURLFromPlayStore(ctx context.Context, packageName string) (string, error) {
	playStoreURL := "https://play.google.com/store/apps/details?id=" + packageName

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, playStoreURL, nil)
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("fetching Play Store page: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return "", ErrPlayStoreNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%w: status %d", ErrPlayStoreNotFound, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading response body: %w", err)
	}

	re := regexp.MustCompile(`https://play-lh\.googleusercontent\.com/[^"'\s]+`)
	matches := re.FindAllString(string(body), -1)
	if len(matches) == 0 {
		return "", ErrNoIconFound
	}

	return matches[0], nil
}

func (c *Cache) saveIconToCache(ctx context.Context, iconURL, iconPath string) error {
	iconReq, err := http.NewRequestWithContext(ctx, http.MethodGet, iconURL, nil)
	if err != nil {
		return fmt.Errorf("creating icon request: %w", err)
	}

	iconResp, err := httpClient.Do(iconReq)
	if err != nil {
		return fmt.Errorf("downloading icon: %w", err)
	}
	defer func() { _ = iconResp.Body.Close() }()

	if err := os.MkdirAll(c.dir, 0750); err != nil {
		return fmt.Errorf("creating cache dir: %w", err)
	}

	file, err := os.Create(iconPath) //nolint:gosec // iconPath is controlled internally
	if err != nil {
		return fmt.Errorf("creating icon file: %w", err)
	}
	defer func() { _ = file.Close() }()

	if _, err := io.Copy(file, iconResp.Body); err != nil {
		return fmt.Errorf("saving icon to file: %w", err)
	}

	return nil
}
