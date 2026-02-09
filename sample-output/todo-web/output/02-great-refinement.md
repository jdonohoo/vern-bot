

# Discovery Analysis: Dark-Themed Todo + Kanban Web App

**TypeScript / React / Vite / Vercel**

---

Allow me to illuminate the optimal approach for this application. What follows is a comprehensive discovery document — the kind of blueprint that makes future developers nod in quiet appreciation.

---

## 1. Problem Space & Goals

You want a focused yet capable todo application supporting two mental models:

1. **Kanban Board** — columns, drag-and-drop, visual workflow
2. **List/Checkbox View** — traditional todo list with completion toggling

Each task carries rich metadata: short description, markdown detail body, dates, and user-defined colored tags. The user chooses their preferred view. Dark theme throughout. Deployed on Vercel.

**Success Criteria:**
- Clean dark UI with proper contrast and legibility
- Unified data model powering both views seamlessly
- Markdown editing as a first-class citizen, not bolted on
- Tags and dates visible at a glance without opening detail
- Fast dev experience (Vite) and instant deployments (Vercel)

---

## 2. Product Scope & UX Model

### 2.1 Core Entities

| Entity | Purpose |
|--------|---------|
| **Task** | Unit of work with title, description, markdown body, status, dates, tags |
| **Column/Status** | Kanban lanes (e.g., Backlog, In Progress, Done) |
| **Tag/Category** | User-defined label with assigned color |

### 2.2 Views

**Board View:**
- Vertical columns representing statuses
- Draggable task cards showing title, short description, tag chips, due date badge
- "Done" is typically the final column
- Empty columns remain visible with drop target

**List View:**
- Checkbox toggles completion state
- Inline metadata: tags, due date, short description
- Grouping: incomplete first, then completed (dimmed/struck-through)
- Optional sort by due date

**Detail Panel (shared by both views):**
- Opens as a slide-over drawer or modal
- Title (editable)
- Short description (editable)
- Markdown editor with live preview toggle
- Tag selector with color swatches
- Date pickers: start date, due date
- Completed date (auto-set on done)
- Created/updated timestamps

### 2.3 UX Edge Cases Worth Addressing

- **No tags, no dates:** Show subtle placeholder text, not empty space
- **Overdue tasks:** Due date badge turns warning color (amber/red)
- **Completed + had due date:** Show completed state, NOT overdue
- **Long descriptions:** Truncate with ellipsis in card/list view
- **Start date after due date:** Warn the user inline
- **Markdown on mobile:** Toggle between edit/preview (no split view)
- **Large markdown bodies:** Debounce preview rendering

---

## 3. Functional Requirements

### 3.1 Task Management
- Create, edit, delete tasks
- Mark done/undone (checkbox or column move)
- Drag between columns (board view)
- Reorder within columns (board view)
- Short description visible in both views
- Detail view with full markdown editor

### 3.2 Tags/Categories
- CRUD operations on tags
- Each tag has a name and hex color
- Multiple tags per task
- Color chips displayed on cards and list items
- Tag management UI (settings or inline)

### 3.3 Dates
- Optional start date and due date via calendar picker
- `completedAt` auto-set when task marked done
- Date badges on cards: "Due Feb 10", "Overdue", "Starts Mar 1"
- ISO string storage, local timezone display

### 3.4 View Toggle
- Persistent preference (localStorage)
- Seamless switch — no data loss, no jarring reloads

---

## 4. Data Model

Observe how elegantly this handles both views from a single source of truth:

### Task
```
id:               string (uuid)
title:            string
shortDescription: string
bodyMarkdown:     string
statusId:         string (references Column)
order:            number (position within column)
tagIds:           string[] (references Tags)
startDate?:       string (ISO 8601)
dueDate?:         string (ISO 8601)
completedAt?:     string (ISO 8601)
createdAt:        string (ISO 8601)
updatedAt:        string (ISO 8601)
```

### Tag
```
id:    string (uuid)
name:  string
color: string (hex, e.g. #FF6B6B)
```

### Column (Status)
```
id:     string (uuid)
name:   string
order:  number
color?: string (optional accent for column header)
```

**Key design decisions:**
- `shortDescription` is separate from `bodyMarkdown` — card views stay clean
- `order` is a numeric field, not derived from array index — stable reordering
- Tags are first-class entities with their own IDs, not inline strings — color management stays sane
- `statusId` maps tasks to columns — moving a card just updates one field

---

## 5. Architecture Recommendations

### 5.1 State Management
- **Zustand** — lightweight, TypeScript-friendly, no boilerplate ceremony
- Single normalized store: tasks map, tags map, columns map
- Persist middleware wrapping localStorage
- Abstracted persistence layer for future backend migration

### 5.2 Persistence Strategy
- **Phase 1 (MVP):** localStorage via Zustand persist
- **Phase 2 (optional):** Vercel serverless + database (Supabase, PlanetScale, or Vercel KV)
- Keep the persistence adapter abstracted so swapping is trivial

### 5.3 Recommended Libraries

