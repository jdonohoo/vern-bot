package llm

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The test binary doubles as a fake LLM CLI: with VERN_FAKE_CLI set it
// prints its arguments (or the env value) and exits.
func TestMain(m *testing.M) {
	switch os.Getenv("VERN_FAKE_CLI") {
	case "":
		os.Exit(m.Run())
	case "args":
		fmt.Print(strings.Join(os.Args[1:], " "))
		os.Exit(0)
	default:
		fmt.Print(os.Getenv("VERN_FAKE_CLI"))
		os.Exit(0)
	}
}

// withFakeProvider registers a provider that runs this test binary.
func withFakeProvider(t *testing.T, p *Provider) {
	t.Helper()
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	p.Binary = exe
	if p.Preamble == nil {
		p.Preamble = func(bool) string { return "" }
	}
	t.Cleanup(Register(p))
}

func fakeArgs(inv Invocation) []string {
	args := []string{"run"}
	if inv.ReadOnly {
		args = append(args, "readonly")
	}
	return modelEffortArgs(args, inv)
}

func TestRunReadOnlyByDefault(t *testing.T) {
	t.Setenv("VERN_LOG", "0")
	t.Setenv("VERN_FAKE_CLI", "args")
	t.Setenv("VERN_ALLOW_WRITES", "")
	SetPolicy(Policy{})
	withFakeProvider(t, &Provider{Name: "fake", ReadOnly: true, Args: fakeArgs})

	res, err := Run(RunOptions{LLM: "fake", Prompt: "hi"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(res.Output, "readonly") || !res.ReadOnly {
		t.Errorf("default run should be read-only, args %q", res.Output)
	}

	res, _ = Run(RunOptions{LLM: "fake", Prompt: "hi", AllowWrites: true})
	if strings.Contains(res.Output, "readonly") {
		t.Errorf("AllowWrites run should not be read-only, args %q", res.Output)
	}

	t.Setenv("VERN_ALLOW_WRITES", "1")
	res, _ = Run(RunOptions{LLM: "fake", Prompt: "hi"})
	if strings.Contains(res.Output, "readonly") {
		t.Errorf("VERN_ALLOW_WRITES=1 should allow writes, args %q", res.Output)
	}
	t.Setenv("VERN_ALLOW_WRITES", "")

	SetPolicy(Policy{AllowWrites: true})
	defer SetPolicy(Policy{})
	res, _ = Run(RunOptions{LLM: "fake", Prompt: "hi"})
	if strings.Contains(res.Output, "readonly") {
		t.Errorf("policy AllowWrites should allow writes, args %q", res.Output)
	}
}

func TestRunModelAndEffort(t *testing.T) {
	t.Setenv("VERN_LOG", "0")
	t.Setenv("VERN_FAKE_CLI", "args")
	SetPolicy(Policy{Tiers: map[string]Tier{
		TierDeep: {Models: map[string]string{"fake": "big-model"}, Effort: "high"},
	}})
	defer SetPolicy(Policy{})
	withFakeProvider(t, &Provider{Name: "fake", ReadOnly: true, Args: fakeArgs})

	res, err := Run(RunOptions{LLM: "fake", Prompt: "hi", Tier: "deep"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(res.Output, "--model big-model --effort high") {
		t.Errorf("tier should set model and effort, args %q", res.Output)
	}

	res, _ = Run(RunOptions{LLM: "fake", Prompt: "hi", Tier: "deep",
		Models: ModelsFor("fake", "step-model"), Effort: "low"})
	if !strings.Contains(res.Output, "--model step-model --effort low") || res.Model != "step-model" {
		t.Errorf("explicit model and effort should win, args %q", res.Output)
	}

	// A model meant for another provider never reaches this one.
	res, _ = Run(RunOptions{LLM: "fake", Prompt: "hi", Models: ModelsFor("codex", "gpt-x")})
	if strings.Contains(res.Output, "gpt-x") {
		t.Errorf("codex model leaked to fake provider, args %q", res.Output)
	}

	if _, err := Run(RunOptions{LLM: "fake", Prompt: "hi", Effort: "ludicrous"}); err == nil {
		t.Error("unknown effort should be an error")
	}
}

func TestRunFailurePrefix(t *testing.T) {
	t.Setenv("VERN_LOG", "0")
	t.Setenv("VERN_FAKE_CLI", "jetski: no output produced — a tool was denied")
	withFakeProvider(t, &Provider{Name: "fake", Args: fakeArgs,
		FailurePrefixes: []string{"jetski: no output produced"}})

	out := filepath.Join(t.TempDir(), "out.md")
	res, err := Run(RunOptions{LLM: "fake", Prompt: "hi", OutputFile: out, QuietStderr: true})
	if err != nil {
		t.Fatal(err)
	}
	if res.ExitCode == 0 || res.Output != "" {
		t.Errorf("failure text with exit 0 should fail the run: exit %d, output %q", res.ExitCode, res.Output)
	}
	if !strings.Contains(res.Stderr, "no output produced") {
		t.Errorf("failure text should move to stderr, got %q", res.Stderr)
	}
}

func TestPersonaTier(t *testing.T) {
	dir := t.TempDir()
	write := func(name, fm string) {
		os.WriteFile(filepath.Join(dir, name+".md"), []byte("---\nname: "+name+"\n"+fm+"\n---\nbody\n"), 0644)
	}
	write("deepone", "model: opus")
	write("mid", "model: sonnet")
	write("quick", "model: haiku")
	write("explicit", "model: haiku\ntier: deep")
	write("none", "model: antigravity")
	want := map[string]string{"deepone": "deep", "mid": "balanced", "quick": "fast", "explicit": "deep", "none": "", "": ""}
	for persona, tier := range want {
		if got := PersonaTier(dir, persona); got != tier {
			t.Errorf("PersonaTier(%q) = %q, want %q", persona, got, tier)
		}
	}
	// Falls back to the built-in personas.
	if got := PersonaTier(t.TempDir(), "architect"); got != TierDeep {
		t.Errorf("embedded architect tier = %q, want deep", got)
	}
}

func TestProviderArgsReadOnlyModelEffort(t *testing.T) {
	inv := Invocation{ReadOnly: true, Model: "m1", Effort: "high", OutputPath: "out.md"}
	cases := map[string][]string{
		"claude":      {"--output-format json", "--permission-mode dontAsk", "--model m1", "--effort high"},
		"codex":       {"--sandbox read-only", "--model m1", "-c model_reasoning_effort=high"},
		"antigravity": {"--mode plan --sandbox --settings", `"command(*)"`, "--model m1", "--effort high"},
		"grok":        {"--model m1", "--effort high", "--always-approve"},
	}
	for name, wants := range cases {
		p, _ := Lookup(name)
		joined := strings.Join(p.Args(inv), " ")
		for _, w := range wants {
			if !strings.Contains(joined, w) {
				t.Errorf("%s read-only args missing %q: %s", name, w, joined)
			}
		}
		if p.ReadOnly && strings.Contains(joined, "dangerously") {
			t.Errorf("%s read-only args still bypass permissions: %s", name, joined)
		}
		write := strings.Join(p.Args(Invocation{OutputPath: "out.md"}), " ")
		if p.ReadOnly && strings.Contains(write, "read-only") {
			t.Errorf("%s write args should not be read-only: %s", name, write)
		}
		if strings.Contains(write, "--model") || strings.Contains(write, "effort") {
			t.Errorf("%s should not pass model or effort when unset: %s", name, write)
		}
	}
}

func TestMissingFlags(t *testing.T) {
	help := `Usage: tool [OPTIONS]
  -p, --print           print and exit
  --model <model>       pick one
  --effort=<level>
  --mode accept-edits|plan`
	got := missingFlags(help, []string{"--print", "--model", "--effort", "--mode", "--sandbox", "--mod"})
	if strings.Join(got, ",") != "--sandbox,--mod" {
		t.Errorf("missingFlags = %v", got)
	}
}

func TestParseModelIDs(t *testing.T) {
	codex := `{"models":[{"slug":"gpt-6-astra","visibility":"list"},{"slug":"gpt-6-luna"}]}`
	if got := strings.Join(parseModelIDs(codex), ","); got != "gpt-6-astra,gpt-6-luna" {
		t.Errorf("codex models = %s", got)
	}
	agy := "Fetching available models...\ngemini-4-argon-high | Gemini 4 Argon (High)\nclaude-opus-5-5-low   Claude Opus 5.5\n"
	if got := strings.Join(parseModelIDs(agy), ","); got != "gemini-4-argon-high,claude-opus-5-5-low" {
		t.Errorf("agy models = %s", got)
	}
}

func TestSignedOut(t *testing.T) {
	if signedOut(`{"loggedIn": true, "authMethod": "claude.ai"}`) {
		t.Error("loggedIn true is signed in")
	}
	if !signedOut(`{"loggedIn": false}`) {
		t.Error("loggedIn false is signed out")
	}
	if !signedOut("Not logged in") {
		t.Error("'Not logged in' is signed out")
	}
	if signedOut("Logged in using ChatGPT") {
		t.Error("'Logged in using ChatGPT' is signed in")
	}
}

func TestParseClaudeJSON(t *testing.T) {
	out := "Warning: this sign-in can end at any time\n" +
		`{"type":"result","subtype":"success","is_error":false,"result":"banana-42","total_cost_usd":0.0123,"usage":{"input_tokens":10,"cache_read_input_tokens":90,"output_tokens":5}}` + "\n"
	text, u, failed := parseClaudeJSON(out, "")
	if text != "banana-42" || failed {
		t.Fatalf("text %q failed %v", text, failed)
	}
	if u.InputTokens != 100 || u.OutputTokens != 5 || u.TotalTokens != 105 || u.CostUSD != 0.0123 {
		t.Errorf("usage = %+v", u)
	}
	if _, _, failed := parseClaudeJSON(`{"type":"result","is_error":true,"result":"API error"}`, ""); !failed {
		t.Error("is_error should fail the run")
	}
	// Anything else passes through as text.
	if text, _, _ := parseClaudeJSON("plain answer {not json}", ""); text != "plain answer {not json}" {
		t.Errorf("plain text changed: %q", text)
	}
}

func TestParseCodexUsage(t *testing.T) {
	_, u, _ := parseCodexUsage("answer", "session id: x\ntokens used\n13,477\n")
	if u.TotalTokens != 13477 {
		t.Errorf("usage = %+v", u)
	}
	_, u, _ = parseCodexUsage("answer", "tokens used 7,501")
	if u.TotalTokens != 7501 {
		t.Errorf("usage = %+v", u)
	}
}
