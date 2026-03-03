# PRD: Project Beavers 🦫

**Title:** Beavers CLI - The Unified Developer Experience Plane  
**Status:** Agile / Draft  
**Target Audience:** Backend Developers, DevOps Engineers  
**Core Tech Stack:** Go (Golang), Cobra, Viper, Bubble Tea  

---

## 1. Vision & Executive Summary
**Beavers** is an industrious CLI tool designed to streamline the management of complex microservice ecosystems. It bridges the gap between physical infrastructure (Monorepos and Standalone Repos) and logical developer workflows. By providing a "single source of truth" for project discovery, local execution, and compliance auditing, Beavers reduces cognitive load and ensures organizational standards are met.

---

## 2. Problem Statement
Backend engineers at the company manage dozens of services across multiple organizations. 
* **Navigational Friction:** Projects are buried deep within monorepo structures (e.g., `src/projects/services/*`), making `cd` commands tedious.
* **Execution Variance:** While most projects use `make`, checking for missing targets or local setup requirements is manual.
* **Configuration Drift:** No automated way to verify if a project follows company standards (e.g., correct PHP versions, presence of a Dockerfile, or specific vendor versions) until CI fails.
* **Repo vs. Project:** A single repository (monorepo) contains multiple independent business projects that need to be treated as first-class entities.

---

## 3. Goals & Objectives
* **Unified Discovery:** Map logical project names (aliases) to physical paths regardless of repo structure.
* **Contextual Command Wrapper:** Execute `make` targets with prefixing and path awareness.
* **Health Transparency:** Instant visibility into Git branch sync status and standards compliance.
* **Extensible Governance:** An evolutive audit engine to enforce file-system and manifest-level rules.

---

## 4. Functional Requirements

### 4.1. Workspace & Multi-Org Management
* **Global Configuration:** Centralized YAML configuration (`~/.beavers/config.yaml`) to avoid repo pollution.
* **Multi-Workspace Support:** Define distinct "Workspaces" (e.g., `Ubitransport`, `Personal`) with different root directories.
* **Granular Discovery:** 
    * **Explicit Typing:** Workspace-level `type` setting (`standalone` or `monorepo-sub`) to force project categorization.
    * **Folder Exclusion:** `excludes` list per workspace to ignore specific directories during discovery.
* **Pattern-Based Discovery:** Use glob patterns (e.g., `services/*`) to resolve sub-projects within monorepo roots.
* **Alias Registry:** Custom short-codes for projects to bypass long FQNs (Fully Qualified Names).

### 4.2. CLI Command Suite
* **`beavers project list`:** Display all discovered projects across all workspaces.
    * *Columns:* Project Name, Alias, Type, Git Branch, Sync Status, Health.
    * *Modifiers:* `--output=json|markdown|html`.
* **`beavers path <alias>`:** Resolve and print the absolute path of a project for shell integration.
* **`beavers svc <install|build|pull> <alias>`:** Change directory to the project and execute the corresponding `make` target.
* **`beavers doctor`:** Validate the local environment (Git access, Go/PHP versions, Docker daemon).

### 4.3. Standards & Audit Engine (Evolutive)
* **Rule Definitions:** Declarative standards in the global config.
* **Checker Modules:**
    * `file_exists`: Verify required files (e.g., `Makefile`, `README`).
    * `makefile_target`: Verify if specific targets exist in the project Makefile.
    * `composer_json_version`: (Advanced) Parse JSON to check for specific library constraints.
* **`beavers svc audit <alias>`:** Run all applicable checks and return a detailed report.

### 4.4. TUI (Terminal User Interface)
* **Dashboard View:** A split-panel interface for monitoring projects.
* **Interactive Navigation:** Filter/Search projects and trigger actions (`i` to install, `a` to audit) via hotkeys.
* **Background Refreshing:** Asynchronous Git fetching to keep sync status current without blocking the UI.

---

## 5. Technical Requirements & Constraints
* **Performance:** Discovery must be cached in `~/.beavers/cache.json`. Any command reading the list must return in < 100ms.
* **Shell Integration:** A helper function (`bgo`) must be provided for `.zshrc/.bashrc` to enable directory jumping.
* **Safety:** The CLI must never perform `git push` or destructive operations automatically.
* **Binary Integrity:** Configuration loading must explicitly look for `.yaml` or `.yml` files to avoid collision with the `beavers` binary.
* **Decoupling:** Discovery logic must not require any "Beavers-specific" files inside the service repositories.

---

## 6. Agile Roadmap

### Phase 1: The Pathfinder (Discovery) ✅
* [x] Cobra/Viper boilerplate.
* [x] Global YAML parsing logic (Flag > Local > Global precedence).
* [x] Folder walker with `standalone` and `monorepo-sub` support.
* [x] Workspace `type` and `excludes` filtering.
* [x] `beavers project list` (pterm table output).
* [x] `beavers path` command for shell `cd` integration.
* [x] Built-in `completion` command.

### Phase 2: The Executor (Workflow) ✅
* [x] `os/exec` wrapper for `make` targets.
* [x] Git branch detection logic.
* [x] Local cache implementation with stale-while-revalidate logic.

### Phase 3: The Auditor (Compliance) ✅
* [x] Core Audit engine.
* [x] File and Makefile checkers.
* [x] Roadmap for `composer.json` parsing.
* [x] `beavers svc audit` report command.

### Phase 4: The TUI Dashboard (Monitoring)
* [ ] Bubble Tea integration.
* [ ] Project list navigation and status icons.

---

## 7. Future Considerations
* **Docker Integration:** Monitoring container status per service.
* **Auto-Fix:** Scaffold missing files (Makefiles/READMEs) based on audit failures.
* **Team Sync:** A way to share the global configuration rules across the team via a central repo.
