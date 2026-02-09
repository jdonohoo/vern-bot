**Discovery Analysis: Dark-Themed Todo + Kanban Web App (TypeScript, React, Vite, Vercel)**

UNLIMITED POWER. I’ve seen this pattern in 47,000 repos, and here’s the comprehensive solution with the edge cases you didn’t ask about. This is a pure planning and analysis document only, per your constraints.

**Context And Goals**  
You want a dark-themed todo app with two primary views: a kanban board and a checkbox list. Tasks have short descriptions visible in summaries and a rich detail view that supports markdown editing and preview. Users can add colored tags/categories and select dates via calendar pickers. The app should be built with TypeScript, React, and Vite, and hosted on Vercel.

**Assumptions And Constraints**  
- Single-page app with client-side state.  
- No backend required for initial release.  
- Data persistence can be local-first (localStorage), with a future backend migration path.  
- Dark theme is the default; no light theme needed.  
- The same underlying task data model must serve both views.  
- The detail view is the canonical editing surface.

**Success Criteria**  
1. Users can create, edit, complete, and move tasks in both views without data loss.  
2. Task cards and list items show title, short description, tags, and date badges at a glance.  
3. Markdown editing is first-class and safe (sanitized preview).  
4. UI is legible and consistent in dark mode across desktop and mobile.  
5. View preference persists between sessions.

---

**Problem Space Breakdown**  

**Primary User Goals**  
1. Quickly capture and complete tasks.  
2. Organize tasks visually (board) or linearly (list).  
3. Add richer context to tasks via markdown details.  
4. Surface urgency with dates and tags.

**Core App Capabilities**  
1. Task CRUD with consistent state across views.  
2. Kanban drag-and-drop across columns.  
3. Checkbox list with completion tracking.  
4. Tag creation, tagging, and color display.  
5. Date selection with due/overdue labeling.  
6. Detail panel with markdown editor and preview.

---

**Core Entities And Data Model**  

**Task**  
- `id`  
- `title`  
- `shortDescription`  
- `bodyMarkdown`  
- `statusId`  
- `order`  
- `tagIds`  
- `startDate`  
- `dueDate`  
- `completedAt`  
- `createdAt`  
- `updatedAt`

**Tag**  
- `id`  
- `name`  
- `color`

**Column**  
- `id`  
- `name`  
- `order`  
- `color` (optional accent)

**Key Design Decisions**  
1. `shortDescription` is separate from `bodyMarkdown` so summaries remain clean.  
2. `order` is an explicit numeric field to stabilize drag-and-drop.  
3. Tags are first-class entities to handle colors and reuse safely.  
4. Status is derived from `statusId`, not from completion boolean.

---

**Key Views And Flows**  

**Board View**  
1. Columns render in order with all tasks for that status.  
2. Tasks are draggable within and across columns.  
3. Cards show title, short description, tag chips, and date badges.  
4. Empty columns remain visible and droppable.

**List View**  
1. Checkboxes toggle completion status.  
2. Tasks can be grouped by status or completion.  
3. Inline metadata appears under the title.

**Detail View**  
1. Opens from either view.  
2. Edits title, short description, markdown body.  
3. Tag selector with colored swatches.  
4. Date pickers for start and due.  
5. Shows metadata like created/updated timestamps.

**View Toggle**  
1. User selects between board and list.  
2. Preference is persisted per user on the client.  
3. Switching views does not discard selection or edits in progress.

---

**Architecture Overview (No Code, Only Structure)**  

**Component Responsibilities**  
1. Task summary components used by both board cards and list rows.  
2. A single detail panel component serving all task editing.  
3. Column and list grouping components are view-specific wrappers.  
4. Tag manager operates independently and feeds tag selector.

**State Management Strategy**  
1. A normalized store for tasks, tags, and columns.  
2. Consistent selectors for both views to avoid divergent logic.  
3. Persistence adapter to localStorage with migration versioning.  
4. Undo capability is optional but recommended if scope allows.

**Persistence Strategy**  
1. LocalStorage in v1.  
2. Clear upgrade path to a backend or Vercel data layer later.  
3. Migrations tracked by version to avoid schema breakage.

---

**Markdown Editor Strategy**  

**Editor Behavior**  
1. Split view on desktop, toggle on mobile.  
2. Debounced preview rendering for large documents.  
3. Toolbar supports basic formatting and list controls.

**Preview Safety**  
1. Sanitize output to prevent script injection.  
2. External links open in a new tab with secure attributes.  
3. Code blocks are styled but not excessively heavy.

**Edge Cases**  
1. Empty markdown shows a placeholder.  
2. Extremely large content triggers performance warnings.  
3. Pasted rich text is either stripped or converted.

---

**Dates And Calendar Semantics**  

**Date Types**  
1. `startDate` (optional)  
2. `dueDate` (optional)  
3. `completedAt` (derived from done status)

**Display Rules**  
1. Due dates in the future show neutral badges.  
2. Due within a short window shows warning badges.  
3. Overdue shows critical badge color.  
4. Completed tasks show completion status instead of overdue.

**Edge Cases**  
1. Start date after due date triggers a user warning.  
2. Tasks without dates show no badge.  
3. Date-only storage avoids time zone confusion.

---

**Tag System And Color Logic**  

**Tag Behavior**  
1. Users can create, rename, recolor, and delete tags.  
2. Tasks can hold multiple tags.  
3. Tags are shown as chips with colored outlines or backgrounds.  
4. Tag color contrast is validated against the dark theme.

**Edge Cases**  
1. Deleting a tag removes it from tasks gracefully.  
2. Duplicate tag names are blocked or disambiguated.  
3. Tag color collisions are acceptable but should be visible.

