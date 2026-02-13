---
id: VTS-001
title: "Label Invalid Source Ref"
complexity: M
status: pending
owner: ""
source: oracle
source_ref: "oracle-architect-breakdown.md (revised 2024-01-15)"
dependencies: []
files:
  - "src/foo/bar.rs"
  - "src/baz/../qux.rs"
---

# Label Invalid Source Ref

This task has a source_ref with characters that would be invalid as a label.
The description-metadata approach should handle this fine.

## Criteria

- source_ref appears in description, not as a label
- No sanitization needed
