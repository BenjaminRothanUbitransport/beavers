package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// discoverProjects walks through configured workspaces and identifies logical projects.
func discoverProjects(cfg *Config) ([]Project, error) {
	var projects []Project

	for _, ws := range cfg.Workspaces {
		wsRoot := ws.Root
		if !filepath.IsAbs(wsRoot) {
			// If not absolute, consider it relative to the workspace file (not implemented here, using current dir)
			abs, err := filepath.Abs(wsRoot)
			if err == nil {
				wsRoot = abs
			}
		}

		if len(ws.Patterns) == 0 {
			// Standalone project: root is the project itself
			name := filepath.Base(wsRoot)
			
			// Skip if excluded
			excluded := false
			for _, ex := range ws.Excludes {
				if ex == name {
					excluded = true
					break
				}
			}
			if excluded {
				continue
			}

			pType := "standalone"
			if ws.Type != "" {
				pType = ws.Type
			}
			branch, sync := detectGitStatus(wsRoot)

			p := Project{
				ID:         fmt.Sprintf("%s-%s", ws.Name, name),
				Name:       name,
				Path:       wsRoot,
				Type:       pType,
				Workspace:  ws.Name,
				GitBranch:  branch,
				SyncStatus: sync,
			}
			// Resolve alias from config (matching by name or absolute path)
			for alias, target := range cfg.Aliases {
				if target == name || target == wsRoot {
					p.Alias = alias
					break
				}
			}
			projects = append(projects, p)
		} else {
			// Monorepo: match sub-projects using glob patterns
			for _, pattern := range ws.Patterns {
				fullPattern := filepath.Join(wsRoot, pattern)
				matches, err := filepath.Glob(fullPattern)
				if err != nil {
					continue
				}

				for _, match := range matches {
					info, err := os.Stat(match)
					if err != nil || !info.IsDir() {
						continue
					}

					name := filepath.Base(match)

					// Skip if excluded
					excluded := false
					for _, ex := range ws.Excludes {
						if ex == name {
							excluded = true
							break
						}
					}
					if excluded {
						continue
					}

					pType := "monorepo-sub"
					if ws.Type != "" {
						pType = ws.Type
					}
					branch, sync := detectGitStatus(match)

					p := Project{
						ID:         fmt.Sprintf("%s-%s", ws.Name, name),
						Name:       name,
						Path:       match,
						Type:       pType,
						Workspace:  ws.Name,
						GitBranch:  branch,
						SyncStatus: sync,
					}
					// Resolve alias from config (matching by name or absolute path)
					for alias, target := range cfg.Aliases {
						if target == name || target == match {
							p.Alias = alias
							break
						}
					}
					projects = append(projects, p)
				}
			}
		}
	}

	return projects, nil
}

// resolveProjectByAlias finds a project by its alias or name.
func resolveProjectByAlias(projects []Project, identifier string) (*Project, error) {
	for _, p := range projects {
		if p.Alias == identifier || p.Name == identifier {
			return &p, nil
		}
	}
	return nil, fmt.Errorf("project with alias or name '%s' not found", identifier)
}
