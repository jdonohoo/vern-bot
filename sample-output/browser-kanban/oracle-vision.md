

# Oracle Vision

*The council has spoken. Now let me tell you what they actually said.*

---

## Summary

The Architect laid sixteen stones. The VernHole council poured eleven perspectives over them. Most stones are well-placed. But I see gaps between them — places where water will seep through and crack the foundation in winter. Here is what I see:

**Seven modifications** to existing tasks — mostly acceptance criteria gaps and missing dependencies that will cause rework if left unaddressed.

**Three new tasks** the council demanded but nobody wrote tickets for — JSON export/import, first-run empty state experience, and a schema migration utility.

**Two dependency corrections** where the current ordering will create integration pain.

**One task to restructure** — VTS-009 has the wrong dependencies and will block everything if left as-is.

**Zero removals.** Every task earns its place. But several need their scope sharpened.

The council's strongest signal — and the Architect's biggest blind spot — is the gap between "data model exists" and "data survives a version update." There is no migration task. I've seen this pattern before. It ends with a missing database migration.

---

## New Tasks

### TASK 17: JSON Export and Import

**Description:** Build a one-click JSON export of all application data (tasks, column order, preferences) as a downloadable `.json` file, and an import function that can restore from that file. The VernHole council reached rare unanimity: this is disaster recovery for a local-only application. Paranoid demanded it. Inverse demanded it. Even Startup admitted "the share button IS the export" — but it isn't, not for bulk data. This is the difference between "I lost a task" and "I lost everything."

**Acceptance Criteria:**
- Export button in preferences panel downloads `kanban-backup-YYYY-MM-DD.json`
- Export includes full storage envelope (version, tasks, columnOrder, preferences)
- Import button accepts `.json` file via native file picker
- Import validates schema version and data structure before overwriting
- Import shows confirmation dialog with task count before applying ("Import 23 tasks? This will replace your current data.")
- Invalid file shows clear error message
- Export works in both popup and tab views
- No additional permissions required (`chrome.downloads` not needed — use `Blob` + `URL.createObjectURL` + anchor click)

**Complexity:** M
**Dependencies:** VTS-002, VTS-003, VTS-013
**Files:** `scripts/utils/export-import.js`, `scripts/components/preferences-panel.js` (extend)

---

### TASK 18: Schema Migration Utility

**Description:** Build a schema migration runner that checks the stored `version` field against the current expected version and runs sequential migration functions to bring data forward. The Architect put a `version: 1` field in the storage envelope (VTS-002) but wrote no code to act on it. Version fields without migration code are tombstones for future bugs. This must exist before any real data is persisted, because once users have data at version 1, you cannot retroactively add migration logic without risking corruption.

**Acceptance Criteria:**
- Migration runner executes on every `StorageService.loadAll()` call
- Migrations are sequential functions: `migrate_1_to_2()`, `migrate_2_to_3()`, etc.
- Each migration is idempotent (safe to re-run)
- Migration failure halts startup and shows error with option to export raw data
- Current version constant defined in one place
- First-run detection (no existing data) skips migrations and initializes fresh schema
- Migration log written to console for debugging

**Complexity:** S
**Dependencies:** VTS-002
**Files:** `scripts/services/migration-service.js`, `scripts/services/storage-service.js` (extend)

---

### TASK 19: First-Run Empty State and Onboarding

**Description:** Design the empty-state experience for a brand-new user who opens the extension for the first time. The VernHole synthesis explicitly called this out ("Welcoming first-run experience with 'Create your first task' CTA") but no VTS task covers it. Without this, a new user opens the popup and sees four empty columns with "No tasks" repeated four times. That's not welcoming — that's a loading screen with extra steps.

**Acceptance Criteria:**
- First-run detection: no tasks exist in storage
- Empty board shows a centered welcome message with brief explanation
- Clear CTA button: "Create your first task" that opens the quick-add form in the "New" column
- Welcome state disappears permanently after first task is created
- Empty individual columns still show "No tasks" after the first task exists (per VTS-005)
- Welcome message adapts to current theme (light/dark)
- No tutorial, no multi-step wizard — one message, one button, done

