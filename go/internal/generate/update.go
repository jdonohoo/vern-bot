package generate

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

// UpdateVMD inserts a new persona into commands/v.md's Specialist Personas table and Routing section.
func UpdateVMD(path, name, desc string, aliases []string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read v.md: %w", err)
	}
	content := string(data)

	aliasStr := strings.Join(aliases, "` / `")
	aliasStr = "`" + aliasStr + "`"

	// Insert table row into Specialist Personas section (alphabetical)
	tableRow := fmt.Sprintf("| %s | `/vern:%s` | %s |", aliasStr, name, desc)
	content, err = insertSortedTableRow(content, "### Specialist Personas", tableRow, name)
	if err != nil {
		return fmt.Errorf("insert table row: %w", err)
	}

	// Insert routing entry (alphabetical among specialist personas)
	routingAliases := strings.Join(aliases, "` / `")
	routingEntry := fmt.Sprintf("- `%s` → invoke `/vern:%s`", routingAliases, name)
	content, err = insertSortedRoutingEntry(content, routingEntry, name)
	if err != nil {
		return fmt.Errorf("insert routing entry: %w", err)
	}

	return os.WriteFile(path, []byte(content), 0644)
}

// insertSortedTableRow inserts a markdown table row into the section starting with sectionHeader,
// maintaining alphabetical order by the persona name.
func insertSortedTableRow(content, sectionHeader, newRow, name string) (string, error) {
	headerIdx := strings.Index(content, sectionHeader)
	if headerIdx == -1 {
		return "", fmt.Errorf("section %q not found", sectionHeader)
	}

	// Find the table rows after the header (skip the header line and the |---|---| line)
	afterHeader := content[headerIdx:]
	lines := strings.Split(afterHeader, "\n")

	// Find table start (first | line after header)
	tableStart := -1
	for i, line := range lines {
		if i == 0 {
			continue // skip header line
		}
		if strings.HasPrefix(strings.TrimSpace(line), "|") {
			tableStart = i
			break
		}
	}
	if tableStart == -1 {
		return "", fmt.Errorf("no table found after %q", sectionHeader)
	}

	// Collect table rows (skip header row and separator row)
	var dataRows []string
	insertionPoint := tableStart + 2 // after header row + separator
	for i := insertionPoint; i < len(lines); i++ {
		trimmed := strings.TrimSpace(lines[i])
		if !strings.HasPrefix(trimmed, "|") {
			break
		}
		dataRows = append(dataRows, lines[i])
	}

	// Add new row and sort
	dataRows = append(dataRows, newRow)
	sort.Slice(dataRows, func(i, j int) bool {
		return extractPersonaName(dataRows[i]) < extractPersonaName(dataRows[j])
	})

	// Rebuild
	var result []string
	result = append(result, lines[:insertionPoint]...)
	result = append(result, dataRows...)
	result = append(result, lines[insertionPoint+len(dataRows)-1:]...)

	return content[:headerIdx] + strings.Join(result, "\n"), nil
}

// extractPersonaName pulls the persona name from a table row like "| `architect` / `arch` / `ar` | `/vern:architect` | ... |"
func extractPersonaName(row string) string {
	// Look for /vern:NAME pattern
	re := regexp.MustCompile(`/vern:([a-z-]+)`)
	m := re.FindStringSubmatch(row)
	if len(m) > 1 {
		return m[1]
	}
	return row
}

// insertSortedRoutingEntry inserts a routing line in alphabetical order.
func insertSortedRoutingEntry(content, newEntry, name string) (string, error) {
	// Find the Routing section
	routingIdx := strings.Index(content, "## Routing")
	if routingIdx == -1 {
		return "", fmt.Errorf("## Routing section not found")
	}

	afterRouting := content[routingIdx:]
	lines := strings.Split(afterRouting, "\n")

	// Find the block of routing entries (lines starting with "- `")
	var routeLines []string
	routeStart := -1
	routeEnd := -1
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "- `") {
			if routeStart == -1 {
				routeStart = i
			}
			routeEnd = i + 1
			routeLines = append(routeLines, line)
		} else if routeStart != -1 && trimmed == "" {
			// Allow blank lines within routing block
			continue
		} else if routeStart != -1 {
			break
		}
	}

	if routeStart == -1 {
		return "", fmt.Errorf("no routing entries found")
	}

	// Add new entry and sort
	routeLines = append(routeLines, newEntry)
	sort.Slice(routeLines, func(i, j int) bool {
		return extractRouteName(routeLines[i]) < extractRouteName(routeLines[j])
	})

	// Rebuild
	var result []string
	result = append(result, lines[:routeStart]...)
	result = append(result, routeLines...)
	result = append(result, lines[routeEnd:]...)

	return content[:routingIdx] + strings.Join(result, "\n"), nil
}

