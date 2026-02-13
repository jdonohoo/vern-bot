package tobeads

import (
	"fmt"
	"strings"

	"github.com/jdonohoo/vern-bot/go/internal/vts"
)

// PreflightReport holds all validation results.
type PreflightReport struct {
	Errors   []string // Hard-fail issues
	Warnings []string // Non-blocking notices
	Order    []string // Topological order of VTS IDs (from Kahn's)
}

func (r *PreflightReport) OK() bool {
	return len(r.Errors) == 0
}

func (r *PreflightReport) String() string {
	var b strings.Builder
	if len(r.Errors) > 0 {
		b.WriteString("ERRORS:\n")
		for _, e := range r.Errors {
			b.WriteString("  - " + e + "\n")
		}
	}
	if len(r.Warnings) > 0 {
		b.WriteString("WARNINGS:\n")
		for _, w := range r.Warnings {
			b.WriteString("  - " + w + "\n")
		}
	}
	return b.String()
}

// Preflight validates a batch of VTS tasks before import.
// Checks: duplicate IDs, unknown statuses, missing dep targets, cycles.
func Preflight(tasks []vts.Task) *PreflightReport {
	report := &PreflightReport{}
	ids := map[string]bool{}

	// Check duplicate IDs
	for _, t := range tasks {
		if ids[t.ID] {
			report.Errors = append(report.Errors, fmt.Sprintf("duplicate VTS ID: %s", t.ID))
		}
		ids[t.ID] = true
	}

	// Check unknown statuses (before normalization)
	for _, t := range tasks {
		if _, ok := StatusMap[strings.ToLower(t.Status)]; !ok {
			report.Errors = append(report.Errors, fmt.Sprintf("%s: unknown status %q", t.ID, t.Status))
		}
	}

	// Check missing dependency targets
	for _, t := range tasks {
		for _, dep := range t.Dependencies {
			if !ids[dep] {
				report.Errors = append(report.Errors, fmt.Sprintf("%s: dependency %s not found in input set", t.ID, dep))
			}
		}
	}

	// Cycle detection via Kahn's algorithm
	order, cycleNodes := kahnsSort(tasks)
	if len(cycleNodes) > 0 {
		report.Errors = append(report.Errors, fmt.Sprintf("dependency cycle detected involving: %s", strings.Join(cycleNodes, ", ")))
	} else {
		report.Order = order
	}

	return report
}

// kahnsSort performs topological sort. Returns ordered IDs and any cycle members.
func kahnsSort(tasks []vts.Task) (order []string, cycleNodes []string) {
	// Build adjacency + in-degree
	inDegree := map[string]int{}
	dependents := map[string][]string{} // dep -> list of tasks that depend on it

	for _, t := range tasks {
		if _, ok := inDegree[t.ID]; !ok {
			inDegree[t.ID] = 0
		}
		for _, dep := range t.Dependencies {
			inDegree[t.ID]++
			dependents[dep] = append(dependents[dep], t.ID)
		}
	}

	// Seed queue with zero in-degree nodes
	var queue []string
	for id, deg := range inDegree {
		if deg == 0 {
			queue = append(queue, id)
		}
	}

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		order = append(order, node)

		for _, dependent := range dependents[node] {
			inDegree[dependent]--
			if inDegree[dependent] == 0 {
				queue = append(queue, dependent)
			}
		}
	}

	// Any nodes still with in-degree > 0 are in cycles
	for id, deg := range inDegree {
		if deg > 0 {
			cycleNodes = append(cycleNodes, id)
		}
	}

	return order, cycleNodes
}
