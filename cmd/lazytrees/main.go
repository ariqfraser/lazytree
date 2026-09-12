package main

import (
	"fmt"

	"lazytrees/internal/git"
)

func main() {
	root := git.GetRepoRoot()
	fmt.Println("Hello, World!" + root)
	fmt.Println("Current branch: " + git.GetCurrentBranch())
}
