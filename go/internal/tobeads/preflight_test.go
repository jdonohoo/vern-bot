package tobeads

import (
	"strings"
	"testing"

	"github.com/jdonohoo/vern-bot/go/internal/vts"
)

func TestPreflight_Valid(t *testing.T) {
	tasks := []vts.Task{
		{ID: "VTS-001", Title: "A", Status: "pending", Dependencies: nil},
		{ID: "VTS-002", Title: "B", Status: "active", Dependencies: []string{"VTS-001"}},
		{ID: "VTS-003", Title: "C", Status: "done", Dependencies: []string{"VTS-002"}},
	}
	report := Preflight(tasks)
	if !report.OK() {
		t.Fatalf("expected OK, got errors: %v", report.Errors)
	}
	if len(report.Order) != 3 {
		t.Errorf("order = %v, want 3 items", report.Order)
	}
	// VTS-001 should come first (no deps)
	if report.Order[0] != "VTS-001" {
		t.Errorf("order[0] = %s, want VTS-001", report.Order[0])
	}
}

func TestPreflight_DuplicateIDs(t *testing.T) {
	tasks := []vts.Task{
		{ID: "VTS-001", Title: "A", Status: "pending"},
		{ID: "VTS-001", Title: "B", Status: "pending"},
	}
	report := Preflight(tasks)
	if report.OK() {
		t.Fatal("expected error for duplicate IDs")
	}
	found := false
	for _, e := range report.Errors {
		if strings.Contains(e, "duplicate") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected duplicate error, got: %v", report.Errors)
	}
}

func TestPreflight_UnknownStatus(t *testing.T) {
	tasks := []vts.Task{
		{ID: "VTS-001", Title: "A", Status: "yolo"},
	}
	report := Preflight(tasks)
	if report.OK() {
		t.Fatal("expected error for unknown status")
	}
}

func TestPreflight_MissingDepTarget(t *testing.T) {
	tasks := []vts.Task{
		{ID: "VTS-001", Title: "A", Status: "pending", Dependencies: []string{"VTS-099"}},
	}
	report := Preflight(tasks)
	if report.OK() {
		t.Fatal("expected error for missing dep target")
	}
	found := false
	for _, e := range report.Errors {
		if strings.Contains(e, "VTS-099") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected missing dep error, got: %v", report.Errors)
	}
}

func TestPreflight_CycleDetection(t *testing.T) {
	tasks := []vts.Task{
		{ID: "VTS-001", Title: "A", Status: "pending", Dependencies: []string{"VTS-002"}},
		{ID: "VTS-002", Title: "B", Status: "pending", Dependencies: []string{"VTS-001"}},
	}
	report := Preflight(tasks)
	if report.OK() {
		t.Fatal("expected error for cycle")
	}
	found := false
	for _, e := range report.Errors {
		if strings.Contains(e, "cycle") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected cycle error, got: %v", report.Errors)
	}
}

func TestPreflight_TransitiveCycle(t *testing.T) {
	tasks := []vts.Task{
		{ID: "VTS-001", Title: "A", Status: "pending", Dependencies: []string{"VTS-003"}},
		{ID: "VTS-002", Title: "B", Status: "pending", Dependencies: []string{"VTS-001"}},
		{ID: "VTS-003", Title: "C", Status: "pending", Dependencies: []string{"VTS-002"}},
	}
	report := Preflight(tasks)
	if report.OK() {
		t.Fatal("expected error for transitive cycle")
	}
}

func TestPreflight_CollectsAllErrors(t *testing.T) {
	tasks := []vts.Task{
		{ID: "VTS-001", Title: "A", Status: "yolo"},
		{ID: "VTS-001", Title: "B", Status: "nope", Dependencies: []string{"VTS-099"}},
	}
	report := Preflight(tasks)
	// Should have: duplicate ID, two unknown statuses, missing dep
	if len(report.Errors) < 3 {
		t.Errorf("expected at least 3 errors, got %d: %v", len(report.Errors), report.Errors)
	}
}
