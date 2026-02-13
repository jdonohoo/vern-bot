package tobeads

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jdonohoo/vern-bot/go/internal/vts"
)

// ImportOptions configures the import pipeline.
type ImportOptions struct {
	VTSDir   string
	BeadsDir string // optional: target Beads repo
	Apply    bool   // false = dry-run (default)
	Sync     bool   // trigger br sync after apply
}

// ImportResult holds the outcome of an import run.
type ImportResult struct {
	Created  int
	Skipped  int // already existed
	Failed   int
	DepsOK   int
	DepsFail int
	Errors   []string
}

// Run executes the full import pipeline: parse → normalize → preflight → execute.
func Run(opts ImportOptions, runner BrRunner) (*ImportResult, error) {
	result := &ImportResult{}

	// 1. Parse VTS files
	tasks, err := vts.ReadDir(opts.VTSDir)
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}
	fmt.Printf("Parsed %d VTS tasks from %s\n", len(tasks), opts.VTSDir)

	// 2. Preflight validation (before normalization, catches statuses early)
	report := Preflight(tasks)
	if !report.OK() {
		fmt.Print(report.String())
		return nil, fmt.Errorf("preflight failed with %d errors", len(report.Errors))
	}

	// 3. Normalize
	specs, normErrs := Normalize(tasks)
	if len(normErrs) > 0 {
		for _, e := range normErrs {
			fmt.Printf("  WARN: %s\n", e)
		}
	}
	if len(specs) == 0 {
		return nil, fmt.Errorf("no valid tasks after normalization")
	}

	// 4. Dry-run: print action plan
	if !opts.Apply {
		printDryRun(specs)
		return result, nil
	}

	// 5. Pass 1: Create issues
	idMap := map[string]string{} // VTS ID -> Beads ID
	for _, spec := range specs {
		beadID, existed, err := runner.Create(spec)
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("create %s: %v", spec.ExternalRef, err))
			fmt.Printf("  FAIL: %s — %v\n", spec.ExternalRef, err)
			continue
		}
		idMap[spec.ExternalRef] = beadID
		if existed {
			result.Skipped++
			fmt.Printf("  SKIP: %s → %s (already exists)\n", spec.ExternalRef, beadID)
		} else {
			result.Created++
			fmt.Printf("  OK:   %s → %s\n", spec.ExternalRef, beadID)
		}
	}

	// Write ID map
	mapPath := filepath.Join(opts.VTSDir, "vts-br-map.json")
	mapData, _ := json.MarshalIndent(idMap, "", "  ")
	if err := os.WriteFile(mapPath, mapData, 0644); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("write map: %v", err))
	} else {
		fmt.Printf("Wrote %s\n", mapPath)
	}

	// 6. Pass 2: Wire dependencies
	fmt.Println("Wiring dependencies...")
	for _, spec := range specs {
		fromBR, ok := idMap[spec.ExternalRef]
		if !ok {
			continue // failed create, skip
		}
		for _, depVTS := range spec.Dependencies {
			toBR, ok := idMap[depVTS]
			if !ok {
				result.DepsFail++
				result.Errors = append(result.Errors, fmt.Sprintf("dep %s→%s: target not in map", spec.ExternalRef, depVTS))
				continue
			}
			if err := runner.DepAdd(fromBR, toBR); err != nil {
				result.DepsFail++
				result.Errors = append(result.Errors, fmt.Sprintf("dep %s→%s: %v", spec.ExternalRef, depVTS, err))
				fmt.Printf("  FAIL: %s → %s: %v\n", fromBR, toBR, err)
			} else {
				result.DepsOK++
				fmt.Printf("  DEP:  %s → %s\n", fromBR, toBR)
			}
		}
	}

	// 7. Optional sync
	if opts.Sync {
		fmt.Println("Running br sync --flush-only...")
		if err := runner.Sync(); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("sync: %v", err))
			fmt.Printf("  WARN: sync failed: %v\n", err)
		} else {
			fmt.Println("  Sync complete.")
		}
	}

	// 8. Summary
	printSummary(result)

	return result, nil
}

func printDryRun(specs []BeadSpec) {
	fmt.Println("\n=== DRY RUN ===")
	fmt.Printf("Would create %d issues:\n\n", len(specs))

	depCount := 0
	for _, s := range specs {
		labels := ""
		if len(s.Labels) > 0 {
			labels = " [" + strings.Join(s.Labels, ", ") + "]"
		}
		deps := ""
		if len(s.Dependencies) > 0 {
			deps = " (deps: " + strings.Join(s.Dependencies, ", ") + ")"
			depCount += len(s.Dependencies)
		}
		fmt.Printf("  %s: %s (%s)%s%s\n", s.ExternalRef, s.Title, s.Status, labels, deps)
	}

	fmt.Printf("\nWould wire %d dependency edges.\n", depCount)
	fmt.Println("Run with --apply to execute.")
}

func printSummary(r *ImportResult) {
	fmt.Println("\n=== SUMMARY ===")
	fmt.Printf("Created: %d\n", r.Created)
	fmt.Printf("Skipped: %d (already existed)\n", r.Skipped)
	fmt.Printf("Failed:  %d\n", r.Failed)
	fmt.Printf("Deps OK: %d\n", r.DepsOK)
	fmt.Printf("Deps Failed: %d\n", r.DepsFail)
	if len(r.Errors) > 0 {
		fmt.Println("Errors:")
		for _, e := range r.Errors {
			fmt.Printf("  - %s\n", e)
		}
	}
}
