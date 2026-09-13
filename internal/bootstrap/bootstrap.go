// Package bootstraps
package bootstrap

import (
	"fmt"
	"os"

	"lazytrees/internal/git"
)

type AppModel struct {
	projectName string
	worktrees   []git.Worktree

	width  int
	height int
}

func Bootstrap() {
	projectName, err := git.GetProjectName()
	if err != nil {
		fmt.Println("FATAL: lazytrees must be called inside a repository")
		os.Exit(1)
	}

	app := AppModel{projectName: projectName}
	fmt.Println(app)

	// fmt.Println("Main Repo Root:", git.GetCommonRoot())
	// fmt.Println("Project:", git.GetProjectName())
	// fmt.Println("Current Root:", git.GetCurrentRoot())
	// fmt.Println("Current branch: " + git.GetCurrentBranch())
	// fmt.Println("")
	// worktrees, err := git.GetWorktrees()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println("---")
	// for i := range worktrees {
	// 	tree := worktrees[i]
	// 	fmt.Println(tree.Alias, tree.Path)
	// }
	// fmt.Println("---")
	// worktree, createErr := git.CreateWorktree()
	// if createErr != nil {
	// 	fmt.Println("Error: creating", worktree.Alias, "\n", createErr)
	// 	return
	// }
	//
	// fmt.Println("Created new worktree:", worktree.Alias)
	//
	// branchErr := worktree.CheckoutBranch("test-branch")
	//
	// if branchErr != nil {
	// 	fmt.Println(branchErr)
	// }
}
