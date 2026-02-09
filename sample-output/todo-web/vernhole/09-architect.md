

# Architect Vern's Analysis: Dark-Themed Todo + Kanban App

*Adjusts reading glasses. Opens whiteboard. Draws a box labeled "the system that will be running when you're asleep."*

Alright. MightyVern did solid work here. The data model is sound, the edge cases are well-cataloged, and the phasing is reasonable. But I'm going to push on the places where "comprehensive" can quietly become "over-engineered," and where a few architectural decisions need sharper edges before anyone writes a line of code.

Let me walk through this the way I'd walk through a design review.

---

## 1. WHAT I AGREE WITH (AND WHY)

**The data model is mostly right.** Normalized entities for Task, Tag, and Column with explicit IDs and an `order` field — that's the correct foundation. Zustand with a normalized store will keep both views honest. No argument there.

**Shared summary component across views.** This is the single most important architectural decision in the whole app. One `TaskSummary` component rendered inside a `BoardCard` wrapper or a `ListRow` wrapper. If someone builds two divergent card components, this project accumulates drift on day two. MightyVern called it. I'm underlining it.

**localStorage with versioned migrations.** For a v1 local-first app, this is the pragmatic call. The migration versioning is critical — without it, you get silent data corruption the first time you change the schema. I've seen this kill apps that were otherwise well-built.

**Sanitized markdown preview.** Non-negotiable. `DOMPurify` or equivalent. Anyone who ships unsanitized user-rendered HTML in 2026 deserves the incident report.

---

## 2. WHERE I PUSH BACK

### 2a. `shortDescription` as a Separate Field — Reconsider

MightyVern recommends separating `shortDescription` from `bodyMarkdown`. I understand the reasoning: keep summaries clean.

But here's the operational reality: **users will not maintain two fields.** They'll fill in the title, skip the short description, write everything in the markdown body, and your cards will be blank. Or they'll duplicate content across both and get confused when edits don't sync.

**My recommendation:** Drop `shortDescription`. Derive the summary from the first N characters (or first line) of `bodyMarkdown`, stripped of formatting. One source of truth. If you want a richer preview, parse the first paragraph. This is less "pure" but dramatically more usable.

If the team strongly wants an explicit subtitle, make it optional and clearly labeled — but don't make it a required part of the card display.

### 2b. The `order` Field Strategy Needs Specifics

"Numeric order field" is correct but incomplete. The devil is in the implementation:

- **Use fractional indexing** (e.g., the `fractional-indexing` library or a simple midpoint calculation). When you drop a task between positions 1.0 and 2.0, assign it 1.5. This avoids renumbering the entire column on every drag.
- **Reindex threshold:** When precision gets too fine (e.g., after 50+ reorderings in the same gap), batch-reindex the column. This is a background operation, not user-facing.
- **Why this matters:** Integer-based reordering means every drag operation writes N updates for N tasks below the insertion point. At 50 tasks per column, that's 50 localStorage writes per drag. Fractional indexing makes it exactly 1 write.

### 2c. Completion Should Be a Status, Not a Boolean

MightyVern's open question #2 — "Should completion be a boolean, a status, or both?" — has a clear answer: **it's a status.**

A task is "Done" when its `statusId` points to a column marked as a completion column. Add a `isCompletionColumn: boolean` flag to the Column entity. When a task moves to that column (by drag or checkbox), `completedAt` gets stamped. When it moves out, `completedAt` gets cleared.

This eliminates the "done without a Done column" edge case entirely. The checkbox in list view is just a shortcut for "move to the first completion column."

No boolean. No dual state. One mechanism, two interfaces.

### 2d. The Detail Panel — Drawer, Not Modal

The plan says "modal or drawer." Pick one: **drawer.** Specifically, a slide-in panel from the right.

Why:
- Modals block context. Users can't see the board while editing.
- A drawer lets users see their kanban columns alongside the detail view on desktop.
- On mobile, a drawer naturally goes full-screen, which is what you want anyway.
- This also eliminates the "switching views discards edits" problem — the drawer is view-independent.

The drawer should be URL-addressable (`?task=<id>`), so deep-linking and browser back work correctly.

