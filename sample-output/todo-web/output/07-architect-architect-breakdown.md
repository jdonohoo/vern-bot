

# Discovery Analysis: Dark-Themed Todo + Kanban Web App

**TypeScript / React / Vite / Vercel**

---

I've been building production systems since before "microservices" was a word anyone said out loud. Let me tell you — a todo app sounds simple until you're debugging drag-and-drop state corruption at 3 AM because someone thought "it's just a list" and didn't think about ordering semantics.

Let's do this right.

---

## Context and Goals

You want a dark-themed todo app with two primary views: a **kanban board** and a **checkbox list**. Tasks have short descriptions visible in summaries and a rich detail view that supports markdown editing and preview. Users can add colored tags/categories and select dates via calendar pickers. The app is built with TypeScript, React, and Vite, and hosted on Vercel.

This is a local-first, single-user app. No backend for v1. That's the right call — ship something real before you add auth and databases.

---

## Assumptions and Constraints

- Single-page app with client-side state
- No backend required for initial release
- Data persistence is local-first (localStorage), with a future backend migration path
- Dark theme is the default; no light theme toggle needed for v1
- The same underlying task data model serves both views — this is non-negotiable
- The detail view is the canonical editing surface
- Vercel handles hosting with zero-config SPA routing

---

## Success Criteria

1. Users can create, edit, complete, and move tasks in both views without data loss
2. Task cards and list items show title, short description, tags, and date badges at a glance
3. Markdown editing is first-class and safe (sanitized preview)
4. UI is legible and consistent in dark mode across desktop and mobile
5. View preference persists between sessions
6. Tag system supports user-defined colors that remain readable against the dark background

---

## Problem Space Breakdown

### Primary User Goals
1. Quickly capture and complete tasks
2. Organize tasks visually (board) or linearly (list)
3. Add richer context to tasks via markdown details
4. Surface urgency with dates and tags

### Core App Capabilities
1. Task CRUD with consistent state across views
2. Kanban drag-and-drop across columns
3. Checkbox list with completion tracking
4. Tag creation, tagging, and color display
5. Date selection with due/overdue labeling
6. Detail panel with markdown editor and preview

---

## Core Entities and Data Model

This is where most todo apps go wrong. They treat the data model as an afterthought. The data model IS the architecture.

### Task
```
id: string (UUID)
title: string
shortDescription: string
bodyMarkdown: string
statusId: string (FK → Column)
order: number
tagIds: string[] (FK → Tag[])
startDate: string | null (ISO date, no time)
dueDate: string | null (ISO date, no time)
completedAt: string | null (ISO datetime)
createdAt: string (ISO datetime)
updatedAt: string (ISO datetime)
```

### Tag
```
id: string (UUID)
name: string
color: string (hex)
```

### Column (Kanban Status)
```
id: string (UUID)
name: string
order: number
color: string | null (optional accent)
```

### Key Design Decisions

1. **`shortDescription` is separate from `bodyMarkdown`** — summaries must remain clean. Never truncate markdown for card display; that's how you get half-rendered headers on a card.
2. **`order` is an explicit numeric field** — this stabilizes drag-and-drop. Relying on array index is how you get phantom reorders.
3. **Tags are first-class entities** — not strings on a task. This handles colors, renames, and deletions safely.
4. **Status is derived from `statusId`, not a boolean** — this means your kanban columns and list grouping use the same source of truth. A "Done" column is just a column, not special-case logic.
5. **Dates are date-only strings** — no timezone confusion. A due date of "2025-03-15" means March 15th everywhere.

---

## Key Views and Flows

### Board View
```
┌──────────────┐  ┌──────────────┐  ┌──────────────┐
│   To Do      │  │  In Progress │  │     Done     │
├──────────────┤  ├──────────────┤  ├──────────────┤
│ ┌──────────┐ │  │ ┌──────────┐ │  │ ┌──────────┐ │
│ │ Task Card│ │  │ │ Task Card│ │  │ │ Task Card│ │
│ │ Title    │ │  │ │ Title    │ │  │ │ ✓ Title  │ │
│ │ Desc...  │ │  │ │ Desc...  │ │  │ │ Desc...  │ │
│ │ [tag][tag]│ │  │ │ [tag]   │ │  │ │ [tag]    │ │
│ │ Due: 3/15│ │  │ │ Due: 3/12│ │  │ │ Done 3/10│ │
│ └──────────┘ │  │ └──────────┘ │  │ └──────────┘ │
│ ┌──────────┐ │  │              │  │              │
│ │ Task Card│ │  │              │  │              │
│ └──────────┘ │  │              │  │              │
│   + Add Task │  │   + Add Task │  │              │
└──────────────┘  └──────────────┘  └──────────────┘
```

