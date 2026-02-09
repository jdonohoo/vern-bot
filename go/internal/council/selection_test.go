package council

import (
	"testing"
)

func testRoster() []Vern {
	return hardcodedRoster()
}

func TestResolveCouncilHammers(t *testing.T) {
	selected, name := ResolveCouncil("hammers", testRoster(), 3)
	if name != "Council of the Three Hammers" {
		t.Errorf("name: got %q", name)
	}
	if len(selected) != 3 {
		t.Fatalf("count: got %d, want 3", len(selected))
	}
	ids := map[string]bool{}
	for _, v := range selected {
		ids[v.ID] = true
	}
	for _, id := range []string{"great", "mediocre", "ketamine"} {
		if !ids[id] {
			t.Errorf("missing core member: %s", id)
		}
	}
}

func TestResolveCouncilConflict(t *testing.T) {
	selected, name := ResolveCouncil("conflict", testRoster(), 3)
	if name != "Max Conflict" {
		t.Errorf("name: got %q", name)
	}
	if len(selected) != 6 {
		t.Fatalf("count: got %d, want 6", len(selected))
	}
}

func TestResolveCouncilFull(t *testing.T) {
	roster := testRoster()
	selected, name := ResolveCouncil("full", roster, 3)
	if name != "The Full Vern Experience" {
		t.Errorf("name: got %q", name)
	}
	if len(selected) != len(roster) {
		t.Errorf("count: got %d, want %d", len(selected), len(roster))
	}
}

func TestResolveCouncilRandom(t *testing.T) {
	roster := testRoster()
	selected, name := ResolveCouncil("random", roster, 3)
	if name != "Fate's Hand" {
		t.Errorf("name: got %q", name)
	}
	if len(selected) < 3 || len(selected) > len(roster) {
		t.Errorf("count: got %d, want 3-%d", len(selected), len(roster))
	}
}

func TestResolveCouncilInner(t *testing.T) {
	selected, name := ResolveCouncil("inner", testRoster(), 3)
	if name != "The Inner Circle" {
		t.Errorf("name: got %q", name)
	}
	if len(selected) < 3 || len(selected) > 5 {
		t.Errorf("count: got %d, want 3-5", len(selected))
	}
	// Core members must be present
	ids := map[string]bool{}
	for _, v := range selected {
		ids[v.ID] = true
	}
	for _, id := range []string{"architect", "inverse", "paranoid"} {
		if !ids[id] {
			t.Errorf("missing core member: %s", id)
		}
	}
}

func TestResolveCouncilBareNumber(t *testing.T) {
	selected, name := ResolveCouncil("5", testRoster(), 3)
	if name != "" {
		t.Errorf("name should be empty for bare number, got %q", name)
	}
	if len(selected) != 5 {
		t.Errorf("count: got %d, want 5", len(selected))
	}
}

func TestResolveCouncilBareNumberClamped(t *testing.T) {
	// Too low: should clamp to min
	selected, _ := ResolveCouncil("1", testRoster(), 3)
	if len(selected) != 3 {
		t.Errorf("count: got %d, want 3 (clamped)", len(selected))
	}
}

func TestScanRoster(t *testing.T) {
	// Test with actual agents dir
	roster := ScanRoster("/Users/justin/projects/jdonohoo/vern-bot/agents")
	if len(roster) < 10 {
		t.Errorf("expected at least 10 agents, got %d", len(roster))
	}
	// Should NOT contain oracle or vernhole-orchestrator
	for _, v := range roster {
		if v.ID == "oracle" || v.ID == "vernhole-orchestrator" {
			t.Errorf("roster should not contain %s", v.ID)
		}
	}
}