### 2e. Column Configuration is Missing

The plan treats columns as static defaults. But users will want to:
- Rename columns
- Reorder columns
- Add/remove columns
- Set which column is the "default" for new tasks

This isn't scope creep — it's table stakes for a kanban app. Without it, you've built a three-column sticky note wall. Add a minimal column management UI (inline rename, drag to reorder, add/delete with confirmation).

---

## 3. ARCHITECTURE I'D PROPOSE

```
src/
├── app/
│   ├── App.tsx                  # Layout shell, routing, drawer host
│   ├── routes.ts                # View routes + task deep-link param
│   └── theme/
│       ├── tokens.ts            # CSS custom properties, dark palette
│       └── global.css           # Reset, typography, token application
│
├── store/
│   ├── index.ts                 # Zustand store, single export
│   ├── slices/
│   │   ├── taskSlice.ts         # Task CRUD, ordering, status changes
│   │   ├── tagSlice.ts          # Tag CRUD, color management
│   │   └── columnSlice.ts       # Column CRUD, ordering
│   ├── selectors.ts             # Memoized selectors for both views
│   ├── persistence.ts           # localStorage adapter + migration runner
│   └── migrations/
│       ├── index.ts             # Migration registry
│       └── v1.ts                # Initial schema
│
├── components/
│   ├── task/
│   │   ├── TaskSummary.tsx      # THE shared summary (tags, date badge, title)
│   │   ├── TaskDetail.tsx       # Drawer content: editor, tags, dates
│   │   └── DateBadge.tsx        # Due/overdue/done badge logic
│   ├── board/
│   │   ├── BoardView.tsx        # Column layout + DnD context
│   │   ├── BoardColumn.tsx      # Single column + droppable
│   │   └── BoardCard.tsx        # Wrapper around TaskSummary
│   ├── list/
│   │   ├── ListView.tsx         # Grouped list + sort controls
│   │   └── ListRow.tsx          # Checkbox + TaskSummary wrapper
│   ├── tags/
│   │   ├── TagChip.tsx          # Single tag display
│   │   ├── TagSelector.tsx      # Multi-select for detail view
│   │   └── TagManager.tsx       # CRUD for tags (settings or inline)
│   ├── dates/
│   │   └── DatePicker.tsx       # Thin wrapper around chosen lib
│   ├── markdown/
│   │   ├── MarkdownEditor.tsx   # Textarea + toolbar
│   │   └── MarkdownPreview.tsx  # Sanitized render
│   └── layout/
│       ├── ViewToggle.tsx       # Board/List switch
│       ├── Drawer.tsx           # Slide-in panel shell
│       └── Header.tsx           # App header, new task, view toggle
│
├── lib/
│   ├── fractionalIndex.ts       # Order calculation utilities
│   ├── markdown.ts              # Parse + sanitize helpers
│   └── dates.ts                 # Date comparison, badge logic
│
└── types/
    └── index.ts                 # Task, Tag, Column, ViewPreference
```

**Key architectural rules:**
1. **Store is the single source of truth.** Components read via selectors, write via actions. No local component state for task data.
2. **TaskSummary is used everywhere.** Board cards and list rows are thin wrappers that add context (drag handle, checkbox) but delegate display to the same component.
3. **The drawer is mounted at the App level,** not inside either view. It persists across view switches.
4. **Persistence is an adapter,** not baked into the store. Swap localStorage for an API client later without touching store logic.

---

## 4. STATE MANAGEMENT — THE ZUSTAND STORE SHAPE

```
{
  tasks: Record<string, Task>,
  taskOrder: Record<string, string[]>,  // columnId -> taskId[] (ordered)
  tags: Record<string, Tag>,
  columns: Record<string, Column>,
  columnOrder: string[],                // ordered column IDs
  
  // UI state
  activeTaskId: string | null,
  viewPreference: 'board' | 'list',
  
  // Persistence metadata
  schemaVersion: number,
}
```

**Why `taskOrder` is separate:** Rather than sorting tasks by their `order` field on every render, maintain an explicit ordered array per column. When a drag completes, splice the array. Derive the fractional index from the array position for persistence. This keeps renders O(1) for ordering and makes DnD operations a simple array splice.

