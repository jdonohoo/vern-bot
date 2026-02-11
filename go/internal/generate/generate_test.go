package generate

import (
	"strings"
	"testing"
)

func TestParseOutput(t *testing.T) {
	output := `Some preamble text

=== AGENT ===
---
name: test
description: Test Vern - A test persona for validation.
model: sonnet
color: blue
---

You are Test Vern.

PERSONALITY:
- Testing things

CATCHPHRASES:
- "Testing, testing, 1-2-3"

SIGN-OFF:
End with a test joke.
=== END AGENT ===

=== COMMAND ===
---
description: Test Vern - A test persona for validation.
argument-hint: [task]
---

# Test Vern

Test this: $ARGUMENTS
=== END COMMAND ===

=== SKILL ===
---
name: test
description: Test Vern - A test persona for validation.
argument-hint: [task]
---

# Test Vern

Test this: $ARGUMENTS
=== END SKILL ===

Here's a dad joke: Why did the test fail? Because it had too many assertions!
`

	gen, err := ParseOutput("test", output)
	if err != nil {
		t.Fatalf("ParseOutput() error: %v", err)
	}

	if gen.Name != "test" {
		t.Errorf("Name = %q, want %q", gen.Name, "test")
	}
	if !strings.Contains(gen.Agent, "name: test") {
		t.Error("Agent missing 'name: test'")
	}
	if !strings.Contains(gen.Command, "description: Test Vern") {
		t.Error("Command missing description")
	}
	if !strings.Contains(gen.Skill, "name: test") {
		t.Error("Skill missing 'name: test'")
	}
}

func TestParseOutput_MissingSection(t *testing.T) {
	output := `=== AGENT ===
---
name: test
description: Test
model: sonnet
color: blue
---
Body
=== END AGENT ===

=== COMMAND ===
---
description: Test
argument-hint: [task]
---
Body
=== END COMMAND ===
`

	_, err := ParseOutput("test", output)
	if err == nil {
		t.Fatal("expected error for missing SKILL section")
	}
	if !strings.Contains(err.Error(), "skill section") {
		t.Errorf("error = %q, want mention of skill section", err.Error())
	}
}

func TestParseOutput_TrailingContent(t *testing.T) {
	output := `=== AGENT ===
---
name: test
description: Test Vern - Testing.
model: sonnet
color: blue
---

You are Test Vern.

PERSONALITY:
- Testing
=== END AGENT ===

=== COMMAND ===
---
description: Test Vern - Testing.
argument-hint: [task]
---

# Test Vern

Test: $ARGUMENTS
=== END COMMAND ===

=== SKILL ===
---
name: test
description: Test Vern - Testing.
argument-hint: [task]
---

# Test Vern

Test: $ARGUMENTS
=== END SKILL ===

This trailing content (like a dad joke) should be silently ignored.
Why did the test pass? Because we wrote it that way!

-- Test Vern (assertions are just vibes)
`

	gen, err := ParseOutput("test", output)
	if err != nil {
		t.Fatalf("ParseOutput() error: %v", err)
	}
	if strings.Contains(gen.Skill, "trailing") {
		t.Error("Skill should not contain trailing content")
	}
}

func TestParseOutput_EmptySection(t *testing.T) {
	output := `=== AGENT ===

=== END AGENT ===

=== COMMAND ===
---
description: Test
argument-hint: [task]
---
Body
=== END COMMAND ===

=== SKILL ===
---
name: test
description: Test
argument-hint: [task]
---
Body
=== END SKILL ===
`

	_, err := ParseOutput("test", output)
	if err == nil {
		t.Fatal("expected error for empty AGENT section")
	}
}

func TestBuildPrompt(t *testing.T) {
	prompt := BuildPrompt("nihilist", "Nothing matters, existential code review", "", "")
	if !strings.Contains(prompt, "nihilist") {
		t.Error("prompt missing name")
	}
	if !strings.Contains(prompt, "Nothing matters") {
		t.Error("prompt missing description")
	}
	if !strings.Contains(prompt, "Choose the model") {
		t.Error("prompt should say 'Choose the model' when no model override")
	}
}

func TestBuildPrompt_WithOverrides(t *testing.T) {
	prompt := BuildPrompt("nihilist", "Nothing matters", "opus", "gray")
	if !strings.Contains(prompt, "Model override") {
		t.Error("prompt missing model override directive")
	}
	if !strings.Contains(prompt, "opus") {
		t.Error("prompt missing opus model")
	}
	if !strings.Contains(prompt, "Color override") {
		t.Error("prompt missing color override directive")
	}
	if !strings.Contains(prompt, "gray") {
		t.Error("prompt missing gray color")
	}
}

func TestValidateName(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"nihilist", false},
		{"code-poet", false},
		{"test123", false},
		{"", true},            // empty
		{"UPPERCASE", true},   // not lowercase
		{"has space", true},   // spaces
		{"has.dot", true},     // dots
		{"123start", true},    // starts with number
		{"-hyphen", true},     // starts with hyphen
		{"mediocre", true},    // conflicts with existing alias
		{"hole", true},        // conflicts with existing workflow
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateName(tt.name, "")
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateName(%q) error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		})
	}
}

func TestComputeAliases(t *testing.T) {
	existing := KnownAliases()

	aliases := ComputeAliases("nihilist", existing)
	if aliases[0] != "nihilist" {
		t.Errorf("first alias should be full name, got %q", aliases[0])
	}
	if len(aliases) < 2 {
		t.Fatal("expected at least 2 aliases")
	}
	if aliases[1] != "nih" {
		t.Errorf("second alias should be 'nih', got %q", aliases[1])
	}

	// Test with a name whose 3-char prefix conflicts
	aliases = ComputeAliases("archive", existing)
	// "arc" conflicts with "arch"? No, "arc" != "arch". Let's check.
	// Actually "arch" is in existing, "arc" is not.
	found3 := false
	for _, a := range aliases {
		if a == "arc" {
			found3 = true
		}
	}
	if !found3 {
		t.Error("expected 'arc' alias for 'archive'")
	}
}

func TestComputeAliases_ConflictingPrefix(t *testing.T) {
	existing := map[string]bool{
		"test":    true,
		"tes":     true,
		"testing": true,
	}

	aliases := ComputeAliases("testing", existing)
	// "testing" itself is in existing, but that's the full name
	// "tes" is also in existing, so it should NOT be added
	for _, a := range aliases {
		if a == "tes" {
			t.Error("should not include conflicting prefix 'tes'")
		}
	}
}
