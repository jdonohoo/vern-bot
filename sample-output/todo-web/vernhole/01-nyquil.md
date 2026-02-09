Plan's solid. Few notes before I fade...

**Data model** — good. `shortDescription` separate from `bodyMarkdown` is right. Don't let markdown bleed into cards.

**Zustand store** — normalize it. One slice per entity. `tasks`, `tags`, `columns`. Selectors do the joins. Don't nest.

**@dnd-kit** — use `sortable` preset. Explicit `order` field yes but use fractional indexing (like `lexorank`) not integers. Avoids reindex storms...

**localStorage** — wrap it. Version key. Schema migration fn on load. `JSON.parse` in a try/catch or you'll corrupt on bad data.

**Markdown** — `react-markdown` + `rehype-sanitize`. Lazy load the editor. Preview is cheap, editing is heavy. Split pane desktop, tab toggle mobile.

**Date picker** — skip massive libs. `input type="date"` works. Style it. Ship it. Add a thin wrapper for overdue badge logic.

**Tags** — enforce min contrast ratio against dark bg programmatically. `wcag-contrast` or roll a 5-line luminance check. Users will pick black on dark gray otherwise.

**What I'd cut from v1:**
- `startDate` — nobody uses it. Just `dueDate`.
- Undo — nice but scope bomb.
- Tag rename — just delete and recreate.

**What I'd add:**
- `Cmd+K` quick-add. Critical for capture speed.
- Persist to `IndexedDB` via `idb-keyval` instead of localStorage. 5MB limit will bite you with markdown bodies.

**Risk the plan undersells:** mobile DnD. `@dnd-kit` touch sensors exist but they're fiddly. Have a fallback — dropdown to move between columns on mobile. Don't ship broken drag on phones.

**Deployment** — Vite + Vercel is zero-config. Just set `base: '/'` and... it works.

done. shipping this is like 2 weekends of focused work...

why did the kanban board break up with the todo list... said they needed more... columns... in the relationship... zzz

-- NyQuil Vern *already asleep on the keyboard*
