package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileExistsChecker(t *testing.T) {
	tempDir := t.TempDir()
	
	// Create a dummy file
	dummyFile := filepath.Join(tempDir, "README.md")
	err := os.WriteFile(dummyFile, []byte("# Dummy"), 0644)
	assert.NoError(t, err)

	checker := &FileExistsChecker{}

	t.Run("File Exists", func(t *testing.T) {
		rule := AuditRule{Type: "file_exists", Params: map[string]string{"filename": "README.md"}}
		result := checker.Audit(tempDir, rule)
		assert.Equal(t, "PASS", result.Status)
	})

	t.Run("File Missing", func(t *testing.T) {
		rule := AuditRule{Type: "file_exists", Params: map[string]string{"filename": "Makefile"}}
		result := checker.Audit(tempDir, rule)
		assert.Equal(t, "FAIL", result.Status)
		assert.Contains(t, result.Message, "not found")
	})

	t.Run("Missing Parameter", func(t *testing.T) {
		rule := AuditRule{Type: "file_exists", Params: map[string]string{}}
		result := checker.Audit(tempDir, rule)
		assert.Equal(t, "FAIL", result.Status)
		assert.Contains(t, result.Message, "Missing or empty")
	})
}

func TestMakefileTargetChecker(t *testing.T) {
	tempDir := t.TempDir()
	
	makefileContent := `
build:
	go build .

test:
	go test ./...

lint :
	golangci-lint run
`
	err := os.WriteFile(filepath.Join(tempDir, "Makefile"), []byte(makefileContent), 0644)
	assert.NoError(t, err)

	checker := &MakefileTargetChecker{}

	t.Run("Target Exists", func(t *testing.T) {
		rule := AuditRule{Type: "makefile_target", Params: map[string]string{"target": "test"}}
		result := checker.Audit(tempDir, rule)
		assert.Equal(t, "PASS", result.Status)
	})

	t.Run("Target with space before colon Exists", func(t *testing.T) {
		rule := AuditRule{Type: "makefile_target", Params: map[string]string{"target": "lint"}}
		result := checker.Audit(tempDir, rule)
		assert.Equal(t, "PASS", result.Status)
	})

	t.Run("Target Missing", func(t *testing.T) {
		rule := AuditRule{Type: "makefile_target", Params: map[string]string{"target": "deploy"}}
		result := checker.Audit(tempDir, rule)
		assert.Equal(t, "FAIL", result.Status)
		assert.Contains(t, result.Message, "not found in Makefile")
	})

	t.Run("Target in Included File", func(t *testing.T) {
		incDir := t.TempDir()
		incFile := filepath.Join(incDir, "inc.mk")
		err := os.WriteFile(incFile, []byte("included_target:\n\t@echo ok\n"), 0644)
		assert.NoError(t, err)

		makeWithInc := "include " + incFile + "\n"
		err = os.WriteFile(filepath.Join(tempDir, "Makefile"), []byte(makeWithInc), 0644)
		assert.NoError(t, err)

		rule := AuditRule{Type: "makefile_target", Params: map[string]string{"target": "included_target"}}
		result := checker.Audit(tempDir, rule)
		assert.Equal(t, "PASS", result.Status)
	})

	t.Run("Makefile Missing", func(t *testing.T) {
		emptyDir := t.TempDir()
		rule := AuditRule{Type: "makefile_target", Params: map[string]string{"target": "build"}}
		result := checker.Audit(emptyDir, rule)
		assert.Equal(t, "FAIL", result.Status)
		assert.Contains(t, result.Message, "Makefile not found")
	})

	t.Run("Missing Parameter", func(t *testing.T) {
		rule := AuditRule{Type: "makefile_target", Params: map[string]string{}}
		result := checker.Audit(tempDir, rule)
		assert.Equal(t, "FAIL", result.Status)
		assert.Contains(t, result.Message, "Missing or empty")
	})
}

func TestRunAudit(t *testing.T) {
	tempDir := t.TempDir()
	err := os.WriteFile(filepath.Join(tempDir, "README.md"), []byte("# Title"), 0644)
	assert.NoError(t, err)

	project := Project{Path: tempDir}
	rules := map[string]AuditRule{
		"check_readme":  {Type: "file_exists", Params: map[string]string{"filename": "README.md"}},
		"check_unknown": {Type: "unknown_type"},
	}

	results := RunAudit(project, rules)
	assert.Len(t, results, 2)

	var foundReadme, foundUnknown bool
	for _, res := range results {
		if res.RuleName == "check_readme" {
			assert.Equal(t, "PASS", res.Status)
			foundReadme = true
		} else if res.RuleName == "check_unknown" {
			assert.Equal(t, "FAIL", res.Status)
			assert.Contains(t, res.Message, "Unknown rule type")
			foundUnknown = true
		}
	}
	assert.True(t, foundReadme)
	assert.True(t, foundUnknown)
}
