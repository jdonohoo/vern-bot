

# Architect Vern — Task Breakdown

*Systems-level decomposition. No files written. Just the blueprint.*

## Task Index

Individual task files: `tasks/`

| # | Task | Complexity | Dependencies |
|---|------|------------|--------------|
| 1 | Project Scaffolding & HTML Shell | S | None — this is the foundation. |
| 2 | File Input & Validation | S | Task 1 (needs the DOM structure to exist) |
| 3 | CSV Parsing with Papa Parse (Streaming) | M | Task 2 (needs the File object), Task 1 (needs error display region) |
| 4 | Table Rendering & Summary Display | M | Task 3 (needs parsed data), Task 1 (needs results area DOM) |
| 5 | Settings Panel (Header Toggle, Delimiter, Encoding) | M | Task 3 (re-triggers parsing), Task 4 (re-triggers rendering), Task 2 (needs the current File object to be retained) |
| 6 | Load More / Chunked Rendering | S | Task 3 (needs the full parsed row buffer), Task 4 (appends to the rendered table) |
| 7 | Error Handling & Edge Cases | S | Tasks 2, 3, 4, 5, 6 (this is a hardening pass over everything) |

## Dependency Graph

```
Task 1 (Scaffold)
  ├── Task 2 (File Input)
  │     └── Task 3 (Parsing)
  │           ├── Task 4 (Rendering)
  │           │     └── Task 6 (Load More)
  │           └── Task 5 (Settings) ──→ re-triggers Task 3 & 4
  └── Task 7 (Error Hardening) ── runs after 2-6 are stable
```

## Build Order Recommendation

| Phase | Tasks | Rationale |
|-------|-------|-----------|
| **Phase 1** | 1 → 2 → 3 → 4 | Get a working upload-to-table flow. This is your "it works" demo. |
| **Phase 2** | 5, 6 (parallel) | Settings and load-more are independent of each other. |
| **Phase 3** | 7 | Harden once the happy path is solid. Don't polish error states before the happy path exists. |

**Total complexity: 2S + 3M + 1S + 1S = 4S + 3M.** This is a single-session build for someone comfortable with vanilla JS. No framework, no build step, no deployment pipeline. One file, one purpose.
