# Agent Instructions: Project Beavers 🦫

This file provides critical context and architectural constraints for AI agents working on the Beavers CLI.

## 🏗 Project Philosophy
- **TDD First:** Every feature MUST start with a test in a corresponding `*_test.go` file. Use `testify` for assertions.
- **Flat Structure:** Maintain a flat root structure for simplicity (no `cmd/` or `internal/` folders unless the project grows significantly).
- **Silent STDOUT:** The `path` command is intended for shell integration (`cd $(beavers path ...)`). NEVER print logs, debug info, or interactive UI elements to STDOUT during this command. Use STDERR for errors or logging if necessary.

## ⚙️ Configuration & Discovery
- **Config Precedence:** CLI Flag > Local `beavers.yaml` > Global `~/.beavers/config.yaml`.
- **Binary Conflict Hazard:** The binary is named `beavers`. When loading configuration, explicitly check for `.yaml` or `.yml` extensions. DO NOT let Viper scan the directory for a generic "beavers" file, as it will attempt to parse the binary and fail with UTF-8 errors.
- **Project Types:**
    - `standalone`: The directory is the project.
    - `monorepo-sub`: The directory contains sub-projects matched by `patterns`.
- **Workspace Context:** Use workspace-level `type` and `excludes` to disambiguate nested repositories.

## 🛠 Tech Stack Constraints
- **Language:** Go 1.24+
- **CLI Framework:** Cobra (use `PersistentPreRunE` on the root command for config loading to ensure flags are parsed first).
- **Config Engine:** Viper.
- **UI/Table Formatting:** PTerm.
- **Testing:** Standard `testing` package + `github.com/stretchr/testify`.

## 🔄 Common Workflows
- **Validation:** Always run `go test -v ./...` before considering a task complete.
- **Manual Check:** Test with the local `beavers.yaml` to ensure discovery logic works against actual filesystem structures.
- **Dependency Management:** Run `go mod tidy` after adding any new imports.
