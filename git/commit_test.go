package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCommit_WritesMessage(t *testing.T) {
	dir := setupRepo(t)

	file := filepath.Join(dir, "hello.txt")
	if err := os.WriteFile(file, []byte("hello\n"), 0644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("git", "add", "hello.txt")
	cmd.Dir = dir
	cmd.Run()

	msg := "feat: add hello\n\nThis adds the hello file."
	if err := Commit(dir, msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	logCmd := exec.Command("git", "log", "--format=%B", "-1")
	logCmd.Dir = dir
	out, _ := logCmd.Output()
	log := strings.TrimSpace(string(out))

	if !strings.Contains(log, "feat: add hello") {
		t.Errorf("subject not found in log: %s", log)
	}
	if !strings.Contains(log, "This adds the hello file.") {
		t.Errorf("body not found in log: %s", log)
	}
}
