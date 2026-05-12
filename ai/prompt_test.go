package ai

import (
	"strings"
	"testing"
)

func TestBuildPrompt_ContainsDiff(t *testing.T) {
	diff := "diff --git a/main.go b/main.go\n+++ added line"
	prompt := BuildPrompt(diff)
	if !strings.Contains(prompt, diff) {
		t.Error("prompt does not contain the diff")
	}
}

func TestBuildPrompt_ContainsJSONInstruction(t *testing.T) {
	prompt := BuildPrompt("some diff")
	if !strings.Contains(prompt, "JSON") {
		t.Error("prompt does not mention JSON output format")
	}
	if !strings.Contains(prompt, "subject") {
		t.Error("prompt does not mention subject field")
	}
}

func TestBuildPrompt_ContainsConventionalCommits(t *testing.T) {
	prompt := BuildPrompt("some diff")
	if !strings.Contains(strings.ToLower(prompt), "conventional") {
		t.Error("prompt does not mention Conventional Commits")
	}
}
