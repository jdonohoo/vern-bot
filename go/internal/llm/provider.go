package llm

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// Provider describes one LLM CLI that vern can drive in headless mode.
//
// Adding a CLI means adding one entry to the registry below; nothing else in
// vern needs to know its flags.
type Provider struct {
	// Name is the canonical name used in config, logs and output.
	Name string
	// Binary is the executable looked up on PATH.
	Binary string
	// Label is the human-readable name shown in the TUI and setup.
	Label string
	// Aliases are other names that resolve to this provider (short forms and
	// retired CLIs whose name still appears in old configs).
	Aliases []string
	// WritesOutputFile is true when the CLI writes its final answer to a file
	// given in Args (codex -o) instead of stdout.
	WritesOutputFile bool
	// Preamble is prepended to every prompt sent to this provider.
	Preamble func(allowFileRead bool) string
	// PromptAsArg is true for CLIs that cannot read a text prompt from stdin
	// (agy). Their prompt goes in Invocation.Prompt; when it is too big for
	// the OS command line it is written to a file and the CLI is told to
	// read it (Invocation.PromptFile). Every other CLI gets the prompt on
	// stdin, so its size is not limited by the OS command-line limits.
	PromptAsArg bool
	// Args builds the argument list.
	Args func(inv Invocation) []string
	// Env is extra environment for the subprocess.
	Env []string
}

// Invocation carries the per-run values a provider needs to build its args.
type Invocation struct {
	Prompt     string // only set when PromptAsArg is true
	PromptFile string // set with Prompt when the full prompt is in a file
	WorkingDir string
	OutputPath string // only set when WritesOutputFile is true
	Timeout    time.Duration
}

const textOnlyPreamble = "IMPORTANT: Output your complete analysis as plain text to stdout. Do NOT create, write, or modify any files. Do NOT use any file-writing tools. Just output your analysis directly as text.\n\n"
const fileReadPreamble = "IMPORTANT: You MAY read files from the filesystem to gather information. Output your complete analysis as plain text to stdout. Do NOT create, write, or modify any files.\n\n"

func textPreamble(allowFileRead bool) string {
	if allowFileRead {
		return fileReadPreamble
	}
	return textOnlyPreamble
}

const codexPreamble = "IMPORTANT: You are acting as a PLANNING and ANALYSIS agent for a discovery pipeline. Write your complete analysis, implementation plan, and recommendations as a detailed markdown document. Do NOT create any code files, project scaffolding, or application code. Do NOT build anything. Your entire output should be a thorough written analysis — problem space, architecture, risks, recommendations — not a built project.\n\n"

var providers = []*Provider{
	{
		Name:     "claude",
		Binary:   "claude",
		Label:    "Claude",
		Aliases:  []string{"c"},
		Preamble: textPreamble,
		// claude -p with no prompt argument reads the prompt from stdin.
		Args: func(inv Invocation) []string {
			return []string{"--dangerously-skip-permissions", "-p"}
		},
		Env: []string{"NODE_OPTIONS=--max-old-space-size=32768"},
	},
	{
		Name:             "codex",
		Binary:           "codex",
		Label:            "Codex",
		Aliases:          []string{"x"},
		WritesOutputFile: true,
		Preamble:         func(bool) string { return codexPreamble },
		// codex exec reads the prompt from stdin when it is "-".
		Args: func(inv Invocation) []string {
			dir := inv.WorkingDir
			if dir == "" {
				dir = "."
			}
			return []string{"exec",
				"--dangerously-bypass-approvals-and-sandbox",
				"--skip-git-repo-check",
				"--cd", dir,
				"-o", inv.OutputPath,
				"-",
			}
		},
	},
	{
		// Google's Antigravity CLI replaced the Gemini CLI in 2026. "gemini"
		// stays as an alias so old configs and pipelines keep working, but it
		// always runs agy.
		Name:        "antigravity",
		Binary:      "agy",
		Label:       "Antigravity",
		Aliases:     []string{"agy", "a", "gemini", "g"},
		Preamble:    textPreamble,
		PromptAsArg: true,
		// agy 1.4 reads a text prompt only from -p (stdin is ignored), and
		// -p must be attached (-p=...) so a prompt starting with "-" is not
		// taken for a flag.
		Args: func(inv Invocation) []string {
			args := []string{"-p=" + inv.Prompt, "--dangerously-skip-permissions"}
			if inv.PromptFile != "" {
				args = append(args, "--add-dir", filepath.Dir(inv.PromptFile))
			}
			if inv.Timeout > 0 {
				args = append(args, "--print-timeout", fmt.Sprintf("%ds", int(inv.Timeout.Seconds())))
			}
			return args
		},
	},
}

