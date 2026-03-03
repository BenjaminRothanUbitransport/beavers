package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

// Checker is the interface for all audit rule implementations.
type Checker interface {
	Audit(projectPath string, rule AuditRule) AuditResult
}

// Registry maps rule types to their Checker implementations.
var checkerRegistry = map[string]Checker{
	"file_exists":     &FileExistsChecker{},
	"makefile_target": &MakefileTargetChecker{},
	// Roadmap for composer.json parsing:
	// "composer_json_version": &ComposerJsonChecker{},
	// This would involve a JSON parser to extract the `require` and `require-dev` blocks,
	// and a semver parsing library (e.g., github.com/Masterminds/semver/v3) to evaluate
	// versions against constraints provided in the rule parameters.
}

// RunAudit orchestrates the execution of multiple audit rules against a project.
func RunAudit(project Project, rules map[string]AuditRule) []AuditResult {
	var results []AuditResult

	for ruleName, rule := range rules {
		checker, ok := checkerRegistry[rule.Type]
		if !ok {
			results = append(results, AuditResult{
				RuleName: ruleName,
				Status:   "FAIL",
				Message:  fmt.Sprintf("Unknown rule type: %s", rule.Type),
			})
			continue
		}

		result := checker.Audit(project.Path, rule)
		result.RuleName = ruleName
		results = append(results, result)
	}

	return results
}

// FileExistsChecker verifies if a specific file exists in the project path.
type FileExistsChecker struct{}

func (c *FileExistsChecker) Audit(projectPath string, rule AuditRule) AuditResult {
	filename, ok := rule.Params["filename"]
	if !ok || filename == "" {
		return AuditResult{Status: "FAIL", Message: "Missing or empty 'filename' parameter"}
	}

	fullPath := filepath.Join(projectPath, filename)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return AuditResult{Status: "FAIL", Message: fmt.Sprintf("File '%s' not found", filename)}
	} else if err != nil {
		return AuditResult{Status: "FAIL", Message: fmt.Sprintf("Error accessing file '%s': %v", filename, err)}
	}

	return AuditResult{Status: "PASS", Message: fmt.Sprintf("File '%s' exists", filename)}
}

// MakefileTargetChecker parses a Makefile to verify if a specific target exists.
type MakefileTargetChecker struct{}

func (c *MakefileTargetChecker) Audit(projectPath string, rule AuditRule) AuditResult {
	target, ok := rule.Params["target"]
	if !ok || target == "" {
		return AuditResult{Status: "FAIL", Message: "Missing or empty 'target' parameter"}
	}

	makefilePath := filepath.Join(projectPath, "Makefile")
	if _, err := os.Stat(makefilePath); os.IsNotExist(err) {
		return AuditResult{Status: "FAIL", Message: "Makefile not found"}
	}

	cmd := exec.Command("make", "-pRrq", "-C", projectPath)
	out, _ := cmd.Output() // make -q typically returns exit code 1 or 2, we just want the stdout

	// Regex to match Makefile targets in make database output (e.g., `target:`)
	pattern := fmt.Sprintf(`(?m)^%s\s*:`, regexp.QuoteMeta(target))
	re, err := regexp.Compile(pattern)
	if err != nil {
		return AuditResult{Status: "FAIL", Message: fmt.Sprintf("Invalid target regex: %v", err)}
	}

	if re.Match(out) {
		return AuditResult{Status: "PASS", Message: fmt.Sprintf("Target '%s' exists in Makefile", target)}
	}

	return AuditResult{Status: "FAIL", Message: fmt.Sprintf("Target '%s' not found in Makefile", target)}
}
