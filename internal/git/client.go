package git

import (
	"bytes"
	"os/exec"
	"strings"
)

type RealGitClient struct{}

func NewClient() *RealGitClient {
	return &RealGitClient{}
}

func (c *RealGitClient) DetectStatus(path string) (string, string) {
	// 1. Get branch name
	cmdBranch := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmdBranch.Dir = path
	var outBranch bytes.Buffer
	cmdBranch.Stdout = &outBranch
	if err := cmdBranch.Run(); err != nil {
		return "", "" // Probably not a git repo or no commits
	}
	branch := strings.TrimSpace(outBranch.String())

	// 2. Get sync status using git status -sb
	cmdStatus := exec.Command("git", "status", "-sb")
	cmdStatus.Dir = path
	var outStatus bytes.Buffer
	cmdStatus.Stdout = &outStatus
	if err := cmdStatus.Run(); err != nil {
		return branch, ""
	}

	statusLines := strings.Split(outStatus.String(), "\n")
	if len(statusLines) == 0 {
		return branch, ""
	}

	firstLine := statusLines[0]
	// firstLine usually looks like "## main...origin/main [ahead 1, behind 2]" or "## main"
	syncStatus := ""
	if strings.Contains(firstLine, "[") && strings.Contains(firstLine, "]") {
		start := strings.Index(firstLine, "[")
		end := strings.Index(firstLine, "]")
		if start < end {
			syncStatus = firstLine[start+1 : end]
		}
	} else if strings.Contains(firstLine, "...") {
		syncStatus = "up-to-date"
	}

	// Check for uncommitted changes
	hasChanges := false
	for _, line := range statusLines[1:] {
		if strings.TrimSpace(line) != "" {
			hasChanges = true
			break
		}
	}

	if hasChanges {
		if syncStatus != "" {
			syncStatus += ", uncommitted"
		} else {
			syncStatus = "uncommitted"
		}
	}

	return branch, syncStatus
}