- Columns render in order with all tasks for that status
- Tasks are draggable within and across columns
- Cards show title, short description, tag chips, and date badges
- Empty columns remain visible and droppable
- Clicking a card opens the detail view

### List View
```
┌──────────────────────────────────────────────┐
│ ☐  Fix login bug                             │
│    Redirects to wrong page after OAuth...    │
│    [bug] [urgent]              Due: Mar 12   │
├──────────────────────────────────────────────┤
│ ☐  Update README                             │
│    Add deployment instructions...            │
│    [docs]                      Due: Mar 15   │
├──────────────────────────────────────────────┤
│ ☑  Set up CI pipeline          Done: Mar 10  │
│    Configure GitHub Actions...               │
│    [devops]                                  │
└──────────────────────────────────────────────┘
```

- Checkboxes toggle completion status
- Tasks grouped by status or completion
- Clicking a row opens the detail view
- Completed tasks are visually dimmed

### Detail View (Drawer/Panel)
```
┌─────────────────────────────────────────┐
│  ← Back                          [Save] │
├─────────────────────────────────────────┤
│  Title: [Fix login bug              ]   │
│  Description: [Redirects to wrong...]   │
│                                         │
│  Status: [In Progress ▼]               │
│  Tags: [bug ×] [urgent ×] [+ Add]      │
│  Start: [📅 Mar 10]  Due: [📅 Mar 12]  │
│                                         │
│  ─── Detail (Markdown) ───────────────  │
│  [Edit] [Preview]                       │
│  ┌─────────────────────────────────┐    │
│  │ ## Steps to Reproduce           │    │
│  │ 1. Login with Google OAuth      │    │
│  │ 2. Observe redirect to `/`      │    │
│  │                                 │    │
│  │ **Expected:** `/dashboard`      │    │
│  │ **Actual:** `/` (home page)     │    │
│  └─────────────────────────────────┘    │
│                                         │
│  Created: Mar 9 · Updated: Mar 11      │
└─────────────────────────────────────────┘
```

- Opens from either view (board card click or list row click)
- Edits title, short description, markdown body
- Tag selector with colored swatches
- Date pickers for start and due
- Shows metadata (created/updated timestamps)
- Split editor/preview on desktop, toggle on mobile

### View Toggle
- User selects between board and list via a persistent toggle
- Preference saved to localStorage
- Switching views does not discard selection or edits in progress

---

## Architecture Overview

### Component Tree (Simplified)
```
App
├── ViewToggle
├── BoardView
│   ├── Column (repeating)
│   │   └── TaskCard (repeating, shared)
│   └── DragDropContext
├── ListView
│   └── TaskRow (repeating, uses shared TaskSummary)
├── TaskDetail (drawer/panel)
│   ├── MarkdownEditor
│   ├── TagSelector
│   └── DatePicker
└── TagManager (settings/modal)
```

### Component Responsibilities
1. **TaskSummary** — shared component used by both board cards and list rows. One component, two consumers. This is the linchpin.
2. **TaskDetail** — single detail panel serving all task editing. The canonical editing surface.
3. **Column / ListView** — view-specific wrappers. They own layout, not data.
4. **TagManager** — operates independently and feeds the tag selector.

### State Management Strategy

I'd recommend **Zustand** here. It's small, TypeScript-friendly, and avoids the boilerplate of Redux without the footguns of raw Context for frequently-updating state.

```
Store Shape:
├── tasks: Record<string, Task>
├── columns: Record<string, Column>
├── tags: Record<string, Tag>
├── columnOrder: string[]
├── taskOrderByColumn: Record<string, string[]>
├── activeTaskId: string | null
├── viewPreference: 'board' | 'list'
└── actions: { create, update, delete, reorder, ... }
```

