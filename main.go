package main

import (
	"flag"
	"fmt"
	"kommit/ai"
	"kommit/git"
	"kommit/picker"
	"os"

	"golang.org/x/term"
)

func waitKey() {
	fd := int(os.Stdin.Fd())
	if term.IsTerminal(fd) {
		old, err := term.MakeRaw(fd)
		if err == nil {
			defer term.Restore(fd, old)
		}
	}
	b := make([]byte, 1)
	os.Stdin.Read(b)
}

func main() {
	aiFlag := flag.String("ai", "", "AI backend to use: claude or codex (default: auto-detect)")
	flag.Parse()

	diff, err := git.GetStagedDiff(".")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if diff == "" {
		fmt.Fprintln(os.Stderr, "nothing staged — run `git add` first")
		fmt.Fprint(os.Stderr, "press any key to exit...")
		waitKey()
		fmt.Fprintln(os.Stderr)
		os.Exit(0)
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
