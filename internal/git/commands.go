// Provides execution of git comands and retrieval of git information.
package git

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

type gitInfo struct {
	root string
}

func getOriginURL() (string, error) {
	stdout, err := exec.Command("git", "config", "--get", "remote.origin.url").Output()
	if err != nil {
		return "", fmt.Errorf("get origin: %s", string(stdout))
	}
	return strings.TrimSpace(string(stdout)), nil
}

func GetProjectName() (string, error) {
	originURL, err := getOriginURL()
	if err != nil {
		return "", err
	}

	return strings.Replace(filepath.Base(originURL), ".git", "", 1), nil
}

func GetCurrentRoot() string {
	root, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return ""
	}
	rootStr := strings.TrimSpace(string(root))
	return filepath.Clean(rootStr)
}

func GetCurrentBranch() string {
	branch, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(branch))
}

func GetCommonRoot() string {
	stdout, err := exec.Command("git", "rev-parse", "--path-format=absolute", "--git-common-dir").Output()
	if err != nil {
		return ""
	}

	gitRoot := filepath.Clean(strings.TrimSpace(string(stdout)))

	if filepath.Base(gitRoot) == ".git" {
		return filepath.Dir(gitRoot)
	}

	return string(gitRoot)
}
