package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectListCommand(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Setup a project
	projDir := filepath.Join(tmpDir, "my-project")
	require.NoError(t, os.MkdirAll(projDir, 0755))

	cfg := &Config{
		Workspaces: []Workspace{
			{Name: "test", Root: projDir},
		},
		Aliases: map[string]string{
			"my-alias": "my-project",
		},
	}

	projects, err := discoverProjects(cfg)
	require.NoError(t, err)

	t.Run("list output contains alias and name", func(t *testing.T) {
		buf := new(bytes.Buffer)
		err := listProjects(buf, projects)
		require.NoError(t, err)

		output := buf.String()
		assert.Contains(t, output, "my-project")
		assert.Contains(t, output, "my-alias")
	})
}

func TestPathCommand(t *testing.T) {
	projects := []Project{
		{Name: "proj1", Alias: "p1", Path: "/abs/path/to/proj1"},
	}

	t.Run("path output is clean for alias", func(t *testing.T) {
		buf := new(bytes.Buffer)
		err := printProjectPath(buf, projects, "p1")
		require.NoError(t, err)

		output := strings.TrimSpace(buf.String())
		assert.Equal(t, "/abs/path/to/proj1", output)
	})

	t.Run("path output is clean for name", func(t *testing.T) {
		buf := new(bytes.Buffer)
		err := printProjectPath(buf, projects, "proj1")
		require.NoError(t, err)

		output := strings.TrimSpace(buf.String())
		assert.Equal(t, "/abs/path/to/proj1", output)
	})

	t.Run("path error for unknown identifier", func(t *testing.T) {
		buf := new(bytes.Buffer)
		err := printProjectPath(buf, projects, "unknown")
		assert.Error(t, err)
	})
}
