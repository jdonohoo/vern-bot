*squints at screen, coffee still brewing*

# Mediocre Reality Check

Vernile wrote a novel. Let me cut it down.

## What's Over-Engineered

**The component tree diagram.** Bro, it's a todo app in a popup. You don't need a formal component architecture document. You need to start coding.

**"Preact + HTM" vs "Lit" vs "Vanilla JS with Web Components" debate.** Pick one. Move on. Preact's fine. We're not writing a thesis.

**The color function got a TypeScript signature and doc comments** for what is literally three if-statements. Just write it when you get there.

**Schema versioning and migration runner.** For a v1 browser extension with zero users. You're building migration infrastructure for a data model that doesn't exist yet. That's a tomorrow problem.

**`columnOrder` redundancy "kept in sync on writes."** You just invented a consistency problem to solve a performance problem you don't have. Hundreds of tasks, `Array.find()` is fine. It's not a database.

**The testing strategy section.** Vitest, integration tests, manual testing checklist — for a popup todo app. Write the app first. Test it by clicking on it. Like a human.

**Accessibility deep-dive in v1 spec.** Important, yes. But Vernile just scoped WCAG AA compliance, ARIA labels, and colorblind secondary signals into an MVP. That's how MVPs die.

**Side panel API consideration.** Scope creep dressed up as a "decision for the builder."

## What Actually Matters

The core is simple. Here's what to build:

1. **Manifest V3 extension, popup, one HTML file**
2. **Preact, esbuild, done** — stop debating
3. **4 columns, cards, click-to-edit detail view**
4. **`chrome.storage.local` with a flat task array** — yes, array. You have < 100 tasks. `find()` is O(nothing).
5. **`<input type="date">`** — Vernile got this right, native picker, zero JS
6. **CSS custom properties for theme** — also correct, that part was solid
7. **Color = 3 if-statements at render time**
8. **Drag and drop: HTML5 native** — agreed, no library
9. **Markdown: `marked` + `DOMPurify`** — or honestly, skip markdown for MVP and add it in week 2
10. **Share = generate string, copy to clipboard, show "Copied!" toast** — Vernile's one-click approach is right, the modal idea was wrong

## What Vernile Got Right

Credit where due:

- Native date picker call — perfect
- One-click share vs modal — correct UX instinct
- No urgency color on Done column — good catch
- `due + T23:59:59` boundary — actually matters
- Minimal extension permissions — `storage` only
- No sync, no backend, offline-only — yes

## What I'd Actually Ship

```
Phase 1 (the only phase that matters):
- Extension scaffold
- Board with 4 columns
- Add/move/delete tasks
- Click task → detail panel with title, description (textarea), date, estimate
- Color indicators
- Dark/light toggle
- Share as markdown copy

Phase "later" (you know when):
- Drag and drop
- Markdown editor
- Anything else
```

## The Real Risk

The real risk isn't architecture. It's that someone reads Vernile's 400-line analysis and spends 3 days planning before writing a single line of code. The app itself is a weekend project. Start with `manifest.json` and go.

The storage model is fine as:
```json
{
  "tasks": [...],
  "theme": "dark"
}
```

That's it. Add fields when you need them. You're not building for scale, you're building for you.

Why did the todo app go to therapy? It had too many unresolved issues and kept putting itself in the "later" column. Ship it.

-- Vern the Mediocre *(overslept, still delivered)*
