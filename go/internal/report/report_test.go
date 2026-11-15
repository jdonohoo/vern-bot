package report

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildFromFilesWithoutManifest(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "input"), 0755)
	os.MkdirAll(filepath.Join(dir, "output", "vts"), 0755)
	os.MkdirAll(filepath.Join(dir, "vernhole"), 0755)
	os.WriteFile(filepath.Join(dir, "input", "prompt.md"), []byte("# Discovery Prompt\n\nA notes site\n"), 0644)
	os.WriteFile(filepath.Join(dir, "output", "01-mighty-analysis.md"), []byte("## Analysis\n\n<script>alert(1)</script>\n\nSolid plan.\n"), 0644)
	os.WriteFile(filepath.Join(dir, "output", "vts", "vts-001-scaffold.md"),
		[]byte("---\nid: VTS-001\ntitle: \"Scaffold\"\ncomplexity: S\nstatus: pending\ndependencies: []\nfiles: []\n---\n\n# Scaffold\n\nSet it up.\n"), 0644)
	os.WriteFile(filepath.Join(dir, "vernhole", "synthesis.md"), []byte("Council agrees.\n"), 0644)

	html, err := Build(dir)
	if err != nil {
		t.Fatal(err)
	}
	s := string(html)
	for _, want := range []string{"A notes site", "VTS-001", "Scaffold", "Solid plan.", "Council agrees.", "VernHole synthesis"} {
		if !strings.Contains(s, want) {
			t.Errorf("report missing %q", want)
		}
	}
	if strings.Contains(s, "<script>alert(1)</script>") {
		t.Error("raw HTML from model output must not reach the report")
	}
	if _, err := Write(dir); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, "output", FileName)); err != nil {
		t.Error("report.html not written")
	}
}

func TestBuildNeedsOutputFolder(t *testing.T) {
	if _, err := Build(t.TempDir()); err == nil {
		t.Error("a folder without output/ should be an error")
	}
}
