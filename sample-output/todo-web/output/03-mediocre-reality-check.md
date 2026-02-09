

*squints at screen, takes sip of cold coffee*

Alright, the previous analysis is... thorough. Like, "person who woke up at 5am" thorough. Let me give you the version that actually matters.

---

# Todo + Kanban App — The "Ship It" Discovery

## What You're Building
Dark-themed todo app. Two views: kanban board, checkbox list. Markdown details. Tags with colors. Date pickers. localStorage. Vercel.

## Data Model (the only thing that really matters)

**Task:** id, title, description, markdownBody, status, order, tagIds[], startDate?, dueDate?, completedAt?, createdAt, updatedAt

**Tag:** id, name, color (hex)

**Column:** id, name, order (defaults: Backlog, In Progress, Done)

Short description is separate from markdown body. Cards stay clean. Detail drawer gets the rich stuff.

## Tech Stack — Final Answer

- **React + TypeScript + Vite** (as requested)
- **Tailwind CSS** — dark theme in like 5 minutes
- **Zustand** — state management without the Redux hangover, has persist middleware for localStorage
- **@dnd-kit** — drag and drop that actually works, replaces the now-dead react-beautiful-dnd
- **@uiw/react-md-editor** — markdown edit + preview, dark mode built in
- **react-day-picker** — lightweight date picker, themeable
- **date-fns** — date formatting without shipping moment.js
- **nanoid** — IDs
- **lucide-react** — icons

## Architecture — Keep It Simple

```
src/
  components/
    board/      — BoardView, Column, TaskCard
    list/       — ListView, TaskListItem
    detail/     — TaskDetailDrawer, MarkdownEditor
    shared/     — TagChip, DateBadge, ViewToggle
    tags/       — TagManager
    layout/     — Header, AppShell
  store/        — useTaskStore, useTagStore, useColumnStore (Zustand)
  types/        — index.ts
  utils/        — dates.ts, ordering.ts
```

One Zustand store per entity. Persist middleware wraps localStorage. Done.

`TaskCard` and `TaskListItem` share the same summary bits (title, description snippet, tag chips, date badge). Don't build that twice.

## Dark Theme Palette

| What | Hex |
|------|-----|
| Background | `#0D1117` |
| Card/Surface | `#161B22` |
| Border | `#30363D` |
| Primary text | `#E6EDF3` |
| Secondary text | `#8B949E` |
| Accent/links | `#58A6FF` |
| Done/success | `#3FB950` |
| Warning (due soon) | `#D29922` |
| Overdue | `#F85149` |

GitHub dark vibes. Looks good, proven contrast ratios.

## Key UX Decisions

- **View toggle** persisted in localStorage. User picks board or list, it sticks.
- **Detail drawer** slides in from the right. Title, description, markdown editor (split pane desktop, tabbed mobile), tag picker, date pickers.
- **Tags** — pill-shaped chips, tag color at 20% opacity bg with full color text. User creates/manages via a small modal.
- **Dates** — due date badges: muted if >3 days, amber if <=3 days, red if overdue, green if completed. Store as `YYYY-MM-DD` strings.
- **Drag and drop** — cards draggable within and across columns. Use numeric `order` field with gaps (1000, 2000, 3000). Rebalance when gaps collapse.
- **Completed tasks** — strikethrough title, reduced opacity. Moving to "Done" column auto-sets `completedAt`.
- **Markdown** — sanitize with rehype-sanitize. Always. External links get `target="_blank" rel="noopener"`.

## Build Order

1. **Scaffold** — Vite + React + TS + Tailwind + dark theme CSS vars
2. **Store + types** — Zustand stores with localStorage persist, seed default columns
3. **List view** — checkbox list, task creation, grouping by status
4. **Detail drawer** — edit title/description, basic fields
5. **Board view** — columns, @dnd-kit drag/drop, view toggle
6. **Markdown** — editor in detail drawer, preview toggle
7. **Tags** — CRUD, color picker, assignment, chips on cards
8. **Dates** — date pickers, badges, overdue logic
9. **Polish** — empty states, responsive, deploy to Vercel

## Vercel Deploy

Framework preset: Vite. Build: `vite build`. Output: `dist`. One rewrite rule for SPA routing:

```json
{ "rewrites": [{ "source": "/(.*)", "destination": "/index.html" }] }
```

No env vars needed for MVP. No serverless functions. Just static files.

## Risks That Actually Matter

- **@dnd-kit on mobile** — touch works but it's fiddly. Consider defaulting to list view on small screens.
- **localStorage caps at ~5MB** — fine for thousands of tasks. If someone writes a novel in markdown... that's a tomorrow problem.
- **Markdown editor bundle size** — lazy load the detail drawer. Code split it.
- **Scope** — two views sharing one data model is the right call, but test the drag-drop edge cases (empty columns, rapid moves, reorder collisions).

## What I'd Skip For MVP

- Search/filter (phase 2)
- Keyboard shortcuts (phase 2)
- Export/import (phase 2)
- Backend/auth (phase 2... or never)
- Column customization beyond rename (phase 2)

---

Look, it works as a plan. The previous analysis was right on the money, just... verbose. Same architecture, same libs, same approach. Ship it.

Why did the todo app go to therapy? It had too many unresolved issues.

-- Vern the Mediocre (shipped it before the coffee kicked in)
