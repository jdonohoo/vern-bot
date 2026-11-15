package tasks

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/jdonohoo/vern-bot/go/internal/vts"
)

const architect = `# Breakdown

### TASK 1: Set up the project

**Description:** Create the repo.

**Acceptance Criteria:**
- builds
- has CI

**Complexity:** S
**Dependencies:** None
**Files:** go.mod, main.go

### TASK 2: Add the parser

**Description:** Parse notes.

**Acceptance Criteria:**
- parses headings

**Complexity:** M
**Dependencies:** Task 1
**Files:** parser.go
`

func TestFromArchitectOutput(t *testing.T) {
	parsed, _, _ := vts.ParseArchitectOutput(architect)
	doc := FromVTS(parsed, "discovery", "notes to site")
	if err := doc.Validate(); err != nil {
		t.Fatal(err)
	}
	if len(doc.Tasks) != 2 || doc.Tasks[1].ID != "VTS-002" {
		t.Fatalf("tasks = %+v", doc.Tasks)
	}
	t2 := doc.Tasks[1]
	if strings.Join(t2.Dependencies, ",") != "VTS-001" || t2.Complexity != "M" || t2.Status != "pending" {
		t.Errorf("task 2 = %+v", t2)
	}
	if len(doc.Tasks[0].AcceptanceCriteria) != 2 || len(doc.Tasks[0].Files) != 2 {
		t.Errorf("task 1 lists = %+v", doc.Tasks[0])
	}
	if t2.VTSFile != "vts/vts-002-add-the-parser.md" {
		t.Errorf("vts file = %s", t2.VTSFile)
	}
}

func TestWriteReadAndLoadFallback(t *testing.T) {
	parsed, _, _ := vts.ParseArchitectOutput(architect)
	out := t.TempDir()
	if err := vts.WriteVTSFiles(parsed, filepath.Join(out, "vts"), "discovery", "arch.md", func(string) {}); err != nil {
		t.Fatal(err)
	}

	// No tasks.json yet: Load reads the VTS files.
	doc, err := Load(out)
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Tasks) != 2 || doc.Tasks[0].Num != 1 || doc.Tasks[1].Dependencies[0] != "VTS-001" {
		t.Fatalf("fallback doc = %+v", doc.Tasks)
	}

	want := FromVTS(parsed, "discovery", "idea")
	if err := want.Write(filepath.Join(out, FileName)); err != nil {
		t.Fatal(err)
	}
	got, err := Load(out)
	if err != nil || got.Idea != "idea" || got.Schema != SchemaID {
		t.Fatalf("Load = %+v, %v", got, err)
	}
}

func TestValidate(t *testing.T) {
	ok := Doc{Version: 1, Groups: []Group{{ID: "E1", Title: "Epic"}},
		Tasks: []Task{{ID: "VTS-001", Title: "a", Group: "E1"}, {ID: "VTS-002", Title: "b", Dependencies: []string{"VTS-001"}, Parent: "E1"}}}
	if err := ok.Validate(); err != nil {
		t.Fatal(err)
	}
	bad := []Doc{
		{Version: 2},
		{Version: 1, Tasks: []Task{{ID: "VTS-001", Title: "a"}, {ID: "VTS-001", Title: "b"}}},
		{Version: 1, Tasks: []Task{{ID: "VTS-001", Title: "a", Dependencies: []string{"VTS-009"}}}},
		{Version: 1, Tasks: []Task{{ID: "VTS-001", Title: "a", Group: "nope"}}},
		{Version: 1, Tasks: []Task{{ID: "VTS-001"}}},
	}
	for i, d := range bad {
		if d.Validate() == nil {
			t.Errorf("bad doc %d validated", i)
		}
	}
}
