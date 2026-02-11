package council

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/jdonohoo/vern-bot/go/internal/embedded"
	"github.com/jdonohoo/vern-bot/go/internal/persona"
)

// Vern represents a persona in the roster.
type Vern struct {
	ID   string
	LLM  string
	Desc string
}

// ScanRoster builds the roster by scanning agents/*.md files.
// Falls back to embedded agent data if the directory isn't available.
// Skips vernhole-orchestrator and oracle (pipeline-only personas).
func ScanRoster(agentsDir string) []Vern {
	skip := map[string]bool{
		"vernhole-orchestrator": true,
		"oracle":               true,
		"historian":            true,
	}

	entries, err := os.ReadDir(agentsDir)
	if err != nil {
		// No agents dir on disk — use embedded data
		return scanEmbeddedRoster(skip)
	}

	var roster []Vern
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		id := strings.TrimSuffix(entry.Name(), ".md")
		if skip[id] {
			continue
		}

		p, err := persona.LoadFile(filepath.Join(agentsDir, entry.Name()))
		if err != nil {
			continue
		}

		llm := persona.ModelToLLM(p.Model)
		desc := persona.ShortDescription(p.Description)

		roster = append(roster, Vern{
			ID:   id,
			LLM:  llm,
			Desc: desc,
		})
	}

	sort.Slice(roster, func(i, j int) bool { return roster[i].ID < roster[j].ID })

	if len(roster) == 0 {
		return scanEmbeddedRoster(skip)
	}
	return roster
}

// scanEmbeddedRoster builds the roster from compiled-in agent data.
func scanEmbeddedRoster(skip map[string]bool) []Vern {
	names := embedded.ListAgents()
	var roster []Vern

	for _, name := range names {
		if skip[name] {
			continue
		}
		content, ok := embedded.GetAgent(name)
		if !ok {
			continue
		}
		p, err := persona.ParseString(content)
		if err != nil {
			continue
		}
		llm := persona.ModelToLLM(p.Model)
		desc := persona.ShortDescription(p.Description)
		roster = append(roster, Vern{ID: name, LLM: llm, Desc: desc})
	}

	if len(roster) == 0 {
		return hardcodedRoster()
	}
	return roster
}

// ResolveCouncil selects Verns based on a council tier name (or bare number).
// Returns the selected Verns and the council display name.
func ResolveCouncil(tierName string, roster []Vern, minVerns int) (selected []Vern, displayName string) {
	if minVerns <= 0 {
		minVerns = 3
	}

	tiers := AllTiers()

	// Handle bare number
	if n, err := strconv.Atoi(tierName); err == nil {
		count := clamp(n, minVerns, len(roster))
		return randomSelect(roster, count), ""
	}

	tier, ok := tiers[tierName]
	if !ok {
		fmt.Fprintf(os.Stderr, "Unknown council: %s\nValid councils: hammers, conflict, inner, round, war, full, random\n", tierName)
		os.Exit(1)
	}

	displayName = tier.Display

	switch tierName {
	case "full":
		return append([]Vern{}, roster...), displayName

	case "random", "":
		count := rand.Intn(len(roster)-minVerns+1) + minVerns
		return randomSelect(roster, count), displayName
	}

	// Build from core members
	var coreVerns []Vern
	coreSet := make(map[string]bool)
	for _, id := range tier.Core {
		if v, ok := findVern(roster, id); ok {
			coreVerns = append(coreVerns, v)
			coreSet[id] = true
		}
	}

	if tier.Fixed {
		return coreVerns, displayName
	}

	// Fill remaining slots randomly
	target := rand.Intn(tier.MaxFill-tier.MinFill+1) + tier.MinFill
	if target <= len(coreVerns) {
		return coreVerns, displayName
	}

	remaining := filterOut(roster, coreSet)
	shuffled := randomSelect(remaining, len(remaining))
	fillCount := target - len(coreVerns)
	if fillCount > len(shuffled) {
		fillCount = len(shuffled)
	}

	selected = append(coreVerns, shuffled[:fillCount]...)
	return selected, displayName
}

func findVern(roster []Vern, id string) (Vern, bool) {
	for _, v := range roster {
		if v.ID == id {
			return v, true
		}
	}
	return Vern{}, false
}

func filterOut(roster []Vern, exclude map[string]bool) []Vern {
	var result []Vern
	for _, v := range roster {
		if !exclude[v.ID] {
			result = append(result, v)
		}
	}
	return result
}

func randomSelect(roster []Vern, count int) []Vern {
	if count >= len(roster) {
		shuffled := make([]Vern, len(roster))
		copy(shuffled, roster)
		rand.Shuffle(len(shuffled), func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })
		return shuffled
	}
	indices := rand.Perm(len(roster))[:count]
	var result []Vern
	for _, i := range indices {
		result = append(result, roster[i])
	}
	return result
}

func clamp(n, min, max int) int {
	if n < min {
		return min
	}
	if n > max {
		return max
	}
	return n
}

func hardcodedRoster() []Vern {
	return []Vern{
		{ID: "academic", LLM: "claude", Desc: "needs more research, cites sources"},
		{ID: "architect", LLM: "claude", Desc: "systems design, blueprints before builds"},
		{ID: "enterprise", LLM: "claude", Desc: "needs 6 meetings and a committee"},
		{ID: "great", LLM: "claude", Desc: "excellence incarnate, elegant solutions"},
		{ID: "inverse", LLM: "claude", Desc: "contrarian takes only"},
		{ID: "ketamine", LLM: "claude", Desc: "multi-dimensional vibes, sees patterns"},
		{ID: "mediocre", LLM: "claude", Desc: "scrappy speed demon, ship it fast"},
		{ID: "mighty", LLM: "claude", Desc: "raw power, comprehensive analysis"},
		{ID: "nyquil", LLM: "claude", Desc: "brilliant but brief, NyQuil kicking in"},
		{ID: "optimist", LLM: "claude", Desc: "everything will be fine!"},
		{ID: "paranoid", LLM: "claude", Desc: "what could possibly go wrong?"},
		{ID: "retro", LLM: "claude", Desc: "we solved this with cron in 2004"},
		{ID: "startup", LLM: "claude", Desc: "MVP or die, move fast"},
		{ID: "ux", LLM: "claude", Desc: "can the user find the button?"},
		{ID: "yolo", LLM: "claude", Desc: "full send, no guardrails"},
	}
}
