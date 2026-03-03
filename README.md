# Beavers 🦫

**Beavers** is an industrious CLI tool designed to bridge the gap between physical infrastructure (Monorepos and Standalone Repos) and logical developer workflows. It provides a "single source of truth" for project discovery, local execution, and compliance auditing.

## 🚀 Key Features (Phase 1: Pathfinder)

- **Unified Discovery:** Map logical project names (aliases) to physical paths regardless of repo structure.
- **Multi-Workspace Support:** Define distinct workspaces (e.g., `Work`, `Personal`) with different root directories.
- **Pattern-Based Discovery:** Supports both standalone repositories and monorepo sub-projects using glob patterns.
- **Granular Control:** Fine-tune discovery with workspace-level `type` (standalone/monorepo-sub) and `excludes` lists.
- **Flexible Configuration:** Uses Viper to handle configuration precedence (Flag > Local `beavers.yaml` > Global `~/.beavers/config.yaml`).
- **Shell Integration:** Optimized `path` command for directory jumping.

## 🛠 Installation

### Prerequisites
- [Go 1.24+](https://go.dev/dl/)

### Build from source
```bash
go build -o beavers .
```

## 📖 Configuration

Beavers looks for a configuration file in the following order:
1.  **Flag:** `--config path/to/config.yaml`
2.  **Local:** `./beavers.yaml`
3.  **Global:** `~/.beavers/config.yaml`

### Example `beavers.yaml`
```yaml
workspaces:
  # 1. Main projects folder (every folder is a standalone project)
  - name: Projects
    root: /Users/name/projects
    patterns:
      - "*"
    type: standalone
    excludes:
      - "monorepo-folder" # Skip the monorepo folder in this workspace

  # 2. A specific monorepo structure
  - name: Monorepo
    root: /Users/name/projects/monorepo/src
    patterns:
      - "services/*"
      - "apps/*"
    type: monorepo-sub

aliases:
  # alias: project_folder_name_or_absolute_path
  api: backend-service
```

## 💻 Usage

### List all discovered projects
```bash
beavers project list
```

### Get a project path
```bash
beavers path <alias_or_name>
```

### Shell Integration (`bgo`)
Add the following to your `.zshrc` or `.bashrc` to jump to projects instantly:
```bash
bgo() {
  local target=$(beavers path "$1")
  if [ -n "$target" ]; then
    cd "$target"
  else
    echo "Project '$1' not found."
  fi
}
```
Usage: `bgo api`

## 🧪 Development

### Running Tests (TDD Approach)
This project follows a Test-Driven Development methodology. All core logic is covered by unit tests.
```bash
go test -v ./...
```

### Project Structure
- `main.go`: Entry point.
- `commands.go`: Cobra CLI command definitions.
- `config.go`: Viper configuration loading logic.
- `discovery.go`: Workspace walking and project identification engine.
- `models.go`: Core data structures.

## 🗺 Roadmap
- **Phase 2: The Executor:** `make` target wrappers and Git status detection.
- **Phase 3: The Auditor:** Declarative project health checks and compliance reporting.
- **Phase 4: The TUI Dashboard:** Interactive monitoring with Bubble Tea.
