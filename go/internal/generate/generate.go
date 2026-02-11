package generate

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/jdonohoo/vern-bot/go/internal/llm"
	"github.com/jdonohoo/vern-bot/go/internal/persona"
)

// Options configures the persona generation.
type Options struct {
	Name        string
	Description string
	Model       string // optional: opus/sonnet/haiku
	Color       string // optional: lipgloss color
	LLM         string // LLM for generation (default: claude)
	DryRun      bool   // print without writing
	NoUpdate    bool   // skip v.md/help/test/roster updates
	RepoRoot    string // detected repo root
	LogFunc     func(string) // optional callback for progress logs
}

func (o Options) log(msg string) {
	if o.LogFunc != nil {
		o.LogFunc(msg)
	} else {
		fmt.Println(msg)
	}
}

// ValidateName checks that a persona name is valid and doesn't conflict with existing agents.
func ValidateName(name, repoRoot string) error {
	if name == "" {
		return fmt.Errorf("name is required")
	}

	// Must be lowercase alphanumeric + hyphens
	validName := regexp.MustCompile(`^[a-z][a-z0-9-]*$`)
	if !validName.MatchString(name) {
		return fmt.Errorf("name must be lowercase alphanumeric with hyphens (e.g. 'nihilist', 'code-poet')")
	}

	if len(name) > 30 {
		return fmt.Errorf("name must be 30 characters or fewer")
	}

	// Check for conflict with existing agents
	if repoRoot != "" {
		agentPath := filepath.Join(repoRoot, "agents", name+".md")
		if _, err := os.Stat(agentPath); err == nil {
			return fmt.Errorf("persona %q already exists at %s", name, agentPath)
		}
	}

	// Check for alias conflicts
	known := KnownAliases()
	if known[name] {
		return fmt.Errorf("name %q conflicts with an existing alias", name)
	}

	return nil
}

