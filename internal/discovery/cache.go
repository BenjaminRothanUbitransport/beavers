package discovery

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/ubitransports/beavers/internal/config"
)

type CacheData struct {
	Projects  []config.Project `json:"projects"`
	UpdatedAt int64            `json:"updated_at"`
}

func getCachePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	beaversDir := filepath.Join(home, ".beavers")
	if err := os.MkdirAll(beaversDir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(beaversDir, "cache.json"), nil
}

func ReadCache() (bool, []config.Project, error) {
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

	return true, cacheData.Projects, nil
}

func WriteCache(projects []config.Project) error {
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
