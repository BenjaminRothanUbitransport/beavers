# Technical Analysis: Project Beavers 🦫

## 1. Architectural Overview
Beavers is designed as a **decoupled discovery and execution layer**. Its primary architecture follows the standard CLI pattern for Go (Cobra/Viper) but adds a discovery engine and a compliance auditor.

### Core Modules:
*   **Config Manager (Viper):** Handles multi-workspace definitions and global audit rules.
*   **Discovery Engine:** Implements a directory walker with glob-pattern support to resolve "logical projects" from "physical paths".
*   **Cache Layer:** Persistent JSON storage to ensure sub-100ms response times for project listing.
*   **Audit Engine:** A pluggable rule-based system for verifying project health and standards.
*   **TUI (Bubble Tea):** An interactive layer for real-time monitoring and navigation.

---

## 2. Component Breakdown

### 2.1. Discovery & Workspace Management
*   **Data Structure:** A `Workspace` defines a `root` and a set of `discovery_patterns`.
*   **Algorithm:** 
    *   Standalone: Simple directory walk.
    *   Monorepo: Recurse into folders matching glob patterns (e.g., `services/*`).
*   **Alias Registry:** A map of short-codes to fully qualified paths, stored in the global configuration or a local cache.

### 2.2. The Audit Engine
The engine must be **declarative and evolutive**. 
*   **Rule Types:**
    *   `FileExistenceRule`: Simple check for `Makefile`, `.gitignore`, etc.
    *   `MakefileTargetRule`: Parses `Makefile` to ensure targets like `install` or `test` exist.
    *   `VersionCheckRule`: (Future) Parses `composer.json` or `go.mod`.
*   **Reporting:** Should support machine-readable formats (JSON) for CI/CD integration and human-readable formats (Table/TUI) for local dev.

### 2.3. Shell Integration (`bgo`)
To enable directory jumping, the CLI itself cannot change the parent shell's directory. 
*   **Mechanism:** `beavers path <alias>` prints the path. 
*   **Helper Function:** A shell script/alias (e.g., `bgo() { cd $(beavers path "$1") }`) must be injected into user profiles.

---

## 3. Data Models (Draft)

```go
type Workspace struct {
    Name     string   `yaml:"name"`
    Root     string   `yaml:"root"`
    Patterns []string `yaml:"patterns"`
}

type Project struct {
    ID       string // Generated from path/name
    Name     string
    Alias    string
    Path     string
    Type     string // monorepo-sub, standalone
    Status   GitStatus
    Compliance HealthStatus
}

type GitStatus struct {
    Branch string
    IsSynced bool
    IsDirty bool
}
```

---

## 4. Key Technical Challenges & Solutions

### 4.1. Performance vs. Accuracy
*   **Challenge:** Walking large monorepos on every command is slow.
*   **Solution:** Use a "Stale-While-Revalidate" cache.
    *   Commands like `beavers project list` read from `cache.json`.
    *   A background process or a `--refresh` flag updates the cache.
    *   File-system watchers (optional) could trigger updates.

### 4.2. Environment Detection
*   **Challenge:** `beavers doctor` needs to verify heterogeneous environments (PHP, Go, Docker).
*   **Solution:** Use `os/exec` to probe binary versions and verify connectivity (e.g., `docker info`).

---

## 5. Execution Strategy

1.  **Phase 1 (The Pathfinder):** Core discovery and shell integration. ✅
2.  **Phase 2 (The Executor):** `make` target wrappers, Git status detection, and performance caching. ✅
3.  **Phase 3 (The Auditor):** Declarative project health checks and compliance reporting.
4.  **Phase 4 (The TUI Dashboard):** Interactive monitoring with Bubble Tea.

---

## 6. Implementation Notes
*   **Concurrency:** Use `sync.WaitGroup` or worker pools for high-performance directory discovery.
*   **Extensibility:** Audit rules should be defined as interfaces to allow easy addition of new checkers.
*   **Caching:** Use standard library `encoding/json` for cache persistence.