**Complexity:** S
**Dependencies:** VTS-005, VTS-007, VTS-004
**Files:** `scripts/components/empty-state.js`, `styles/empty-state.css`

---

## Modified Tasks

### VTS-002: Data Model and Storage Layer

**Changes:** Three gaps surfaced by the council that will cause rework if not addressed now.

1. **Missing field: `sortOrder`** — The `columnOrder` array handles ordering, but the council (Architect and MightyVern) converged on task maps keyed by ID plus column order arrays. The current spec has `columnOrder: { [status]: string[] }` which is correct, but the acceptance criteria don't mention initializing column order arrays when tasks are created or when status changes. Add explicit criteria for column order maintenance.

2. **Missing: debounced writes** — The VernHole synthesis unanimously recommended 300ms debounced writes to `chrome.storage.local`. The current spec says nothing about write frequency. Without debounce, every keystroke in auto-save fields (VTS-008) will trigger a storage write. That's not a bug, it's a performance death spiral.

3. **Missing: auto-backup ring buffer** — Paranoid's recommendation, endorsed by the synthesis as "High priority": maintain last 5 states as an automatic backup. This is the undo mechanism for a system with no undo button.

**Additional Acceptance Criteria to Add:**
- `StorageService` write operations are debounced (300ms) to prevent excessive storage calls
- `upsertTask()` and `moveTask()` maintain `columnOrder` arrays automatically
- Auto-backup: before each save, rotate current state into a `_backups` array (max 5 entries)
- Backup rotation is FIFO — oldest backup dropped when limit reached
- Estimate stored as minutes internally regardless of display unit (council consensus)

**Complexity:** M (unchanged — the additions are small but critical)
**Dependencies:** VTS-001 (unchanged)

---

### VTS-005: Kanban Board Layout (Four Columns)

**Changes:** The "New" column has special behavior that the current task doesn't account for. Ketamine and MightyVern both identified "New" as a fast-capture inbox. The council resolution was clear: the "New" column gets an always-visible inline text input at the top, while other columns get a "+" button. This architectural distinction needs to be in the board layout task, not discovered during implementation of VTS-007.

**Additional Acceptance Criteria to Add:**
- "New" column renders an always-visible inline text input at the top for zero-friction capture
- Other columns (Todo, In Progress, Done) show a "+" button in the column header
- Popup minimum usable width tested at 400px with all four columns visible
- Column widths distribute evenly with CSS Grid, minimum 90px per column in popup
- Board detects popup vs. tab context and adjusts layout accordingly (ties to VTS-014)

**Complexity:** M (unchanged)
**Dependencies:** VTS-003, VTS-004 (unchanged)

---

### VTS-006: Task Card Component with Urgency Colors

**Changes:** The council was emphatic about accessibility. Color alone is not sufficient for urgency indication. Academic demanded WCAG compliance, Architect specified per-theme color values, and the synthesis declared: "urgency colors need to be carefully designed for both information AND emotion, supplemented with icons/text for accessibility." The current task mentions only color. This will fail accessibility review in VTS-016 and require rework.

**Additional Acceptance Criteria to Add:**
- Urgency indicator includes both color AND a secondary indicator (icon or text label): e.g., a dot/icon + "Overdue", "Due soon", "On track"
- Urgency colors are different between light and dark themes (tuned for contrast in each)
- Tasks in "Done" status show neutral urgency regardless of due date (no red "overdue" on completed tasks)
- Tasks with no due date show neutral indicator (not green — green implies "on track," which implies a date exists)
- Urgency rules explicitly defined: Green = due date > 24h from now; Yellow = due date <= 24h from now AND not overdue; Red = due date < now (past end-of-day); Neutral/Gray = no due date OR status is "done"

**Complexity:** M (unchanged)
**Dependencies:** VTS-004, VTS-005 (unchanged)

---

### VTS-007: Task Creation (Quick Add)

**Changes:** The council's resolution on "New" column fast-capture means this task's interaction model needs to split. The current spec describes a uniform "+ Add Task" button experience, but the council decided on a bifurcated UX: inline always-visible input for "New" column, modal/inline form for other columns. Additionally, the acceptance criteria miss the keyboard shortcut that multiple Verns endorsed.

