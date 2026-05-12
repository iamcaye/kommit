package picker

import (
	"fmt"
	"kommit/ai"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// Pick opens an fzf picker for the given messages and returns the selected
// commit message text (possibly edited by the user via ctrl-e).
// Returns ("", nil) if the user aborts without selecting.
func Pick(messages []ai.CommitMessage) (string, error) {
	if _, err := exec.LookPath("fzf"); err != nil {
		return "", fmt.Errorf("install fzf (https://github.com/junegunn/fzf)")
	}

	tmpDir, err := os.MkdirTemp("", "kommit-*")
	if err != nil {
		return "", fmt.Errorf("could not create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	var lines []string
	for i, msg := range messages {
		path := filepath.Join(tmpDir, strconv.Itoa(i))
		content := msg.Subject
		if msg.Body != "" {
			content += "\n\n" + msg.Body
		}
		if err := os.WriteFile(path, []byte(content), 0600); err != nil {
			return "", fmt.Errorf("could not write message file: %w", err)
		}
		lines = append(lines, fmt.Sprintf("%d\t%s", i, msg.Subject))
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}

	previewCmd := fmt.Sprintf("cat %s/{1}", tmpDir)
	bindCmd := fmt.Sprintf("ctrl-e:execute(%s %s/{1})", editor, tmpDir)

	cmd := exec.Command("fzf",
		"--with-nth=2",
		"--delimiter=\t",
		"--preview="+previewCmd,
		"--preview-window=right:60%:wrap",
		"--bind="+bindCmd,
		"--prompt=commit> ",
		"--header=Enter to commit  •  ctrl-e to edit",
	)
	cmd.Stdin = strings.NewReader(strings.Join(lines, "\n"))
	cmd.Stderr = os.Stderr

	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 130 {
			// User aborted with ctrl-c or ESC — not an error
			return "", nil
		}
		return "", fmt.Errorf("fzf error: %w", err)
	}

	selected := strings.TrimSpace(string(out))
	if selected == "" {
		return "", nil
	}

	parts := strings.SplitN(selected, "\t", 2)
	idx, err := strconv.Atoi(parts[0])
	if err != nil {
		return "", fmt.Errorf("unexpected fzf output: %q", selected)
	}

	// Read the temp file — may have been edited by ctrl-e
	content, err := os.ReadFile(filepath.Join(tmpDir, strconv.Itoa(idx)))
	if err != nil {
		return "", fmt.Errorf("could not read selected message: %w", err)
	}
	return strings.TrimSpace(string(content)), nil
}
