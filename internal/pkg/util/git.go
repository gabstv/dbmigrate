package util

import (
	"os/exec"
	"strings"
)

// GitRoot retrieves the root folder of the project from Git
func GitRoot() (string, error) {
	path, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	return strings.TrimSpace(string(path)), err
}
