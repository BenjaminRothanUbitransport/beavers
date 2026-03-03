package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// getCachePath returns the path to the cache file.
func getCachePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	// The config uses ~/.beavers/config.yaml, so cache goes into the same directory
	beaversDir := filepath.Join(home, ".beavers")
	if err := os.MkdirAll(beaversDir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(beaversDir, "cache.json"), nil
}

// readCache attempts to read the cache file.
// Returns a boolean indicating if the cache is valid and fresh enough, and the cached projects.
func readCache() (bool, []Project, error) {
	cachePath, err := getCachePath()
	if err != nil {
		return false, nil, err
	}

	data, err := os.ReadFile(cachePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil, nil
		}
		return false, nil, err
	}

	var cacheData CacheData
	if err := json.Unmarshal(data, &cacheData); err != nil {
		return false, nil, err
	}

	// Stale-while-revalidate means we return the cache even if it's "stale",
	// but we might want to know if it's totally missing vs stale.
	// Returning true means we found it.
	return true, cacheData.Projects, nil
}

// writeCache saves the discovered projects to the cache file.
func writeCache(projects []Project) error {
	cachePath, err := getCachePath()
	if err != nil {
		return err
	}

	cacheData := CacheData{
		Projects:  projects,
		UpdatedAt: time.Now().Unix(),
	}

	data, err := json.MarshalIndent(cacheData, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(cachePath, data, 0644)
}
