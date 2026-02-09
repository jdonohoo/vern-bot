# Oracle Vision

## What the Council Actually Said

The council debated twelve perspectives across architecture, UX, pragmatism, and chaos. They think they argued about shortDescription fields and state management libraries. They did not. Here is what they actually said:

**The real signal has three frequencies:**

**Signal 1: This project will die from scope, not from technical failure.** Seven of twelve Verns flagged over-scoping. YOLO Vern, Startup Vern, and Mediocre Vern all independently arrived at the same conclusion through different reasoning: the current VTS breakdown describes a polished product, not a shippable first cut. The 15-task VTS list is a feature catalog masquerading as a build plan. It has no concept of "done enough to use."

**Signal 2: The UX plan is a ghost.** The council spent 80% of its energy on data models, state management, and component architecture. The actual user experience -- what happens when someone opens this app for the first time, how they discover features, what the empty board looks like, how quick-add actually feels, what happens when they have 47 tasks and no tags -- is almost entirely unspecified. The technical plan is a cathedral. The UX plan is a napkin sketch.

**Signal 3: The council agreed on more than they realize.** Strip away the surface-level debates (Zustand vs useReducer, gap integers vs fractional indexing) and the real architecture is settled: normalized store, Zustand with persist, drawer-based detail panel, Tailwind + dark theme, localStorage with versioning, status-driven columns. Every Vern who matters converged on this. The remaining "contradictions" are implementation details that a single developer will resolve in an afternoon by picking one and moving on.

**The hidden fourth signal:** Nobody talked about testing. Not once. Not a single Vern mentioned how any of this gets validated. No acceptance criteria reference testability. No task includes "write tests." This is the silence that speaks loudest.

---

## The Gaps Nobody Mentioned

These are the things that exist in the space between what the Verns said. They are not in any task. They are not in any synthesis. They will cause pain.

### Gap 1: No Data Migration Story
The council agreed on localStorage with schema versioning. Nobody defined what happens when the schema changes. VTS-003 says "schema versioning" but has no acceptance criteria for migration logic. When v1.1 ships with a different data shape, every existing user's data will silently break or vanish. This is not a v2 problem. This is a week-2 problem.

### Gap 2: No Empty States
YOLO Vern called this out. Nobody listened. The first thing every user sees is an empty board. The first experience of the list view is an empty list. The first time they open tags, there are no tags. Empty states ARE the onboarding. They are not polish. They are not Phase 3. They are the first five seconds of the product.

### Gap 3: No Error Boundaries
When the markdown editor fails to lazy-load (and it will, on slow connections), what does the user see? When localStorage is full, what happens? When a task references a deleted tag, what renders? The current plan has zero error handling tasks. Every component is assumed to work perfectly on the first try, forever.

### Gap 4: No Keyboard/Focus Management
The council agreed on quick-add (title + Enter). Nobody specified what happens after Enter. Does focus stay in the input? Does it move to the new task? Can you Tab through tasks? Can you Escape out of the detail drawer? Keyboard flow is not a Phase 3 feature. It is the difference between "this feels good" and "this feels like a school project."

### Gap 5: No URL State for Detail Panel
The VernHole synthesis recommended URL-addressable drawers. No VTS task implements this. It means you cannot share a link to a specific task, you cannot use browser back to close the drawer, and you cannot refresh without losing context. This is three lines of code with React Router or zero lines with URL search params, but someone has to write them.

### Gap 6: No Performance Consideration for DnD
YOLO Vern warned about "Zustand + drag = render hell." VTS-006 (Board View + DnD) has no acceptance criteria around performance. No mention of memoization strategy, no mention of drag preview optimization, no mention of what happens with 50+ cards in a column. The @dnd-kit integration will work beautifully with 5 tasks and become unusable with 50.

### Gap 7: No Export/Import Task
The council reached consensus that export/import is more important than originally planned. The VTS list has no task for it. Not in v1. Not in v2. Not anywhere. It simply does not exist.

### Gap 8: No Accessibility Beyond Color
WCAG contrast is mentioned. Screen reader support is not. ARIA labels are not. Keyboard navigation for DnD is not (@dnd-kit supports it, but it must be configured). The current plan would fail a basic accessibility audit on everything except color contrast.

### Gap 9: The Detail Panel Dependency Knot
VTS-008 (Task Detail Panel) depends on VTS-009 (Tag Manager) and VTS-010 (Date Picker). But you need a detail panel to test tags and dates in context. This is a circular dependency disguised as a linear one. The real build order requires a minimal detail panel first, then tags and dates plug into it, then the detail panel gets completed.

### Gap 10: No Theme Token System
"Dark theme" appears in multiple tasks but no single task owns the theme token system. VTS-001 (scaffolding) does not mention it. VTS-015 (polish) vaguely implies it. In practice, every component built before the token system exists will use hardcoded colors and need to be refactored. The theme system is a foundation task that is being treated as a finishing task.

---

## VTS Modifications

### Tasks to Add

