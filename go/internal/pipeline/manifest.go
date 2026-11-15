package pipeline

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/jdonohoo/vern-bot/go/internal/llm"
	"github.com/jdonohoo/vern-bot/go/internal/tasks"
)

// ManifestFile is the run manifest a discovery run keeps in output/.
const ManifestFile = "run.json"

// Version is vern's version, recorded in run manifests. main sets it.
var Version = "dev"

// Manifest describes a discovery run for agents and tools: what ran, on
// which LLM and model, how long it took, what it cost where the CLIs say,
// and where the results are. Paths are relative to the discovery folder.
type Manifest struct {
	Version     int          `json:"manifest_version"`
	Kind        string       `json:"kind"` // discovery
	VernVersion string       `json:"vern_version"`
	Status      string       `json:"status"` // running, ok, partial, failed
	Idea        string       `json:"idea"`
	Dir         string       `json:"dir"`
	Pipeline    string       `json:"pipeline"`
	LLMMode     string       `json:"llm_mode"`
	ReadOnly    bool         `json:"read_only"`
	StartedAt   string       `json:"started_at"`
	FinishedAt  string       `json:"finished_at,omitempty"`
	DurationMS  int64        `json:"duration_ms"`
	Steps       []StepResult `json:"steps"`
	FailedSteps []int        `json:"failed_steps,omitempty"`
	Tasks       *TaskSummary `json:"tasks,omitempty"`
	VernHole    *Extra       `json:"vernhole,omitempty"`
	Oracle      *Extra       `json:"oracle,omitempty"`
	Usage       *llm.Usage   `json:"usage,omitempty"` // summed over steps that reported it
}

// TaskSummary points at a run's tasks.
type TaskSummary struct {
	Count  int    `json:"count"`
	File   string `json:"file"`    // output/tasks.json
	VTSDir string `json:"vts_dir"` // output/vts
}

// Extra records an optional stage.
type Extra struct {
	Status string `json:"status"` // ok, failed
	Path   string `json:"path"`
}

// ReadManifest loads output/run.json from a discovery folder.
func ReadManifest(discoveryDir string) (*Manifest, error) {
	data, err := os.ReadFile(filepath.Join(discoveryDir, "output", ManifestFile))
	if err != nil {
		return nil, err
	}
	m := &Manifest{}
	if err := json.Unmarshal(data, m); err != nil {
		return nil, err
	}
	return m, nil
}

// manifest builds the current manifest. final marks the end of the run.
func (p *Pipeline) manifest(final bool) *Manifest {
	dir := p.opts.DiscoveryDir
	rel := func(path string) string {
		if r, err := filepath.Rel(dir, path); err == nil {
			return filepath.ToSlash(r)
		}
		return filepath.ToSlash(path)
	}
	m := &Manifest{
		Version:     1,
		Kind:        "discovery",
		VernVersion: Version,
		Status:      "running",
		Idea:        p.opts.Idea,
		Dir:         dir,
		Pipeline:    p.mode,
		LLMMode:     p.cfg.LLMMode,
		ReadOnly:    !llm.WritesAllowed(),
		StartedAt:   p.startTime.UTC().Format(time.RFC3339),
		DurationMS:  time.Since(p.startTime).Milliseconds(),
		FailedSteps: p.failed,
		Steps:       []StepResult{},
	}
	if abs, err := filepath.Abs(dir); err == nil {
		m.Dir = abs
	}
	var total llm.Usage
	for _, r := range p.results {
		if r.Name == "" {
			continue
		}
		if r.OutputFile != "" {
			r.OutputFile = rel(r.OutputFile)
		}
		if r.Usage != nil {
			total = total.Add(*r.Usage)
		}
		m.Steps = append(m.Steps, r)
	}
	if !total.IsZero() {
		m.Usage = &total
	}

	outputDir := filepath.Join(dir, "output")
	if doc, err := tasks.Load(outputDir); err == nil && len(doc.Tasks) > 0 {
		m.Tasks = &TaskSummary{Count: len(doc.Tasks), File: "output/" + tasks.FileName, VTSDir: "output/vts"}
	}
	if p.vernhole != "" {
		m.VernHole = &Extra{Status: p.vernhole, Path: "vernhole"}
	}
	if fileSize(filepath.Join(dir, "oracle-vision.md")) > 0 {
		m.Oracle = &Extra{Status: "ok", Path: "oracle-vision.md"}
	}

	if final {
		m.FinishedAt = time.Now().UTC().Format(time.RFC3339)
		switch {
		case len(m.Steps) > 0 && len(p.failed) == len(m.Steps):
			m.Status = "failed"
		case len(p.failed) > 0 || p.vernhole == "failed" || m.Tasks == nil:
			m.Status = "partial"
		default:
			m.Status = "ok"
		}
	}
	return m
}

// writeManifest saves output/run.json and returns what it wrote.
func (p *Pipeline) writeManifest(final bool) *Manifest {
	m := p.manifest(final)
	data, err := json.MarshalIndent(m, "", "  ")
	if err == nil {
		os.WriteFile(filepath.Join(p.opts.DiscoveryDir, "output", ManifestFile), append(data, '\n'), 0644)
	}
	return m
}
