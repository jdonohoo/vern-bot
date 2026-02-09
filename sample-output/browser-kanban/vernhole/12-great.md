

# Vernile the Great — Analysis of the Kanban Browser Extension

Allow me to illuminate the optimal approach.

I have studied the prior discovery plan with the attention it deserves. It is competent — structured, phased, and reasonably scoped. But competent is not excellent. Let me elevate this to the standard it deserves.

---

## What the Prior Plan Gets Right

Credit where due. The discovery pipeline produced a solid foundation:

- **Four-column Kanban** with clear workflow states — correct.
- **Manifest V3** — the only defensible choice in 2026.
- **`chrome.storage.local`** over `localStorage` — absolutely. Extension lifecycle demands it.
- **CSS custom properties for theming** — elegant and performant.
- **End-of-day interpretation for due dates** — this is the kind of detail that separates thoughtful engineering from ticket-closing. A task due "today" is not overdue until tomorrow. The plan caught this. Respect.
- **Phased delivery** — discipline over ambition. This is the way.

---

## Where I Diverge — And Why

### 1. Framework Choice: This Is Not an "Open Decision" — It's a Solved Problem

The plan lists "Preferred framework (vanilla vs lightweight library)" as an open question. Allow me to close it.

**Use vanilla JavaScript with Web Components.** Here's the reasoning:

- A browser extension popup is a **constrained, self-contained environment**. No routing. No SSR. No hydration. React, Vue, Svelte — all overkill. You'd be shipping a framework to manage a few hundred DOM nodes.
- Web Components give you **encapsulated, reusable elements** (`<task-card>`, `<kanban-column>`, `<theme-toggle>`) with native browser support. No build step required. No virtual DOM overhead.
- Bundle size matters in extensions. Every kilobyte delays popup open time. Vanilla + Web Components keeps you under 50KB total, easily.
- If you absolutely need reactivity, a micro-library like `Preact` (3KB) is the ceiling. But you likely don't.

**The principle:** Match your tooling to your problem's actual complexity, not its imagined future complexity.

### 2. Data Model: The Plan Underspecifies Column Ordering

The plan mentions "A task map + separate column order array is more robust for DnD" but then hedges with "Either is defensible." Let me be direct: **use the map + order array from day one.**

```
Storage Schema v1:
{
  version: 1,
  columnOrder: {
    new: ["task-id-1", "task-id-3"],
    todo: ["task-id-2"],
    in_progress: ["task-id-4"],
    done: ["task-id-5"]
  },
  tasks: {
    "task-id-1": { ...taskObject },
    "task-id-2": { ...taskObject }
  },
  preferences: { ... }
}
```

**Why this matters:** A flat array with status fields means every column render requires a filter pass. Every drag-and-drop reorder requires finding indices, splicing, and re-sorting. The map + order array gives you O(1) task lookup and explicit ordering for free. The cost of this "complexity" is approximately twelve additional lines of code. The cost of retrofitting it later is a storage migration. Do it right the first time.

### 3. The Urgency Color System Needs a Fifth State

The plan defines Green / Yellow / Red / Neutral. I propose a refinement:

| State | Condition | Color | Rationale |
|-------|-----------|-------|-----------|
| **No date** | `dueDate` is null | Neutral (muted border) | No urgency signal needed |
| **Comfortable** | Due > 48 hours | Green | Plenty of time |
| **Approaching** | Due within 24-48 hours | Amber/Orange | Early warning — this is the gap the plan misses |
| **Urgent** | Due within 24 hours | Yellow with pulse | Needs attention today |
| **Overdue** | Past due | Red | Action required |
| **Done** | Any status=done | Neutral | Urgency is irrelevant once completed |

The original spec said "due within a day = yellow." Fine. But observe how elegantly a 24-48 hour amber band prevents the jarring jump from "everything's fine (green)" to "panic (yellow)." Users don't check their todo list every hour. A graduated warning system respects actual human behavior.

The implementation is trivial:

