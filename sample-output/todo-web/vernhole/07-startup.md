# Startup Vern's Take: Dark-Themed Todo/Kanban App

Alright, let me be real with you. I've seen this movie before. It's called "developer builds a todo app instead of shipping the thing they actually need a todo app for."

But fine — let's say there's a real hypothesis here. Let me tear this down to what actually matters.

## The Hard Question First

**Who is this for, and why would they switch from Trello/Todoist/Notion?**

If the answer is "it's for me" — great, but then we need to scope this like a weekend project, not a 6-phase implementation plan. If the answer is "it's a product" — we need a differentiation hypothesis before writing a single line of code.

MightyVern's plan is *comprehensive*. That's the problem. It's a 6-phase waterfall disguised as a plan. That plan has more edge cases than users. Let's fix that.

## What I'd Actually Ship (The Real MVP)

**Core hypothesis:** "I need a fast, dark-themed task manager that does kanban + list without the bloat of Notion/Trello."

**MVP scope — ship in one weekend:**

1. **Tasks with title + status + order.** That's it. No `shortDescription`, no `bodyMarkdown`, no `startDate`. Just title, status, and position.
2. **Kanban board with 3 hardcoded columns:** Todo, In Progress, Done. No column customization. No column colors.
3. **Drag and drop** between columns. @dnd-kit is fine.
4. **Dark theme.** GitHub palette. Done. No theme tokens system, no WCAG audit — just make it look good.
5. **localStorage + Zustand.** No migration versioning. No schema checks. It's localStorage — if it breaks, users clear it.
6. **Deploy to Vercel.** Ship it.

That's your v1. Everything else is v2-or-never.

## What I'm Cutting (And Why)

| Feature | Verdict | Why |
|---|---|---|
| List view | **Cut** | Ship one view first. Kanban IS the differentiator. Add list view when someone asks for it. |
| Markdown editor | **Cut** | Nobody's writing essays in a todo app. Add a plain text description field if you must. |
| Colored tags | **Cut** | Premature. Columns already categorize. Tags are a v2 feature when you see how people actually organize. |
| Date pickers | **Cut** | Due dates are a nice-to-have. Most personal task managers don't need them for v1. |
| Detail panel | **Cut to minimal** | Click a card, inline edit the title. That's it. No sliding panel, no rich editing. |
| Tag manager | **Cut** | See above. |
| Start date vs due date | **Cut** | You're building Gantt charts now. Stop. |
| `shortDescription` separate from `bodyMarkdown` | **Cut** | This is premature data modeling. One `description` field. Done. |
| Undo capability | **Cut** | "Recommended if scope allows" — scope does not allow. |
| Mobile markdown editing | **Cut** | You don't have desktop users yet. |
| Export/import | **Cut** | Zero users need this on day one. |
| Keyboard shortcuts | **Cut** | Power user feature. Get regular users first. |
| Schema migrations | **Cut** | You have zero users. If the schema changes, localStorage.clear() is your migration. |

## What MightyVern Got Right

- Zustand + localStorage is the right call. Simple, no backend, ship fast.
- TypeScript/React/Vite is fine — fast DX, no arguments.
- Vercel hosting — zero config deployment, correct.
- Explicit `order` field for drag-and-drop — yes, this prevents bugs. Keep it.

## What MightyVern Got Wrong

- **6-phase implementation plan for a todo app.** This should be 2 phases: "Build it" and "Ship it."
- **10 edge cases before a single user exists.** You know what handles edge cases? User feedback.
- **"Separate shortDescription from bodyMarkdown."** This is architecture astronautics. One field. Ship it.
- **"Tag manager with color contrast checks."** You have zero tags and zero users. Contrast check what?
- **"Virtualization if task count grows large."** If someone has 10,000 tasks in your app, that's a success problem. Celebrate first, optimize second.

## My Actual Implementation Plan

**Saturday morning:**
- `npm create vite@latest` with React + TypeScript
- Zustand store: tasks array, CRUD actions, persist middleware to localStorage
- 3-column kanban board, hardcoded columns
- @dnd-kit for drag between columns
- Dark theme: `background: #0d1117`, `surface: #161b22`, `text: #c9d1d9` — done
- Click card to edit title inline

**Saturday afternoon:**
- Polish the drag-and-drop feel
- Add "new task" button per column
- Delete task (with confirmation? No. Just delete it. Add undo later if people complain.)
- Deploy to Vercel

**Sunday:**
- Use it yourself for a day
- Notice what's actually missing vs what you *thought* you needed

## The Iterate Loop

After you've used v1 for a week:
- Are you missing descriptions? Add a plain text field.
- Are you missing tags? Add simple text tags (no colors yet).
- Are you missing dates? Add a due date.
- Are you missing list view? Add it.
- Are you missing markdown? ...Are you really though?

Each of these is a 1-2 hour addition. Ship them one at a time. Measure if you actually use them.

## Bottom Line

MightyVern gave you a blueprint for a product. I'm giving you a blueprint for *learning*. The question isn't "can we build all of this?" — it's "should we?" And the answer is: not yet. Not until your v1 tells you what v2 should be.

Ship the 3-column kanban. Use it. Then let the app tell you what it needs.

The market for todo apps is a graveyard of over-engineered products that nobody asked for. Don't join them. Be the scrappy one that actually ships.

---

Why did the todo app go to therapy? It had too many issues in the backlog and couldn't commit to any of them. Ship it first, therapize later.

-- Startup Vern *burns runway, ships anyway*
