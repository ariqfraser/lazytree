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

type Worktree struct {
	Alias    string
	Path     string
	Head     string
	Branch   string
	Detached bool
}

func GetWorktrees() ([]Worktree, error) {
	if _, err := CleanWorktreeMetadata(); err != nil {
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

	var worktrees []Worktree
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

			worktrees = append(worktrees, Worktree{
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

func CreateWorktree() (Worktree, error) {
	alias, path, candidateErr := generateWorktreeCandidate()

	newTree := Worktree{}

	if candidateErr != nil {
		return newTree, candidateErr
	}

	cmd := exec.Command("git", "worktree", "add", "--detach", path)
	cmd.Dir = GetCommonRoot()
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		message := strings.TrimSpace(string(stdout))
		return newTree, fmt.Errorf(
			"create worktree %q: %w: %s",
			path,
			err,
			message,
		)
	}

	newTree.Alias = alias
	newTree.Path = path
	newTree.Detached = true

	return newTree, nil
}

func (tree Worktree) CheckoutBranch(branch string) error {
	cmd := exec.Command("git", "checkout", branch)
	cmd.Dir = tree.Path

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("checkout err %q", string(out))
	}

	return nil
}

func (tree Worktree) Delete() error {
	if tree.Path == "" {
		return fmt.Errorf("delete worktree: path is empty")
	}

	mainRoot := GetCommonRoot()
	if mainRoot == "" {
		return fmt.Errorf("delete worktree: could not determine repository root")
	}

	cmd := exec.Command("git", "worktree", "remove", tree.Path)
	cmd.Dir = mainRoot

	output, err := cmd.CombinedOutput()
	if err != nil {
		message := strings.TrimSpace(string(output))
		if message != "" {
			return fmt.Errorf("delete worktree %q: %w: %s", tree.Path, err, message)
		}
		return fmt.Errorf("delete worktree %q: %w", tree.Path, err)
	}

	return nil
}
