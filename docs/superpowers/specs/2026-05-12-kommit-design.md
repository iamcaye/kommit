# kommit — Design Spec

**Date:** 2026-05-12
**Status:** Approved

## Overview

`kommit` is a terminal CLI written in Go that generates git commit messages using an AI CLI (`claude` or `codex`). It reads the staged diff, asks the AI to produce 1–4 commit message options scaled to the complexity of the changes, presents them in an `fzf` picker with a preview pane, and performs the commit when the user selects one.

## Architecture

Single Go binary. No subcommands. Flags:

- `--ai <claude|codex>` — force a specific AI backend (default: auto-detect)

### Package Layout

```
kommit/
├── main.go
├── git/
│   ├── git.go        — run git diff --staged, validate repo state
│   └── commit.go     — run git commit -F <tmpfile>
├── ai/
│   ├── ai.go         — detect available AI CLI, invoke it, return raw output
│   ├── prompt.go     — build the structured prompt sent to the AI
│   └── parser.go     — parse AI JSON output into []CommitMessage
└── picker/
    └── picker.go     — build fzf invocation, handle ctrl+e edit shortcut
```

### Core Types

```go
type CommitMessage struct {
    Subject string
    Body    string
}
```

## Data Flow

```
kommit
  │
  ├─ git diff --staged          → empty? exit "nothing staged"
  │
  ├─ detect claude / codex      → neither found? exit with install hint
  │
  ├─ call AI CLI                → 30s timeout, spinner while waiting
  │                             → parse failure? show raw output + exit
  │
  ├─ fzf picker
  │   ├─ preview pane: full message (subject + blank line + body)
  │   └─ ctrl+e: open in $EDITOR, commit with edited result
  │
  └─ git commit -F <tmpfile>    → failure? surface git's error message
```

## AI Integration

### Backend Detection

Auto-detect order: `claude` first, `codex` second. Detection is a simple `exec.LookPath`. `--ai` flag overrides.

### Invocation

- **claude:** `claude --print "<prompt>"`
- **codex:** `codex "<prompt>"`

Both are called with a 30-second timeout. A spinner is shown on stderr while waiting.

### Prompt

The prompt instructs the AI to:

1. Analyze the diff and assess complexity (trivial / moderate / complex)
2. Return a JSON array of commit message options:
   - 1 option for trivial changes
   - 2–3 for moderate
   - 3–4 for complex/multi-file
3. Follow Conventional Commits format (`type(scope): subject`)
4. Include a body only when the change warrants explanation

**Output format:**

```json
[
  {"subject": "fix: correct off-by-one in pagination", "body": ""},
  {"subject": "fix(pagination): correct off-by-one error in page count", "body": "The page count was calculated before slicing, causing the last page\nto always appear empty when results were an exact multiple of page size."}
]
```

### Parsing

Primary: JSON unmarshal into `[]CommitMessage`.

Fallback: if JSON parse fails, split output on blank lines and treat each block's first line as `Subject` and the rest as `Body`. If fallback also fails, print raw AI output and exit with an error.

## Picker UI

Each `CommitMessage` is written to a numbered temp file (`0`, `1`, `2`, …) in a temp directory before `fzf` is launched. The format of each file is:

```
{subject}

{body}
```

`fzf` is invoked with:

- Input: lines of the form `{index}\t{subject}`
- `--with-nth=2` — display only the subject column
- `--preview "cat <tmpdir>/{1}"` — preview reads the temp file for the highlighted index
- `--bind "ctrl+e:execute($EDITOR <tmpdir>/{1})+abort"` — opens the selected message's file in `$EDITOR`; `kommit` detects the `ctrl+e` abort code, reads the edited file, and commits it

If the user aborts normally (`ctrl+c` / `ESC`), `fzf` exits non-zero and `kommit` exits cleanly with no commit performed. The `ctrl+e` path uses a distinct exit code (`130` vs `1`) to distinguish an edit-then-commit from a plain abort.

## Commit

The selected (or edited) message is written to a temp file and passed to:

```
git commit -F <tmpfile>
```

Multi-line messages (subject + body) are preserved correctly this way.

## Error Handling

| Condition | Behavior |
|---|---|
| No staged changes | Exit with message: "nothing staged — run `git add` first" |
| Not in a git repo | Exit with git's own error message |
| No AI CLI found | Exit with: "install claude or codex and ensure it is in PATH" |
| AI times out (>30s) | Exit with: "AI timed out after 30s" |
| AI returns unparseable output | Print raw output, exit with error |
| `fzf` not found | Exit with: "install fzf (https://github.com/junegunn/fzf)" |
| User aborts picker | Exit 0, no commit |
| `git commit` fails | Surface git's stderr, exit with error |

## Testing

- **Unit — `ai/parser.go`:** valid JSON, malformed JSON triggering fallback, fully unparseable output
- **Unit — `ai/prompt.go`:** prompt contains the diff text, option count instruction is present
- **Integration — `git/git.go`:** uses a temp git repo with staged changes to verify diff capture
- **AI output:** no live AI calls in tests; fixture files provide canned responses

## Out of Scope

- Unstaged changes (always requires explicit `git add`)
- Amending existing commits
- Signing commits (`-S`)
- Any GUI or web interface
