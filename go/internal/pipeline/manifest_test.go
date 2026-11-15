package pipeline

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jdonohoo/vern-bot/go/internal/llm"
	"github.com/jdonohoo/vern-bot/go/internal/packs"
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

const fakeDeliverables = `# Story deliverables

### DELIVERABLE 1: Outline
**Type:** outline
**Description:** Three acts.
**Acceptance Criteria:**
- each act turns
**Dependencies:** None

### DELIVERABLE 2: Beat sheet
**Type:** beats
**Description:** Beats per act.
**Acceptance Criteria:**
- 15 beats
**Dependencies:** Deliverable 1
`

// A creative pack pipeline yields deliverables with types in tasks.json.
func TestCreativePipelineWritesDeliverables(t *testing.T) {
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(llm.Register(&llm.Provider{
		Name: "fake", Binary: exe, Label: "Fake", ReadOnly: true,
		Preamble: func(bool) string { return "" },
		Args:     func(llm.Invocation) []string { return nil },
	}))
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Setenv("VERN_LOG", "0")
	t.Setenv("VERN_PATH", "")
	t.Setenv("VERN_FAKE_CLI", fakeDeliverables)

	pack := filepath.Join(t.TempDir(), "story-pack")
	os.MkdirAll(filepath.Join(pack, "pipelines"), 0755)
	os.WriteFile(filepath.Join(pack, "pipelines", "tale.json"), []byte(`{
	  "name": "tale", "kind": "creative", "steps": [
	    {"step": 1, "name": "Premise", "persona": "mighty", "llm": "fake", "context_mode": "prompt_only", "prompt_prefix": "Premise."},
	    {"step": 2, "name": "Bible", "persona": "mighty", "llm": "fake", "context_mode": "all_previous", "prompt_prefix": "Consolidate."},
	    {"step": 3, "name": "Deliverables", "persona": "mighty", "llm": "fake", "context_mode": "consolidation", "prompt_prefix": "Deliverables."}
	  ]}`), 0644)
	if _, _, err := packs.Setup([]string{pack}, t.TempDir(), home); err != nil {
		t.Fatal(err)
	}
	defer packs.Reset()

	dir := filepath.Join(t.TempDir(), "tale")
	m, err := RunWithManifest(Options{
		Idea: "a lighthouse keeper's last night", DiscoveryDir: dir, BatchMode: true,
		SkipHistorian: true, Pipeline: "tale", OnLog: func(string) {},
	})
	if err != nil {
		t.Fatal(err)
	}
	if m.PipelineKind != "creative" || m.Status != "ok" {
		t.Fatalf("manifest = %+v", m)
	}
	doc, err := tasks.Read(filepath.Join(dir, "output", tasks.FileName))
	if err != nil || len(doc.Tasks) != 2 {
		t.Fatalf("tasks.json = %+v, %v", doc, err)
	}
	d2 := doc.Tasks[1]
	if d2.Kind != "deliverable" || d2.Type != "beats" || len(d2.Dependencies) != 1 {
		t.Errorf("deliverable 2 = %+v", d2)
	}
}

// A rerun whose breakdown has no tasks must not report the last run's.
func TestRerunWithoutTasksIgnoresStaleTasks(t *testing.T) {
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(llm.Register(&llm.Provider{
		Name: "fake", Binary: exe, Label: "Fake", ReadOnly: true,
		Preamble: func(bool) string { return "" },
		Args:     func(llm.Invocation) []string { return nil },
	}))
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Setenv("VERN_LOG", "0")
	cfgDir := filepath.Join(home, ".config", "vern")
	os.MkdirAll(cfgDir, 0755)
	os.WriteFile(filepath.Join(cfgDir, "config.json"), []byte(`{
	  "discovery_pipelines": {"default": [
	    {"step": 1, "name": "Architect Breakdown", "persona": "architect", "llm": "fake", "context_mode": "prompt_only", "prompt_prefix": "Break it down."}
	  ]}
	}`), 0644)
	dir := filepath.Join(t.TempDir(), "rerun")
	opts := Options{Idea: "x", DiscoveryDir: dir, BatchMode: true, SkipHistorian: true, OnLog: func(string) {}}

	t.Setenv("VERN_FAKE_CLI", fakeBreakdown)
	if m, err := RunWithManifest(opts); err != nil || m.Tasks == nil {
		t.Fatalf("first run: %+v, %v", m, err)
	}
	t.Setenv("VERN_FAKE_CLI", "No tasks here, just an essay about the idea.")
	m, err := RunWithManifest(opts)
	if err != nil {
		t.Fatal(err)
	}
	if m.Tasks != nil || m.Status == "ok" {
		t.Errorf("rerun reported stale tasks: status %s, tasks %+v", m.Status, m.Tasks)
	}
}