// extractRouteName pulls the first alias from a routing entry like "- `architect` / `arch` / `ar` → invoke `/vern:architect`"
func extractRouteName(line string) string {
	re := regexp.MustCompile("→ invoke `/vern:([a-z-]+)`")
	m := re.FindStringSubmatch(line)
	if len(m) > 1 {
		return m[1]
	}
	return line
}

// UpdateHelpMD inserts a persona line into commands/help.md's SPECIALIST PERSONAS section
// and adds its short alias to the aliases line.
func UpdateHelpMD(path, name, model, shortDesc, shortAlias string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read help.md: %w", err)
	}
	content := string(data)

	// Map model to display name
	modelDisplay := strings.Title(model)
	switch strings.ToLower(model) {
	case "opus":
		modelDisplay = "Opus   "
	case "sonnet":
		modelDisplay = "Sonnet "
	case "haiku":
		modelDisplay = "Haiku  "
	}

	// Build the new line matching existing format:
	// "  /vern:%-17s %-8s- %s"
	newLine := fmt.Sprintf("  /vern:%-17s%-8s- %s", name+" <task>", modelDisplay, shortDesc)

	// Find SPECIALIST PERSONAS section and insert alphabetically
	specIdx := strings.Index(content, "SPECIALIST PERSONAS")
	if specIdx == -1 {
		return fmt.Errorf("SPECIALIST PERSONAS section not found in help.md")
	}

	lines := strings.Split(content, "\n")
	var specLines []int // indices of specialist persona lines
	inSpec := false
	for i, line := range lines {
		if strings.Contains(line, "SPECIALIST PERSONAS") {
			inSpec = true
			continue
		}
		if inSpec {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "/vern:") {
				specLines = append(specLines, i)
			} else if trimmed != "" && !strings.HasPrefix(trimmed, "/vern:") {
				break
			}
		}
	}

	if len(specLines) == 0 {
		return fmt.Errorf("no specialist persona lines found")
	}

	// Collect existing lines, add new one, sort
	var existingLines []string
	for _, idx := range specLines {
		existingLines = append(existingLines, lines[idx])
	}
	existingLines = append(existingLines, newLine)
	sort.Slice(existingLines, func(i, j int) bool {
		return extractHelpName(existingLines[i]) < extractHelpName(existingLines[j])
	})

	// Replace the old lines with the new sorted set
	firstIdx := specLines[0]
	lastIdx := specLines[len(specLines)-1]

	var newLines []string
	newLines = append(newLines, lines[:firstIdx]...)
	newLines = append(newLines, existingLines...)
	newLines = append(newLines, lines[lastIdx+1:]...)
	content = strings.Join(newLines, "\n")

	// Add short alias to the Aliases line
	aliasIdx := strings.Index(content, "Aliases:")
	if aliasIdx != -1 {
		// Find the end of the aliases line(s)
		afterAlias := content[aliasIdx:]
		aliasEnd := strings.Index(afterAlias, "\n\n")
		if aliasEnd == -1 {
			aliasEnd = len(afterAlias)
		}
		aliasSection := afterAlias[:aliasEnd]

		// Insert the new alias before the last entry
		if !strings.Contains(aliasSection, shortAlias) {
			// Add to end of alias list
			newAliasSection := strings.TrimRight(aliasSection, " \n") + ", " + shortAlias
			content = content[:aliasIdx] + newAliasSection + content[aliasIdx+aliasEnd:]
		}
	}

	return os.WriteFile(path, []byte(content), 0644)
}

// extractHelpName pulls the persona name from a help line like "  /vern:inverse <task>     Sonnet ..."
func extractHelpName(line string) string {
	re := regexp.MustCompile(`/vern:([a-z-]+)`)
	m := re.FindStringSubmatch(line)
	if len(m) > 1 {
		return m[1]
	}
	return line
}

// UpdateEmbeddedTest adds a name to the expectedAgents slice in embedded_test.go.
func UpdateEmbeddedTest(path, name string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read embedded_test.go: %w", err)
	}
	content := string(data)

	// Find the expectedAgents slice
	sliceStart := strings.Index(content, "var expectedAgents = []string{")
	if sliceStart == -1 {
		return fmt.Errorf("expectedAgents slice not found")
	}

	sliceEnd := strings.Index(content[sliceStart:], "}")
	if sliceEnd == -1 {
		return fmt.Errorf("expectedAgents slice end not found")
	}
	sliceEnd += sliceStart

	sliceContent := content[sliceStart : sliceEnd+1]

	// Extract existing names
	re := regexp.MustCompile(`"([a-z-]+)"`)
	matches := re.FindAllStringSubmatch(sliceContent, -1)
	var names []string
	for _, m := range matches {
		names = append(names, m[1])
	}

	// Check if already present
	for _, n := range names {
		if n == name {
			return nil // already there
		}
	}

	// Add and sort
	names = append(names, name)
	sort.Strings(names)

	// Rebuild the slice with the same formatting (5 per line)
	var sb strings.Builder
	sb.WriteString("var expectedAgents = []string{\n")
	for i := 0; i < len(names); i += 5 {
		end := i + 5
		if end > len(names) {
			end = len(names)
		}
		sb.WriteString("\t")
		for j := i; j < end; j++ {
			sb.WriteString(fmt.Sprintf("%q", names[j]))
			if j < len(names)-1 {
				sb.WriteString(", ")
			} else {
				sb.WriteString(",")
			}
		}
		sb.WriteString("\n")
	}
	sb.WriteString("}")

	content = content[:sliceStart] + sb.String() + content[sliceEnd+1:]
	return os.WriteFile(path, []byte(content), 0644)
}

