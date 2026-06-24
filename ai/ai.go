package ai

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Backend string

const (
	BackendClaude Backend = "claude"
	BackendCodex  Backend = "codex"
)

func DetectBackend(flag string) (Backend, error) {
	switch flag {
	case "claude":
		return BackendClaude, nil
	case "codex":
		return BackendCodex, nil
	case "":
		if _, err := exec.LookPath("claude"); err == nil {
			return BackendClaude, nil
		}
		if _, err := exec.LookPath("codex"); err == nil {
			return BackendCodex, nil
		}
		return "", fmt.Errorf("install claude or codex and ensure it is in PATH")
	default:
		return "", fmt.Errorf("unknown AI backend %q: use claude or codex", flag)
	}
}

// Generate runs the full pipeline: build prompt → invoke AI CLI → parse response.
func Generate(diff string, backend Backend) ([]CommitMessage, error) {
	prompt := BuildPrompt(diff)

	done := make(chan struct{})
	go showSpinner(done)
	raw, err := invoke(backend, prompt)
	close(done)

	if err != nil {
		return nil, err
	}

	msgs, err := ParseResponse(raw)
	if err != nil {
		return nil, fmt.Errorf("could not parse AI response:\n%s\n\nerror: %w", raw, err)
	}
	return msgs, nil
}

func invoke(backend Backend, prompt string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var cmd *exec.Cmd
	switch backend {
	case BackendClaude:
		cmd = exec.CommandContext(ctx, "claude", "--model", "claude-haiku-4-5-20251001", "--print", "--no-session-persistence")
	case BackendCodex:
		cmd = exec.CommandContext(ctx, "codex", "exec", "-m", "gpt-4o-mini", "--ephemeral")
	}

	cmd.Stdin = strings.NewReader(prompt)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("AI timed out after 60s")
		}
		return "", fmt.Errorf("AI CLI error: %s", strings.TrimSpace(stderr.String()))
	}

	return strings.TrimSpace(stdout.String()), nil
}

func showSpinner(done <-chan struct{}) {
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	i := 0
	for {
		select {
		case <-done:
			fmt.Fprint(os.Stderr, "\r\033[K")
			return
		default:
			fmt.Fprintf(os.Stderr, "\r%s Generating commit messages...", frames[i%len(frames)])
			time.Sleep(80 * time.Millisecond)
			i++
		}
	}
}
