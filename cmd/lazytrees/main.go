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
	for i := range worktrees {
		tree := worktrees[i]
		fmt.Println(tree.Alias, tree.Path)
	}
}
