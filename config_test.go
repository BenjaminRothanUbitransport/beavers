package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	tmpDir := t.TempDir()

	// 1. Test Default/Empty Config
	t.Run("empty config", func(t *testing.T) {
		viper.Reset()
		cfg, err := loadConfig("")
		require.NoError(t, err)
		assert.Empty(t, cfg.Workspaces)
	})

	// 2. Test Local Config (beavers.yaml)
	t.Run("local config beavers.yaml", func(t *testing.T) {
		viper.Reset()
		localFile := filepath.Join(tmpDir, "beavers.yaml")
		err := os.WriteFile(localFile, []byte(`
workspaces:
  - name: local
    root: /tmp/local
`), 0644)
		require.NoError(t, err)

		cfg, err := loadConfig(localFile)
		require.NoError(t, err)
		require.Len(t, cfg.Workspaces, 1)
		assert.Equal(t, "local", cfg.Workspaces[0].Name)
	})

	// 3. Test Global Config (Simulation)
	t.Run("global config simulation", func(t *testing.T) {
		viper.Reset()
		globalFile := filepath.Join(tmpDir, ".beavers", "config.yaml")
		require.NoError(t, os.MkdirAll(filepath.Dir(globalFile), 0755))
		err := os.WriteFile(globalFile, []byte(`
workspaces:
  - name: global
    root: /tmp/global
`), 0644)
		require.NoError(t, err)

		// Set the Home environment for viper to find the config
		oldHome := os.Getenv("HOME")
		os.Setenv("HOME", tmpDir)
		defer os.Setenv("HOME", oldHome)

		cfg, err := loadConfig("")
		require.NoError(t, err)
		require.Len(t, cfg.Workspaces, 1)
		assert.Equal(t, "global", cfg.Workspaces[0].Name)
	})

	// 4. Test Precedence: Flag over Global
	t.Run("precedence flag over global", func(t *testing.T) {
		viper.Reset()
		// Global setup
		globalDir := filepath.Join(tmpDir, "precedence", ".beavers")
		require.NoError(t, os.MkdirAll(globalDir, 0755))
		globalFile := filepath.Join(globalDir, "config.yaml")
		err := os.WriteFile(globalFile, []byte(`
workspaces:
  - name: global
`), 0644)
		require.NoError(t, err)

		// Flag setup
		flagFile := filepath.Join(tmpDir, "precedence", "custom.yaml")
		err = os.WriteFile(flagFile, []byte(`
workspaces:
  - name: flag
`), 0644)
		require.NoError(t, err)

		oldHome := os.Getenv("HOME")
		os.Setenv("HOME", filepath.Join(tmpDir, "precedence"))
		defer os.Setenv("HOME", oldHome)

		cfg, err := loadConfig(flagFile)
		require.NoError(t, err)
		require.Len(t, cfg.Workspaces, 1)
		assert.Equal(t, "flag", cfg.Workspaces[0].Name)
	})
}
