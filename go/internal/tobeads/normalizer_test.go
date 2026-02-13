package tobeads

import (
	"strings"
	"testing"

	"github.com/jdonohoo/vern-bot/go/internal/vts"
)

func TestNormalize_ValidTask(t *testing.T) {
	tasks := []vts.Task{{
		ID:           "VTS-001",
		Title:        "Test Task",
		Status:       "pending",
		Complexity:   "M",
		Source:       "oracle",
		SourceRef:    "breakdown.md",
		Description:  "Do the thing.",
		Dependencies: []string{"VTS-002"},
		Files:        []string{"foo.go", "bar.go"},
		Criteria:     []string{"It works", "It doesn't break"},
	}}

	specs, errs := Normalize(tasks)
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(specs) != 1 {
		t.Fatalf("got %d specs, want 1", len(specs))
	}

	s := specs[0]
	if s.ExternalRef != "VTS-001" {
		t.Errorf("ExternalRef = %q", s.ExternalRef)
	}
	if s.Status != "open" {
		t.Errorf("Status = %q, want open", s.Status)
	}
	if len(s.Labels) != 2 {
		t.Errorf("Labels = %v, want 2", s.Labels)
	}
	if s.Labels[0] != "complexity:M" {
		t.Errorf("Labels[0] = %q", s.Labels[0])
	}
	if s.Labels[1] != "source:oracle" {
		t.Errorf("Labels[1] = %q", s.Labels[1])
	}
	if !strings.Contains(s.Description, "## Files") {
		t.Error("description missing files section")
	}
	if !strings.Contains(s.Description, "## Acceptance Criteria") {
		t.Error("description missing criteria section")
	}
	if !strings.Contains(s.Description, "breakdown.md") {
		t.Error("description missing source_ref metadata")
	}
}

func TestNormalize_AllStatuses(t *testing.T) {
	for vtsStatus, beadsStatus := range StatusMap {
		tasks := []vts.Task{{
			ID:     "VTS-001",
			Title:  "Test",
			Status: vtsStatus,
		}}
		specs, errs := Normalize(tasks)
		if len(errs) > 0 {
			t.Errorf("status %q: unexpected error: %v", vtsStatus, errs)
		}
		if len(specs) == 1 && specs[0].Status != beadsStatus {
			t.Errorf("status %q → %q, want %q", vtsStatus, specs[0].Status, beadsStatus)
		}
	}
}

func TestNormalize_UnknownStatus(t *testing.T) {
	tasks := []vts.Task{{
		ID:     "VTS-001",
		Title:  "Test",
		Status: "yolo",
	}}
	_, errs := Normalize(tasks)
	if len(errs) == 0 {
		t.Fatal("expected error for unknown status")
	}
}

func TestNormalize_UnknownComplexity(t *testing.T) {
	tasks := []vts.Task{{
		ID:         "VTS-001",
		Title:      "Test",
		Status:     "pending",
		Complexity: "XXXL",
	}}
	_, errs := Normalize(tasks)
	if len(errs) == 0 {
		t.Fatal("expected error for unknown complexity")
	}
}

func TestNormalize_EmptyOwner(t *testing.T) {
	tasks := []vts.Task{{
		ID:     "VTS-001",
		Title:  "Test",
		Status: "pending",
		Owner:  "",
	}}
	specs, _ := Normalize(tasks)
	if len(specs) == 1 && specs[0].Assignee != "" {
		t.Errorf("Assignee = %q, want empty", specs[0].Assignee)
	}
}
