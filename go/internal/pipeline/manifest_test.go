package pipeline

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jdonohoo/vern-bot/go/internal/llm"
	"github.com/jdonohoo/vern-bot/go/internal/tasks"
)

// With VERN_FAKE_CLI set, the test binary acts as an LLM CLI that prints it.
func TestMain(m *testing.M) {
	if out := os.Getenv("VERN_FAKE_CLI"); out != "" {
		fmt.Print(out)
		os.Exit(0)
	}
	os.Exit(m.Run())
}

const fakeBreakdown = `# Plan

### TASK 1: Scaffold

**Description:** Set it up.
**Acceptance Criteria:**
- builds
**Complexity:** S
**Dependencies:** None
**Files:** main.go

### TASK 2: Render notes

**Description:** Turn markdown into HTML.
**Acceptance Criteria:**
- renders headings
**Complexity:** M
**Dependencies:** Task 1
**Files:** render.go
`

func TestDiscoveryWritesManifestAndTasks(t *testing.T) {
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(llm.Register(&llm.Provider{
		Name: "fake", Binary: exe, Label: "Fake", ReadOnly: true,
		Preamble: func(bool) string { return "" },
		Args:     func(llm.Invocation) []string { return nil },
		ParseOutput: func(stdout, _ string) (string, llm.Usage, bool) {
			return stdout, llm.Usage{TotalTokens: 10, CostUSD: 0.5}, false
		},
	}))

	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Setenv("VERN_LOG", "0")
	t.Setenv("VERN_FAKE_CLI", fakeBreakdown)
	cfgDir := filepath.Join(home, ".config", "vern")
	os.MkdirAll(cfgDir, 0755)
	os.WriteFile(filepath.Join(cfgDir, "config.json"), []byte(`{
	  "discovery_pipelines": {"default": [
	    {"step": 1, "name": "Analysis", "persona": "mighty", "llm": "fake", "context_mode": "prompt_only", "prompt_prefix": "Analyze.", "model": "m-big", "effort": "high"},
	    {"step": 2, "name": "Architect Breakdown", "persona": "architect", "llm": "fake", "context_mode": "previous", "prompt_prefix": "Break it down."}
	  ]}
	}`), 0644)

	dir := filepath.Join(t.TempDir(), "notes-site")
	m, err := RunWithManifest(Options{
		Idea: "notes to a static site", DiscoveryDir: dir, BatchMode: true,
		SkipHistorian: true, OnLog: func(string) {},
	})
	if err != nil {
		t.Fatal(err)
	}
	if m == nil || m.Status != "ok" || len(m.Steps) != 2 {
		t.Fatalf("manifest = %+v", m)
	}
	s1 := m.Steps[0]
	if s1.LLMUsed != "fake" || s1.Model != "m-big" || s1.Effort != "high" || !s1.ReadOnly || s1.Persona != "mighty" {
		t.Errorf("step 1 = %+v", s1)
	}
	if !strings.HasPrefix(s1.OutputFile, "output/01-") {
		t.Errorf("step output should be relative: %s", s1.OutputFile)
	}
	if m.Usage == nil || m.Usage.TotalTokens != 20 || m.Usage.CostUSD != 1.0 {
		t.Errorf("usage = %+v", m.Usage)
	}
	if m.Tasks == nil || m.Tasks.Count != 2 {
		t.Fatalf("tasks = %+v", m.Tasks)
	}

	onDisk, err := ReadManifest(dir)
	if err != nil || onDisk.Status != "ok" || onDisk.FinishedAt == "" {
		t.Fatalf("run.json = %+v, %v", onDisk, err)
	}
	doc, err := tasks.Read(filepath.Join(dir, "output", tasks.FileName))
	if err != nil || len(doc.Tasks) != 2 || doc.Idea != "notes to a static site" {
		t.Fatalf("tasks.json = %+v, %v", doc, err)
	}
}
