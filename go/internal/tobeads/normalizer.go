package tobeads

import (
	"fmt"
	"strings"

	"github.com/jdonohoo/vern-bot/go/internal/vts"
)

// Normalize transforms a slice of parsed VTS tasks into BeadSpecs.
// Returns specs and any normalization errors (non-fatal, collected).
func Normalize(tasks []vts.Task) ([]BeadSpec, []string) {
	var specs []BeadSpec
	var errs []string

	for _, t := range tasks {
		spec, err := normalizeTask(t)
		if err != nil {
			errs = append(errs, fmt.Sprintf("%s: %v", t.ID, err))
			continue
		}
		specs = append(specs, spec)
	}

	return specs, errs
}

func normalizeTask(t vts.Task) (BeadSpec, error) {
	spec := BeadSpec{
		ExternalRef:  t.ID,
		Title:        t.Title,
		Dependencies: t.Dependencies,
	}

	// Map status
	beadsStatus, ok := StatusMap[strings.ToLower(t.Status)]
	if !ok {
		return spec, fmt.Errorf("unknown status %q", t.Status)
	}
	spec.Status = beadsStatus

	// Build labels
	var labels []string
	cx := strings.ToUpper(strings.TrimSpace(t.Complexity))
	if ValidComplexity[cx] {
		labels = append(labels, "complexity:"+cx)
	} else if cx != "" && cx != "?" {
		return spec, fmt.Errorf("unknown complexity %q", t.Complexity)
	}
	if t.Source != "" {
		labels = append(labels, "source:"+t.Source)
	}
	spec.Labels = labels

	// Assignee
	if t.Owner != "" {
		spec.Assignee = t.Owner
	}

	// Build description: body + files list + source_ref metadata + criteria
	spec.Description = buildDescription(t)

	return spec, nil
}

func buildDescription(t vts.Task) string {
	var parts []string

	if t.Description != "" {
		parts = append(parts, t.Description)
	}

	if len(t.Files) > 0 {
		var fileLines []string
		fileLines = append(fileLines, "## Files")
		for _, f := range t.Files {
			fileLines = append(fileLines, "- "+f)
		}
		parts = append(parts, strings.Join(fileLines, "\n"))
	}

	if t.SourceRef != "" {
		parts = append(parts, fmt.Sprintf("---\n_Source: %s | Ref: %s_", t.Source, t.SourceRef))
	}

	if len(t.Criteria) > 0 {
		var criteriaLines []string
		criteriaLines = append(criteriaLines, "## Acceptance Criteria")
		for _, c := range t.Criteria {
			criteriaLines = append(criteriaLines, "- "+c)
		}
		parts = append(parts, strings.Join(criteriaLines, "\n"))
	}

	return strings.Join(parts, "\n\n")
}
