package main

import (
	"flag"
	"fmt"
	"kommit/ai"
	"kommit/git"
	"kommit/picker"
	"os"
)

func main() {
	aiFlag := flag.String("ai", "", "AI backend to use: claude or codex (default: auto-detect)")
	flag.Parse()

	diff, err := git.GetStagedDiff(".")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	backend, err := ai.DetectBackend(*aiFlag)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	msgs, err := ai.Generate(diff, backend)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	message, err := picker.Pick(msgs)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if message == "" {
		os.Exit(0)
	}

	if err := git.Commit(".", message); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