Key principles:
1. **Normalized store** — tasks, tags, and columns stored by ID, not nested
2. **Consistent selectors** for both views to avoid divergent logic
3. **Persistence middleware** to localStorage with migration versioning
4. **Undo capability** is optional but recommended if scope allows

### Persistence Strategy
1. localStorage in v1 with JSON serialization
2. Schema version number stored alongside data
3. Migration functions keyed by version for forward compatibility
4. Clear upgrade path to a backend (Supabase, Vercel KV, etc.) later
5. Warn if storage approaches browser limits (~5MB)

---

## Library Recommendations

| Concern | Recommendation | Rationale |
|---------|---------------|-----------|
| State | Zustand | Small, TS-native, middleware for persistence |
| Drag & Drop | @dnd-kit/core | Modern, accessible, tree-shakeable, actively maintained |
| Markdown Editor | @uiw/react-md-editor or react-markdown + textarea | Depends on how rich you want the editing toolbar |
| Markdown Sanitization | DOMPurify + rehype-sanitize | Non-negotiable for XSS prevention |
| Date Picker | react-day-picker | Lightweight, accessible, easy to theme dark |
| Date Utils | date-fns | Tree-shakeable, no global mutation (unlike moment) |
| Styling | Tailwind CSS or CSS Modules | Tailwind pairs well with Vite and dark themes |
| IDs | crypto.randomUUID() or nanoid | UUID generation for entities |

**What I would NOT use:** Any heavy component library (MUI, Ant Design) for this project. You want control over the dark theme, and fighting a component library's theme system is a time sink you don't need.

---

## Markdown Editor Strategy

### Editor Behavior
- Split view (editor left, preview right) on desktop
- Toggle between edit and preview on mobile
- Debounced preview rendering — 300ms delay for large documents
- Toolbar supports: bold, italic, headers, lists, code blocks, links

### Preview Safety
- **Sanitize ALL output with DOMPurify** — this is not optional
- External links get `target="_blank" rel="noopener noreferrer"`
- Code blocks are syntax-highlighted but don't load heavy dependencies
- No inline HTML execution in preview

### Edge Cases
- Empty markdown shows a helpful placeholder ("Add details...")
- Extremely large content (>50KB) triggers a performance warning
- Pasted rich text is stripped to plain text to avoid rendering surprises

---

## Dates and Calendar Semantics

### Date Types
- `startDate` — optional, when work begins
- `dueDate` — optional, when it's due
- `completedAt` — auto-set when task moves to Done status

### Display Rules
- Due date > 3 days away: neutral badge (muted)
- Due date within 3 days: warning badge (amber)
- Overdue: critical badge (red)
- Completed tasks show "Done [date]" instead of overdue, even if they were late

### Edge Cases
- Start date after due date: show inline warning, don't block save
- Tasks without dates: no badge, no date section in summary
- Date-only storage (no time component) avoids timezone headaches entirely
- Clearing a date is always allowed

---

## Tag System and Color Logic

### Tag Behavior
- Users can create, rename, recolor, and delete tags
- Tasks can hold multiple tags
- Tags shown as colored chips (colored border + tinted background in dark mode)
- Maximum tag name length: ~24 characters (prevents layout breakage)

### Color Handling in Dark Mode
This is where most dark-theme apps get lazy. Tag colors picked by users (bright green, yellow) can be illegible against dark backgrounds.

**Strategy:** Store the user's chosen color as-is, but render chips using a computed low-opacity background with the full-opacity color as a border/text. This keeps user intent while maintaining readability.

```
Chip rendering:
- Background: user color at 15% opacity
- Border: user color at 60% opacity  
- Text: user color at 90% opacity (or white if contrast ratio < 4.5:1)
```

### Edge Cases
- Deleting a tag removes it from all tasks (cascade removal from `tagIds`)
- Duplicate tag names are blocked at creation time
- Tag color collisions are acceptable — users own their palette

---

## Kanban Board Behavior

### Ordering Strategy

This is where I've seen the most bugs in production. Here's how to do it right:

1. Tasks have a numeric `order` field (float, not integer)
2. When dropping between two tasks, new order = midpoint of neighbors
3. When dropping at top or bottom, order = neighbor ± 1.0
4. Reindex all orders in a column when the gap between adjacent orders drops below 0.001

