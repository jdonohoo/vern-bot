package packs

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func write(t *testing.T, path, content string) {
	t.Helper()
	os.MkdirAll(filepath.Dir(path), 0755)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestSetupOrderAndLookup(t *testing.T) {
	defer Reset()
	root := t.TempDir()
	flagPack := filepath.Join(root, "flagpack")
	envPack := filepath.Join(root, "envpack")
	home := filepath.Join(root, "home")
	project := filepath.Join(root, "repo")
	sub := filepath.Join(project, "src", "deep")
	os.MkdirAll(sub, 0755)

	write(t, filepath.Join(flagPack, "agents", "editor.md"), "---\nname: editor\n---\nflag editor\n")
	write(t, filepath.Join(flagPack, "pack.json"), `{"name": "story-pack", "description": "Story personas", "version": "1.0.0"}`)
	write(t, filepath.Join(project, ".vern", "agents", "editor.md"), "---\nname: editor\n---\nproject editor\n")
	write(t, filepath.Join(project, ".vern", "agents", "lorekeeper.md"), "---\nname: lorekeeper\n---\nlore\n")
	write(t, filepath.Join(project, ".vern", "config.json"), `{"sink": "exec:cat", "packs": ["../shared"], "pipeline": "story"}`)
	write(t, filepath.Join(project, "shared", "pipelines", "story.json"), `{"name": "story", "kind": "creative", "steps": [{"step": 1}]}`)
	write(t, filepath.Join(envPack, "pipelines", "story.json"), `{"name": "story", "steps": []}`)
	write(t, filepath.Join(envPack, "pipelines", "unnamed.json"), `{"steps": [{"step": 1}]}`)
	write(t, filepath.Join(envPack, "pipelines", "broken.json"), `{not json`)
	write(t, filepath.Join(home, ".config", "vern", "packs", "mine", "agents", "solo.md"), "---\nname: solo\n---\n")
	t.Setenv("VERN_PATH", envPack)

	got, proj, err := Setup([]string{flagPack}, sub, home)
	if err != nil {
		t.Fatal(err)
	}
	if proj == nil || proj.Sink != "exec:cat" || proj.Pipeline != "story" {
		t.Fatalf("project = %+v", proj)
	}
	var names []string
	for _, p := range got {
		names = append(names, p.Name+":"+p.Source)
	}
	want := []string{"story-pack:flag", "repo/.vern:project", "shared:project", "envpack:VERN_PATH", "mine:user"}
	if len(names) != len(want) {
		t.Fatalf("packs = %v", names)
	}
	for i := range want {
		if names[i] != want[i] {
			t.Errorf("pack %d = %s, want %s", i, names[i], want[i])
		}
	}

	// First pack wins.
	path, ok := FindAgent("", "editor")
	if !ok || filepath.Dir(filepath.Dir(path)) != flagPack {
		t.Errorf("editor resolved to %s", path)
	}
	if _, ok := FindAgent("", "lorekeeper"); !ok {
		t.Error("project persona not found")
	}
	if _, ok := FindAgent("", "../etc/passwd"); ok {
		t.Error("names with path separators must not resolve")
	}

	// vern's own agents dir comes after the packs.
	own := filepath.Join(root, "agents")
	write(t, filepath.Join(own, "mighty.md"), "---\nname: mighty\n---\n")
	write(t, filepath.Join(own, "editor.md"), "---\nname: editor\n---\nbuilt-in\n")
	files := AgentFiles(own)
	if filepath.Dir(filepath.Dir(files["editor"])) != flagPack || files["mighty"] == "" || files["solo"] == "" {
		t.Errorf("AgentFiles = %v", files)
	}

	pls, errs := Pipelines()
	if len(errs) != 1 {
		t.Errorf("want one broken pipeline error, got %v", errs)
	}
	byName := map[string]Pipeline{}
	for _, p := range pls {
		byName[p.Name] = p
	}
	if byName["story"].Pack != "shared" || byName["story"].Kind != "creative" {
		t.Errorf("story should come from the project's shared pack: %+v", byName["story"])
	}
	if _, ok := byName["unnamed"]; !ok {
		t.Error("a pipeline without a name is named after its file")
	}
}

func TestFindProjectNone(t *testing.T) {
	p, err := FindProject(t.TempDir())
	if p != nil || err != nil {
		t.Errorf("FindProject in an empty tree = %+v, %v", p, err)
	}
}

func TestFindProjectBadConfig(t *testing.T) {
	dir := t.TempDir()
	write(t, filepath.Join(dir, ".vern", "config.json"), "{oops")
	p, err := FindProject(dir)
	if err == nil || p == nil {
		t.Errorf("a broken config should return the folder and an error: %+v, %v", p, err)
	}
}

// The example pack in the repo must load: personas with tiers, pipelines
// that parse and end in a consolidation step.
func TestCreativeExamplePack(t *testing.T) {
	defer Reset()
	dir := filepath.Join("..", "..", "..", "packs", "creative-example")
	if _, _, err := Setup([]string{dir}, t.TempDir(), t.TempDir()); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"editor", "worldbuilder", "skeptic-reader", "lorekeeper"} {
		path, ok := FindAgent("", name)
		if !ok {
			t.Errorf("persona %s not found", name)
			continue
		}
		data, _ := os.ReadFile(path)
		if !strings.Contains(string(data), "name: "+name+"\n") || !strings.Contains(string(data), "\ntier: ") {
			t.Errorf("%s needs a matching name and a tier", name)
		}
	}
	pls, errs := Pipelines()
	if len(errs) > 0 {
		t.Fatal(errs)
	}
	if len(pls) != 3 {
		t.Fatalf("want 3 pipelines, got %d", len(pls))
	}
	for _, pl := range pls {
		var steps []struct {
			Persona     string `json:"persona"`
			ContextMode string `json:"context_mode"`
		}
		if err := json.Unmarshal(pl.Steps, &steps); err != nil || len(steps) == 0 {
			t.Errorf("%s: bad steps: %v", pl.Name, err)
			continue
		}
		if pl.Kind != "creative" {
			t.Errorf("%s: kind %q, want creative", pl.Name, pl.Kind)
		}
		if last := steps[len(steps)-1]; last.ContextMode != "consolidation" {
			t.Errorf("%s: last step should be consolidation, got %q", pl.Name, last.ContextMode)
		}
	}
}
