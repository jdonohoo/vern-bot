

# Discovery Analysis: Todo Web App (Kanban + List)

## Let's Cut to the Chase

You want a todo app. Cool. The world has 10,000 of them. That means the market is validated — but it also means you need to ship FAST and find your angle, or you're just building a portfolio project. Let's scope this tight.

## Core Hypothesis

**"Users want a simple, dark-themed task manager that bridges the gap between a basic checklist and a full project management tool — with markdown-powered detail views."**

That's your wedge. Most todo apps are too simple (Apple Reminders) or too complex (Jira). You're aiming for the sweet spot.

## MVP Feature Breakdown (Ruthlessly Prioritized)

### MUST-HAVE (Ship This First)

| Feature | Why It Makes Money | Effort |
|---|---|---|
| Dark theme | Table stakes in 2026. Non-negotiable. | Low — use a component library with dark mode baked in |
| List view (checkbox toggle) | Core UX. This is the "todo" in todo app. | Low |
| Kanban board view | This is what differentiates you from Apple Reminders | Medium |
| View toggle (list ↔ kanban) | User choice = user retention | Low |
| Task card with title + description preview | Users need to scan tasks at a glance | Low |
| Task detail panel/modal | Where the real work happens | Medium |
| Markdown editor/previewer in detail | THIS is your differentiator. Ship it. | Medium (use a library) |
| Tags/categories with colors | Visual organization = dopamine = retention | Medium |
| Due dates with calendar picker | Everyone needs deadlines | Low-Medium |
| Local storage persistence | You need data to persist. No backend yet. | Low |

### CUT FROM V1 (Nice-to-Have, Iterate Later)

- User auth / accounts — localStorage first, auth later
- Backend / database — ship client-only, add Supabase or Firebase when you validate
- Drag-and-drop between kanban columns — use buttons/dropdowns to move cards initially
- Search / filtering — cut it, ship without it
- Recurring tasks — not core
- Subtasks / checklists within tasks — v2
- Notifications / reminders — v2
- Export / import — v2
- Collaboration / sharing — way v2
- Mobile-optimized responsive design — desktop-first, responsive later

## Tech Stack Decisions (Buy, Don't Build)

```
Framework:       React + TypeScript + Vite  ✅ (as requested)
Hosting:         Vercel                      ✅ (as requested)
UI Components:   shadcn/ui                   (dark mode built-in, copy-paste components, no vendor lock)
Styling:         Tailwind CSS                (ships with shadcn, fast to iterate)
Markdown:        @uiw/react-md-editor        (editor + preview in one, actively maintained)
  OR:            react-markdown + textarea   (simpler, less polished, faster to ship)
Calendar:        shadcn date picker          (built on react-day-picker, already in the ecosystem)
State:           Zustand                     (tiny, simple, perfect for local-first apps)
Persistence:     localStorage via Zustand middleware  (zero backend, instant ship)
Icons:           Lucide React                (ships with shadcn)
Drag-and-drop:   SKIP FOR V1                (move tasks via dropdown/buttons)
```

**Why shadcn/ui over Material UI or Chakra?** Copy-paste model means no bloated dependency. Dark mode is trivial. Components are yours to modify. Perfect for MVPs that need to look good without fighting a design system.

**Why Zustand over Redux or Context?** Three lines of code to set up. Persistence middleware for localStorage in one line. You're not building a state management thesis, you're shipping a todo app.

## Data Model (Keep It Lean)

```typescript
interface Task {
  id: string;                  // crypto.randomUUID()
  title: string;
  description: string;         // Short preview text
  detail: string;              // Markdown content
  status: 'todo' | 'in-progress' | 'done';
  tags: string[];              // References to tag IDs
  dueDate: string | null;      // ISO date string
  createdAt: string;
  updatedAt: string;
}

interface Tag {
  id: string;
  name: string;
  color: string;               // Hex color
}

interface AppState {
  tasks: Task[];
  tags: Tag[];
  viewMode: 'list' | 'kanban';
  // CRUD operations...
}
```

That's it. No user model. No project model. No workspace model. Ship it.

## Architecture (One Page App)

