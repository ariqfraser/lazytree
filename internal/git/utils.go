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
	repoName, err := GetProjectName()
	if err != nil {
		return "", "", err
	}
	alias := generateWorktreeAlias()
	path := filepath.Join(getLazytreesRoot(), repoName, alias, repoName)

	// [TODO] Validate alias

	return alias, path, nil
}

func CleanWorktreeMetadata() (string, error) {
	mainRoot := GetCommonRoot()
	if mainRoot == "" {
		return "", fmt.Errorf("could not determine repository root")
	}

	dryRunCmd := exec.Command("git", "worktree", "prune", "--dry-run")
	dryRunCmd.Dir = mainRoot

	dryRunOut, err := dryRunCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("dry run worktree prune: %w", err)
	}

	pruneCmd := exec.Command("git", "worktree", "prune")
	pruneCmd.Dir = mainRoot

	pruneOutput, err := pruneCmd.CombinedOutput()
	if err != nil {
		pruneOutput := strings.TrimSpace(string(pruneOutput))
		if pruneOutput != "" {
			return "", fmt.Errorf("prune error: %s", pruneOutput)
		}
	}

	return strings.TrimSpace(string(dryRunOut)), nil
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