// retired names map to a reason shown when someone still asks for them.
var retired = map[string]string{
	"copilot": "Copilot support was removed in v3",
	"p":       "Copilot support was removed in v3",
}

// Providers returns the registered providers in display order.
func Providers() []*Provider {
	out := make([]*Provider, len(providers))
	copy(out, providers)
	return out
}

// Names returns the canonical provider names in display order.
func Names() []string {
	names := make([]string, 0, len(providers))
	for _, p := range providers {
		names = append(names, p.Name)
	}
	return names
}

// Lookup returns the provider for a name or alias (case-insensitive).
func Lookup(name string) (*Provider, bool) {
	n := strings.ToLower(strings.TrimSpace(name))
	for _, p := range providers {
		if p.Name == n {
			return p, true
		}
		for _, a := range p.Aliases {
			if a == n {
				return p, true
			}
		}
	}
	return nil, false
}

// Canonical returns the canonical provider name for a name or alias, or ""
// when the name is unknown or retired.
func Canonical(name string) string {
	if p, ok := Lookup(name); ok {
		return p.Name
	}
	return ""
}

// IsRetired reports whether name is a provider vern no longer supports, and why.
func IsRetired(name string) (string, bool) {
	reason, ok := retired[strings.ToLower(strings.TrimSpace(name))]
	return reason, ok
}

// Available reports whether the provider's CLI is on PATH.
func (p *Provider) Available() bool {
	_, err := exec.LookPath(p.Binary)
	return err == nil
}

// Detect returns availability for every registered provider, keyed by
// canonical name.
func Detect() map[string]bool {
	out := make(map[string]bool, len(providers))
	for _, p := range providers {
		out[p.Name] = p.Available()
	}
	return out
}

// resolveLLM normalizes an LLM name to a canonical provider, falling back to
// claude (with a warning) when the CLI is missing, retired or unknown.
func resolveLLM(name string) string {
	p, ok := Lookup(name)
	if !ok {
		if reason, gone := IsRetired(name); gone {
			fmt.Fprintf(os.Stderr, "[vern-run] Warning: %s; using claude\n", reason)
			return "claude"
		}
		// Unknown names pass through so Run can report them.
		return strings.ToLower(name)
	}
	if p.Name == "claude" {
		return p.Name
	}
	if !p.Available() {
		fmt.Fprintf(os.Stderr, "[vern-run] Warning: %s CLI (%s) not found, falling back to claude\n", p.Name, p.Binary)
		return "claude"
	}
	return p.Name
}

// maxArgPrompt is the largest prompt passed as a single command-line
// argument. Linux caps one argument at 128 KB and Windows caps the whole
// command line at 32 K characters; leave room for the other arguments.
func maxArgPrompt() int {
	if runtime.GOOS == "windows" {
		return 24 * 1024
	}
	return 100 * 1024
}

// promptFilePointer is the short prompt sent instead of a prompt too big for
// the command line.
func promptFilePointer(path string) string {
	return "Your complete instructions are in the file " + path +
		". Read that entire file now, then follow the instructions in it exactly. " +
		"Do not modify or delete that file. Do not mention the file in your answer."
}
