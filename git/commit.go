package git

import (
	"fmt"
	"os"
	"os/exec"
)

// Commit writes message to a temp file and runs `git commit -F <tmpfile>` from dir.
func Commit(dir, message string) error {
	f, err := os.CreateTemp("", "kommit-*.txt")
	if err != nil {
		return fmt.Errorf("could not create temp file: %w", err)
	}
	defer os.Remove(f.Name())

	if _, err := f.WriteString(message); err != nil {
		f.Close()
		return fmt.Errorf("could not write commit message: %w", err)
	}
	f.Close()

	cmd := exec.Command("git", "commit", "-F", f.Name())
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git commit failed")
	}
	return nil
}