| Concern | Library | Rationale |
|---------|---------|-----------|
| State | **Zustand** | Minimal API, great TS support, persist middleware |
| Drag & Drop | **@dnd-kit/core** + **@dnd-kit/sortable** | Modern, accessible, well-maintained, tree-shakeable |
| Markdown Editor | **@uiw/react-md-editor** or **Milkdown** | Rich editing with preview, dark theme support |
| Markdown Render | **react-markdown** + **rehype-sanitize** | Safe HTML output, plugin ecosystem |
| Date Picker | **react-day-picker** | Lightweight, themeable, no heavy dependencies |
| Date Utils | **date-fns** | Tree-shakeable, immutable, no moment.js bloat |
| Icons | **lucide-react** | Clean icon set, tree-shakeable |
| Styling | **Tailwind CSS** | Rapid dark theme, utility-first, Vite plugin |
| ID Generation | **nanoid** | Fast, URL-safe, tiny |

### 5.4 Component Architecture

```
App
├── Header (view toggle, settings)
├── ViewContainer
│   ├── BoardView
│   │   ├── Column (droppable)
│   │   │   └── TaskCard (draggable)
│   │   └── AddColumnButton
│   └── ListView
│       ├── TaskListGroup (by status)
│       │   └── TaskListItem (checkbox + metadata)
│       └── AddTaskInput
├── TaskDetailDrawer
│   ├── TitleEditor
│   ├── DescriptionEditor
│   ├── MarkdownEditor (with preview toggle)
│   ├── TagSelector
│   ├── DatePickers
│   └── MetadataFooter
├── TagManager (modal/popover)
└── Providers (DndContext, StoreProvider)
```

**Key insight:** `TaskCard` and `TaskListItem` share a `TaskSummary` sub-component for the common elements (title, description, tags, date badge). This is the way.

---

## 6. Dark Theme Design System

### 6.1 Color Palette

| Token | Hex | Usage |
|-------|-----|-------|
| `--bg-primary` | `#0D1117` | Page background |
| `--bg-surface` | `#161B22` | Cards, panels |
| `--bg-elevated` | `#1C2129` | Modals, dropdowns |
| `--border` | `#30363D` | Dividers, card borders |
| `--text-primary` | `#E6EDF3` | Headings, primary text |
| `--text-secondary` | `#8B949E` | Descriptions, metadata |
| `--text-muted` | `#484F58` | Placeholders, disabled |
| `--accent` | `#58A6FF` | Primary actions, links |
| `--success` | `#3FB950` | Completed states |
| `--warning` | `#D29922` | Due soon |
| `--danger` | `#F85149` | Overdue |

### 6.2 Typography
- **Headings:** Inter or system-ui, bold weight
- **Body:** Same family, regular weight, high readability
- **Monospace (code blocks):** JetBrains Mono or Fira Code

### 6.3 Visual Treatment
- Cards: 1px border with `--border`, subtle shadow on hover
- Tag chips: Pill shape, tag color at ~20% opacity background with full-color text/border
- Date badges: Small, rounded, color-coded by urgency
- Completed tasks: Reduced opacity, strikethrough on title
- Drag ghost: Slight rotation, elevated shadow

---

## 7. Markdown Editor Specification

### 7.1 Desktop Behavior
- Split pane: editor left, preview right (resizable)
- Toolbar: bold, italic, headers, code, link, list, checkbox
- Live preview with debounced rendering (~200ms)

### 7.2 Mobile Behavior
- Tab toggle: "Edit" | "Preview"
- Full-width editor or preview (not split)

### 7.3 Safety
- **Always sanitize** rendered HTML via `rehype-sanitize`
- External links open in new tab (`target="_blank"`, `rel="noopener"`)
- Code blocks get syntax highlighting via `rehype-highlight`

### 7.4 Edge Cases
- Empty body: Show placeholder "Add details in markdown..."
- Extremely large documents: Virtual scrolling or warn at threshold
- Paste rich text: Convert to markdown or strip formatting

---

## 8. Kanban Board Specifics

### 8.1 Drag & Drop Behavior
- **Within column:** Reorder (update `order` field)
- **Across columns:** Move + reorder (update `statusId` + `order`)
- **Visual feedback:** Drop placeholder line between cards
- **Empty columns:** Show drop zone with dashed border

### 8.2 Default Columns
- Backlog, In Progress, Done (user can rename/add/remove)
- "Done" column has special behavior: auto-sets `completedAt`

### 8.3 Ordering Strategy
- Use fractional indexing or gap-based ordering (e.g., orders of 1000, 2000, 3000)
- Rebalance when gaps collapse
- Never derive order from array position alone

---

## 9. Date & Calendar Handling

