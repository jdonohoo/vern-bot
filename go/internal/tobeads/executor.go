package tobeads

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// BrRunner abstracts br CLI operations for testability.
type BrRunner interface {
	Create(spec BeadSpec) (beadID string, alreadyExisted bool, err error)
	DepAdd(fromID, toID string) error
	Sync() error
}

// RealBrRunner shells out to the br CLI binary.
type RealBrRunner struct {
	WorkDir string // --beads-dir: sets exec.Cmd.Dir
}

func (r *RealBrRunner) Create(spec BeadSpec) (string, bool, error) {
	args := []string{"create",
		"--title", spec.Title,
		"--status", spec.Status,
		"--external-ref", spec.ExternalRef,
		"--json",
	}

	if spec.Description != "" {
		args = append(args, "--description", spec.Description)
	}

	if len(spec.Labels) > 0 {
		args = append(args, "--labels", strings.Join(spec.Labels, ","))
	}

	if spec.Assignee != "" {
		args = append(args, "--assignee", spec.Assignee)
	}

	stdout, stderr, err := r.runBr(args...)
	if err != nil {
		// Check for duplicate external_ref error
		if strings.Contains(stderr, "external_ref") || strings.Contains(stderr, "already exists") || strings.Contains(stderr, "unique") {
			// Try to find existing bead by external ref
			existingID, lookupErr := r.lookupByExternalRef(spec.ExternalRef)
			if lookupErr == nil && existingID != "" {
				return existingID, true, nil
			}
			return "", false, fmt.Errorf("duplicate external_ref %s: %s", spec.ExternalRef, stderr)
		}
		return "", false, fmt.Errorf("br create failed: %s (stderr: %s)", err, stderr)
	}

	// Parse JSON output for Beads ID
	beadID, parseErr := parseCreateJSON(stdout)
	if parseErr != nil {
		return "", false, fmt.Errorf("parse br create output: %w", parseErr)
	}

	return beadID, false, nil
}

func (r *RealBrRunner) DepAdd(fromID, toID string) error {
	_, stderr, err := r.runBr("dep", "add", fromID, toID)
	if err != nil {
		return fmt.Errorf("br dep add %s %s: %s (stderr: %s)", fromID, toID, err, stderr)
	}
	return nil
}

func (r *RealBrRunner) Sync() error {
	_, stderr, err := r.runBr("sync", "--flush-only")
	if err != nil {
		return fmt.Errorf("br sync: %s (stderr: %s)", err, stderr)
	}
	return nil
}

func (r *RealBrRunner) lookupByExternalRef(extRef string) (string, error) {
	stdout, _, err := r.runBr("list", "--json")
	if err != nil {
		return "", err
	}

	var items []map[string]interface{}
	if err := json.Unmarshal([]byte(stdout), &items); err != nil {
		return "", err
	}

	for _, item := range items {
		if ref, ok := item["external_ref"].(string); ok && ref == extRef {
			if id, ok := item["id"].(string); ok {
				return id, nil
			}
		}
	}
	return "", fmt.Errorf("not found")
}

func (r *RealBrRunner) runBr(args ...string) (stdout, stderr string, err error) {
	cmd := exec.Command("br", args...)
	if r.WorkDir != "" {
		cmd.Dir = r.WorkDir
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err = cmd.Run()
	return stdoutBuf.String(), stderrBuf.String(), err
}

// parseCreateJSON extracts the Beads issue ID from br create --json output.
func parseCreateJSON(output string) (string, error) {
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		return "", fmt.Errorf("invalid JSON: %w", err)
	}

	// Try common field names for the ID
	for _, key := range []string{"id", "ID", "issue_id", "bead_id"} {
		if v, ok := result[key]; ok {
			switch id := v.(type) {
			case string:
				return id, nil
			case float64:
				return fmt.Sprintf("%d", int(id)), nil
			}
		}
	}

	return "", fmt.Errorf("no id field found in JSON output: %s", output)
}

// MockBrRunner is a test double for BrRunner.
type MockBrRunner struct {
	Created    []BeadSpec
	Deps       [][2]string
	Synced     bool
	NextID     int
	ExistingRefs map[string]string // external_ref -> bead ID (simulates already-created)
	FailCreate map[string]error   // external_ref -> error
	FailDep    error
}

func NewMockBrRunner() *MockBrRunner {
	return &MockBrRunner{
		NextID:       1,
		ExistingRefs: map[string]string{},
		FailCreate:   map[string]error{},
	}
}

func (m *MockBrRunner) Create(spec BeadSpec) (string, bool, error) {
	if err, ok := m.FailCreate[spec.ExternalRef]; ok {
		return "", false, err
	}
	if id, ok := m.ExistingRefs[spec.ExternalRef]; ok {
		return id, true, nil
	}
	id := fmt.Sprintf("BR-%03d", m.NextID)
	m.NextID++
	m.Created = append(m.Created, spec)
	m.ExistingRefs[spec.ExternalRef] = id
	return id, false, nil
}

func (m *MockBrRunner) DepAdd(fromID, toID string) error {
	if m.FailDep != nil {
		return m.FailDep
	}
	m.Deps = append(m.Deps, [2]string{fromID, toID})
	return nil
}

func (m *MockBrRunner) Sync() error {
	m.Synced = true
	return nil
}
