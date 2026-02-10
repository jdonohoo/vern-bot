package embedded

import (
	"encoding/json"
	"strings"
	"testing"
)

var expectedAgents = []string{
	"academic", "architect", "enterprise", "great", "inverse",
	"ketamine", "mediocre", "mighty", "nyquil", "optimist",
	"oracle", "paranoid", "retro", "startup", "ux",
	"vernhole-orchestrator", "yolo",
}

func TestListAgentsReturnsAll(t *testing.T) {
	names := ListAgents()
	if len(names) != len(expectedAgents) {
		t.Fatalf("ListAgents() returned %d agents, want %d", len(names), len(expectedAgents))
	}
	for i, name := range names {
		if name != expectedAgents[i] {
			t.Errorf("ListAgents()[%d] = %q, want %q", i, name, expectedAgents[i])
		}
	}
}

func TestListAgentsIsSorted(t *testing.T) {
	names := ListAgents()
	for i := 1; i < len(names); i++ {
		if names[i] < names[i-1] {
			t.Errorf("ListAgents() not sorted: %q before %q", names[i-1], names[i])
		}
	}
}

func TestGetAgentReturnsContent(t *testing.T) {
	for _, name := range expectedAgents {
		content, ok := GetAgent(name)
		if !ok {
			t.Errorf("GetAgent(%q) not found", name)
			continue
		}
		if content == "" {
			t.Errorf("GetAgent(%q) returned empty content", name)
			continue
		}
		// Every agent should have YAML frontmatter
		if !strings.Contains(content, "---") {
			t.Errorf("GetAgent(%q) missing YAML frontmatter delimiter", name)
		}
		if !strings.Contains(content, "name:") {
			t.Errorf("GetAgent(%q) missing name field in frontmatter", name)
		}
	}
}

func TestGetAgentNotFound(t *testing.T) {
	_, ok := GetAgent("nonexistent-agent")
	if ok {
		t.Error("GetAgent(\"nonexistent-agent\") should return false")
	}
}

func TestGetDefaultConfigIsValidJSON(t *testing.T) {
	raw := GetDefaultConfig()
	if raw == "" {
		t.Fatal("GetDefaultConfig() returned empty string")
	}

	var cfg map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &cfg); err != nil {
		t.Fatalf("GetDefaultConfig() is not valid JSON: %v", err)
	}

	// Check expected top-level keys
	for _, key := range []string{"version", "timeout_seconds", "llms", "llm_mode", "llm_modes", "discovery_pipelines"} {
		if _, ok := cfg[key]; !ok {
			t.Errorf("GetDefaultConfig() missing key %q", key)
		}
	}
}

func TestSortStrings(t *testing.T) {
	input := []string{"zebra", "apple", "mango", "banana"}
	sortStrings(input)
	want := []string{"apple", "banana", "mango", "zebra"}
	for i, v := range input {
		if v != want[i] {
			t.Errorf("sortStrings result[%d] = %q, want %q", i, v, want[i])
		}
	}
}
