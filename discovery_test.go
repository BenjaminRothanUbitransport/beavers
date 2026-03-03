package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiscoverProjects(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup a standalone project
	standaloneDir := filepath.Join(tmpDir, "standalone")
	require.NoError(t, os.MkdirAll(standaloneDir, 0755))

	// Setup a monorepo
	monorepoDir := filepath.Join(tmpDir, "monorepo")
	require.NoError(t, os.MkdirAll(filepath.Join(monorepoDir, "services/svc1"), 0755))
	require.NoError(t, os.MkdirAll(filepath.Join(monorepoDir, "services/svc2"), 0755))
	require.NoError(t, os.MkdirAll(filepath.Join(monorepoDir, "apps/app1"), 0755))

	cfg := &Config{
		Workspaces: []Workspace{
			{
				Name: "standalone",
				Root: standaloneDir,
			},
			{
				Name:     "monorepo",
				Root:     monorepoDir,
				Patterns: []string{"services/*", "apps/*"},
			},
		},
		Aliases: map[string]string{
			"my-svc1": "svc1",
		},
	}

	t.Run("discover projects", func(t *testing.T) {
		projects, err := discoverProjects(cfg)
		require.NoError(t, err)

		// Expecting:
		// 1. Standalone project
		// 2. monorepo/services/svc1
		// 3. monorepo/services/svc2
		// 4. monorepo/apps/app1
		assert.Len(t, projects, 4)

		// Check standalone
		var standalone *Project
		for _, p := range projects {
			if p.Workspace == "standalone" {
				standalone = &p
				break
			}
		}
		require.NotNil(t, standalone)
		assert.Equal(t, "standalone", standalone.Name)
		assert.Equal(t, "standalone", standalone.Type)

		// Check monorepo sub-project with alias
		var svc1 *Project
		for _, p := range projects {
			if p.Name == "svc1" {
				svc1 = &p
				break
			}
		}
		require.NotNil(t, svc1)
		assert.Equal(t, "my-svc1", svc1.Alias)
		assert.Equal(t, "monorepo-sub", svc1.Type)
	})
}
