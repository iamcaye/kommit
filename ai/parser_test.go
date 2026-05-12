package ai

import (
	"os"
	"testing"
)

func TestParseResponse_ValidJSON(t *testing.T) {
	raw, err := os.ReadFile("testdata/valid.json")
	if err != nil {
		t.Fatal(err)
	}
	msgs, err := ParseResponse(string(raw))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(msgs) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(msgs))
	}
	if msgs[0].Subject != "fix: correct off-by-one in pagination" {
		t.Errorf("unexpected subject: %q", msgs[0].Subject)
	}
	if msgs[0].Body != "" {
		t.Errorf("expected empty body, got %q", msgs[0].Body)
	}
	if msgs[1].Subject != "fix(pagination): correct off-by-one error in page count" {
		t.Errorf("unexpected subject: %q", msgs[1].Subject)
	}
}

func TestParseResponse_Fallback(t *testing.T) {
	raw, err := os.ReadFile("testdata/fallback.txt")
	if err != nil {
		t.Fatal(err)
	}
	msgs, err := ParseResponse(string(raw))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(msgs) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(msgs))
	}
	if msgs[0].Subject != "fix: correct off-by-one in pagination" {
		t.Errorf("unexpected subject: %q", msgs[0].Subject)
	}
	if msgs[1].Subject != "fix(pagination): correct off-by-one error in page count" {
		t.Errorf("unexpected subject: %q", msgs[1].Subject)
	}
	if msgs[1].Body != "The page count was calculated before slicing." {
		t.Errorf("unexpected body: %q", msgs[1].Body)
	}
}

func TestParseResponse_Unparseable(t *testing.T) {
	// Only whitespace — both JSON and fallback parser produce no messages
	_, err := ParseResponse("   \n\n   \n\n   ")
	if err == nil {
		t.Fatal("expected error for empty/whitespace input")
	}
}
