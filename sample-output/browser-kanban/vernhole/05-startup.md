

# Startup Vern's Take: Kanban Browser Extension

Alright, let me cut through the enterprise-grade planning doc above and get real about what actually matters here.

## The Core Hypothesis

**"People want a lightweight, in-browser task board that's faster to open than Trello/Notion/Jira."**

That's it. That's the bet. The value prop is **zero friction** — click the extension icon, you're managing tasks. No login. No cloud sync. No account creation. No loading spinner.

## What I'd Actually Ship in v1 (The Real MVP)

The discovery plan above is thorough. Too thorough. It's a 13-section planning doc for a todo app. Let me take a machete to it.

### MUST-HAVE (Week 1 ship)

- **3 columns**: Todo | In Progress | Done. Cut "New" — it's the same as "Todo" and users don't need the cognitive overhead of a 4th column. If they disagree, we add it back in 20 minutes.
- **Task cards** with title only on the board
- **Click to expand** → title, plain text description, due date, estimate
- **Color indicators** (green/yellow/red) on cards — this is the killer differentiator over a plain list
- **Move between columns** via dropdown/buttons (not drag-and-drop yet)
- **Dark/light theme** — yes, keep this in v1, it's table stakes in 2026 and takes 30 minutes with CSS variables
- **chrome.storage.local** persistence

### CUT FROM v1 (Nice-to-haves that don't validate the hypothesis)

- **Markdown editor** — Cut it. Plain text with link pasting works fine. Nobody's writing essays in a browser extension popup. If users scream for it, add it in week 2.
- **Drag and drop** — Cut it. Buttons to move cards work. DnD is a time sink with edge cases. Ship buttons, see if people even use the board first.
- **Share as markdown** — Cut it. Cool feature, zero users are asking for it yet. Copy-paste exists. Add it when someone requests it.
- **"Open in tab" view** — Cut it. Solve the popup constraint problem when you actually have it.
- **Estimate unit selection** (min/hour/day) — Cut the unit picker. Just do hours. One field, one unit. Simplify.
- **Export/import, search, filters, recurring tasks** — Phase 3? More like Phase Never-Unless-There's-Traction.

## The Honest Math

The discovery plan has 3 phases. Here's my counter:

| Phase | What | Why |
|-------|------|-----|
| **v0.1** | Board + cards + color urgency + theme + persistence | Validates: do people actually use this daily? |
| **v0.2** | Drag-and-drop + markdown (if requested) | Polish based on real feedback |
| **v0.3** | Share + whatever users are yelling about | Feature expansion driven by data |

## What I'd Actually Measure

Before writing a single line of code beyond v0.1:

1. **Daily active usage** — Are people opening this more than twice a day?
2. **Task completion rate** — Are tasks moving to Done, or is this a graveyard?
3. **Retention at day 7** — Still using it after a week?

If those numbers are dead, pivot. Maybe it becomes a "quick notes" extension. Maybe it becomes a meeting agenda tool. Let the data talk.

## Architecture Hot Take

The discovery plan suggests considering "vanilla vs lightweight library." Here's the answer: **vanilla JS + HTML + CSS**. It's a popup todo board. You don't need React. You don't need Vue. You don't need a build system. You need:

- One `popup.html`
- One `popup.js` (or a couple small modules)
- One `styles.css` with CSS variables for theming
- A `manifest.json`

Total codebase should be under 1,000 lines for v0.1. If it's more, you're overbuilding.

## The Real Risk Nobody Mentioned

The discovery plan lists "popup size constraints" and "markdown XSS" as risks. Those are engineering risks. The actual risk is:

**Nobody uses it.** The browser extension market for todo apps is a graveyard. Momentum, Taskade, Todoist — they all have extensions. The question isn't "can we build this?" It's "why would anyone switch?"

The answer better be: "Because it's stupidly simple and the color-coded urgency view is instantly useful." If that's not enough, no amount of markdown rendering saves you.

## My Recommendation

1. Build v0.1 in a weekend
2. Put it on the Chrome Web Store (free)
3. Share it in 3 communities where people complain about task management
4. Watch the numbers for 2 weeks
5. If people use it → iterate. If not → pivot or kill it.

Stop planning. Start shipping. The market will tell you what to build next faster than any discovery document.

---

Why did the lean startup founder only build half a todo app? Because the other half was a *pivot table*.

-- Startup Vern (shipped it yesterday, iterating today)