**Why floats?** Because integer reordering requires updating every task below the insertion point. Floats let you insert between any two items with a single update. When gaps get too small (after many reorders), you reindex the whole column — but that's rare.

### Drag-and-Drop Edge Cases
- Touch input on mobile: @dnd-kit handles this, but test thoroughly
- Fast drag across columns: use stable IDs, never rely on position
- Empty columns must accept drops (show a drop zone placeholder)
- Dragging the last task out of a column: column stays visible
- Provide keyboard-based reordering as an accessibility fallback

---

## List View Behavior

### Completion Logic
1. Checkbox toggles `completedAt` timestamp
2. Checking a task also moves it to the "Done" column (statusId change)
3. Unchecking a task moves it back to the first non-Done column (or "To Do" default)
4. This keeps board and list views in sync — completion IS a status change

### Sorting Defaults
1. Incomplete first, completed last
2. Within incomplete: sorted by due date (soonest first), then by creation date
3. Optional user-controlled sort (alphabetical, manual, etc.) as a future enhancement

---

## Theming and Visual System

### Dark Theme Tokens
```
--bg-base:        #0d1117    (GitHub dark-style base)
--bg-surface:     #161b22    (cards, panels)
--bg-elevated:    #1c2128    (dropdowns, modals)
--bg-hover:       #21262d    (interactive hover)

--text-primary:   #e6edf3    (high-emphasis text)
--text-secondary: #8b949e    (low-emphasis, metadata)
--text-muted:     #484f58    (disabled, placeholder)

--border:         #30363d    (dividers, card borders)
--border-active:  #58a6ff    (focus rings, active states)

--accent:         #58a6ff    (primary actions)
--success:        #3fb950    (completed, positive)
--warning:        #d29922    (due soon)
--danger:         #f85149    (overdue, destructive)
```

### Guidance
- Maintain WCAG AA contrast ratios (4.5:1 for normal text)
- Cards use subtle borders, not drop shadows (shadows look wrong in dark mode)
- Focus states are visible and high-contrast
- Tag chips use the computed opacity technique described above

---

## Accessibility Considerations

This isn't optional. It's how you build software that works for everyone.

1. **Keyboard navigation** for view toggle, task selection, and detail panel
2. **Focus trap** in the detail drawer when open
3. **ARIA labels** for drag-and-drop handles (e.g., "Move task: Fix login bug")
4. **Live regions** for status updates ("Task moved to In Progress")
5. **Clear focus indicators** on all interactive elements
6. **Reduced motion** support — disable drag animations if user prefers

---

## Performance and Bundle Strategy

1. **Lazy-load the markdown editor** — it's the heaviest dependency and not needed on initial render
2. **Lazy-load the date picker** — same reasoning
3. **Virtualize task lists** if count exceeds ~100 (react-window or @tanstack/virtual)
4. **Debounce markdown preview** rendering (300ms)
5. **Debounce localStorage writes** (500ms) — don't serialize state on every keystroke
6. **Code-split by view** — board and list views can be separate chunks

---

## Security and Data Integrity

1. **Sanitize markdown preview** with DOMPurify — no exceptions
2. **Validate tag names and colors** at creation (length limits, valid hex)
3. **Schema validation on localStorage read** — malformed data should fall back to defaults, not crash the app
4. **All state mutations go through the store** — no direct localStorage manipulation from components
5. **No eval, no dangerouslySetInnerHTML without sanitization**

---

## Risks and Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| Scope creep from dual views | High | Shared data model, shared TaskSummary component |
| Drag-and-drop instability | Medium | Explicit float `order` field, stable IDs, reindexing |
| Markdown editor bundle size | Medium | Lazy-load editor, use lightweight library |
| Dark theme tag readability | Medium | Computed opacity chip rendering, contrast checks |
| LocalStorage limits (~5MB) | Low | Storage size monitoring, warn at 80% |
| State corruption on schema change | Medium | Versioned migrations, fallback to defaults |

---

## Suggested Project Structure

