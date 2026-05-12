package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func setupRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	run := func(args ...string) {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("command %v failed: %s", args, out)
		}
	}
	run("git", "init")
	run("git", "config", "user.email", "test@test.com")
	run("git", "config", "user.name", "Test")
	return dir
}

func TestGetStagedDiff_WithStagedChanges(t *testing.T) {
	dir := setupRepo(t)
	file := filepath.Join(dir, "hello.txt")
	if err := os.WriteFile(file, []byte("hello world\n"), 0644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("git", "add", "hello.txt")
	cmd.Dir = dir
	cmd.Run()

	diff, err := GetStagedDiff(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if diff == "" {
		t.Fatal("expected non-empty diff")
	}
	if !strings.Contains(diff, "hello world") {
		t.Errorf("diff does not contain staged content:\n%s", diff)
	}
}

func TestGetStagedDiff_NothingStaged(t *testing.T) {
	dir := setupRepo(t)
	_, err := GetStagedDiff(dir)
	if err == nil {
		t.Fatal("expected error when nothing is staged")
	}
	if !strings.Contains(err.Error(), "nothing staged") {
		t.Errorf("unexpected error message: %v", err)
	}
}
