package llm

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestLookupAliases(t *testing.T) {
	tests := map[string]string{
		"claude": "claude", "c": "claude", "Claude": "claude",
		"codex": "codex", "x": "codex",
		"antigravity": "antigravity", "agy": "antigravity", "a": "antigravity",
		// The Gemini CLI was replaced by Antigravity: the old name runs agy.
		"gemini": "antigravity", "g": "antigravity",
	}
	for in, want := range tests {
		if got := Canonical(in); got != want {
			t.Errorf("Canonical(%q) = %q, want %q", in, got, want)
		}
	}
	if p, _ := Lookup("gemini"); p.Binary != "agy" {
		t.Errorf("gemini should run the agy binary, got %q", p.Binary)
	}
}

func TestCopilotRetired(t *testing.T) {
	if Canonical("copilot") != "" {
		t.Error("copilot should no longer be a provider")
	}
	if _, gone := IsRetired("copilot"); !gone {
		t.Error("copilot should be reported as retired")
	}
	if got := resolveLLM("copilot"); got != "claude" {
		t.Errorf("resolveLLM(copilot) = %q, want claude", got)
	}
}

func TestNamesOrder(t *testing.T) {
	got := strings.Join(Names(), ",")
	if got != "claude,codex,antigravity" {
		t.Errorf("Names() = %s", got)
	}
}

func TestAntigravityArgs(t *testing.T) {
	p, _ := Lookup("antigravity")
	if !p.PromptAsArg {
		t.Fatal("agy only reads a text prompt from -p")
	}
	args := p.Args(Invocation{Prompt: "-starts with a dash", Timeout: 20 * time.Minute})
	if args[0] != "-p=-starts with a dash" {
		t.Errorf("prompt must be attached to -p, got %q", args[0])
	}
	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "--print-timeout 1200s") {
		t.Errorf("print timeout should match vern's timeout: %s", joined)
	}
	if strings.Contains(joined, "--add-dir") {
		t.Errorf("no --add-dir without a prompt file: %s", joined)
	}
	args = p.Args(Invocation{Prompt: "read it", PromptFile: "/tmp/vern-prompt-1/prompt.md"})
	if !strings.Contains(strings.Join(args, " "), "--add-dir /tmp/vern-prompt-1") {
		t.Errorf("prompt file dir should be readable: %v", args)
	}
}

func TestCodexReadsPromptFromStdin(t *testing.T) {
	p, _ := Lookup("codex")
	args := p.Args(Invocation{OutputPath: "/tmp/out.md"})
	if args[len(args)-1] != "-" {
		t.Errorf("codex prompt should be read from stdin (-), got %v", args)
	}
}

// fakeCLI puts a shell script named name on PATH for the test.
func fakeCLI(t *testing.T, name, script string) {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("fake CLIs are shell scripts")
	}
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, name), []byte("#!/bin/sh\n"+script), 0755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("VERN_LOG", "0")
}

func TestRunSendsLargePromptOnStdin(t *testing.T) {
	fakeCLI(t, "claude", "cat\n")
	// Bigger than Linux's 128 KB single-argument limit.
	big := strings.Repeat("x", 300*1024) + "END-MARKER"
	res, err := Run(RunOptions{LLM: "claude", Prompt: big, Timeout: 30 * time.Second, QuietStderr: true})
	if err != nil {
		t.Fatal(err)
	}
	if res.ExitCode != 0 || !strings.Contains(res.Output, "END-MARKER") {
		t.Errorf("prompt did not arrive on stdin (exit %d, %d bytes)", res.ExitCode, len(res.Output))
	}
}

func TestRunAntigravityViaGeminiAlias(t *testing.T) {
	// Echo the -p= argument back.
	fakeCLI(t, "agy", `for a in "$@"; do case "$a" in -p=*) printf '%s' "${a#-p=}";; esac; done`+"\n")
	res, err := Run(RunOptions{LLM: "gemini", Prompt: "hello agy", Timeout: 30 * time.Second, QuietStderr: true})
	if err != nil {
		t.Fatal(err)
	}
	if res.LLMUsed != "antigravity" {
		t.Errorf("LLMUsed = %q, want antigravity", res.LLMUsed)
	}
	if !strings.Contains(res.Output, "hello agy") {
		t.Errorf("prompt not passed to agy: %q", res.Output)
	}
}

func TestRunAntigravityLargePromptUsesFile(t *testing.T) {
	// Print the pointer prompt, then the file it points at.
	fakeCLI(t, "agy", `for a in "$@"; do case "$a" in -p=*) p="${a#-p=}";; esac; done
f=$(printf '%s' "$p" | sed -n 's/.*in the file \([^ ]*prompt\.md\).*/\1/p')
printf 'POINTER:%s\n' "$p" | head -c 200; echo; cat "$f"
`)
	big := strings.Repeat("y", maxArgPrompt()+1024) + "FILE-MARKER"
	res, err := Run(RunOptions{LLM: "antigravity", Prompt: big, Timeout: 30 * time.Second, QuietStderr: true})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(res.Output, "POINTER:Your complete instructions are in the file") {
		t.Errorf("large prompt should be replaced by a file pointer, got %.120q", res.Output)
	}
	if !strings.Contains(res.Output, "FILE-MARKER") {
		t.Error("prompt file did not hold the full prompt")
	}
}