```
src/
├── components/
│   ├── board/
│   │   ├── BoardView.tsx
│   │   ├── Column.tsx
│   │   └── TaskCard.tsx
│   ├── list/
│   │   ├── ListView.tsx
│   │   └── TaskRow.tsx
│   ├── detail/
│   │   ├── TaskDetail.tsx
│   │   ├── MarkdownEditor.tsx
│   │   └── DatePickerField.tsx
│   ├── shared/
│   │   ├── TaskSummary.tsx
│   │   ├── TagChip.tsx
│   │   ├── DateBadge.tsx
│   │   └── ViewToggle.tsx
│   └── tags/
│       ├── TagManager.tsx
│       └── TagSelector.tsx
├── store/
│   ├── index.ts
│   ├── taskSlice.ts
│   ├── tagSlice.ts
│   ├── columnSlice.ts
│   └── persistence.ts
├── types/
│   └── index.ts
├── utils/
│   ├── dates.ts
│   ├── ordering.ts
│   └── markdown.ts
├── theme/
│   └── tokens.css
├── App.tsx
└── main.tsx
```

---

## Implementation Plan

## VTS Task Index

Individual VTS files: `vts/`

*Revised per Oracle Vision review. 20 tasks total (VTS-008 split into 008a/008b, four new tasks 016-019).*

| # | ID | Task | Complexity | Dependencies |
|---|-----|------|------------|--------------|
| 1 | VTS-001 | Project Scaffolding and Tooling | S | None |
| 2 | VTS-002 | Data Model and Type Definitions | S | None |
| 3 | VTS-003 | State Store with Persistence | M | VTS-002, VTS-016 |
| 4 | VTS-004 | Shared Task Summary Components | M | VTS-002, VTS-003 |
| 5 | VTS-005 | List View with Checkbox Completion | M | VTS-003, VTS-004 |
| 6 | VTS-006 | Board View with Drag-and-Drop | L | VTS-003, VTS-004 |
| 7 | VTS-007 | View Toggle with Persistence | S | VTS-003, VTS-005 |
| 8 | VTS-008a | Detail Panel Shell | S | VTS-003, VTS-004 |
| 9 | VTS-008b | Detail Panel Integration | S | VTS-008a, VTS-009, VTS-010, VTS-011 |
| 10 | VTS-009 | Tag Manager and Tag Selector | M | VTS-003 |
| 11 | VTS-010 | Date Picker Fields | M | VTS-003 |
| 12 | VTS-011 | Markdown Editing and Preview | M | VTS-008a |
| 13 | VTS-012 | Task Creation Flow | S | VTS-003, VTS-005 |
| 14 | VTS-013 | Task Deletion | S | VTS-003, VTS-004 |
| 15 | VTS-014 | Responsive Layout and Mobile Polish | M | VTS-005, VTS-006, VTS-007, VTS-008a |
| 16 | VTS-015 | Deploy Pipeline and Final Polish | M | All other tasks |
| 17 | VTS-016 | Theme Token System and Dark Mode Foundation | S | VTS-001 |
| 18 | VTS-017 | Empty States and Onboarding | S | VTS-005, VTS-006 |
| 19 | VTS-018 | JSON Export/Import | M | VTS-003 |
| 20 | VTS-019 | Error Boundaries and Fallbacks | S | VTS-004 |

## Open Questions for You

1. **Multiple boards or single board?** Right now this is scoped as one board with multiple columns. Do you want named boards/projects?
2. **Completion semantics** — should checking the checkbox always move to "Done" column, or should completion be independent of status?
3. **Priorities** — do you want a priority field (P1/P2/P3) beyond what tags and dates provide?
4. **Collaboration** — any plans for sharing or multi-user? This affects the persistence strategy significantly.
5. **Markdown checklists** — should `- [ ]` checkboxes in markdown body sync with task completion state?

---

## Bottom Line

This plan delivers a coherent, maintainable todo app that supports two mental models (board and list) without splitting your data or UI logic. The architecture is intentionally simple and local-first, with a clean path to future upgrades. The detail view is the single source of truth, while the summary views stay fast, readable, and consistent.

The shared TaskSummary component is the most important architectural decision here — it ensures your two views never diverge. The float-based ordering is the second — it prevents the drag-and-drop bugs that haunt every kanban implementation that "just uses array index."

Ship the list view first. It's simpler and proves the data model. Then layer on the board view — by then, your store and shared components are battle-tested.

---

Why did the todo app break up with the kanban board? Because it felt like their relationship was always "In Progress" and never "Done." ...I'll see myself out.

-- Architect Vern (measure twice, deploy once)
