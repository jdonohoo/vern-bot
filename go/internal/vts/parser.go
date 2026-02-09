package vts

import (
	"fmt"
	"regexp"
	"strings"
)

// Task represents a parsed VTS task from architect output.
type Task struct {
	Num          int
	Title        string
	Body         string // Full raw body including header
	Description  string
	Complexity   string
	Dependencies []string // e.g. ["VTS-001", "VTS-002"]
	Criteria     []string
	Files        []string
}

var taskPattern = regexp.MustCompile(`(?im)^#{2,3}\s+Task\s+(\d+)\s*[:\.—]\s*(.+)`)

// ParseArchitectOutput extracts VTS tasks from architect breakdown markdown.
// Returns tasks, the header (text before first task), and the footer (text after last task).
func ParseArchitectOutput(content string) (tasks []Task, header string, footer string) {
	matches := taskPattern.FindAllStringSubmatchIndex(content, -1)
	if len(matches) == 0 {
		return nil, content, ""
	}

	// Header: everything before the first task
	header = strings.TrimRight(content[:matches[0][0]], "\n- ")

	// Find footer: first ## header (h2, not h3) after the last task
	lastTaskPos := matches[len(matches)-1][0]
	remaining := content[lastTaskPos:]
	footerPattern := regexp.MustCompile(`(?m)^## [^#]`)
	footerMatch := footerPattern.FindStringIndex(remaining)

	var footerAbs int
	if footerMatch != nil {
		footerAbs = lastTaskPos + footerMatch[0]
		// Strip trailing --- separator before footer
		pre := strings.TrimRight(content[:footerAbs], " \n\t")
		if strings.HasSuffix(pre, "---") {
			footerAbs = len(pre) - 3
		}
		footer = strings.TrimLeft(content[footerAbs:], "-")
		footer = strings.TrimSpace(footer)
	} else {
		footerAbs = len(content)
	}

	// Extract individual tasks
	for idx, match := range matches {
		taskNum := 0
		fmt.Sscanf(content[match[2]:match[3]], "%d", &taskNum)
		taskTitle := strings.TrimSpace(content[match[4]:match[5]])

		start := match[0]
		var end int
		if idx+1 < len(matches) {
			end = matches[idx+1][0]
		} else {
			end = footerAbs
		}

		taskBody := strings.TrimRight(content[start:end], " \n\t")
		// Remove trailing --- separator between tasks
		if strings.HasSuffix(taskBody, "---") {
			taskBody = strings.TrimRight(taskBody[:len(taskBody)-3], " \n\t")
		}

		task := Task{
			Num:          taskNum,
			Title:        taskTitle,
			Body:         taskBody,
			Description:  extractDescription(taskBody),
			Complexity:   normalizeComplexity(extractField(taskBody, "Complexity")),
			Dependencies: extractDependencies(taskBody),
			Criteria:     extractList(taskBody, "Acceptance Criteria"),
			Files:        extractFiles(taskBody),
		}
		tasks = append(tasks, task)
	}

	return tasks, header, footer
}

// extractField extracts a **Field:** value from markdown.
func extractField(body, field string) string {
	pattern := regexp.MustCompile(`\*\*` + regexp.QuoteMeta(field) + `:\*\*\s*(.+)`)
	m := pattern.FindStringSubmatch(body)
	if m != nil {
		return strings.TrimSpace(m[1])
	}
	return ""
}

// extractList extracts a **Field:** followed by comma-separated values or bullet list.
func extractList(body, field string) []string {
	// Try bullet list after the header first (more specific match)
	pattern := regexp.MustCompile(`\*\*` + regexp.QuoteMeta(field) + `:\*\*\s*\n((?:\s*[-*]\s+.+\n?)+)`)
	m := pattern.FindStringSubmatch(body)
	if m != nil {
		return parseBulletList(m[1])
	}

	// Fall back to inline comma-separated
	val := extractField(body, field)
	if val != "" && !isEmptyValue(val) {
		items := splitListValues(val)
		if len(items) > 0 {
			return items
		}
	}

	return nil
}

func extractDescription(body string) string {
	// Try **Description:** field first
	desc := extractField(body, "Description")
	if desc != "" {
		return desc
	}

	// Fall back to text between title and first bold field
	pattern := regexp.MustCompile(`(?im)^#{2,3}\s+Task\s+\d+.*\n+(.*?)(?:\n\*\*|\z)`)
	m := pattern.FindStringSubmatch(body)
	if m != nil {
		return strings.TrimSpace(m[1])
	}
	return ""
}

func extractDependencies(body string) []string {
	raw := extractField(body, "Dependencies")
	if raw == "" || isEmptyValue(raw) {
		return nil
	}
	// Extract task numbers like "Task 1", "Task 2, Task 3", "#1, #2"
	pattern := regexp.MustCompile(`(?i)(?:Task\s+|#)(\d+)`)
	matches := pattern.FindAllStringSubmatch(raw, -1)
	var deps []string
	for _, m := range matches {
		num := 0
		fmt.Sscanf(m[1], "%d", &num)
		deps = append(deps, fmt.Sprintf("VTS-%03d", num))
	}
	return deps
}

func extractFiles(body string) []string {
	for _, field := range []string{"Files", "Files Touched", "Files Affected", "Key Files"} {
		files := extractList(body, field)
		if len(files) > 0 {
			return files
		}
	}
	return nil
}

func normalizeComplexity(cx string) string {
	upper := strings.ToUpper(strings.Trim(cx, "*` "))
	switch upper {
	case "S", "M", "L", "XL":
		return upper
	default:
		if cx != "" {
			return cx
		}
		return "?"
	}
}

func isEmptyValue(val string) bool {
	lower := strings.ToLower(strings.TrimSpace(val))
	return lower == "none" || lower == "n/a" || lower == "-" || lower == ""
}

func splitListValues(val string) []string {
	parts := regexp.MustCompile(`[,;]`).Split(val, -1)
	var items []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		p = strings.TrimLeft(p, "- ")
		if p != "" {
			items = append(items, p)
		}
	}
	return items
}

func parseBulletList(text string) []string {
	var items []string
	for _, line := range strings.Split(strings.TrimSpace(text), "\n") {
		line = strings.TrimSpace(line)
		line = strings.TrimLeft(line, "-* ")
		if line != "" {
			items = append(items, line)
		}
	}
	return items
}