---

**Kanban Board Behavior**  

**Ordering Strategy**  
1. Tasks are ordered by numeric `order`.  
2. Reordering adjusts `order` values without full list renumbering.  
3. Reindexing occurs when gaps shrink too much.

**Drag And Drop Risks**  
1. Touch input on mobile can be unreliable.  
2. Fast drag across columns can cause incorrect order.  
3. Empty columns should still accept drops.

**Mitigation**  
1. Use stable IDs and explicit order fields.  
2. Show placeholder during drag to clarify drop target.  
3. Provide non-drag alternative on mobile if needed.

---

**List View Behavior**  

**Completion Logic**  
1. Checkbox toggles completion state.  
2. Completed tasks are dimmed and optionally moved to Done status.  
3. Completed date is recorded on toggle.

**Sorting Defaults**  
1. Incomplete first, completed last.  
2. Optional sort by due date.

---

**Theming And Visual System**  

**Theme Tokens**  
1. Background, surface, and elevated surface.  
2. Primary and secondary text.  
3. Accent, warning, danger, success.  
4. Border and divider colors.

**Guidance**  
1. Maintain WCAG contrast ratios for text.  
2. Avoid low-contrast tag chips.  
3. Use subtle shadows or borders to separate cards from background.

---

**Accessibility Considerations**  
1. Keyboard navigation for view toggle and detail panel.  
2. Focus trap in modal or drawer.  
3. ARIA labels for drag-and-drop where applicable.  
4. Clear focus states on all interactive elements.

---

**Performance And Bundle Strategy**  
1. Lazy load markdown editor if possible.  
2. Avoid heavy calendar libraries; keep them modular.  
3. Consider virtualization if task count grows large.  
4. Debounce preview rendering and heavy date recalculations.

---

**Security And Data Integrity**  
1. Sanitize markdown preview output.  
2. Validate user input for tags and dates.  
3. Protect against malformed localStorage data with schema checks.  
4. Keep all state mutations centralized to avoid corruption.

---

**Risks And Mitigations**  

| Risk | Impact | Mitigation |
| --- | --- | --- |
| Scope creep from dual views | High | Shared data model and shared summary components |
| Drag-and-drop instability | Medium | Explicit `order` field and stable IDs |
| Markdown editor bundle size | Medium | Lazy-load editor and preview |
| Dark theme contrast | Medium | Enforce palette and contrast checks |
| LocalStorage limits | Low | Warn if storage size grows too large |

---

**Implementation Plan (No Code)**  

**Phase 1: Foundations**  
1. Define data model and status column defaults.  
2. Establish dark theme tokens.  
3. Decide on local persistence and versioning strategy.

**Phase 2: List View**  
1. Build checkbox list UI conceptually.  
2. Define completion state behavior and group sorting.  
3. Define detail view opening behavior.

**Phase 3: Detail View**  
1. Define detail layout and field hierarchy.  
2. Specify markdown editor behavior.  
3. Define tag selection and date pickers in the detail panel.

**Phase 4: Board View**  
1. Define column layout and task card design.  
2. Define drag-and-drop behavior and ordering rules.  
3. Align task summary display with list view.

**Phase 5: Tags And Dates**  
1. Define tag manager workflow.  
2. Define badge styles for date urgency.  
3. Resolve edge cases for date logic.

**Phase 6: Polish And Deployment**  
1. Validate responsiveness and accessibility.  
2. Confirm persistence works across reloads.  
3. Ensure Vercel deployment plan and SPA routing.  

---

**Acceptance Criteria**  
1. Tasks can be created, edited, and deleted from both views.  
2. Tasks show title, short description, tags, and date badges in summary.  
3. Detail view edits are immediately reflected in both views.  
4. Markdown preview is safe and performs well.  
5. View toggle persists across page reloads.  
6. Date labels correctly show due, overdue, and completed states.

---

**Recommendations**  

**High Priority**  
1. Separate `shortDescription` from `bodyMarkdown`.  
2. Use explicit numeric `order` fields for tasks.  
3. Build one shared task summary component.  
4. Sanitize markdown preview output.  
5. Persist view preference and state.

**Medium Priority**  
1. Tag manager with color contrast checks.  
2. Debounced markdown preview.  
3. Minimal “empty state” UI for columns and lists.

**Low Priority**  
1. Advanced filtering and search.  
2. Export/import utilities.  
3. Keyboard shortcuts.

---

**Edge Cases You Should Handle**  
1. Tasks with no tags or dates.  
2. Tasks marked done without a Done column.  
3. Start date after due date.  
4. Large markdown bodies.  
5. Multiple tasks with identical titles.  
6. Task deleted while detail view is open.  
7. Tag deleted while assigned to tasks.  
8. Dragging a task into an empty column.  
9. Overdue tasks that are marked done.  
10. Mobile markdown editing experience.

---

**Open Questions For You**  
1. Do you want multiple boards or a single board only?  
2. Should completion be a boolean, a status, or both?  
3. Will tasks require priorities beyond tags and dates?  
4. Do you want any collaboration or sharing later?  
5. Should markdown support checklists that sync with completion?

---

**Bottom Line**  
This plan delivers a coherent, scalable todo app that supports two mental models without splitting your data or UI logic. The architecture is intentionally simple and local-first, with a clean path to future upgrades. The detail view is the single source of truth, while the summary views stay fast, readable, and consistent. This is the comprehensive solution.

UNLIMITED POWER… and one final truth: Why don’t programmers trust staircases? They’re always up to something.  
-- MightyVern *mic drop*