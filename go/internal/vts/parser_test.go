package vts

import (
	"os"
	"path/filepath"
	"testing"
)

const sampleArchitectOutput = `# Architect Breakdown

Here is the task breakdown for the project.

### TASK 1: Set Up Project Structure

**Description:** Initialize the project with proper directory layout and build configuration.
**Acceptance Criteria:**
- Go module initialized
- Directory structure created
- CI configuration added
**Complexity:** M
**Dependencies:** None
**Files:** go.mod, cmd/main.go, internal/

---

### TASK 2: Implement Config Loader

**Description:** Load configuration from JSON file with 3-tier fallback chain.
**Acceptance Criteria:**
- User config loaded first
- Falls back to default config
- Hardcoded defaults as last resort
**Complexity:** S
**Dependencies:** Task 1
**Files:** internal/config/config.go, config.default.json

---

### TASK 3: Build CLI Interface

**Description:** Create cobra-based CLI with subcommands.
**Acceptance Criteria:**
- Run subcommand works
- Discovery subcommand works
- Hole subcommand works
**Complexity:** L
**Dependencies:** Task 1, Task 2
**Files:** cmd/vern/main.go, cmd/vern/run.go

## Dependency Graph

Task 1 -> Task 2 -> Task 3

## Total Estimate

Medium effort, ~2 weeks.
`

func TestParseArchitectOutput(t *testing.T) {
	tasks, header, footer := ParseArchitectOutput(sampleArchitectOutput)

	if len(tasks) != 3 {
		t.Fatalf("expected 3 tasks, got %d", len(tasks))
	}

	// Check header
	if header == "" {
		t.Error("header should not be empty")
	}

	// Check footer contains dependency graph
	if footer == "" {
		t.Error("footer should not be empty")
	}

	// Task 1
	t1 := tasks[0]
	if t1.Num != 1 {
		t.Errorf("task 1 num: got %d", t1.Num)
	}
	if t1.Title != "Set Up Project Structure" {
		t.Errorf("task 1 title: got %q", t1.Title)
	}
	if t1.Complexity != "M" {
		t.Errorf("task 1 complexity: got %q", t1.Complexity)
	}
	if len(t1.Dependencies) != 0 {
		t.Errorf("task 1 deps: got %v", t1.Dependencies)
	}
	if len(t1.Criteria) != 3 {
		t.Errorf("task 1 criteria: got %d, want 3", len(t1.Criteria))
	}

	// Task 2 - has dependency on Task 1
	t2 := tasks[1]
	if len(t2.Dependencies) != 1 || t2.Dependencies[0] != "VTS-001" {
		t.Errorf("task 2 deps: got %v, want [VTS-001]", t2.Dependencies)
	}
	if t2.Complexity != "S" {
		t.Errorf("task 2 complexity: got %q", t2.Complexity)
	}

	// Task 3 - has dependencies on Task 1 and 2
	t3 := tasks[2]
	if len(t3.Dependencies) != 2 {
		t.Errorf("task 3 deps: got %v, want [VTS-001, VTS-002]", t3.Dependencies)
	}
	if t3.Complexity != "L" {
		t.Errorf("task 3 complexity: got %q", t3.Complexity)
	}
}

func TestParseNoTasks(t *testing.T) {
	tasks, _, _ := ParseArchitectOutput("Just some text with no tasks.")
	if len(tasks) != 0 {
		t.Errorf("expected 0 tasks, got %d", len(tasks))
	}
}

func TestWriteVTSFiles(t *testing.T) {
	dir := t.TempDir()
	tasks := []Task{
		{
			Num:         1,
			Title:       "Test Task",
			Description: "A test task",
			Complexity:  "S",
			Criteria:    []string{"It works", "Tests pass"},
			Files:       []string{"main.go"},
		},
	}

	err := WriteVTSFiles(tasks, dir, "discovery", "architect.md")
	if err != nil {
		t.Fatal(err)
	}

	// Check file was created
	files, _ := filepath.Glob(filepath.Join(dir, "vts-*.md"))
	if len(files) != 1 {
		t.Fatalf("expected 1 VTS file, got %d", len(files))
	}

	// Check content
	content, _ := os.ReadFile(files[0])
	s := string(content)
	if !contains(s, "id: VTS-001") {
		t.Error("missing VTS ID")
	}
	if !contains(s, `title: "Test Task"`) {
		t.Error("missing title")
	}
	if !contains(s, "complexity: S") {
		t.Error("missing complexity")
	}
	if !contains(s, "source: discovery") {
		t.Error("missing source")
	}
}

func TestWriteSummary(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "architect.md")

	tasks := []Task{
		{Num: 1, Title: "First", Complexity: "S"},
		{Num: 2, Title: "Second", Complexity: "M", Dependencies: []string{"VTS-001"}},
	}

	err := WriteSummary(tasks, file, "# Header", "## Footer", "")
	if err != nil {
		t.Fatal(err)
	}

	content, _ := os.ReadFile(file)
	s := string(content)
	if !contains(s, "VTS-001") {
		t.Error("missing VTS-001 in summary")
	}
	if !contains(s, "VTS-002") {
		t.Error("missing VTS-002 in summary")
	}
	if !contains(s, "## VTS Task Index") {
		t.Error("missing task index header")
	}
}

func TestSlugify(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Set Up Project Structure", "set-up-project-structure"},
		{"Implement Config Loader", "implement-config-loader"},
		{"Hello  World!!!", "hello-world"},
	}
	for _, tt := range tests {
		got := Slugify(tt.input)
		if got != tt.want {
			t.Errorf("Slugify(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsStr(s, substr))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
