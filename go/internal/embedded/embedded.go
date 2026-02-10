// Package embedded provides compiled-in agent personas and default config
// so the standalone binary works without the agents/ directory on disk.
//
// Regenerate after changing agents/*.md or config.default.json:
//
//	cd go && go generate ./internal/embedded/
package embedded

//go:generate go run gen.go

// GetAgent returns the full markdown content for a persona by name.
// Returns empty string and false if not found.
func GetAgent(name string) (string, bool) {
	content, ok := AgentData[name]
	return content, ok
}

// ListAgents returns the sorted names of all embedded agents.
func ListAgents() []string {
	names := make([]string, 0, len(AgentData))
	for name := range AgentData {
		names = append(names, name)
	}
	// AgentData is generated in sorted order, but map iteration isn't
	// deterministic, so sort here.
	sortStrings(names)
	return names
}

// GetDefaultConfig returns the embedded config.default.json content.
func GetDefaultConfig() string {
	return DefaultConfigJSON
}

func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}