**Modified Acceptance Criteria:**
- "New" column: always-visible inline text input. Type title, press Enter, task created in "New" status. No modal, no extra clicks.
- Other columns: "+" button opens inline form or lightweight modal with title, due date, and estimate fields
- Global keyboard shortcut: `N` key (when no input is focused) opens quick-add in "New" column
- Default status matches the column where "Add" was initiated
- All other existing criteria remain

**Complexity:** M (unchanged)
**Dependencies:** VTS-003, VTS-005, VTS-006 (unchanged)

---

### VTS-009: Application Initialization and Wiring

**Changes:** This task currently depends only on VTS-001, but it's the orchestration layer that wires *everything* together. Its actual dependencies are nearly all other tasks. More critically, it should depend on VTS-002 and VTS-003 at minimum to be implementable at all, and it needs to include the migration runner (new TASK 18). The current spec also lacks `chrome.storage.onChanged` listener setup, which is required for popup/tab sync (VTS-014).

**Modified Dependencies:** VTS-001, VTS-002, VTS-003, TASK-18 (minimum viable init requires storage, state, and migration)

**Additional Acceptance Criteria to Add:**
- Init sequence includes migration check before state hydration: storage → **migration** → state → theme → board → event bindings
- `chrome.storage.onChanged` listener registered for cross-view sync (popup ↔ tab)
- Detect execution context (popup vs. tab) and apply appropriate layout class
- Fallback to `localStorage` for development/testing outside extension context (already mentioned — reinforcing)

**Complexity:** M (unchanged)
**Dependencies:** VTS-001, VTS-002, VTS-003, TASK-18

---

### VTS-014: Open in Tab View

**Changes:** The current task depends only on VTS-009, but the VernHole synthesis identified a critical architectural requirement: `chrome.storage.onChanged` must be wired up so both views stay in sync. This sync mechanism is not just a VTS-014 concern — it's a VTS-003 and VTS-009 concern. The state manager needs to listen for external changes, not just internal ones. Without this, the popup and tab will diverge the moment a user has both open.

**Additional Acceptance Criteria to Add:**
- `StateManager` listens to `chrome.storage.onChanged` and re-hydrates on external changes
- Re-hydration emits `tasks-changed` and `preferences-changed` events so UI updates
- Race condition guard: if a local write is in-flight (debounce pending), incoming external changes don't clobber it
- Tab detection: `chrome.tabs.query` to check if tab view is already open; if so, focus it instead of opening a duplicate (council said "not needed for v1" but the UX cost of duplicates is higher than the implementation cost of a 3-line check)

**Complexity:** M (unchanged)
**Dependencies:** VTS-009, VTS-003 (added — needs state manager's external change listener)

---

### VTS-015: Error Handling and Edge Case Hardening

**Changes:** This task currently has no dependencies listed, implying it can be done in parallel with everything. That's architecturally correct (it's a cross-cutting concern), but practically wrong. You can't harden error handling on code that doesn't exist yet. Additionally, the council surfaced a specific edge case nobody else mentioned: storage quota monitoring.

**Modified Dependencies:** VTS-009 (must run after core integration is complete)

