# Mediocre Analysis: Kanban Browser Extension

*squints at screen, sips cold coffee*

Ok look, I read the plan. It's... thorough. Too thorough. 47 sections for a todo app. Let me cut through this.

## What Actually Matters

Build a Chrome extension. Popup with 4 columns. Cards move between them. Done.

## My Take on the Plan

**What's right:**
- Manifest V3, `chrome.storage.local`, CSS variables for theming. Yep. Standard stuff.
- The data model is fine. Task with id, title, description, status, dueDate, estimate. Ship it.
- Urgency color logic is correct. Green/yellow/red/neutral. End-of-day interpretation for due dates. Good call.

**What's over-engineered:**
- "Storage Envelope with versioning" — it's v1. There is no v2. We'll refactor later.
- "TaskMap vs flat array" debate — use an array. `Array.filter()` exists. Next.
- "Phase 3: Recurring tasks" — nobody asked for this. Stop.
- 10 edge cases listed like it's a NASA launch. It's a todo app. Handle nulls and move on.
- "Open in Tab view" — that's a tomorrow problem.

**What's missing from the plan:**
- Actually picking a framework. I'll say it: **vanilla JS**. It's a popup. You don't need React for 4 columns and some cards. A lightweight lib like Preact is acceptable if you're soft. But vanilla works.
- The share-as-markdown feature is trivial — string template, `navigator.clipboard.writeText()`. Done in 20 lines.
- For markdown rendering, just use `marked` + `DOMPurify`. Two dependencies. Don't build a parser.

## What I'd Actually Build

```
manifest.json
popup.html
popup.js
styles.css
```

That's it. Four files. Maybe a `marked.min.js` and `purify.min.js` vendored in.

**Storage:** One key in `chrome.storage.local`. JSON blob. Array of tasks + preferences object.

**Drag and drop:** HTML5 native `draggable` attribute. No library. It's janky but it works.

**Date picker:** `<input type="date">`. Native. Free. Done.

**Theme:** `data-theme="dark"` on body, CSS variables swap. 15 lines of extra CSS.

**Estimate field:** `<input type="number">` + a `<select>` for min/hour/day. Don't overthink it.

## The "Open Decisions" Are Not Open

- Framework: Vanilla JS
- Storage: Array
- Markdown: Phase 1, it's just two libs
- Drag and drop: Phase 1, it's native HTML5
- Popup only (tab view is scope creep)

## Realistic Timeline

A decent dev ships this in a weekend. Phase 1 and Phase 2 from the plan are one phase. There's no reason to separate "CRUD" from "drag and drop" — you're building the board once.

## Verdict

The discovery plan is solid but needs a diet. Cut the phases down, skip the architecture astronautics, use vanilla JS, and ship it. The core idea is good — a local-first kanban in your browser is genuinely useful. Just don't let planning become procrastination.

Look, it works.

---

Why did the developer quit his job? Because he didn't get arrays. Ship it.

-- Vern the Mediocre *(5 more minutes...)*