### 9.1 Picker UX
- Calendar popover from date badge click
- Clear button to remove dates
- Start date picker constrains due date (can't be before start)

### 9.2 Display Logic

| Condition | Badge Style |
|-----------|-------------|
| Due > 3 days away | Default/muted |
| Due within 3 days | Warning (amber) |
| Overdue (past due, not done) | Danger (red) |
| Completed | Success (green), shows completion date |
| No due date | No badge shown |

### 9.3 Timezone
- Store as ISO strings (date-only for simplicity: `YYYY-MM-DD`)
- Display using browser locale
- No timezone conversion needed for date-only values

---

## 10. Project Structure

```
todo-app/
├── public/
├── src/
│   ├── components/
│   │   ├── board/
│   │   │   ├── BoardView.tsx
│   │   │   ├── Column.tsx
│   │   │   └── TaskCard.tsx
│   │   ├── list/
│   │   │   ├── ListView.tsx
│   │   │   └── TaskListItem.tsx
│   │   ├── detail/
│   │   │   ├── TaskDetailDrawer.tsx
│   │   │   └── MarkdownEditor.tsx
│   │   ├── shared/
│   │   │   ├── TaskSummary.tsx
│   │   │   ├── TagChip.tsx
│   │   │   ├── DateBadge.tsx
│   │   │   └── ViewToggle.tsx
│   │   ├── tags/
│   │   │   └── TagManager.tsx
│   │   └── layout/
│   │       ├── Header.tsx
│   │       └── AppShell.tsx
│   ├── store/
│   │   ├── useTaskStore.ts
│   │   ├── useTagStore.ts
│   │   └── useColumnStore.ts
│   ├── types/
│   │   └── index.ts
│   ├── utils/
│   │   ├── dates.ts
│   │   ├── ordering.ts
│   │   └── markdown.ts
│   ├── hooks/
│   │   ├── useViewPreference.ts
│   │   └── useOverdueCheck.ts
│   ├── styles/
│   │   └── theme.css
│   ├── App.tsx
│   └── main.tsx
├── index.html
├── tailwind.config.ts
├── tsconfig.json
├── vite.config.ts
├── vercel.json
└── package.json
```

---

## 11. Vercel Deployment

### 11.1 Configuration
- Framework preset: Vite
- Build command: `vite build`
- Output directory: `dist`
- No server-side functions needed for MVP

### 11.2 `vercel.json`
```json
{
  "rewrites": [{ "source": "/(.*)", "destination": "/index.html" }]
}
```

### 11.3 Environment
- No environment variables needed for localStorage-only MVP
- Add when backend integration begins

---

## 12. Risks & Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| Scope creep from dual views | High | Unified data model, shared components |
| Markdown editor weight | Medium | Lazy-load editor, code-split detail drawer |
| Drag-and-drop edge cases | Medium | Use battle-tested @dnd-kit, explicit order field |
| Dark theme contrast issues | Low | Define palette upfront, test with WCAG checker |
| localStorage limits (~5MB) | Low | Sufficient for thousands of tasks; warn on approach |
| Mobile drag-and-drop UX | Medium | @dnd-kit touch support; consider list-only on small screens |

---

## 13. Implementation Phases

### Phase 1: Foundation
- Vite + React + TypeScript + Tailwind setup
- Dark theme CSS variables
- Data types and Zustand store with localStorage persist
- Default columns seeded

### Phase 2: List View
- Task creation (title + short description)
- Checkbox completion
- Task list with grouping by status
- Basic task detail drawer (no markdown yet)

### Phase 3: Board View
- Column rendering
- @dnd-kit integration (within + across columns)
- Task card component with shared summary
- View toggle with persistence

### Phase 4: Rich Features
- Markdown editor in detail drawer
- Tag CRUD + color picker
- Tag assignment to tasks
- Tag chips on cards/list items

### Phase 5: Dates & Polish
- Date pickers (start, due)
- Date badges with overdue logic
- Empty states, loading states
- Responsive adjustments
- Vercel deployment

### Phase 6: Optional Enhancements
- Search/filter by tag or status
- Keyboard shortcuts
- Export/import JSON
- "Today" / "Upcoming" filtered views

---

## 14. Key Architectural Decisions Summary

1. **Separate `shortDescription` from `bodyMarkdown`** — cards stay clean, detail stays rich
2. **Numeric `order` field** — stable drag-and-drop reordering without array gymnastics
3. **Tags as first-class entities** — color management, reuse, and filtering become trivial
4. **Zustand over Redux/Context** — right-sized for this app, zero boilerplate
5. **@dnd-kit over react-beautiful-dnd** — actively maintained, accessible, performant
6. **Tailwind CSS** — rapid dark theme iteration, utility-first composability
7. **localStorage-first persistence** — zero backend cost for MVP, clean migration path

Excellence is not negotiable. This architecture gives you a robust, maintainable todo app that serves both the "quick checkbox" user and the "visual kanban" user from a single, elegant data model. The markdown detail view elevates it beyond a typical todo app into a genuine productivity tool. Every component has a clear responsibility, every edge case has been considered, and future maintainers will understand the intent immediately.

This is the way.

---

And now, a moment of levity befitting our success: Why did the Kanban board break up with the todo list? Because it needed more columns in the relationship.

-- Vernile the Great (excellence is the only option)
