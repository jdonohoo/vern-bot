package persona

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jdonohoo/vern-bot/go/internal/embedded"
)

// Persona represents a parsed agent file.
type Persona struct {
	Name        string
	Description string
	Model       string
	Color       string
	Body        string // Everything after the YAML frontmatter
}

// Load reads an agent file and parses its YAML frontmatter and body.
// agentsDir is the path to the agents/ directory, name is the persona ID (e.g. "mighty").
// Falls back to embedded agent data if the file doesn't exist on disk.
func Load(agentsDir, name string) (*Persona, error) {
	path := filepath.Join(agentsDir, name+".md")
	p, err := LoadFile(path)
	if err == nil {
		return p, nil
	}

	// Fall back to embedded data
	return LoadEmbedded(name)
}

// LoadFile reads and parses a single agent markdown file.
func LoadFile(path string) (*Persona, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open agent file: %w", err)
	}
	defer f.Close()

	p := &Persona{}
	scanner := bufio.NewScanner(f)

	// State machine: look for opening ---, parse frontmatter, then collect body
	state := 0 // 0=before frontmatter, 1=in frontmatter, 2=body
	var bodyLines []string

	for scanner.Scan() {
		line := scanner.Text()

		switch state {
		case 0:
			if strings.TrimSpace(line) == "---" {
				state = 1
			}
		case 1:
			if strings.TrimSpace(line) == "---" {
				state = 2
				continue
			}
			// Simple YAML key: value parsing
			if idx := strings.Index(line, ":"); idx > 0 {
				key := strings.TrimSpace(line[:idx])
				val := strings.TrimSpace(line[idx+1:])
				switch key {
				case "name":
					p.Name = val
				case "description":
					p.Description = val
				case "model":
					p.Model = val
				case "color":
					p.Color = val
				}
			}
		case 2:
			bodyLines = append(bodyLines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read agent file: %w", err)
	}

	p.Body = strings.Join(bodyLines, "\n")

	// Default model to claude
	if p.Model == "" {
		p.Model = "claude"
	}

	return p, nil
}

// LoadEmbedded loads a persona from the compiled-in embedded agent data.
func LoadEmbedded(name string) (*Persona, error) {
	content, ok := embedded.GetAgent(name)
	if !ok {
		return nil, fmt.Errorf("agent %q not found in embedded data", name)
	}
	return ParseString(content)
}

// ParseString parses persona markdown content from a string (same format as .md files).
func ParseString(content string) (*Persona, error) {
	p := &Persona{}
	state := 0
	var bodyLines []string

	for _, line := range strings.Split(content, "\n") {
		switch state {
		case 0:
			if strings.TrimSpace(line) == "---" {
				state = 1
			}
		case 1:
			if strings.TrimSpace(line) == "---" {
				state = 2
				continue
			}
			if idx := strings.Index(line, ":"); idx > 0 {
				key := strings.TrimSpace(line[:idx])
				val := strings.TrimSpace(line[idx+1:])
				switch key {
				case "name":
					p.Name = val
				case "description":
					p.Description = val
				case "model":
					p.Model = val
				case "color":
					p.Color = val
				}
			}
		case 2:
			bodyLines = append(bodyLines, line)
		}
	}

	p.Body = strings.Join(bodyLines, "\n")
	if p.Model == "" {
		p.Model = "claude"
	}
	return p, nil
}

// ModelToLLM maps a model identifier from agent frontmatter to the LLM engine name.
func ModelToLLM(model string) string {
	switch strings.ToLower(model) {
	case "opus", "sonnet", "haiku":
		return "claude"
	default:
		return "claude"
	}
}

// ShortDescription extracts the short descriptor from a full description string.
// E.g. "MightyVern / Codex Vern - Raw computational power." -> "Raw computational power"
func ShortDescription(desc string) string {
	if idx := strings.Index(desc, " - "); idx >= 0 {
		short := desc[idx+3:]
		if dot := strings.Index(short, ". "); dot >= 0 {
			return short[:dot]
		}
		return strings.TrimSuffix(short, ".")
	}
	if dot := strings.Index(desc, ". "); dot >= 0 {
		return desc[:dot]
	}
	return strings.TrimSuffix(desc, ".")
}
