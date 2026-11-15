// Package tasks defines vern's structured task output (schema v1): the
// JSON form of a discovery run's VTS tasks that agents and task sinks read
// instead of parsing markdown.
package tasks

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jdonohoo/vern-bot/go/internal/vts"
)

// SchemaID identifies schema v1; the schema itself is schemas/tasks.v1.json.
const SchemaID = "https://github.com/jdonohoo/vern-bot/blob/master/schemas/tasks.v1.json"

// FileName is the task file a discovery run writes next to output/vts/.
const FileName = "tasks.json"

// Doc is a run's task graph.
type Doc struct {
	Schema  string  `json:"$schema"`
	Version int     `json:"version"`
	Source  string  `json:"source"`         // discovery, oracle
	Idea    string  `json:"idea,omitempty"` // the request the tasks answer
	Groups  []Group `json:"groups,omitempty"`
	Tasks   []Task  `json:"tasks"`
}

// Group is an optional grouping above tasks (an epic, a milestone). Sinks
// that support hierarchy create groups first and parent tasks to them.
type Group struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Parent string `json:"parent,omitempty"`
	Body   string `json:"body,omitempty"`
}

// Task is one unit of work.
type Task struct {
	ID                 string   `json:"id"` // VTS-001
	Num                int      `json:"num"`
	Title              string   `json:"title"`
	Kind               string   `json:"kind,omitempty"` // task (default) or deliverable
	Type               string   `json:"type,omitempty"` // a deliverable's type: outline, beats, character...
	Description        string   `json:"description,omitempty"`
	AcceptanceCriteria []string `json:"acceptance_criteria,omitempty"`
	Complexity         string   `json:"complexity,omitempty"` // S, M, L, XL
	Dependencies       []string `json:"dependencies,omitempty"`
	Files              []string `json:"files,omitempty"`
	Group              string   `json:"group,omitempty"`  // a Group id
	Parent             string   `json:"parent,omitempty"` // a Task or Group id
	Status             string   `json:"status,omitempty"` // pending, active, done...
	Body               string   `json:"body,omitempty"`   // the task's full markdown
	VTSFile            string   `json:"vts_file,omitempty"`
}

// ID returns the VTS id for a task number.
func ID(num int) string { return fmt.Sprintf("VTS-%03d", num) }

// FromVTS converts parsed VTS tasks (architect output or files on disk).
func FromVTS(in []vts.Task, source, idea string) Doc {
	doc := Doc{Schema: SchemaID, Version: 1, Source: source, Idea: idea, Tasks: []Task{}}
	for _, t := range in {
		id := t.ID
		if id == "" {
			id = ID(t.Num)
		}
		if t.Num == 0 {
			fmt.Sscanf(id, "VTS-%d", &t.Num)
		}
		status := t.Status
		if status == "" {
			status = "pending"
		}
		doc.Tasks = append(doc.Tasks, Task{
			ID:                 id,
			Num:                t.Num,
			Title:              t.Title,
			Kind:               "task",
			Type:               t.Type,
			Description:        t.Description,
			AcceptanceCriteria: t.Criteria,
			Complexity:         t.Complexity,
			Dependencies:       t.Dependencies,
			Files:              t.Files,
			Status:             status,
			Body:               strings.TrimSpace(t.Body),
			VTSFile:            fmt.Sprintf("vts/vts-%03d-%s.md", t.Num, vts.Slugify(t.Title)),
		})
	}
	return doc
}

// AsDeliverables marks every task as a deliverable (creative pipelines).
func (d Doc) AsDeliverables() Doc {
	for i := range d.Tasks {
		d.Tasks[i].Kind = "deliverable"
	}
	return d
}

// Validate checks the rules schema v1 cannot express: unique ids, and
// dependencies, groups and parents that point at something in the doc.
func (d Doc) Validate() error {
	if d.Version != 1 {
		return fmt.Errorf("unsupported task schema version %d", d.Version)
	}
	ids := map[string]bool{}
	groups := map[string]bool{}
	for _, g := range d.Groups {
		if g.ID == "" || groups[g.ID] {
			return fmt.Errorf("group id %q is empty or repeated", g.ID)
		}
		groups[g.ID] = true
	}
	for _, t := range d.Tasks {
		if t.ID == "" || t.Title == "" {
			return fmt.Errorf("task %d needs an id and a title", t.Num)
		}
		if ids[t.ID] {
			return fmt.Errorf("task id %s is repeated", t.ID)
		}
		ids[t.ID] = true
	}
	for _, t := range d.Tasks {
		for _, dep := range t.Dependencies {
			if !ids[dep] {
				return fmt.Errorf("%s depends on unknown task %s", t.ID, dep)
			}
		}
		if t.Group != "" && !groups[t.Group] {
			return fmt.Errorf("%s is in unknown group %s", t.ID, t.Group)
		}
		if t.Parent != "" && !ids[t.Parent] && !groups[t.Parent] {
			return fmt.Errorf("%s has unknown parent %s", t.ID, t.Parent)
		}
	}
	return nil
}

// Write saves the doc as indented JSON.
func (d Doc) Write(path string) error {
	data, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0644)
}

// Read loads a task file.
func Read(path string) (Doc, error) {
	var d Doc
	data, err := os.ReadFile(path)
	if err != nil {
		return d, err
	}
	if err := json.Unmarshal(data, &d); err != nil {
		return d, fmt.Errorf("parse %s: %w", path, err)
	}
	return d, nil
}

// Load returns the tasks for a run's output folder: tasks.json when it
// exists, else the VTS files in output/vts/ (runs from before v3.1).
func Load(outputDir string) (Doc, error) {
	if d, err := Read(filepath.Join(outputDir, FileName)); err == nil {
		return d, nil
	}
	vtsTasks, err := vts.ReadDir(filepath.Join(outputDir, "vts"))
	if err != nil {
		return Doc{}, err
	}
	return FromVTS(vtsTasks, "discovery", ""), nil
}
