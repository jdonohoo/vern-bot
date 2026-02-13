---
id: VTS-002
title: "Build the Parser"
complexity: S
status: pending
owner: ""
source: discovery
source_ref: "architect-breakdown.md"
dependencies:
  - VTS-001
files:
  - "parser/parser.go"
  - "parser/parser_test.go"
---

# Build the Parser

Implement the file parser with frontmatter support.

## Criteria

- Parses YAML frontmatter correctly
- Handles missing optional fields
- Unit tests cover edge cases
