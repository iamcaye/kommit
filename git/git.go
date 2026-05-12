package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// GetStagedDiff returns `git diff --staged` output from dir.
func GetStagedDiff(dir string) (string, error) {
	cmd := exec.Command("git", "diff", "--staged")
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%s", strings.TrimSpace(stderr.String()))
	}

	diff := strings.TrimSpace(stdout.String())
	if diff == "" {
		return "", fmt.Errorf("nothing staged — run `git add` first")
	}
	return diff, nil
}
