package ai

import (
	"encoding/json"
	"errors"
	"strings"
)

type CommitMessage struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// ParseResponse parses the raw AI output into commit messages.
// Tries JSON first; falls back to blank-line-separated blocks.
func ParseResponse(raw string) ([]CommitMessage, error) {
	raw = strings.TrimSpace(raw)

	// Strip markdown code fences the AI may wrap around JSON
	if start := strings.Index(raw, "["); start != -1 {
		if end := strings.LastIndex(raw, "]"); end > start {
			raw = raw[start : end+1]
		}
	}

	var msgs []CommitMessage
	if err := json.Unmarshal([]byte(raw), &msgs); err == nil && len(msgs) > 0 {
		return msgs, nil
	}

	return parseFallback(raw)
}

func parseFallback(raw string) ([]CommitMessage, error) {
	blocks := strings.Split(strings.TrimSpace(raw), "\n\n")
	var msgs []CommitMessage
	for _, block := range blocks {
		block = strings.TrimSpace(block)
		if block == "" {
			continue
		}
		lines := strings.SplitN(block, "\n", 2)
		msg := CommitMessage{Subject: strings.TrimSpace(lines[0])}
		if len(lines) == 2 {
			msg.Body = strings.TrimSpace(lines[1])
		}
		if msg.Subject != "" {
			msgs = append(msgs, msg)
		}
	}
	if len(msgs) == 0 {
		return nil, errors.New("could not parse AI response")
	}
	return msgs, nil
}
