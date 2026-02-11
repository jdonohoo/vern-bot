package generate

import (
	"fmt"
	"strings"

	"github.com/jdonohoo/vern-bot/go/internal/persona"
)

// GeneratedPersona holds the parsed output from the LLM.
type GeneratedPersona struct {
	Name    string
	Agent   string // Full content for agents/{name}.md
	Command string // Full content for commands/{name}.md
	Skill   string // Full content for skills/{name}/SKILL.md
}

// ParseOutput extracts the 3 delimited sections from LLM output.
// Content after the last === END SKILL === (e.g. dad jokes) is ignored.
func ParseOutput(name, output string) (*GeneratedPersona, error) {
	agent, err := extractSection(output, "AGENT")
	if err != nil {
		return nil, fmt.Errorf("agent section: %w", err)
	}
	command, err := extractSection(output, "COMMAND")
	if err != nil {
		return nil, fmt.Errorf("command section: %w", err)
	}
	skill, err := extractSection(output, "SKILL")
	if err != nil {
		return nil, fmt.Errorf("skill section: %w", err)
	}

	// Validate agent frontmatter parses correctly
	p, err := persona.ParseString(agent)
	if err != nil {
		return nil, fmt.Errorf("agent frontmatter parse: %w", err)
	}
	if p.Name == "" {
		return nil, fmt.Errorf("agent frontmatter missing name field")
	}
	if p.Model == "" {
		return nil, fmt.Errorf("agent frontmatter missing model field")
	}

	return &GeneratedPersona{
		Name:    name,
		Agent:   agent,
		Command: command,
		Skill:   skill,
	}, nil
}

// extractSection pulls content between === SECTION === and === END SECTION === delimiters.
func extractSection(output, section string) (string, error) {
	startDelim := "=== " + section + " ==="
	endDelim := "=== END " + section + " ==="

	startIdx := strings.Index(output, startDelim)
	if startIdx == -1 {
		return "", fmt.Errorf("missing start delimiter %q", startDelim)
	}

	contentStart := startIdx + len(startDelim)
	// Skip the newline after the delimiter
	if contentStart < len(output) && output[contentStart] == '\n' {
		contentStart++
	}

	remaining := output[contentStart:]
	endIdx := strings.Index(remaining, endDelim)
	if endIdx == -1 {
		return "", fmt.Errorf("missing end delimiter %q", endDelim)
	}

	content := strings.TrimRight(remaining[:endIdx], "\n")
	if strings.TrimSpace(content) == "" {
		return "", fmt.Errorf("section %q is empty", section)
	}

	return content + "\n", nil
}
