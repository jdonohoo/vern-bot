package vts

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestParseVTSFile_Valid(t *testing.T) {
	content := `---
id: VTS-001
title: "Freeze Field Mapping"
complexity: XS
status: pending
owner: ""
source: oracle
source_ref: "oracle-breakdown.md"
dependencies:
  - VTS-008
files:
  - "mapping.yaml"
  - "schema.rs"
---

# Freeze Field Mapping

Create the mapping specification.

## Criteria

- Mapping file exists
- Status allowlist covers all cases
`
	task, err := ParseVTSFile(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if task.ID != "VTS-001" {
		t.Errorf("ID = %q, want VTS-001", task.ID)
	}
	if task.Num != 1 {
		t.Errorf("Num = %d, want 1", task.Num)
	}
	if task.Title != "Freeze Field Mapping" {
		t.Errorf("Title = %q, want 'Freeze Field Mapping'", task.Title)
	}
	if task.Complexity != "XS" {
		t.Errorf("Complexity = %q, want XS", task.Complexity)
	}
	if task.Status != "pending" {
		t.Errorf("Status = %q, want pending", task.Status)
	}
	if task.Owner != "" {
		t.Errorf("Owner = %q, want empty", task.Owner)
	}
	if task.Source != "oracle" {
		t.Errorf("Source = %q, want oracle", task.Source)
	}
	if task.SourceRef != "oracle-breakdown.md" {
		t.Errorf("SourceRef = %q, want oracle-breakdown.md", task.SourceRef)
	}
	if len(task.Dependencies) != 1 || task.Dependencies[0] != "VTS-008" {
		t.Errorf("Dependencies = %v, want [VTS-008]", task.Dependencies)
	}
	if len(task.Files) != 2 {
		t.Errorf("Files = %v, want 2 files", task.Files)
	}
	if task.Description != "Create the mapping specification." {
		t.Errorf("Description = %q", task.Description)
	}
	if len(task.Criteria) != 2 {
		t.Errorf("Criteria = %v, want 2 items", task.Criteria)
	}
}

func TestParseVTSFile_EmptyArrays(t *testing.T) {
	content := `---
id: VTS-008
title: "Verify external_ref"
complexity: XS
status: pending
owner: ""
source: oracle
source_ref: "breakdown.md"
dependencies: []
files: []
---

# Verify external_ref

Run experiments.
`
	task, err := ParseVTSFile(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(task.Dependencies) != 0 {
		t.Errorf("Dependencies = %v, want empty", task.Dependencies)
	}
	if len(task.Files) != 0 {
		t.Errorf("Files = %v, want empty", task.Files)
	}
}

func TestParseVTSFile_MissingOptionalFields(t *testing.T) {
	content := `---
id: VTS-002
title: "Build Parser"
complexity: S
status: pending
dependencies: []
files: []
---

# Build Parser

Parse VTS files.
`
	task, err := ParseVTSFile(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if task.Owner != "" {
		t.Errorf("Owner = %q, want empty", task.Owner)
	}
	if task.Source != "" {
		t.Errorf("Source = %q, want empty", task.Source)
	}
	if task.SourceRef != "" {
		t.Errorf("SourceRef = %q, want empty", task.SourceRef)
	}
}

func TestParseVTSFile_NoCriteria(t *testing.T) {
	content := `---
id: VTS-003
title: "No Criteria Task"
complexity: M
status: pending
dependencies: []
files: []
---

# No Criteria Task

Just a description, no criteria section.
`
	task, err := ParseVTSFile(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(task.Criteria) != 0 {
		t.Errorf("Criteria = %v, want empty", task.Criteria)
	}
	if task.Description != "Just a description, no criteria section." {
		t.Errorf("Description = %q", task.Description)
	}
}

func TestParseVTSFile_MissingID(t *testing.T) {
	content := `---
title: "No ID"
complexity: S
status: pending
dependencies: []
files: []
---

# No ID
`
	_, err := ParseVTSFile(content)
	if err == nil {
		t.Fatal("expected error for missing id")
	}
}

func TestParseVTSFile_MissingTitle(t *testing.T) {
	content := `---
id: VTS-001
complexity: S
status: pending
dependencies: []
files: []
---

# Untitled
`
	_, err := ParseVTSFile(content)
	if err == nil {
		t.Fatal("expected error for missing title")
	}
}

func TestParseVTSFile_MalformedFrontmatter(t *testing.T) {
	content := `---
id: VTS-001
title: "Bad"
no closing delimiter
`
	_, err := ParseVTSFile(content)
	if err == nil {
		t.Fatal("expected error for missing closing ---")
	}
}

func TestReadDir(t *testing.T) {
	dir := t.TempDir()

	// Write 3 VTS files
	for i, id := range []string{"VTS-001", "VTS-002", "VTS-003"} {
		content := fmt.Sprintf(`---
id: %s
title: "Task %d"
complexity: S
status: pending
dependencies: []
files: []
---

# Task %d

Description for task %d.
`, id, i+1, i+1, i+1)
		err := os.WriteFile(filepath.Join(dir, id+".md"), []byte(content), 0644)
		if err != nil {
			t.Fatalf("write fixture: %v", err)
		}
	}

	tasks, err := ReadDir(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(tasks) != 3 {
		t.Fatalf("got %d tasks, want 3", len(tasks))
	}

	// Verify sorted by ID
	if tasks[0].ID != "VTS-001" || tasks[2].ID != "VTS-003" {
		t.Errorf("tasks not sorted: %s, %s, %s", tasks[0].ID, tasks[1].ID, tasks[2].ID)
	}
}

func TestReadDir_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	_, err := ReadDir(dir)
	if err == nil {
		t.Fatal("expected error for empty dir")
	}
}
