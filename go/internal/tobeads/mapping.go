package tobeads

// Status mapping: VTS status → Beads status.
// Unknown statuses cause a hard preflight failure.
var StatusMap = map[string]string{
	"pending":  "open",
	"active":   "in_progress",
	"blocked":  "blocked",
	"done":     "closed",
	"complete": "closed",
	"deferred": "deferred",
}

// Valid complexity codes. Mapped to labels like "complexity:XS".
var ValidComplexity = map[string]bool{
	"XS": true,
	"S":  true,
	"M":  true,
	"L":  true,
	"XL": true,
}

// BeadSpec is a normalized task ready for Beads creation.
type BeadSpec struct {
	ExternalRef  string   // VTS ID (e.g. "VTS-001")
	Title        string
	Description  string   // Body + files + source_ref metadata + acceptance criteria
	Status       string   // Beads status (mapped)
	Labels       []string // ["complexity:M", "source:oracle"]
	Assignee     string   // From owner, empty = unset
	Dependencies []string // VTS IDs (resolved to Beads IDs in executor)
}
