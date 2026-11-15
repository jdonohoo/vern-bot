package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// Status is what vern doctor reports for one provider.
type Status struct {
	Name         string   `json:"name"`
	Label        string   `json:"label"`
	Binary       string   `json:"binary"`
	Installed    bool     `json:"installed"`
	Path         string   `json:"path,omitempty"`
	Version      string   `json:"version,omitempty"`
	SignedIn     *bool    `json:"signed_in,omitempty"` // nil when the CLI has no check
	ReadOnly     bool     `json:"read_only"`
	MissingFlags []string `json:"missing_flags,omitempty"`
	Install      string   `json:"install,omitempty"`
	Error        string   `json:"error,omitempty"`
}

// OK reports whether vern can drive this provider as configured.
func (s Status) OK() bool {
	return s.Installed && len(s.MissingFlags) == 0 && (s.SignedIn == nil || *s.SignedIn)
}

// probeTimeout bounds every doctor and models subprocess.
var probeTimeout = 20 * time.Second

// probe runs bin with args and returns combined output and the exit code.
func probe(bin string, args ...string) (string, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), probeTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return out.String(), ExitTimeout, ctx.Err()
	}
	return out.String(), exitCodeFromErr(err), nil
}

// Check inspects one provider's CLI without spending tokens: whether it is
// installed, its version, whether it is signed in (when the CLI can say),
// and whether its help still documents every flag vern passes.
func (p *Provider) Check() Status {
	st := Status{Name: p.Name, Label: p.Label, Binary: p.Binary, ReadOnly: p.ReadOnly, Install: p.Install}
	path, err := exec.LookPath(p.Binary)
	if err != nil {
		return st
	}
	st.Installed = true
	st.Path = path

	if len(p.VersionArgs) > 0 {
		out, code, _ := probe(path, p.VersionArgs...)
		if code == 0 {
			st.Version = FirstLine(out)
		}
	}
	if len(p.AuthArgs) > 0 {
		out, code, _ := probe(path, p.AuthArgs...)
		signed := code == 0 && !signedOut(out)
		st.SignedIn = &signed
	}
	if len(p.HelpArgs) > 0 && len(p.Flags) > 0 {
		out, code, perr := probe(path, p.HelpArgs...)
		if perr != nil || (code != 0 && out == "") {
			st.Error = "could not read " + p.Binary + " " + strings.Join(p.HelpArgs, " ")
		} else {
			st.MissingFlags = missingFlags(out, p.Flags)
		}
	}
	return st
}

// signedOut catches auth commands that exit 0 but say the user is signed
// out (claude auth status prints JSON with loggedIn).
func signedOut(out string) bool {
	var v struct {
		LoggedIn *bool `json:"loggedIn"`
	}
	if json.Unmarshal([]byte(strings.TrimSpace(out)), &v) == nil && v.LoggedIn != nil {
		return !*v.LoggedIn
	}
	l := strings.ToLower(out)
	return strings.Contains(l, "not logged in") || strings.Contains(l, "not signed in")
}

// missingFlags returns the flags that do not appear as whole words in help.
func missingFlags(help string, flags []string) []string {
	var missing []string
	for _, f := range flags {
		re := regexp.MustCompile(`(^|[\s,\[(])` + regexp.QuoteMeta(f) + `($|[\s,=<\[\])])`)
		if !re.MatchString(help) {
			missing = append(missing, f)
		}
	}
	return missing
}

// ListModels runs the provider's model list command. It returns the model
// ids it could parse and the raw output. ok is false when the CLI has no
// list command or is not installed.
func (p *Provider) ListModels() (ids []string, raw string, ok bool) {
	if len(p.ModelsArgs) == 0 || !p.Available() {
		return nil, "", false
	}
	out, code, err := probe(p.Binary, p.ModelsArgs...)
	if err != nil || code != 0 {
		return nil, out, false
	}
	return parseModelIDs(out), out, true
}

var slugRe = regexp.MustCompile(`"(?:slug|id)"\s*:\s*"([^"]+)"`)
var modelLineRe = regexp.MustCompile(`^\s*([A-Za-z0-9][A-Za-z0-9._:/-]*[A-Za-z0-9])\s*(\||$|\s{2,})`)

// parseModelIDs reads model ids from JSON ("slug"/"id" fields, as codex
// debug models prints) or from lines that start with an id (agy models).
func parseModelIDs(out string) []string {
	seen := map[string]bool{}
	var ids []string
	add := func(id string) {
		if !seen[id] {
			seen[id] = true
			ids = append(ids, id)
		}
	}
	trimmed := strings.TrimSpace(out)
	if strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[") {
		for _, m := range slugRe.FindAllStringSubmatch(out, -1) {
			add(m[1])
		}
		return ids
	}
	for _, line := range strings.Split(out, "\n") {
		m := modelLineRe.FindStringSubmatch(line)
		if m == nil || !strings.ContainsAny(m[1], "-.") {
			continue // headings and prose, not ids
		}
		add(m[1])
	}
	return ids
}
