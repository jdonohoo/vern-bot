package vts

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Slugify converts a string to a URL-friendly slug.
func Slugify(s string) string {
	s = strings.ToLower(s)
	s = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(s, "-")
	s = regexp.MustCompile(`-+`).ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}

// WriteVTSFiles writes individual VTS task files to the given directory.
// source identifies the origin (e.g. "discovery" or "oracle").
func WriteVTSFiles(tasks []Task, dir string, source string, sourceRef string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create VTS dir: %w", err)
	}

	for _, task := range tasks {
		taskID := fmt.Sprintf("VTS-%03d", task.Num)

		var lines []string
		lines = append(lines, "---")
		lines = append(lines, fmt.Sprintf("id: %s", taskID))
		lines = append(lines, fmt.Sprintf("title: %q", task.Title))
		lines = append(lines, fmt.Sprintf("complexity: %s", task.Complexity))
		lines = append(lines, "status: pending")
		lines = append(lines, `owner: ""`)
		lines = append(lines, fmt.Sprintf("source: %s", source))
		lines = append(lines, fmt.Sprintf("source_ref: %q", sourceRef))

		if len(task.Dependencies) > 0 {
			lines = append(lines, "dependencies:")
			for _, d := range task.Dependencies {
				lines = append(lines, fmt.Sprintf("  - %s", d))
			}
		} else {
			lines = append(lines, "dependencies: []")
		}

		if len(task.Files) > 0 {
			lines = append(lines, "files:")
			for _, f := range task.Files {
				lines = append(lines, fmt.Sprintf("  - %q", f))
			}
		} else {
			lines = append(lines, "files: []")
		}

		lines = append(lines, "---")
		lines = append(lines, "")
		lines = append(lines, fmt.Sprintf("# %s", task.Title))
		lines = append(lines, "")

		if task.Description != "" {
			lines = append(lines, task.Description)
			lines = append(lines, "")
		}

		if len(task.Criteria) > 0 {
			lines = append(lines, "## Criteria")
			lines = append(lines, "")
			for _, c := range task.Criteria {
				lines = append(lines, fmt.Sprintf("- %s", c))
			}
			lines = append(lines, "")
		}

		filename := fmt.Sprintf("vts-%03d-%s.md", task.Num, Slugify(task.Title))
		path := filepath.Join(dir, filename)
		content := strings.Join(lines, "\n")
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return fmt.Errorf("write VTS file %s: %w", filename, err)
		}
		fmt.Printf("  Created: vts/%s\n", filename)
	}

	return nil
}

// WriteSummary rewrites the architect file as a summary with a VTS task index.
func WriteSummary(tasks []Task, architectFile string, header string, footer string, indexTitle string) error {
	if indexTitle == "" {
		indexTitle = "VTS Task Index"
	}

	var lines []string
	lines = append(lines, "| ID | Task | Complexity | Dependencies |")
	lines = append(lines, "|----|------|------------|--------------|")

	for _, task := range tasks {
		taskID := fmt.Sprintf("VTS-%03d", task.Num)
		depStr := "None"
		if len(task.Dependencies) > 0 {
			depStr = strings.Join(task.Dependencies, ", ")
		}
		lines = append(lines, fmt.Sprintf("| %s | %s | %s | %s |", taskID, task.Title, task.Complexity, depStr))
	}

	taskIndex := strings.Join(lines, "\n")

	var content strings.Builder
	content.WriteString(header)
	content.WriteString("\n\n")
	content.WriteString(fmt.Sprintf("## %s\n\n", indexTitle))
	content.WriteString("Individual VTS files: `vts/`\n\n")
	content.WriteString(taskIndex)
	content.WriteString("\n")

	if footer != "" {
		content.WriteString("\n")
		content.WriteString(footer)
		content.WriteString("\n")
	}

	if err := os.WriteFile(architectFile, []byte(content.String()), 0644); err != nil {
		return fmt.Errorf("write summary: %w", err)
	}

	fmt.Printf("  Rewrote %s as summary (%d VTS tasks extracted)\n",
		filepath.Base(architectFile), len(tasks))
	return nil
}