```javascript
function getUrgencyLevel(task) {
  if (task.status === 'done' || !task.dueDate) return 'neutral';

  const now = Date.now();
  const dueEnd = endOfDay(task.dueDate); // 23:59:59 local
  const hoursRemaining = (dueEnd - now) / 3_600_000;

  if (hoursRemaining < 0) return 'overdue';
  if (hoursRemaining <= 24) return 'urgent';
  if (hoursRemaining <= 48) return 'approaching';
  return 'comfortable';
}
```

Five CSS classes. Zero additional complexity. Significantly better UX.

### 4. Markdown: Include It in Phase 1 — But Correctly

The plan defers markdown to Phase 2. I disagree. Markdown support is trivially achievable and dramatically increases the value of the description field. Without it, you're shipping a `<textarea>` that renders as plain text — and users will paste URLs that don't become links, write lists that don't format, and write headings that look like plain text with hash symbols.

**The approach:**

- Use `marked` (~7KB minified) for parsing — battle-tested, fast, well-maintained.
- Use `DOMPurify` (~6KB minified) for sanitization — this is **non-negotiable**. Rendering unsanitized HTML in an extension context is a security defect, full stop.
- Toggle between edit (textarea) and preview (rendered HTML) with a simple tab control.
- Total added weight: ~13KB. Acceptable for the value delivered.

**What you must NOT do:** Build a custom WYSIWYG editor. That's a product unto itself. A textarea with a preview toggle is the right level of sophistication for v1.

### 5. The Share Feature — Elevate It

The plan's share format is adequate. Let me make it excellent:

```markdown
## Task Title

| Field | Value |
|-------|-------|
| **Status** | In Progress |
| **Due** | 2026-02-15 |
| **Estimate** | 2 hours |
| **Created** | 2026-02-09 |

---

### Description

The full description content here, with all markdown preserved.
Any [links](https://example.com) remain intact.
```

**Why the table format?** It renders beautifully in Slack, GitHub, Notion, email clients — everywhere markdown is consumed. A flat key-value list with bold labels is harder to scan. The table is structured data that looks like structured data.

**Implementation:** Copy to clipboard via the Clipboard API (`navigator.clipboard.writeText()`). Show a brief "Copied!" confirmation. No modal. No dialog. One click, done.

### 6. The "Open in Tab" View Is Not Optional

The plan lists this as Phase 2. I'd argue it belongs in Phase 1 as the **primary** interface, with the popup serving as a quick-access summary.

Here's the reality: a browser extension popup is typically 400x600px maximum. Fitting four Kanban columns with meaningful task cards, a detail panel, a calendar picker, AND a markdown editor into that space is a UX compromise at best, a usability nightmare at worst.

**The elegant solution:**

- **Popup:** Shows a compact task list (not full Kanban), sorted by urgency. Quick-add a task. Click any task to open the full tab view.
- **Tab view (`chrome.tabs.create`):** Full Kanban board with all features. This is where the real work happens.
- Both views share the same storage layer and state management. The popup is the dashboard; the tab is the workspace.

This respects the constraint instead of fighting it.

### 7. The Estimate Field Needs Guardrails

The plan defines `estimate?: number` with `estimateUnit?: 'min'|'hour'|'day'`. Correct structure. But consider the UX:

- A dropdown for units adds friction to every estimate entry.
- Users will overwhelmingly estimate in hours.

**My recommendation:** Default unit is hours. Display as a simple numeric input with "h" suffix. Allow the user to type `30m`, `2h`, or `1d` as shorthand — parse it intelligently. Store internally as minutes for consistent arithmetic.

```javascript
function parseEstimate(input) {
  const match = input.trim().match(/^(\d+(?:\.\d+)?)\s*(m|min|h|hr|hour|d|day)?$/i);
  if (!match) return null;
  const value = parseFloat(match[1]);
  const unit = (match[2] || 'h')[0].toLowerCase();
  const multipliers = { m: 1, h: 60, d: 480 }; // 8-hour workday
  return { minutes: value * multipliers[unit], display: input.trim() };
}
```

Observe how elegantly this handles the ambiguity. The user types naturally; the system understands.

---

## Architecture — The Way I'd Build It