---

## 5. FAILURE MODES — THE 3 AM TEST

| Scenario | What Happens | Mitigation |
|---|---|---|
| localStorage is full | Silent data loss on next save | Check `storage` quota before write. Show persistent warning banner if within 80%. |
| Corrupted localStorage JSON | App crashes on load | Wrap parse in try/catch. If parse fails, offer "reset" or "export raw data" option. Never silently discard. |
| User drags task during slow render | Stale drop target, wrong column | Use `@dnd-kit`'s collision detection with closestCenter. Debounce is the enemy of drag — keep it synchronous. |
| Two tabs open, both editing | Last-write-wins data loss | Listen to `storage` events. Show "data updated in another tab" banner. Reload or merge. |
| Detail drawer open, task deleted (e.g., keyboard shortcut) | Drawer shows stale/missing data | Drawer watches `activeTaskId` — if task no longer exists in store, close drawer with a toast notification. |
| Markdown body is 50KB+ | localStorage bloats, preview lags | Debounce preview rendering (300ms). Warn at 10KB. Hard limit at 100KB with clear messaging. |

---

## 6. WHAT I'D CUT FROM V1

The plan is solid but ambitious. For a shippable v1, I'd defer:

- **Undo/redo.** Nice to have, adds significant complexity to the store. Ship without it, add it when users ask.
- **Advanced filtering and search.** V1 is for small task counts. If you need search, you have too many tasks for a personal todo app.
- **Export/import.** Useful, but not launch-blocking.
- **Start date.** Most users only care about due dates. One date field is simpler. Add start date in v2 if needed.
- **Split markdown editor (side-by-side).** Start with a tab toggle (edit/preview). Side-by-side is a layout headache on smaller screens.

**What I would NOT cut:**
- Column management (rename, reorder, add/delete). Without this, the kanban board is decorative.
- Multi-tab awareness. Two-tab data loss is a support nightmare.
- Schema migrations. Get this right from day one or pay for it forever.

---

## 7. TECHNOLOGY NOTES

- **@dnd-kit** — Correct choice. It's the most maintained React DnD library, handles keyboard DnD for accessibility, and works with Zustand cleanly.
- **Markdown** — Use `react-markdown` with `remark-gfm` for rendering, `DOMPurify` for sanitization. For the editor, a plain `<textarea>` with a formatting toolbar is sufficient. Don't reach for a WYSIWYG editor — they're heavy and fragile.
- **Date picker** — `react-day-picker` is lightweight and themeable. Avoid `react-datepicker` (large bundle, hard to style for dark themes).
- **Zustand middleware** — Use the `persist` middleware with a custom storage adapter that handles versioned migrations. This is built into Zustand and well-documented.
- **CSS approach** — CSS Modules or vanilla-extract for scoping. Tailwind is fine if the team knows it, but for a dark-theme-first app with custom tokens, CSS custom properties give you more control and less fighting.

---

## 8. VTS TASK BREAKDOWN

### TASK 1: Project Scaffolding and Theme Foundation

**Description:** Initialize the Vite + React + TypeScript project. Configure the dark theme token system using CSS custom properties. Set up the app shell layout (header, main content area, drawer mount point). Configure Vercel deployment.
**Acceptance Criteria:**
- Vite project builds and deploys to Vercel
- Dark theme tokens defined and applied globally
- App shell renders with header, content area, and empty drawer slot
- TypeScript strict mode enabled
**Complexity:** M
**Dependencies:** None
**Files:** `vite.config.ts`, `src/app/App.tsx`, `src/app/theme/tokens.ts`, `src/app/theme/global.css`, `src/types/index.ts`, `vercel.json`

### TASK 2: Data Model and Zustand Store

