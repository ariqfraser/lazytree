// Provides execution of git comands and retrieval of git information.
package git

import (
	"os/exec"
	"strings"
)

func GetRepoRoot() string {
	root, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(root))
}

func GetCurrentBranch() string {
	branch, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(branch))
}