```
src/
├── components/
│   ├── layout/
│   │   └── AppShell.tsx          # Dark theme wrapper, header, view toggle
│   ├── kanban/
│   │   ├── KanbanBoard.tsx       # Three columns: todo, in-progress, done
│   │   └── KanbanCard.tsx        # Card with title, desc preview, tags, due date
│   ├── list/
│   │   ├── ListView.tsx          # Checkbox list grouped by status
│   │   └── ListItem.tsx          # Row with checkbox, title, tags, due date
│   ├── task/
│   │   ├── TaskDetail.tsx        # Modal/panel with full detail + markdown editor
│   │   ├── TaskForm.tsx          # Create/edit form (title, desc, tags, date)
│   │   └── TagPicker.tsx         # Tag selector with color dots
│   └── ui/                       # shadcn components (button, dialog, badge, etc.)
├── store/
│   └── useStore.ts               # Zustand store with localStorage persistence
├── lib/
│   └── utils.ts                  # shadcn utility (cn function)
├── App.tsx
└── main.tsx
```

No routing. No pages directory. It's a single-page todo app. Don't over-engineer it.

## Key UI Decisions

**View Toggle:** Simple segmented control in the header. Two buttons: "List" and "Board." Done.

**Kanban Columns:** Three fixed columns — Todo, In Progress, Done. Don't let users create custom columns in v1. That's scope creep wearing a trench coat.

**Task Detail:** Opens as a side panel (sheet) on desktop. Modal on mobile. The markdown editor lives here.

**Tags:** Inline badge components with colored backgrounds. Tag management via a simple dialog — name + color picker. Pre-populate 4-5 default tags (Bug, Feature, Urgent, Personal, Work).

**Due Dates:** Date picker in the task form. Show relative dates on cards ("Tomorrow", "Overdue", "Dec 15"). Color-code overdue dates red.

**Dark Theme:** Set dark mode as default. Don't even build a light mode toggle for v1. Ship dark, iterate later.

## Shipping Plan (Get to Vercel in Days, Not Weeks)

### Day 1: Foundation
- `npm create vite@latest` with React + TS
- Install shadcn/ui, Tailwind, Zustand
- Set up dark theme
- Build the Zustand store with localStorage
- Create the AppShell with view toggle

### Day 2: List View
- ListView + ListItem components
- Checkbox to mark done
- "Add task" inline form
- Basic CRUD working

### Day 3: Kanban View
- KanbanBoard with three columns
- KanbanCard component
- Status change via dropdown on each card
- View toggle working between list and kanban

### Day 4: Detail + Markdown
- TaskDetail side panel
- Integrate markdown editor
- Description preview on cards
- Calendar date picker on task form

### Day 5: Tags + Polish
- Tag management (create, color pick)
- Tags on cards (colored badges)
- Due date display with overdue indicators
- Deploy to Vercel

### Day 6: Ship It
- Vercel deployment config (just push to main, Vercel handles it)
- Add a landing-ish header with app name
- Test on mobile (don't fix everything, just don't break)
- **SHIP IT**

## What "Done" Looks Like for MVP

A user can:
1. Toggle between list and kanban view
2. Create a task with title, description, due date, and tags
3. Open a task and edit rich detail in markdown
4. Check off tasks as done / move between statuses
5. Create colored tags and apply them to tasks
6. See all their data persist on reload (localStorage)

That's it. That validates whether anyone wants this thing.

## What You're NOT Building (Yet)

- Backend / API / database
- Authentication
- Collaboration
- Drag-and-drop reordering
- Search or filtering
- Mobile app
- Notifications
- Analytics

All of those are "scaling problems." You need to earn those problems first by getting users.

## Next Steps After Ship

1. **Put it in front of 10 people.** Watch them use it. Don't explain anything.
2. **Measure:** Do they come back? What do they try to do that they can't?
3. **Decision point:** If retention > 0, add Supabase for persistence + auth. If not, pivot the concept.

## Risk Callouts

- **Markdown editor bundle size** — `@uiw/react-md-editor` is ~200KB. For an MVP, who cares. Optimize later.
- **localStorage limits** — ~5MB. That's thousands of tasks. A good problem to have.
- **No backend = no sharing** — Correct. Ship anyway. Sharing is a v2 feature you add when you know someone cares.

---

Ship the list view first. If you catch yourself designing a database schema before you have a working checkbox, you've gone too far.

Why did the lean startup founder refuse to add a settings page? Because "later" is the best feature you'll never build. Ship it!

-- Startup Vern *(pivoting my way to prod)*