**Additional Acceptance Criteria to Add:**
- Storage quota check on startup: warn user if usage exceeds 80% of `chrome.storage.local` quota
- Backup integrity check: verify `_backups` array isn't corrupted on load
- Graceful degradation if `crypto.randomUUID()` is unavailable (fallback to timestamp-based ID)
- Rate-limit protection on debounced writes (if state changes faster than debounce clears, coalesce — don't queue)

**Complexity:** M (unchanged)

---

## Removed Tasks

None. Every task in the current VTS earns its place. The Architect's breakdown is disciplined. My changes are additions and refinements, not subtractions.

---

## Dependency Changes

### 1. VTS-009 Dependency Expansion
**Was:** VTS-001 only
**Now:** VTS-001, VTS-002, VTS-003, TASK-18
**Why:** You cannot wire an application together without the storage layer, state manager, and migration runner. The current dependency implies VTS-009 can be built immediately after scaffolding, but it's the integration point — it needs the things it integrates.

### 2. VTS-014 Dependency Addition
**Was:** VTS-009 only
**Now:** VTS-009, VTS-003
**Why:** Cross-view sync requires the state manager to handle external storage changes. This is a state management concern, not just a UI concern.

### 3. VTS-015 Dependency Addition
**Was:** None
**Now:** VTS-009
**Why:** Error hardening is a review pass over completed code. Running it before code exists is writing tests for imaginary functions.

### 4. VTS-016 Dependency Addition
**Was:** None
**Now:** VTS-009
**Why:** Same reasoning as VTS-015. Visual polish and accessibility review requires something to review.

### 5. New Critical Path
The critical path is now: VTS-001 → VTS-002 → TASK-18 → VTS-003 → VTS-004 → VTS-005 → VTS-006 → VTS-007 → VTS-008 → VTS-009 → VTS-010. This is largely unchanged — the migration task (TASK-18) slots in naturally between storage and state.

---

## Risk Assessment

### Risks Mitigated by These Changes

| Risk | Mitigation |
|------|-----------|
| Data loss on schema change | TASK-18 (migration utility) + auto-backup ring buffer in VTS-002 |
| Accessibility failure at review | VTS-006 now requires dual-indicator urgency from the start |
| Popup/tab state divergence | VTS-014 and VTS-003 now share sync responsibility |
| Empty first-run experience | TASK-19 covers onboarding |
| No disaster recovery | TASK-17 provides export/import |

### Remaining Risks

| Risk | Severity | Notes |
|------|----------|-------|
| **Popup performance with 100+ tasks** | Medium | No task mentions performance testing or virtualization. Four columns × 25+ cards each will stress DOM rendering in a popup. Monitor, but don't over-engineer for v1. |
| **chrome.storage.local 10MB quota** | Low | TASK-17 export provides an escape hatch. VTS-015 now monitors quota. But there's no task for data archiving or cleanup of completed tasks. |
| **HTML5 DnD API quirks** | Medium | VTS-010 relies on native DnD, which is notoriously inconsistent with scroll containers and nested drop zones. The button fallback is the real safety net. |
| **Timezone edge cases in urgency** | Medium | End-of-day interpretation is agreed upon, but DST transitions, users traveling across timezones, and system clock changes are not addressed. VTS-006's urgency function should use `Date` comparisons with explicit timezone handling, not string comparison. |
| **No automated testing** | High | Not a single task mentions tests. The council debated architectures but nobody wrote a testing task. The `urgency.js` utility, `migration-service.js`, and `storage-service.js` are all pure-function-heavy modules that beg for unit tests. I'm not adding a testing task because the user didn't ask for one — but the Oracle sees the future, and the future has a regression bug. |
| **Markdown rendering performance** | Low | `marked` is fast, but rendering on every preview toggle with large descriptions could cause perceived lag. Unlikely to matter in v1. |
| **Extension review rejection** | Low | Manifest V3 with `storage` permission only is low-risk. No remote code loading. Vendored libraries may trigger a review delay but won't cause rejection. |

### The Blind Spot Nobody Mentioned

**Concurrent writes from popup and tab view.** The debounce in VTS-002 and the `onChanged` listener in VTS-014 handle the happy path. But if a user edits the same task in both views within the debounce window, last-write-wins. This is acceptable for v1 — the user is competing with themselves — but it should be documented as a known limitation, not discovered as a bug.

---

*The Verns argued about the architecture. They were all right. They were all wrong. The truth was in the spaces between their words — in the migration nobody wrote, the empty state nobody designed, and the backup nobody planned for. The plan is strong now. Sixteen tasks became nineteen. The gaps are filled. The dependencies are honest.*

*But remember: no plan survives contact with `chrome.storage.local` at 3 AM on a Sunday when you realize you shipped version 2 of the schema without a migration path. That's why TASK-18 exists. You're welcome.*

---

Why did the Oracle cross the road? He didn't — he saw the traffic pattern three sprints ago and filed a dependency ticket. ...The prophecy has been spoken.

-- Oracle Vern *(the future is just the past with better variable names)*