// Run executes the full persona generation workflow.
func Run(opts Options) error {
	if opts.LLM == "" {
		opts.LLM = "claude"
	}

	// 1. Validate
	if err := ValidateName(opts.Name, opts.RepoRoot); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}

	opts.log(fmt.Sprintf("Generating persona %q...", opts.Name))
	opts.log(fmt.Sprintf("  LLM: %s | Timeout: 5m", opts.LLM))

	// 2. Build prompt
	prompt := BuildPrompt(opts.Name, opts.Description, opts.Model, opts.Color)

	// 3. Call LLM
	opts.log("\nWaiting for LLM response...")

	result, err := llm.Run(llm.RunOptions{
		Ctx:     context.Background(),
		LLM:     opts.LLM,
		Prompt:  prompt,
		Timeout: 5 * time.Minute,
	})
	if err != nil {
		return fmt.Errorf("LLM call failed: %w", err)
	}
	if result.Output == "" {
		return fmt.Errorf("LLM returned empty output")
	}

	opts.log(fmt.Sprintf("  Response received (%.1fs, %d chars)", result.Duration.Seconds(), len(result.Output)))

	// 4. Parse output
	opts.log("\nParsing output...")
	gen, err := ParseOutput(opts.Name, result.Output)
	if err != nil {
		return fmt.Errorf("parse output: %w", err)
	}

	// Extract model and color from the generated agent for display/updates
	p, _ := persona.ParseString(gen.Agent)
	modelName := "sonnet"
	colorName := ""
	if p != nil {
		if p.Model != "" {
			modelName = p.Model
		}
		if p.Color != "" {
			colorName = p.Color
		}
	}

	opts.log(fmt.Sprintf("  Agent:   OK (model=%s, color=%s)", modelName, colorName))
	opts.log("  Command: OK")
	opts.log("  Skill:   OK")

	// 5. Dry run: print and exit
	if opts.DryRun {
		opts.log("\n--- DRY RUN: Generated content ---")
		opts.log("\n=== agents/" + opts.Name + ".md ===")
		opts.log(gen.Agent)
		opts.log("\n=== commands/" + opts.Name + ".md ===")
		opts.log(gen.Command)
		opts.log("\n=== skills/" + opts.Name + "/SKILL.md ===")
		opts.log(gen.Skill)
		return nil
	}

	// 6. Write files
	opts.log("\nWriting files...")

	agentPath := filepath.Join(opts.RepoRoot, "agents", opts.Name+".md")
	if err := os.WriteFile(agentPath, []byte(gen.Agent), 0644); err != nil {
		return fmt.Errorf("write agent file: %w", err)
	}
	opts.log(fmt.Sprintf("  Created %s", agentPath))

	commandPath := filepath.Join(opts.RepoRoot, "commands", opts.Name+".md")
	if err := os.WriteFile(commandPath, []byte(gen.Command), 0644); err != nil {
		return fmt.Errorf("write command file: %w", err)
	}
	opts.log(fmt.Sprintf("  Created %s", commandPath))

	skillDir := filepath.Join(opts.RepoRoot, "skills", opts.Name)
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		return fmt.Errorf("create skill dir: %w", err)
	}
	skillPath := filepath.Join(skillDir, "SKILL.md")
	if err := os.WriteFile(skillPath, []byte(gen.Skill), 0644); err != nil {
		return fmt.Errorf("write skill file: %w", err)
	}
	opts.log(fmt.Sprintf("  Created %s", skillPath))

	// 7. Update registrations
	if !opts.NoUpdate {
		opts.log("\nUpdating registrations...")

		aliases := ComputeAliases(opts.Name, KnownAliases())
		shortDesc := persona.ShortDescription(p.Description)

		// v.md
		vmdPath := filepath.Join(opts.RepoRoot, "commands", "v.md")
		if err := UpdateVMD(vmdPath, opts.Name, shortDesc, aliases); err != nil {
			opts.log(fmt.Sprintf("  Warning: v.md update failed: %v", err))
		} else {
			shortAlias := opts.Name
			if len(aliases) > 1 {
				shortAlias = aliases[1]
			}
			opts.log(fmt.Sprintf("  Updated commands/v.md (aliases: %s)", shortAlias))
		}

		// help.md
		helpPath := filepath.Join(opts.RepoRoot, "commands", "help.md")
		shortAlias := opts.Name
		if len(aliases) > 1 {
			shortAlias = aliases[1]
		}
		if err := UpdateHelpMD(helpPath, opts.Name, modelName, shortDesc, shortAlias); err != nil {
			opts.log(fmt.Sprintf("  Warning: help.md update failed: %v", err))
		} else {
			opts.log("  Updated commands/help.md")
		}

		// embedded_test.go
		testPath := filepath.Join(opts.RepoRoot, "go", "internal", "embedded", "embedded_test.go")
		if err := UpdateEmbeddedTest(testPath, opts.Name); err != nil {
			opts.log(fmt.Sprintf("  Warning: embedded_test.go update failed: %v", err))
		} else {
			opts.log("  Updated go/internal/embedded/embedded_test.go")
		}

		// hardcodedRoster
		selectionPath := filepath.Join(opts.RepoRoot, "go", "internal", "council", "selection.go")
		if err := UpdateHardcodedRoster(selectionPath, opts.Name, shortDesc); err != nil {
			opts.log(fmt.Sprintf("  Warning: selection.go update failed: %v", err))
		} else {
			opts.log("  Updated go/internal/council/selection.go")
		}

		// 8. Regenerate embedded assets
		opts.log("\nRegenerating embedded assets...")
		goDir := filepath.Join(opts.RepoRoot, "go")
		cmd := exec.Command("go", "generate", "./internal/embedded/")
		cmd.Dir = goDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			opts.log(fmt.Sprintf("  Warning: go generate failed: %v", err))
			opts.log("  Run manually: cd go && go generate ./internal/embedded/")
		} else {
			opts.log("  Done")
		}
	}

	opts.log(fmt.Sprintf("\nPersona %q is ready.", opts.Name))
	return nil
}

// DetectRepoRoot walks up from cwd looking for agents/ + go/ coexisting.
func DetectRepoRoot() (string, error) {
	// Check VERN_ROOT env var first
	if root := os.Getenv("VERN_ROOT"); root != "" {
		if hasRepoMarkers(root) {
			return root, nil
		}
	}

	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get working directory: %w", err)
	}

	for {
		if hasRepoMarkers(dir) {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("could not find vern-bot repo root (looking for agents/ + go/ dirs)")
}

func hasRepoMarkers(dir string) bool {
	agentsInfo, err := os.Stat(filepath.Join(dir, "agents"))
	if err != nil || !agentsInfo.IsDir() {
		return false
	}
	goInfo, err := os.Stat(filepath.Join(dir, "go"))
	if err != nil || !goInfo.IsDir() {
		return false
	}
	return true
}

// ExtractModel extracts the model from generated agent content.
func ExtractModel(agentContent string) string {
	p, err := persona.ParseString(agentContent)
	if err != nil {
		return "sonnet"
	}
	if p.Model == "" {
		return "sonnet"
	}
	return p.Model
}

// ExtractColor extracts the color from generated agent content.
func ExtractColor(agentContent string) string {
	p, err := persona.ParseString(agentContent)
	if err != nil {
		return ""
	}
	return p.Color
}

// ExtractShortDesc extracts the short description from generated agent content.
func ExtractShortDesc(agentContent string) string {
	p, err := persona.ParseString(agentContent)
	if err != nil {
		return ""
	}
	return strings.TrimSuffix(persona.ShortDescription(p.Description), ".")
}
