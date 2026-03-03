package discovery

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ubitransports/beavers/internal/app"
	"github.com/ubitransports/beavers/internal/config"
)

func DiscoverProjects(appCtx *app.App) ([]config.Project, error) {
	var projects []config.Project

	for _, ws := range appCtx.Config.Workspaces {
		wsRoot := ws.Root
		if !filepath.IsAbs(wsRoot) {
			abs, err := filepath.Abs(wsRoot)
			if err == nil {
				wsRoot = abs
			}
		}

		if len(ws.Patterns) == 0 {
			p := discoverStandalone(appCtx, ws, wsRoot)
			if p != nil {
				projects = append(projects, *p)
			}
		} else {
			ps := discoverMonorepo(appCtx, ws, wsRoot)
			projects = append(projects, ps...)
		}
	}

	return projects, nil
}

func discoverStandalone(appCtx *app.App, ws config.Workspace, wsRoot string) *config.Project {
	name := filepath.Base(wsRoot)
	if isExcluded(name, ws.Excludes) {
		return nil
	}

	pType := "standalone"
	if ws.Type != "" {
		pType = ws.Type
	}
	branch, sync := appCtx.Git.DetectStatus(wsRoot)

	p := config.Project{
		ID:         fmt.Sprintf("%s-%s", ws.Name, name),
		Name:       name,
		Path:       wsRoot,
		Type:       pType,
		Workspace:  ws.Name,
		GitBranch:  branch,
		SyncStatus: sync,
	}
	p.Alias = resolveAlias(appCtx.Config.Aliases, name, wsRoot)
	return &p
}

func discoverMonorepo(appCtx *app.App, ws config.Workspace, wsRoot string) []config.Project {
	var projects []config.Project

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
			if isExcluded(name, ws.Excludes) {
				continue
			}

			pType := "monorepo-sub"
			if ws.Type != "" {
				pType = ws.Type
			}
			branch, sync := appCtx.Git.DetectStatus(match)

			p := config.Project{
				ID:         fmt.Sprintf("%s-%s", ws.Name, name),
				Name:       name,
				Path:       match,
				Type:       pType,
				Workspace:  ws.Name,
				GitBranch:  branch,
				SyncStatus: sync,
			}
			p.Alias = resolveAlias(appCtx.Config.Aliases, name, match)
			projects = append(projects, p)
		}
	}
	return projects
}

func isExcluded(name string, excludes []string) bool {
	for _, ex := range excludes {
		if ex == name {
			return true
		}
	}
	return false
}

func resolveAlias(aliases map[string]string, name, path string) string {
	for alias, target := range aliases {
		if target == name || target == path {
			return alias
		}
	}
	return ""
}

// ResolveProjectByAlias finds a project by its alias or name.
func ResolveProjectByAlias(projects []config.Project, identifier string) (*config.Project, error) {
	for _, p := range projects {
		if p.Alias == identifier || p.Name == identifier {
			return &p, nil
		}
	}
	return nil, fmt.Errorf("project with alias or name '%s' not found", identifier)
}