// UpdateHardcodedRoster inserts a new Vern entry into the hardcodedRoster() function in selection.go.
func UpdateHardcodedRoster(path, name, desc string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read selection.go: %w", err)
	}
	content := string(data)

	// Find the hardcodedRoster function and its return slice
	funcIdx := strings.Index(content, "func hardcodedRoster()")
	if funcIdx == -1 {
		return fmt.Errorf("hardcodedRoster() function not found")
	}

	// Find return []Vern{
	returnIdx := strings.Index(content[funcIdx:], "return []Vern{")
	if returnIdx == -1 {
		return fmt.Errorf("return []Vern{ not found in hardcodedRoster")
	}
	returnIdx += funcIdx

	// Find closing }
	braceDepth := 0
	sliceStart := returnIdx + len("return []Vern{")
	sliceEnd := -1
	for i := sliceStart; i < len(content); i++ {
		switch content[i] {
		case '{':
			braceDepth++
		case '}':
			if braceDepth == 0 {
				sliceEnd = i
				goto found
			}
			braceDepth--
		}
	}
found:
	if sliceEnd == -1 {
		return fmt.Errorf("could not find end of hardcodedRoster slice")
	}

	// Extract existing entries
	sliceContent := content[sliceStart:sliceEnd]
	entryRe := regexp.MustCompile(`\{ID: "([^"]+)", LLM: "([^"]+)", Desc: "([^"]+)"\}`)
	matches := entryRe.FindAllStringSubmatch(sliceContent, -1)

	type entry struct {
		ID, LLM, Desc string
	}
	var entries []entry
	for _, m := range matches {
		entries = append(entries, entry{ID: m[1], LLM: m[2], Desc: m[3]})
	}

	// Check if already present
	for _, e := range entries {
		if e.ID == name {
			return nil
		}
	}

	// Add and sort
	entries = append(entries, entry{ID: name, LLM: "claude", Desc: desc})
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].ID < entries[j].ID
	})

	// Rebuild
	var sb strings.Builder
	sb.WriteString("\n")
	for _, e := range entries {
		sb.WriteString(fmt.Sprintf("\t\t{ID: %q, LLM: %q, Desc: %q},\n", e.ID, e.LLM, e.Desc))
	}
	sb.WriteString("\t")

	content = content[:sliceStart] + sb.String() + content[sliceEnd:]
	return os.WriteFile(path, []byte(content), 0644)
}

// ComputeAliases generates aliases for a persona name.
// Returns the full name + a 3-letter prefix. Avoids conflicts with existing aliases.
func ComputeAliases(name string, existingAliases map[string]bool) []string {
	aliases := []string{name}

	// 3-letter prefix
	prefix := name
	if len(prefix) > 3 {
		prefix = prefix[:3]
	}
	if !existingAliases[prefix] && prefix != name {
		aliases = append(aliases, prefix)
	}

	return aliases
}

// KnownAliases returns the set of aliases already in use by existing personas.
func KnownAliases() map[string]bool {
	return map[string]bool{
		// Core personas
		"mediocre": true, "med": true, "m": true,
		"great": true, "vernile": true, "g": true,
		"nyquil": true, "nq": true, "n": true,
		"ketamine": true, "ket": true, "k": true,
		"yolo": true, "y": true,
		"mighty": true, "codex": true, "c": true,
		// Specialists
		"architect": true, "arch": true, "ar": true,
		"inverse": true, "inv": true, "i": true,
		"paranoid": true, "para": true, "p": true,
		"optimist": true, "opt": true, "o": true,
		"academic": true, "acad": true, "a": true,
		"startup": true, "su": true, "s": true,
		"enterprise": true, "ent": true, "e": true,
		"ux": true, "u": true,
		"retro": true, "ret": true, "r": true,
		"oracle": true, "orc": true, "ora": true,
		// Workflows
		"setup": true,
		"new-idea": true, "new": true, "ni": true,
		"discovery": true, "disco": true, "d": true,
		"hole": true, "khole": true, "vh": true,
		"generate": true, "gen": true,
		"help": true, "h": true, "?": true,
	}
}
