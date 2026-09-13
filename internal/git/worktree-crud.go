// Package git provides CRUD operations for Git worktrees.
package git

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

type WorktreeStatus struct {
	Alias    string
	Path     string
	Head     string
	Branch   string
	Detached bool
}

func GetWorktrees() ([]WorktreeStatus, error) {
	if err := cleanWorktreeMetadata(); err != nil {
		return nil, err
	}

	mainRoot := GetCommonRoot()
	if mainRoot == "" {
		return nil, fmt.Errorf("could not determine repository root")
	}

	cmd := exec.Command("git", "worktree", "list", "--porcelain")
	cmd.Dir = mainRoot

	stdout, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("list worktrees: %w", err)
	}

	var worktrees []WorktreeStatus
	scanner := bufio.NewScanner(bytes.NewReader(stdout))
	current := -1

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		key, value, _ := strings.Cut(line, " ")

		switch key {
		case "worktree":
			path := filepath.Clean(value)

			worktrees = append(worktrees, WorktreeStatus{
				Path:  path,
				Alias: getWorktreeAlias(path, mainRoot),
			})

			current = len(worktrees) - 1

		case "HEAD":
			if current >= 0 {
				worktrees[current].Head = value
			}

		case "branch":
			if current >= 0 {
				worktrees[current].Branch = value
			}

		case "detached":
			if current >= 0 {
				worktrees[current].Detached = true
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("parse worktrees: %w", err)
	}

	return worktrees, nil
}

func CreateWorktree() (string, error) {
	alias, path, candidateErr := generateWorktreeCandidate()

	if candidateErr != nil {
		return alias, candidateErr
	}

	cmd := exec.Command("git", "worktree", "add", "--detach", path)
	cmd.Dir = GetCommonRoot()
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		message := strings.TrimSpace(string(stdout))
		return alias, fmt.Errorf(
			"create worktree %q: %w: %s",
			path,
			err,
			message,
		)
	}

	return alias, nil
}
