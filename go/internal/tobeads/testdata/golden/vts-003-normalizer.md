---
id: VTS-003
title: "Build Normalizer"
complexity: M
status: pending
owner: ""
source: discovery
source_ref: "architect-breakdown.md"
dependencies:
  - VTS-001
  - VTS-002
files:
  - "normalizer/normalizer.go"
---

# Build Normalizer

Transform parsed data into normalized output format.

## Criteria

- All status values map correctly
- Unknown statuses produce validation error
- Complexity codes produce correct labels
