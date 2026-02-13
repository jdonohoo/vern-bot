---
id: VTS-004
title: "Implement Executor"
complexity: M
status: active
owner: "vern"
source: oracle
source_ref: "oracle-vision (revised 2024-01-15).md"
dependencies:
  - VTS-002
  - VTS-003
files:
  - "executor/executor.go"
  - "executor/executor_test.go"
---

# Implement Executor

Build the execution layer that creates issues from normalized data.

## Criteria

- Creates issues via CLI with all mapped fields
- Captures returned ID from JSON output
- Handles duplicate gracefully
- Produces ID mapping file