**Description:** Define TypeScript types for Task, Tag, and Column. Implement the Zustand store with normalized entity maps, task ordering per column, and all CRUD actions. Include the persistence adapter with versioned migration support.
**Acceptance Criteria:**
- All entity types defined with explicit fields
- Store supports create, read, update, delete for tasks, tags, and columns
- Task ordering uses fractional indexing
- Persistence writes to localStorage on every mutation
- Schema version tracked; migration runner executes on load
- Corrupted localStorage handled gracefully (try/catch with recovery option)
**Complexity:** L
**Dependencies:** Task 1
**Files:** `src/store/index.ts`, `src/store/slices/taskSlice.ts`, `src/store/slices/tagSlice.ts`, `src/store/slices/columnSlice.ts`, `src/store/selectors.ts`, `src/store/persistence.ts`, `src/store/migrations/index.ts`, `src/store/migrations/v1.ts`, `src/types/index.ts`, `src/lib/fractionalIndex.ts`

### TASK 3: Shared Task Summary Component

**Description:** Build the `TaskSummary` component that renders title, derived description (first line of markdown body), tag chips, and date badges. This single component will be used by both the board card and list row.
**Acceptance Criteria:**
- Renders title, truncated first-line description, tag chips, and date badge
- Date badge shows correct state (due, soon, overdue, completed)
- Tag chips display with correct colors against dark background
- Component accepts no view-specific props — it's purely data-driven
**Complexity:** M
**Dependencies:** Task 2
**Files:** `src/components/task/TaskSummary.tsx`, `src/components/task/DateBadge.tsx`, `src/components/tags/TagChip.tsx`, `src/lib/dates.ts`

### TASK 4: Board View with Drag-and-Drop

**Description:** Implement the kanban board view using `@dnd-kit`. Columns render in order, task cards are draggable within and across columns. Empty columns accept drops. Column management (rename, reorder, add, delete) is included.
**Acceptance Criteria:**
- Columns render in configured order with task cards
- Tasks are draggable within a column (reorder) and across columns (status change)
- Empty columns show droppable placeholder
- Column headers support inline rename
- Columns can be added, deleted (with confirmation if tasks exist), and reordered
- Keyboard drag-and-drop works for accessibility
- Drop updates store immediately with fractional index recalculation
**Complexity:** XL
**Dependencies:** Task 2, Task 3
**Files:** `src/components/board/BoardView.tsx`, `src/components/board/BoardColumn.tsx`, `src/components/board/BoardCard.tsx`

### TASK 5: List View with Completion Toggle

**Description:** Implement the checkbox list view. Tasks are grouped by status, with completed tasks dimmed and sorted to the bottom. Checkbox toggles move tasks to/from the completion column.
**Acceptance Criteria:**
- Tasks render as checkbox rows grouped by column/status
- Checking a task moves it to the completion column and stamps `completedAt`
- Unchecking moves it back to previous column and clears `completedAt`
- Completed tasks are visually dimmed
- Sort options: manual order, due date
- Each row wraps `TaskSummary` with a checkbox prepended
**Complexity:** L
**Dependencies:** Task 2, Task 3
**Files:** `src/components/list/ListView.tsx`, `src/components/list/ListRow.tsx`

### TASK 6: Detail Drawer with Markdown Editor

**Description:** Build the slide-in detail drawer. Includes editable title, markdown editor with edit/preview toggle, tag selector, date pickers, and metadata display. Drawer is mounted at App level and driven by `activeTaskId` in the store.
**Acceptance Criteria:**
- Drawer slides in from right on desktop, full-screen on mobile
- Title is inline-editable
- Markdown editor is a textarea with formatting toolbar
- Preview mode renders sanitized markdown with `DOMPurify`
- Tag selector allows adding/removing tags from the task
- Date pickers for due date (start date deferred)
- Drawer URL-addressable via query parameter (`?task=<id>`)
- Drawer closes gracefully if active task is deleted
- Focus trap when drawer is open
**Complexity:** XL
**Dependencies:** Task 2, Task 7, Task 8
**Files:** `src/components/layout/Drawer.tsx`, `src/components/task/TaskDetail.tsx`, `src/components/markdown/MarkdownEditor.tsx`, `src/components/markdown/MarkdownPreview.tsx`, `src/lib/markdown.ts`

### TASK 7: Tag Management System