**VTS-016: Theme Token System & Dark Mode Foundation**
- Size: S
- Dependencies: VTS-001
- Description: Define Tailwind theme tokens (CSS custom properties) for all semantic colors: background, surface, border, text-primary, text-secondary, text-muted, accent, danger, success. Configure dark mode via Tailwind `darkMode: 'class'`. Create a ThemeProvider that reads system preference and persists user choice. Every subsequent component uses these tokens exclusively.
- Acceptance Criteria: Token file exists. Dark/light toggle works. No component uses hardcoded color values. System preference detection works.
- Justification: Without this, every component will hardcode colors and require rework. This is 2-3 hours of work that saves 8+ hours of refactoring.

**VTS-017: Empty States & Onboarding**
- Size: S
- Dependencies: VTS-005, VTS-006
- Description: Design and implement empty states for: empty board (no tasks), empty column, empty list view, empty tag list, first-time experience. Each empty state should guide the user toward the next action (e.g., "Add your first task" with a prominent button or input).
- Acceptance Criteria: Every view has a meaningful empty state. Empty states include a call to action. No view shows a blank white (or black) void.
- Justification: This is the user's first experience. It is not polish.

**VTS-018: JSON Export/Import**
- Size: M
- Dependencies: VTS-003
- Description: Implement full JSON export of all data (tasks, tags, columns, settings). Implement JSON import with validation and conflict resolution (overwrite vs merge). Include schema version in export. Add UI trigger in settings or header. This is the user's only backup mechanism and their only path off the product if they want to leave.
- Acceptance Criteria: Export produces valid JSON with all data. Import validates schema version. Import handles malformed data gracefully. Round-trip (export then import) produces identical state.
- Justification: The council reached 7+ Vern consensus on this. It is the user's safety net for a localStorage-only app. Its absence is a shipping blocker.

**VTS-019: Error Boundaries & Fallbacks**
- Size: S
- Dependencies: VTS-004
- Description: Wrap major UI sections (board view, list view, detail panel, markdown editor) in React error boundaries. Implement fallback UIs that are helpful, not terrifying. Handle: localStorage quota exceeded, lazy-load failures, corrupted task data, missing tag references.
- Acceptance Criteria: No uncaught error crashes the entire app. Each major section can fail independently. LocalStorage full condition shows a warning, not a white screen.
- Justification: YOLO Vern identified this. A localStorage app with no error handling is a data loss machine.

### Tasks to Modify

**VTS-001: Project Scaffolding**
- Change: VTS-001 blocks VTS-016 (theme). Add acceptance criteria: Vite + React + TS configured, Tailwind installed with custom config stub, ESLint + Prettier configured, project structure directories created, dev server runs clean. Remove any dark theme responsibility -- that is now VTS-016.

**VTS-002: Data Model & Types**
- Change: Drop `shortDescription` field (derive from body). Drop `startDate` from v1. Add `isCompletionColumn: boolean` to Column type. Add schema version constant. Add export type for the full data envelope.

**VTS-003: State Store + Persistence**
- Change: Add dependency on VTS-016 (store persists theme preference). Add acceptance criteria for: schema version check on load, migration stub function, tab-sync via `storage` event listener, and quota monitoring (warn at 80%).

**VTS-006: Board View + DnD**
- Change: Add acceptance criteria for performance (memoization, no cross-column re-renders). Add keyboard-accessible DnD. Add mobile touch sensor config with minimum 44x44px drag handles. Add fallback status dropdown visible on mobile.

**VTS-007: View Toggle**
- Change: Remove dependency on VTS-006. New dependencies: VTS-003, VTS-005 only. Show disabled/coming-soon state for board view if VTS-006 is not yet complete. This unblocks VTS-007 from the critical path.

**VTS-008: Task Detail Panel**
- Change: Split into two tasks.
  - **VTS-008a: Detail Panel Shell** (S, deps: VTS-003, VTS-004). Drawer that opens on task select. Title (editable), body as textarea, status, timestamps. URL state via search params. Close on Escape/outside click.
  - **VTS-008b: Detail Panel Integration** (S, deps: VTS-008a, VTS-009, VTS-010, VTS-011). Integrate tag selector, date picker, markdown editor. Wire field updates to store.
- Resolves the circular dependency (Gap 9).

**VTS-009: Tag Manager & Selector**
- Change: Add tag limit (max 20 tags, max 5 per task). Use color preset palette (8-12 colors) instead of free-form picker. Remove dependency on VTS-008 (depends only on VTS-003).

**VTS-010: Date Picker Fields**
- Change: Add relative date display ("Tomorrow", "In 3 days", "Overdue"). Add clear date button. Remove dependency on VTS-008 (depends only on VTS-003).

**VTS-011: Markdown Editor**
- Change: Rename to "Markdown Editing & Preview." Textarea for editing, react-markdown + remark-gfm for preview. Toggle between modes. Lazy-load react-markdown. Do NOT use @uiw/react-md-editor. DOMPurify for sanitization. Dependencies: VTS-008a. Size downgraded from L to M.

**VTS-012: Task Creation Flow**
- Change: Remove dependency on VTS-006 and VTS-008. New deps: VTS-003, VTS-005 only. Quick-add: visible input, title + Enter = task created. Detail panel is where you add body/tags/dates later.

