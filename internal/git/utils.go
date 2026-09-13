// Package git : utils for git actions
package git

import (
	"fmt"
	"math/rand"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	prefixes = []string{"fattened", "plump", "iron", "happy", "crying", "sad"}
	nouns    = []string{"oyster", "kitty", "pumpkin", "giraffe", "donky", "firefly"}
)

const lazytreePattern = `[.]lazytrees[\\/][A-Za-z0-9_-]+[\\/]([A-Za-z0-9_-]+)(?:[\\/]|$)`

const LazyTreeDir = ".lazytrees"

var aliasRegexp = regexp.MustCompile(lazytreePattern)

func generateWorktreeAlias() string {
	fmt.Println("Max alias combos:", len(prefixes)*len(nouns))
	return prefixes[rand.Intn(len(prefixes))] + "-" + nouns[rand.Intn(len(nouns))]
}

func getLazytreesRoot() string {
	return filepath.Join(filepath.Dir(GetCommonRoot()), LazyTreeDir)
}

// returns Alias, Path, error
func generateWorktreeCandidate() (string, string, error) {
	repoName := GetProjectName()
	alias := generateWorktreeAlias()
	path := filepath.Join(getLazytreesRoot(), repoName, alias, repoName)

	// [TODO] Validate alias

	return alias, path, nil
}

func cleanWorktreeMetadata() error {
	mainRoot := GetCommonRoot()
	if mainRoot == "" {
		return fmt.Errorf("could not determine repository root")
	}

	dryRunCmd := exec.Command("git", "worktree", "prune", "--dry-run")
	dryRunCmd.Dir = mainRoot

	_, err := dryRunCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("dry run worktree prune: %w", err)
	}

	pruneCmd := exec.Command("git", "worktree", "prune")
	pruneCmd.Dir = mainRoot

	output, err := pruneCmd.CombinedOutput()
	if err != nil {
		if output := strings.TrimSpace(string(output)); output != "" {
			fmt.Println(output)
		}

		return fmt.Errorf("prune worktrees: %w", err)
	}

	return nil
}

func getWorktreeAlias(path, mainRoot string) string {
	if path == mainRoot {
		return "main"
	}

	if match := aliasRegexp.FindStringSubmatch(path); len(match) > 1 {
		return match[1]
	}

	parent := filepath.Base(filepath.Dir(path))
	name := filepath.Base(path)

	return filepath.Join(parent, name)
}
