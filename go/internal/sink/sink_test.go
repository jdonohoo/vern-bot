package sink

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/jdonohoo/vern-bot/go/internal/tasks"
	"github.com/jdonohoo/vern-bot/go/internal/vts"
)

const breakdown = `### TASK 1: Scaffold
**Description:** Set up.
**Complexity:** S
**Dependencies:** None
**Files:** main.go

### TASK 2: Build it
**Description:** Build.
**Complexity:** M
**Dependencies:** Task 1
**Files:** app.go
`

func runDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	parsed, _, _ := vts.ParseArchitectOutput(breakdown)
	if err := vts.WriteVTSFiles(parsed, filepath.Join(dir, "output", "vts"), "discovery", "x.md", func(string) {}); err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestValidate(t *testing.T) {
	for _, ok := range []string{"", "vts", "json", "json:out/t.json", "beads", "exec:cat", "EXEC:cat"} {
		if err := Validate(ok); err != nil {
			t.Errorf("Validate(%q) = %v", ok, err)
		}
	}
	for _, bad := range []string{"exec:", "exec:  ", "jira", "trello:x"} {
		if Validate(bad) == nil {
			t.Errorf("Validate(%q) should fail", bad)
		}
	}
}

func TestVTSAndJSONSinks(t *testing.T) {
	dir := runDir(t)
	info := Run(Options{Spec: "vts", DiscoveryDir: dir})
	if info.Status != "ok" {
		t.Fatalf("vts sink = %+v", info)
	}
	// A run from before tasks.json gets one written.
	if _, err := tasks.Read(filepath.Join(dir, "output", tasks.FileName)); err != nil {
		t.Errorf("tasks.json not written: %v", err)
	}
	dest := filepath.Join(t.TempDir(), "plan", "tasks.json")
	info = Run(Options{Spec: "json:" + dest, DiscoveryDir: dir})
	doc, err := tasks.Read(dest)
	if info.Status != "ok" || err != nil || len(doc.Tasks) != 2 {
		t.Errorf("json sink = %+v, %v", info, err)
	}
}

func TestExecSink(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("uses sh")
	}
	dir := runDir(t)
	out := filepath.Join(t.TempDir(), "got.json")
	info := Run(Options{Spec: `exec:cat > "` + out + `"; echo "run=$VERN_RUN_DIR"; test -f "$VERN_TASKS_FILE"`,
		DiscoveryDir: dir, Stderr: io.Discard})
	if info.Status != "ok" {
		t.Fatalf("exec sink = %+v", info)
	}
	doc, err := tasks.Read(out)
	if err != nil || len(doc.Tasks) != 2 || doc.Tasks[1].Dependencies[0] != "VTS-001" {
		t.Errorf("exec sink stdin = %+v, %v", doc, err)
	}
	log, _ := os.ReadFile(filepath.Join(dir, "output", LogFile))
	if !strings.Contains(string(log), "run=") || info.Log != "output/sink.log" {
		t.Errorf("sink log = %q", log)
	}

	info = Run(Options{Spec: "exec:exit 7", DiscoveryDir: dir, Stderr: io.Discard})
	if info.Status != "failed" || info.ExitCode != 7 {
		t.Errorf("failing exec sink = %+v", info)
	}
}

func TestNoTasksIsSkipped(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "output"), 0755)
	if info := Run(Options{Spec: "exec:cat", DiscoveryDir: dir}); info.Status != "skipped" {
		t.Errorf("run without tasks = %+v", info)
	}
}

func TestBeadsSinkUsesCallback(t *testing.T) {
	dir := runDir(t)
	var got string
	info := Run(Options{Spec: "beads", DiscoveryDir: dir, Beads: func(v string) error { got = v; return nil }})
	if info.Status != "ok" || filepath.Base(got) != "vts" {
		t.Errorf("beads sink = %+v, dir %s", info, got)
	}
}
