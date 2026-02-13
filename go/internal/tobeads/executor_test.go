package tobeads

import (
	"fmt"
	"testing"
)

func TestMockBrRunner_Create(t *testing.T) {
	mock := NewMockBrRunner()
	spec := BeadSpec{ExternalRef: "VTS-001", Title: "Test", Status: "open"}

	id, existed, err := mock.Create(spec)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if existed {
		t.Error("expected new, got already existed")
	}
	if id != "BR-001" {
		t.Errorf("id = %q, want BR-001", id)
	}
	if len(mock.Created) != 1 {
		t.Errorf("created = %d, want 1", len(mock.Created))
	}
}

func TestMockBrRunner_CreateDuplicate(t *testing.T) {
	mock := NewMockBrRunner()
	spec := BeadSpec{ExternalRef: "VTS-001", Title: "Test", Status: "open"}

	id1, _, _ := mock.Create(spec)
	id2, existed, err := mock.Create(spec)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !existed {
		t.Error("expected already existed on second create")
	}
	if id1 != id2 {
		t.Errorf("ids differ: %s vs %s", id1, id2)
	}
}

func TestMockBrRunner_CreateFailure(t *testing.T) {
	mock := NewMockBrRunner()
	mock.FailCreate["VTS-001"] = fmt.Errorf("br exploded")

	_, _, err := mock.Create(BeadSpec{ExternalRef: "VTS-001"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestMockBrRunner_DepAdd(t *testing.T) {
	mock := NewMockBrRunner()
	err := mock.DepAdd("BR-002", "BR-001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(mock.Deps) != 1 {
		t.Fatalf("deps = %d, want 1", len(mock.Deps))
	}
	if mock.Deps[0] != [2]string{"BR-002", "BR-001"} {
		t.Errorf("dep = %v", mock.Deps[0])
	}
}

func TestMockBrRunner_Sync(t *testing.T) {
	mock := NewMockBrRunner()
	err := mock.Sync()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !mock.Synced {
		t.Error("expected synced = true")
	}
}

func TestParseCreateJSON(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		wantID string
		wantOK bool
	}{
		{"string id", `{"id": "abc-123"}`, "abc-123", true},
		{"numeric id", `{"id": 42}`, "42", true},
		{"ID key", `{"ID": "def-456"}`, "def-456", true},
		{"no id", `{"title": "test"}`, "", false},
		{"bad json", `not json`, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := parseCreateJSON(tt.input)
			if tt.wantOK && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tt.wantOK && err == nil {
				t.Fatal("expected error")
			}
			if id != tt.wantID {
				t.Errorf("id = %q, want %q", id, tt.wantID)
			}
		})
	}
}
