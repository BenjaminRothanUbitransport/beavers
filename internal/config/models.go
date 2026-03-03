package config

// Workspace defines a root directory and discovery patterns for projects.
type Workspace struct {
	Name     string   `mapstructure:"name" yaml:"name"`
	Root     string   `mapstructure:"root" yaml:"root"`
	Patterns []string `mapstructure:"patterns" yaml:"patterns"`
	Type     string   `mapstructure:"type" yaml:"type"`         // Default type for projects in this workspace
	Excludes []string `mapstructure:"excludes" yaml:"excludes"` // Folders to ignore
}

// Config represents the global beavers configuration.
type Config struct {
	Workspaces []Workspace          `mapstructure:"workspaces" yaml:"workspaces"`
	Aliases    map[string]string    `mapstructure:"aliases" yaml:"aliases"`
	AuditRules map[string]AuditRule `mapstructure:"audit_rules" yaml:"audit_rules"`
}

// AuditRule defines a single compliance check.
type AuditRule struct {
	Type   string            `mapstructure:"type" yaml:"type"`
	Params map[string]string `mapstructure:"params" yaml:"params"`
}

// AuditResult represents the outcome of an audit rule check.
type AuditResult struct {
	RuleName string
	Status   string // "PASS" or "FAIL"
	Message  string
}

// Project represents a discovered logical project.
type Project struct {
	ID         string
	Name       string
	Alias      string
	Path       string
	Type       string // "standalone" or "monorepo-sub"
	Workspace  string
	GitBranch  string
	SyncStatus string
}