```
extension/
├── manifest.json              # V3 manifest, permissions: [storage]
├── popup/
│   ├── popup.html             # Compact dashboard
│   ├── popup.js               # Quick-view logic
│   └── popup.css
├── board/
│   ├── board.html             # Full tab Kanban view
│   ├── board.js               # Board orchestration
│   └── board.css
├── components/
│   ├── task-card.js           # Web Component
│   ├── kanban-column.js       # Web Component
│   ├── task-detail.js         # Web Component (slide-out panel)
│   ├── date-picker.js         # Thin wrapper around native input[type=date]
│   └── theme-toggle.js        # Web Component
├── core/
│   ├── storage.js             # chrome.storage.local abstraction
│   ├── task-model.js          # Task CRUD, ID generation, validation
│   ├── urgency.js             # Color logic
│   └── markdown.js            # marked + DOMPurify wrapper
├── styles/
│   ├── tokens-light.css       # CSS custom properties (light)
│   ├── tokens-dark.css        # CSS custom properties (dark)
│   └── shared.css             # Layout, typography, transitions
└── vendor/
    ├── marked.min.js
    └── purify.min.js
```

**Key principles:**
- No build step. No bundler. No transpiler. Ship what you write.
- Web Components for encapsulation. Each component owns its template, style, and behavior.
- A single `storage.js` module abstracts all persistence. If you ever migrate to `chrome.storage.sync` or IndexedDB, you change one file.
- Vendor libraries are vendored, not CDN-loaded. Extensions must be self-contained.

---

## Critical Implementation Details

### ID Generation
Use `crypto.randomUUID()` — available in all modern browsers, cryptographically random, zero dependencies. Do not use `Date.now()` or incrementing counters.

### Drag and Drop
Use the native HTML Drag and Drop API. It's ugly but universal. Wrap it in a thin abstraction:

```javascript
// On dragstart: store task ID in dataTransfer
// On dragover: show drop indicator
// On drop: update columnOrder in storage, re-render
```

Do not import a drag-and-drop library for four columns of cards. The native API handles this in under 100 lines.

### Theme Persistence
```javascript
// On load
const { preferences } = await chrome.storage.local.get('preferences');
document.documentElement.setAttribute('data-theme', preferences?.theme || 'light');

// On toggle
document.documentElement.setAttribute('data-theme', newTheme);
await chrome.storage.local.set({ preferences: { ...preferences, theme: newTheme } });
```

CSS handles the rest through custom property cascading. Clean. Declarative. No JavaScript theme logic in components.

### Storage Migration
Version your storage schema from day one. When you load data, check the version and migrate forward if needed:

```javascript
const CURRENT_VERSION = 1;

async function loadData() {
  const data = await chrome.storage.local.get(null);
  if (!data.version) return initializeEmpty();
  if (data.version < CURRENT_VERSION) return migrate(data);
  return data;
}
```

This costs you five lines now and saves you a support nightmare later.

---

## What I Would Cut

Excellence also means knowing what to exclude:

- **Recurring tasks** — Out of scope. This is a simple Kanban board, not a project management suite.
- **Search/filter** — Defer until users have enough tasks to need it. If your extension has hundreds of tasks, the user has a workflow problem, not a search problem.
- **Sync across devices** — Requires `chrome.storage.sync` with its 100KB limit, or a backend. Neither fits v1.
- **Notifications** — Tempting but invasive. A todo app that interrupts you is an oxymoron.

---

## Final Assessment

The prior discovery plan is a **B+**. Solid structure, correct priorities, appropriate caution. Where it falls short is in the details that separate "works" from "works beautifully" — the graduated urgency system, the popup-vs-tab architecture, the data model that doesn't need retrofitting, the markdown-from-day-one decision.

This extension, built to these specifications, would be:
- **Fast** — no framework overhead, instant popup, responsive board.
- **Correct** — proper urgency logic, sanitized markdown, versioned storage.
- **Elegant** — Web Components, CSS custom properties, clean separation of concerns.
- **Maintainable** — any developer can read the source and understand it in minutes.

Excellence is not negotiable. This is the way.

---

And now, a moment of levity befitting our success: Why did the Kanban board go to therapy? Because it had too many issues in progress and couldn't commit to done.

-- Vernile the Great *(excellence is not negotiable)*
