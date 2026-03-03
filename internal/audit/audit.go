package audit

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/ubitransports/beavers/internal/app"
	"github.com/ubitransports/beavers/internal/config"
)

type Checker interface {
	Audit(appCtx *app.App, projectPath string, rule config.AuditRule) config.AuditResult
}

func getRegistry() map[string]Checker {
	return map[string]Checker{
		"file_exists":     &FileExistsChecker{},
		"makefile_target": &MakefileTargetChecker{},
	}
}

func RunAudit(appCtx *app.App, project config.Project, rules map[string]config.AuditRule) []config.AuditResult {
	var results []config.AuditResult
	registry := getRegistry()

	for ruleName, rule := range rules {
		checker, ok := registry[rule.Type]
		if !ok {
			results = append(results, config.AuditResult{
				RuleName: ruleName,
				Status:   "FAIL",
				Message:  fmt.Sprintf("Unknown rule type: %s", rule.Type),
			})
			continue
		}

		result := checker.Audit(appCtx, project.Path, rule)
		result.RuleName = ruleName
		results = append(results, result)
	}

	return results
}

type FileExistsChecker struct{}

func (c *FileExistsChecker) Audit(appCtx *app.App, projectPath string, rule config.AuditRule) config.AuditResult {
	filename, ok := rule.Params["filename"]
	if !ok || filename == "" {
		return config.AuditResult{Status: "FAIL", Message: "Missing or empty 'filename' parameter"}
	}

	fullPath := filepath.Join(projectPath, filename)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return config.AuditResult{Status: "FAIL", Message: fmt.Sprintf("File '%s' not found", filename)}
	} else if err != nil {
		return config.AuditResult{Status: "FAIL", Message: fmt.Sprintf("Error accessing file '%s': %v", filename, err)}
	}

	return config.AuditResult{Status: "PASS", Message: fmt.Sprintf("File '%s' exists", filename)}
}

type MakefileTargetChecker struct{}

func (c *MakefileTargetChecker) Audit(appCtx *app.App, projectPath string, rule config.AuditRule) config.AuditResult {
	target, ok := rule.Params["target"]
	if !ok || target == "" {
		return config.AuditResult{Status: "FAIL", Message: "Missing or empty 'target' parameter"}
	}

	makefilePath := filepath.Join(projectPath, "Makefile")
	if _, err := os.Stat(makefilePath); os.IsNotExist(err) {
		return config.AuditResult{Status: "FAIL", Message: "Makefile not found"}
	}

	out, _ := appCtx.Exec.Run(projectPath, "make", "-pRrq")

	pattern := fmt.Sprintf(`(?m)^%s\s*:`, regexp.QuoteMeta(target))
	re, err := regexp.Compile(pattern)
	if err != nil {
		return config.AuditResult{Status: "FAIL", Message: fmt.Sprintf("Invalid target regex: %v", err)}
	}

	if re.Match(out) {
		return config.AuditResult{Status: "PASS", Message: fmt.Sprintf("Target '%s' exists in Makefile", target)}
	}

	return config.AuditResult{Status: "FAIL", Message: fmt.Sprintf("Target '%s' not found in Makefile", target)}
}
