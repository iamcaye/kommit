package ai

import (
	"testing"
)

func TestDetectBackend_FlagOverrideClaude(t *testing.T) {
	backend, err := DetectBackend("claude")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if backend != BackendClaude {
		t.Errorf("expected claude, got %q", backend)
	}
}

func TestDetectBackend_FlagOverrideCodex(t *testing.T) {
	backend, err := DetectBackend("codex")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if backend != BackendCodex {
		t.Errorf("expected codex, got %q", backend)
	}
}

func TestDetectBackend_InvalidFlag(t *testing.T) {
	_, err := DetectBackend("gpt99")
	if err == nil {
		t.Fatal("expected error for unknown backend")
	}
}

func TestDetectBackend_AutoDetect(t *testing.T) {
	backend, err := DetectBackend("")
	if err != nil {
		// Acceptable — neither CLI installed in this environment
		t.Logf("no AI CLI found (expected in CI): %v", err)
		return
	}
	if backend != BackendClaude && backend != BackendCodex {
		t.Errorf("unexpected backend: %q", backend)
	}
}