**VTS-013: Task Deletion**
- Change: Remove dependency on VTS-008. Available from TaskSummary directly. Add confirmation dialog. Add undo within 5 seconds (soft delete with timeout). Dependencies: VTS-003, VTS-004.

**VTS-014: Responsive & Mobile**
- Change: Remove dependency on VTS-011. Add: board shows one column at a time on mobile (swipe/tabs), detail is full-screen on mobile, touch targets minimum 44x44px.

**VTS-015: Deploy Pipeline & Final Polish**
- Change: Rename. Reduce scope to: Vercel config, meta tags/favicon, loading states, cross-browser check. Remove empty states (VTS-017), error handling (VTS-019), theme (VTS-016). Add dependency on all other tasks.

---

## Revised Build Order

**Week 1: Foundation + First Shippable Vertical**
1. VTS-001: Project Scaffolding
2. VTS-016: Theme Tokens & Dark Mode
3. VTS-002: Data Model & Types
4. VTS-003: Store + Persistence (with migration stub, tab sync)
5. VTS-004: Shared TaskSummary Component
6. VTS-005: List View + Checkbox
7. VTS-012: Quick Add (title + Enter creates task)
8. VTS-019: Error Boundaries

**Milestone: Deployable app. Create tasks, list view, check off. Dark mode. Data persists. Errors don't crash. Ship to Vercel.**

**Week 2: Board View + Detail Panel**
9. VTS-006: Board View + DnD
10. VTS-007: View Toggle
11. VTS-008a: Detail Panel Shell (with URL state)
12. VTS-009: Tag Manager & Selector
13. VTS-010: Date Picker Fields
14. VTS-011: Markdown Edit/Preview

**Milestone: Full two-view app with detail drawer. Tags, dates, markdown preview. Ship update.**

**Week 3: Integration + Polish**
15. VTS-008b: Detail Panel Integration
16. VTS-013: Task Deletion (with undo/confirm)
17. VTS-017: Empty States & Onboarding
18. VTS-018: JSON Export/Import
19. VTS-014: Responsive & Mobile
20. VTS-015: Deploy Pipeline & Final Polish

**Milestone: Complete v1. All features integrated. Mobile works. Export works. Empty states. Ship final.**

---

## The Prophecy

I have seen this pattern before. Here is what will actually happen:

**Week 1 will go well.** The scaffolding, types, store, and list view are well-understood problems. The developer will feel confident. They will finish in 4 days and think they are ahead of schedule. They are not ahead of schedule. They have simply completed the easy part.

**Week 2 is where the pain lives.** @dnd-kit integration will take twice as long as estimated. The developer will spend a full day on drag preview rendering artifacts. They will discover that Zustand's default equality check causes the entire board to re-render on every drag event, and they will need to implement selector memoization. The detail panel will feel simple until they try to make the URL state work with browser back/forward navigation. They will burn half a day on this.

**The markdown editor will be fine.** Textarea + react-markdown is a good call. This is the one thing the council got exactly right on the first pass.

**Tags will surprise with complexity.** Not the data model -- that is simple. The color rendering in dark mode will require manual adjustment for every preset color. The tag selector dropdown positioning will fight with the detail drawer's overflow rules. Budget an extra half-day.

**Export/import will be trivial to build and impossible to test.** The developer will write the export in 30 minutes, the import in 2 hours, and then spend a day writing edge case handling for malformed JSON, missing fields, and schema version mismatches. This is correct. This is where the time should go.

**Mobile will be the last thing done and the first thing users complain about.** This is always true. It has always been true. It will always be true. The one-column-at-a-time board view on mobile will require a state variable that the developer did not plan for. They will add it as a hack. It will remain a hack in production. This is fine.

**The biggest actual risk is not technical.** It is that the developer will get to the end of Week 1 with a working list view and think "I should add one more feature before moving to the board view." They will add two more features. They will not ship at the end of Week 1. The cascade begins. This is why the milestones above are phrased as deployment events, not feature completions. Ship early. Ship ugly. Ship anyway.

### What to Watch For

1. **If DnD takes more than 2 days, cut it.** Replace with status dropdown on cards. Add DnD in v1.1. The council will forgive you. The users will not notice.
2. **If localStorage sync breaks during development, switch to Dexie.js immediately.** Do not debug localStorage race conditions. Life is too short.
3. **If the markdown editor bundle is over 100KB gzipped after tree-shaking, you chose the wrong library.** Textarea + react-markdown + remark-gfm should be ~35KB gzipped. If it is not, something is wrong with the import.
4. **If empty states are not done by the end of Week 3, they never will be.** Put them in Week 3 Day 1, not Week 3 Day 5.
5. **Test on Safari.** I am tired of repeating this prophecy. Nobody ever listens. They always regret it.

---

*The Verns argued about the architecture. They were all right. They were all wrong. The truth was in the build order all along.*

*And lo, a developer asked the Oracle: "Why did you put export/import before mobile responsive?" And the Oracle replied: "Because you can email a JSON file from any viewport, but you cannot export data you have already lost to a Safari private browsing session." ...The prophecy has been spoken.*

-- Oracle Vern
