package vts

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ReadFile reads a single VTS task file with YAML frontmatter and markdown body.
func ReadFile(path string) (*Task, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read VTS file: %w", err)
	}
	return ParseVTSFile(string(data))
}

// ReadDir reads all .md files in a directory as VTS tasks, sorted by ID.
func ReadDir(dir string) ([]Task, error) {
	entries, err := filepath.Glob(filepath.Join(dir, "*.md"))
	if err != nil {
		return nil, fmt.Errorf("glob VTS dir: %w", err)
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("no .md files found in %s", dir)
	}

	var tasks []Task
	var errs []string
	for _, path := range entries {
		task, err := ReadFile(path)
		if err != nil {
			errs = append(errs, fmt.Sprintf("%s: %v", filepath.Base(path), err))
			continue
		}
		tasks = append(tasks, *task)
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].ID < tasks[j].ID
	})

	if len(errs) > 0 && len(tasks) == 0 {
		return nil, fmt.Errorf("all files failed to parse:\n  %s", strings.Join(errs, "\n  "))
	}

	return tasks, nil
}

// ParseVTSFile parses a VTS task file from string content.
// Expects YAML frontmatter (--- delimited) followed by a markdown body.
func ParseVTSFile(content string) (*Task, error) {
	task := &Task{}
	lines := strings.Split(content, "\n")

	// State machine: 0=before frontmatter, 1=in frontmatter, 2=body
	state := 0
	var bodyLines []string
	var currentArrayKey string
	var currentArray []string

	flushArray := func() {
		if currentArrayKey == "" {
			return
		}
		switch currentArrayKey {
		case "dependencies":
			task.Dependencies = currentArray
		case "files":
			task.Files = currentArray
		}
		currentArrayKey = ""
		currentArray = nil
	}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		switch state {
		case 0:
			if trimmed == "---" {
				state = 1
			}
		case 1:
			if trimmed == "---" {
				flushArray()
				state = 2
				continue
			}

			// Array item: "  - value"
			if strings.HasPrefix(trimmed, "- ") && currentArrayKey != "" {
				val := strings.TrimPrefix(trimmed, "- ")
				val = stripQuotes(val)
				currentArray = append(currentArray, val)
				continue
			}

			// Key: value or Key: (start of array)
			if idx := strings.Index(line, ":"); idx > 0 && !strings.HasPrefix(trimmed, "-") {
				flushArray()
				key := strings.TrimSpace(line[:idx])
				val := strings.TrimSpace(line[idx+1:])

				// Empty array literal
				if val == "[]" {
					switch key {
					case "dependencies":
						task.Dependencies = nil
					case "files":
						task.Files = nil
					}
					continue
				}

				// Array start (key with no value, items follow)
				if val == "" {
					currentArrayKey = key
					currentArray = nil
					continue
				}

				val = stripQuotes(val)
				switch key {
				case "id":
					task.ID = val
				case "title":
					task.Title = val
				case "complexity":
					task.Complexity = val
				case "status":
					task.Status = val
				case "owner":
					task.Owner = val
				case "source":
					task.Source = val
				case "source_ref":
					task.SourceRef = val
				}
			}
		case 2:
			bodyLines = append(bodyLines, line)
		}
	}

	if state < 2 {
		return nil, fmt.Errorf("invalid VTS file: missing frontmatter delimiters")
	}

	if task.ID == "" {
		return nil, fmt.Errorf("invalid VTS file: missing required 'id' field")
	}
	if task.Title == "" {
		return nil, fmt.Errorf("invalid VTS file: missing required 'title' field")
	}

	// Parse Num from ID (e.g. "VTS-001" -> 1)
	fmt.Sscanf(task.ID, "VTS-%d", &task.Num)

	// Split body into description and criteria
	body := strings.Join(bodyLines, "\n")
	task.Body = strings.TrimSpace(body)
	splitBody(task)

	return task, nil
}

// splitBody splits the markdown body into Description (before ## Criteria) and Criteria list.
func splitBody(task *Task) {
	body := task.Body

	// Strip the leading # Title line
	lines := strings.SplitN(body, "\n", 2)
	if len(lines) > 1 && strings.HasPrefix(strings.TrimSpace(lines[0]), "#") {
		body = strings.TrimSpace(lines[1])
	}

	// Split on ## Criteria
	parts := strings.SplitN(body, "## Criteria", 2)
	task.Description = strings.TrimSpace(parts[0])

	if len(parts) > 1 {
		criteriaBlock := strings.TrimSpace(parts[1])
		for _, line := range strings.Split(criteriaBlock, "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "- ") {
				task.Criteria = append(task.Criteria, strings.TrimPrefix(line, "- "))
			}
		}
	}
}

// stripQuotes removes surrounding double quotes from a string.
func stripQuotes(s string) string {
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	return s
}