**Description:** Build the tag CRUD system: create tags with name and color, edit, delete (with cascading removal from tasks). Include the tag selector for the detail view and a tag manager accessible from settings or inline.
**Acceptance Criteria:**
- Tags can be created with a name and color picker
- Tags can be renamed, recolored, and deleted
- Deleting a tag removes it from all tasks
- Duplicate tag names are prevented
- Tag colors have sufficient contrast against dark background (validation or curated palette)
- Tag selector in detail view supports multi-select
**Complexity:** M
**Dependencies:** Task 2
**Files:** `src/components/tags/TagManager.tsx`, `src/components/tags/TagSelector.tsx`, `src/components/tags/TagChip.tsx`, `src/store/slices/tagSlice.ts`

### TASK 8: Date Picker and Due Date Logic

**Description:** Integrate a lightweight date picker (`react-day-picker`), implement due date badge logic (neutral, soon, overdue, completed), and wire it into the detail drawer.
**Acceptance Criteria:**
- Date picker styled for dark theme
- Due date badge shows correct urgency level
- Completed tasks show "done" state regardless of due date
- Tasks without dates show no badge
- Date stored as ISO date string (no time component)
**Complexity:** M
**Dependencies:** Task 2
**Files:** `src/components/dates/DatePicker.tsx`, `src/components/task/DateBadge.tsx`, `src/lib/dates.ts`

### TASK 9: View Toggle and Preference Persistence

**Description:** Implement the board/list view toggle in the header. Persist the user's preference in the store (and thus localStorage). Ensure switching views does not close the detail drawer or lose context.
**Acceptance Criteria:**
- Toggle switches between board and list views
- Preference persists across page reloads
- Detail drawer remains open during view switch
- Active task selection is preserved during switch
**Complexity:** S
**Dependencies:** Task 4, Task 5, Task 6
**Files:** `src/components/layout/ViewToggle.tsx`, `src/components/layout/Header.tsx`, `src/app/App.tsx`

### TASK 10: Multi-Tab Awareness and Storage Resilience

**Description:** Listen for `storage` events to detect changes from other tabs. Show a notification banner when external changes are detected. Implement storage quota monitoring with warnings.
**Acceptance Criteria:**
- `storage` event listener detects external changes
- Banner prompts user to reload when external changes detected
- Storage usage is checked before writes; warning shown at 80% capacity
- Corrupted data on load shows recovery option, not a crash
**Complexity:** M
**Dependencies:** Task 2
**Files:** `src/store/persistence.ts`, `src/components/layout/StorageBanner.tsx`

---

## 9. WHAT I'D WATCH

1. **The TaskSummary component will be the most-modified file in the project.** Design it for composition (slots/render props for the checkbox, drag handle, etc.) rather than conditional props. If you see `if (view === 'board')` inside TaskSummary, you've gone wrong.

2. **Fractional indexing is simple until it isn't.** Write unit tests for the index calculation — specifically for the "insert at beginning," "insert at end," and "reindex" cases. This is the kind of utility that breaks silently.

3. **The Zustand persist middleware is convenient but opaque.** Write your own thin persistence layer on top of it so you control serialization, migration, and error handling. The middleware's default behavior on parse failure is... not great.

4. **Vercel SPA routing.** You'll need a `vercel.json` rewrite rule to serve `index.html` for all routes. This is a one-liner but it's the kind of thing that wastes an hour if you forget it.

---

## 10. FINAL ASSESSMENT

MightyVern's plan is comprehensive and largely correct. My refinements are:

1. **Drop `shortDescription`** — derive from markdown body
2. **Use fractional indexing** — not integer reordering
3. **Completion is a status** — not a boolean
4. **Drawer, not modal** — URL-addressable, view-independent
5. **Column management is required** — not optional
6. **Cut start date, undo, and split editor** from v1
7. **Add multi-tab awareness** — prevent silent data loss
8. **Write your own persistence layer** — don't trust Zustand's default

The architecture is straightforward. The risk is scope, not complexity. Build the store and TaskSummary first — everything else is UI wiring on top of a solid foundation.

If someone asks "where does X happen?" the answer should always be obvious from the file tree. If it isn't, the architecture has failed.

---

Why did the systems architect add a circuit breaker to their marriage? Because after three failed requests, it's healthier to stop retrying and fall back to the default state. ...I'll see myself out.

-- Architect Vern (measure twice, deploy once)
