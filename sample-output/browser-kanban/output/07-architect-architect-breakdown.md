

# Architect Vern's VTS Breakdown: Kanban Browser Extension

Alright. I've read the master plan. It's thorough — MightyVern did good work on the vision. Now let me do what I do: turn vision into buildable, shippable tasks that a developer can pick up on the worst day of their life and still make progress.

Let me look at what we're working with first.

---

## System Architecture Overview

```
┌─────────────────────────────────────────────────┐
│                  Browser Extension               │
│                  (Manifest V3)                   │
├─────────────────────────────────────────────────┤
│  UI Layer                                        │
│  ┌──────────┐ ┌──────────┐ ┌──────────────────┐ │
│  │ Kanban   │ │ Task     │ │ Preferences      │ │
│  │ Board    │ │ Detail   │ │ Panel            │ │
│  │ (4 cols) │ │ Panel    │ │ (theme, settings)│ │
│  └────┬─────┘ └────┬─────┘ └────────┬─────────┘ │
│       │             │                │           │
├───────┴─────────────┴────────────────┴───────────┤
│  State Management Layer                          │
│  (in-memory state + event dispatch)              │
├──────────────────────────────────────────────────┤
│  Storage Layer                                   │
│  (chrome.storage.local, versioned JSON schema)   │
├──────────────────────────────────────────────────┤
│  Utilities                                       │
│  (urgency calc, markdown render, share export)   │
└──────────────────────────────────────────────────┘
```

---

## Phase 1: Core Foundation

## VTS Task Index

Individual VTS files: `vts/`

| ID | Task | Complexity | Dependencies |
|----|------|------------|--------------|
| VTS-001 | Extension Scaffold and Manifest V3 Setup | S | None |
| VTS-002 | Data Model and Storage Layer | M | VTS-001 |
| VTS-003 | State Management Layer | M | VTS-002 |
| VTS-004 | Theme System (CSS Custom Properties + Toggle) | M | VTS-001, VTS-003 |
| VTS-005 | Kanban Board Layout (Four Columns) | M | VTS-003, VTS-004 |
| VTS-006 | Task Card Component with Urgency Colors | M | VTS-004, VTS-005 |
| VTS-007 | Task Creation (Quick Add) | M | VTS-003, VTS-005, VTS-006 |
| VTS-008 | Task Detail Panel (View + Edit) | L | VTS-003, VTS-006, VTS-007 |
| VTS-009 | Application Initialization and Wiring | M | VTS-001 |
| VTS-010 | Drag and Drop Between Columns | L | VTS-005, VTS-006 |
| VTS-011 | Markdown Editor and Renderer | L | VTS-008 |
| VTS-012 | Share Task as Markdown | S | VTS-008 |
| VTS-013 | Preferences Panel | S | VTS-004, VTS-011 |
| VTS-014 | Open in Tab View | M | VTS-009 |
| VTS-015 | Error Handling and Edge Case Hardening | M | None |
| VTS-016 | Visual Polish and Accessibility Pass | M | None |

## Dependency Graph

```
TASK 1 (Scaffold)
  ├── TASK 2 (Storage)
  │     └── TASK 3 (State)
  │           ├── TASK 4 (Theme) ──────────────────┐
  │           │     └── TASK 5 (Board) ────────────┤
  │           │           ├── TASK 6 (Cards) ──────┤
  │           │           │     └── TASK 10 (DnD)  │
  │           │           └── TASK 7 (Quick Add)   │
  │           │                 └── TASK 8 (Detail) │
  │           │                       ├── TASK 11 (Markdown)
  │           │                       │     └── TASK 13 (Prefs)
  │           │                       └── TASK 12 (Share)
  │           └── TASK 9 (Wiring) ─── TASK 14 (Tab View)
  │
  └── TASK 15 (Error Hardening) ← All tasks
      └── TASK 16 (Polish + A11y) ← All tasks
```

## Summary

| Phase | Tasks | Complexity Budget |
|-------|-------|-------------------|
| Phase 1: Core | Tasks 1–9 | 2S + 6M + 1L |
| Phase 2: Enhanced | Tasks 10–14 | 2S + 2M + 2L |
| Phase 3: Polish | Tasks 15–16 | 2M |
| **Total** | **16 tasks** | **4S + 10M + 3L** |

## Key Architectural Decisions

1. **No framework.** A popup extension doesn't need React. Vanilla JS with a clean component pattern keeps the bundle small and the mental model simple.
2. **Task map + column order arrays** over flat array. Pays for itself the moment drag-and-drop lands.
3. **CSS custom properties for theming.** One attribute toggle, zero JS style manipulation.
4. **Native date picker.** No date library dependency for a date field.
5. **Bundled markdown library.** No CDN, no external fetches, no additional permissions.
6. **Auto-save on blur.** No save buttons. Fewer clicks, fewer lost edits, fewer "did I save?" moments.
7. **Storage versioning from day one.** Migration path built in before you need it.

---

How will this fail at 3 AM? It won't — because the storage layer has explicit error handling, the urgency calc is a pure function with no side effects, and the state manager is the single source of truth. The worst case is a storage write failure, and even then the in-memory state stays consistent and the user sees a clear error message.

The next developer might be having the worst day of their life. Every file in this plan is under 300 lines, every function does one thing, and every component has exactly one job.

Why did the browser extension architect insist on four Kanban columns? Because three-tier architecture just wasn't enough layers of separation. ...I'll see myself out.

-- Architect Vern (measure twice, deploy once)
