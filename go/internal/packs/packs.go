// Package packs finds extra personas and pipelines outside vern's own
// agents/ folder, and reads a project's .vern/config.json.
//
// A pack is a directory with any of:
//
//	pack.json          {"name", "description", "version"} (optional)
//	agents/*.md        personas, same format as vern's agents/
//	pipelines/*.json   discovery pipelines: {"name", "description", "steps": [...]}
//
// Packs are searched in order, first match wins: --pack flags, the
// project's .vern/ folder and the packs its config lists, VERN_PATH, then
// ~/.config/vern/packs/*. vern's own agents/ come after every pack, and
// the personas built into the binary last.
package packs

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

// ProjectDirName is the per-project folder: a pack plus config.json.
const ProjectDirName = ".vern"

// Pack is one pack directory.
type Pack struct {
	Dir         string `json:"dir"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Version     string `json:"version,omitempty"`
	Source      string `json:"source"` // flag, project, VERN_PATH, user
}

// Project is a project's .vern/config.json. Every field is optional; set
// fields override the user config for runs in that project.
type Project struct {
	// Sink is the default task sink for discovery runs (vts, json, beads,
	// exec:<command>).
	Sink string `json:"sink,omitempty"`
	// Packs are more pack directories, relative to the .vern folder.
	Packs []string `json:"packs,omitempty"`
	// Pipeline is the default discovery pipeline name.
	Pipeline string `json:"pipeline,omitempty"`
	// LLMMode is the default LLM mode.
	LLMMode string `json:"llm_mode,omitempty"`
	// AllowWrites lets the LLM CLIs edit files in this project.
	AllowWrites *bool `json:"allow_writes,omitempty"`
	// Tiers override tier models and effort (same shape as the user config).
	Tiers json.RawMessage `json:"tiers,omitempty"`

	// Dir is the .vern folder this was read from (not serialized).
	Dir string `json:"-"`
}

var (
	mu      sync.RWMutex
	active  []Pack
	project *Project
)

// Setup resolves the pack search path. flagDirs come from --pack, cwd is
// where the project's .vern/ is looked for (walking up), home is the
// user's home directory. It returns the packs found and the project
// config, if any, and makes them active for this process.
func Setup(flagDirs []string, cwd, home string) ([]Pack, *Project, error) {
	var out []Pack
	seen := map[string]bool{}
	add := func(dir, source string) {
		abs, err := filepath.Abs(dir)
		if err != nil {
			return
		}
		if info, err := os.Stat(abs); err != nil || !info.IsDir() || seen[abs] {
			return
		}
		seen[abs] = true
		out = append(out, readPack(abs, source))
	}

	for _, d := range flagDirs {
		add(d, "flag")
	}

	proj, err := FindProject(cwd)
	if proj != nil {
		add(proj.Dir, "project")
		for _, d := range proj.Packs {
			if !filepath.IsAbs(d) {
				d = filepath.Join(proj.Dir, d)
			}
			add(d, "project")
		}
	}

	for _, d := range filepath.SplitList(os.Getenv("VERN_PATH")) {
		if d != "" {
			add(d, "VERN_PATH")
		}
	}

	if home != "" {
		userPacks := filepath.Join(home, ".config", "vern", "packs")
		if entries, rerr := os.ReadDir(userPacks); rerr == nil {
			for _, e := range entries {
				if e.IsDir() {
					add(filepath.Join(userPacks, e.Name()), "user")
				}
			}
		}
	}

	mu.Lock()
	active = out
	project = proj
	mu.Unlock()
	return out, proj, err
}

// Active returns the packs Setup found.
func Active() []Pack {
	mu.RLock()
	defer mu.RUnlock()
	return append([]Pack{}, active...)
}

// ActiveProject returns the project config Setup found, or nil.
func ActiveProject() *Project {
	mu.RLock()
	defer mu.RUnlock()
	return project
}

// Reset clears the active packs and project (for tests).
func Reset() {
	mu.Lock()
	active, project = nil, nil
	mu.Unlock()
}

func readPack(dir, source string) Pack {
	p := Pack{Dir: dir, Name: filepath.Base(dir), Source: source}
	if source == "project" && filepath.Base(dir) == ProjectDirName {
		p.Name = filepath.Base(filepath.Dir(dir)) + "/" + ProjectDirName
	}
	if data, err := os.ReadFile(filepath.Join(dir, "pack.json")); err == nil {
		var meta struct{ Name, Description, Version string }
		if json.Unmarshal(data, &meta) == nil {
			if meta.Name != "" {
				p.Name = meta.Name
			}
			p.Description, p.Version = meta.Description, meta.Version
		}
	}
	return p
}

// FindProject walks up from dir to the nearest .vern/ folder and reads its
// config.json (an empty config when the folder has none). It returns nil
// when there is no .vern/ folder, and an error with the config it could
// not parse.
func FindProject(dir string) (*Project, error) {
	if dir == "" {
		return nil, nil
	}
	dir, _ = filepath.Abs(dir)
	home, _ := os.UserHomeDir()
	for {
		candidate := filepath.Join(dir, ProjectDirName)
		// ~/.vern would make every folder under home a project; skip it.
		if info, err := os.Stat(candidate); err == nil && info.IsDir() && dir != home {
			p := &Project{Dir: candidate}
			data, err := os.ReadFile(filepath.Join(candidate, "config.json"))
			if err == nil {
				if jerr := json.Unmarshal(data, p); jerr != nil {
					return &Project{Dir: candidate}, jerr
				}
				p.Dir = candidate
			}
			return p, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return nil, nil
		}
		dir = parent
	}
}

// FindAgent returns the path of a persona file: the first pack that has
// agents/<name>.md, else agentsDir/<name>.md. ok is false when neither
// has it (callers then use the built-in personas).
func FindAgent(agentsDir, name string) (string, bool) {
	if name == "" || strings.ContainsAny(name, `/\`) {
		return "", false
	}
	for _, p := range Active() {
		path := filepath.Join(p.Dir, "agents", name+".md")
		if fileExists(path) {
			return path, true
		}
	}
	if agentsDir != "" {
		path := filepath.Join(agentsDir, name+".md")
		if fileExists(path) {
			return path, true
		}
	}
	return "", false
}

// AgentFiles maps each persona name to its file across the packs and
// agentsDir, with packs taking precedence.
func AgentFiles(agentsDir string) map[string]string {
	out := map[string]string{}
	dirs := []string{}
	for _, p := range Active() {
		dirs = append(dirs, filepath.Join(p.Dir, "agents"))
	}
	if agentsDir != "" {
		dirs = append(dirs, agentsDir)
	}
	for _, d := range dirs {
		files, _ := filepath.Glob(filepath.Join(d, "*.md"))
		sort.Strings(files)
		for _, f := range files {
			name := strings.TrimSuffix(filepath.Base(f), ".md")
			if _, ok := out[name]; !ok {
				out[name] = f
			}
		}
	}
	return out
}

// PackAgentFiles maps persona names to files for the packs only.
func PackAgentFiles() map[string]string { return AgentFiles("") }

// Pipeline is a pipeline file from a pack. Steps are left raw so this
// package does not depend on config.
type Pipeline struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Kind        string          `json:"kind,omitempty"` // software (default) or creative
	Steps       json.RawMessage `json:"steps"`
	Pack        string          `json:"-"`
	File        string          `json:"-"`
}

// Pipelines returns the pipelines the active packs define, first pack
// wins for a name. A file without a name is named after the file.
func Pipelines() ([]Pipeline, []error) {
	var out []Pipeline
	var errs []error
	seen := map[string]bool{}
	for _, p := range Active() {
		files, _ := filepath.Glob(filepath.Join(p.Dir, "pipelines", "*.json"))
		sort.Strings(files)
		for _, f := range files {
			data, err := os.ReadFile(f)
			if err != nil {
				errs = append(errs, err)
				continue
			}
			var pl Pipeline
			if err := json.Unmarshal(data, &pl); err != nil {
				errs = append(errs, &FileError{File: f, Err: err})
				continue
			}
			if pl.Name == "" {
				pl.Name = strings.TrimSuffix(filepath.Base(f), ".json")
			}
			if seen[pl.Name] {
				continue
			}
			seen[pl.Name] = true
			pl.Pack, pl.File = p.Name, f
			out = append(out, pl)
		}
	}
	return out, errs
}

// FileError is a pack file that could not be read.
type FileError struct {
	File string
	Err  error
}

func (e *FileError) Error() string { return e.File + ": " + e.Err.Error() }

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
