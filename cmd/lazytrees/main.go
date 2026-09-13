package main

import (
	"fmt"

	"lazytrees/internal/git"
)

func main() {
	fmt.Println("Main Repo Root:", git.GetCommonRoot())
	fmt.Println("Project:", git.GetProjectName())
	fmt.Println("Current Root:", git.GetCurrentRoot())
	fmt.Println("Current branch: " + git.GetCurrentBranch())
	fmt.Println("")
	worktrees, err := git.GetWorktrees()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("---")
	for i := range worktrees {
		tree := worktrees[i]
		fmt.Println(tree.Alias, tree.Path)
	}
	fmt.Println("---")
	worktree, createErr := git.CreateWorktree()
	if createErr != nil {
		fmt.Println("Error: creating", worktree.Alias, "\n", createErr)
		return
	}

	fmt.Println("Created new worktree:", worktree.Alias)

	branchErr := git.CheckoutBranch(worktree, "test-branch")

	if branchErr != nil {
		fmt.Println(branchErr)
	}
}
