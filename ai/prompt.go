package ai

import "fmt"

func BuildPrompt(diff string) string {
	return fmt.Sprintf(`You are an expert at writing git commit messages.

Analyze this staged git diff and generate commit message options.

Rules:
1. Use Conventional Commits format: type(scope): subject
   Valid types: feat, fix, refactor, chore, docs, test, perf, style
2. Scale the number of options to diff complexity:
   - Trivial (single line, typo fix): 1 option
   - Moderate (one feature or bug fix): 2-3 options
   - Complex (multiple files or subsystems): 3-4 options
3. Include a body only when the change warrants explanation (why, not what)
4. Keep subjects under 72 characters

Respond ONLY with a JSON array. No markdown, no explanation, no code fences.

Format:
[
  {"subject": "type(scope): concise description", "body": ""},
  {"subject": "type: alternative description", "body": "Optional multi-line\nbody here"}
]

Git diff:
%s`, diff)
}
