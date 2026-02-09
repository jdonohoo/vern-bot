

# Updated Task Breakdown — Architect Vern's Revised Blueprint

*The Oracle saw the gaps between the stones. Now the Architect fills them with load-bearing concrete.*

---

## Architecture Overview

The Oracle's vision is sound. Three new tasks, seven modifications, and dependency corrections that prevent integration pain. I've incorporated every recommendation, sharpened scope where needed, and ensured the critical path is honest about what depends on what.

**Critical Path:** VTS-001 → VTS-002 → TASK-18 → VTS-003 → VTS-004 → VTS-005 → VTS-006 → VTS-007 → VTS-008 → VTS-009 → VTS-010

**Parallelizable after VTS-009:** VTS-011, VTS-012, VTS-013, VTS-014, TASK-17, TASK-19, VTS-015, VTS-016

## VTS Task Index

Individual VTS files: `vts/`

| ID | Task | Complexity | Dependencies |
|----|------|------------|--------------|
| VTS-001 | Extension Scaffold and Manifest V3 Setup | S | None |
| VTS-002 | Data Model and Storage Layer | M | VTS-001 |
| VTS-003 | Schema Migration Utility | S | VTS-002 |
| VTS-004 | State Management Layer | M | VTS-002, VTS-003 |
| VTS-005 | Theme System (CSS Custom Properties + Toggle) | M | VTS-001, VTS-004 |
| VTS-006 | Kanban Board Layout (Four Columns) | M | VTS-004, VTS-005 |
| VTS-007 | Task Card Component with Urgency Colors | M | VTS-005, VTS-006 |
| VTS-008 | Task Creation (Quick Add) | M | VTS-004, VTS-006, VTS-007 |
| VTS-009 | Task Detail Panel (View + Edit) | L | VTS-004, VTS-007, VTS-008 |
| VTS-010 | Application Initialization and Wiring | M | VTS-001, VTS-002, VTS-003, VTS-004 |
| VTS-011 | Drag and Drop Between Columns | L | VTS-006, VTS-007 |
| VTS-012 | Markdown Editor and Renderer | L | VTS-009 |
| VTS-013 | Share Task as Markdown | S | VTS-009 |
| VTS-014 | Preferences Panel | S | VTS-005, VTS-012 |
| VTS-015 | Open in Tab View | M | VTS-010, VTS-004 |
| VTS-016 | Error Handling and Edge Case Hardening | M | VTS-010 |
| VTS-017 | Visual Polish and Accessibility Pass | M | VTS-010 |
| VTS-018 | JSON Export and Import | M | VTS-002, VTS-004, VTS-014 |
| VTS-019 | First-Run Empty State and Onboarding | S | VTS-006, VTS-008, VTS-005 |

## Dependency Graph (ASCII)

```
TASK 1 (Scaffold)
  │
  ├──→ TASK 2 (Storage Layer)
  │      │
  │      ├──→ TASK 3 (Migration Utility)
  │      │      │
  │      │      └──→ TASK 4 (State Manager)
  │      │             │
  │      │             ├──→ TASK 5 (Theme System) ──→ TASK 6 (Board Layout)
  │      │             │                                │
  │      │             │                                ├──→ TASK 7 (Task Cards)
  │      │             │                                │      │
  │      │             │                                │      ├──→ TASK 8 (Quick Add)
  │      │             │                                │      │      │
  │      │             │                                │      │      └──→ TASK 9 (Detail Panel)
  │      │             │                                │      │             │
  │      │             │                                │      │             ├──→ TASK 12 (Markdown)
  │      │             │                                │      │             ├──→ TASK 13 (Share)
  │      │             │                                │      │             └──→ TASK 14 (Preferences)
  │      │             │                                │      │                    │
  │      │             │                                │      │                    └──→ TASK 18 (Export/Import)
  │      │             │                                │      │
  │      │             │                                │      └──→ TASK 11 (Drag & Drop)
  │      │             │                                │
  │      │             │                                └──→ TASK 19 (Empty State)
  │      │             │
  │      │             └──→ TASK 10 (Init & Wiring)
  │      │                    │
  │      │                    ├──→ TASK 15 (Tab View)
  │      │                    ├──→ TASK 16 (Error Hardening)
  │      │                    └──→ TASK 17 (Accessibility)
```

---

## Summary of Changes from Original VTS

| Change | What | Why |
|--------|------|-----|
| **New** | Task 3: Schema Migration Utility | Version fields without migration code are tombstones for future bugs |
| **New** | Task 18: JSON Export and Import | Disaster recovery for a local-only application |
| **New** | Task 19: First-Run Empty State | Four empty columns saying "No tasks" is not onboarding |
| **Modified** | Task 2: Added debounce, auto-backup, column order maintenance | Performance death spiral prevention and undo mechanism |
| **Modified** | Task 4: Added `chrome.storage.onChanged` listener | Cross-view sync is a state management concern |
| **Modified** | Task 6: Added "New" column inline input, popup width testing | Architectural UX distinction belongs in layout task |
| **Modified** | Task 7: Added dual urgency indicators, per-theme colors, explicit rules | Color alone fails accessibility |
| **Modified** | Task 8: Bifurcated UX, keyboard shortcut | Council's resolution on fast-capture |
| **Modified** | Task 10: Expanded dependencies, migration in init sequence | Can't wire what doesn't exist |
| **Modified** | Task 15: Added sync via StateManager, duplicate tab detection | Prevent popup/tab divergence |
| **Modified** | Task 16: Added quota monitoring, backup integrity, UUID fallback | Hardening without these is incomplete |
| **Modified** | Task 17: Added dependency on Task 10 | Can't review code that doesn't exist |
| **Renumbered** | Original VTS-003 through VTS-016 shifted to accommodate Task 3 insertion | Sequential numbering maintained |

---

## Known Limitations (Documented, Not Bugs)

1. **Concurrent writes from popup and tab within debounce window:** Last-write-wins. Acceptable for v1 — the user is competing with themselves.
2. **No automated tests:** The Oracle sees the regression bug. So do I. But the user didn't ask for tests, and I don't add scope the user didn't request.
3. **No data archiving or cleanup of completed tasks:** Export (Task 18) provides the escape hatch. Quota monitoring (Task 16) provides the warning.

---

Why did the Architect add a migration utility to the task breakdown? Because a `version: 1` field without migration code is like a fire extinguisher with no pin — it looks responsible until the day you actually need it. ...Measure twice, deploy once.

-- Architect Vern *(the gaps between the stones are where the water gets in)*
