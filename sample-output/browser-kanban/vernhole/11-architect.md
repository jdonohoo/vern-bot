

# Architect Vern's Systems Analysis: Kanban Browser Extension

*Sits down, opens a blank whiteboard, and asks the question nobody wants to hear:*

**"Before we write a single line of code — how will this fail at 3 AM?"**

---

## 1. Assessment of the Prior Discovery Plan

The MightyVern plan is solid groundwork. Good requirements capture, reasonable phasing, correct identification of the popup constraint problem. But it reads like a feature list wearing an architecture costume. Let me put on the hard hat and talk about **systems**.

What's missing:

- **Component boundaries and responsibility mapping** — who owns what?
- **State management strategy** — how does data flow, and where are the mutation points?
- **Migration strategy** — storage schema *will* change; how do we handle that on day one?
- **Error boundaries** — what happens when `chrome.storage.local` silently fails?
- **Rendering pipeline** — markdown is not "bonus," it's a security surface area decision that affects architecture from the start.

Let me fix that.

---

## 2. Architecture: The Blueprint

### 2.1 Technology Decision

**Vanilla JS + Web Components. No framework.**

Here's why:
- A browser extension popup is ~400x600px of UI. React/Vue/Svelte adds build tooling, bundle size, and cognitive overhead for a problem that doesn't need it.
- Web Components give us encapsulation, reusability, and native browser support — no build step required if we want to keep it simple.
- CSS custom properties handle theming natively.
- If the team later *proves* they need a framework, the component boundaries we define now make migration straightforward.

**Trade-off acknowledged:** No framework means hand-rolling some state management. That's fine — the state surface is small enough that a simple pub/sub store is cleaner than importing 40KB of React for a todo list.

### 2.2 Component Architecture

```
┌─────────────────────────────────────────────────┐
│                  Extension Shell                 │
│  (popup.html / tab.html — same app, two hosts)  │
├─────────────────────────────────────────────────┤
│                                                  │
│  ┌──────────────────────────────────────────┐   │
│  │            <kanban-board>                 │   │
│  │  ┌────────┐┌────────┐┌────────┐┌──────┐ │   │
│  │  │<kanban-││<kanban-││<kanban-││<kanban│ │   │
│  │  │column> ││column> ││column> ││column>│ │   │
│  │  │ "New"  ││ "Todo" ││"In Pr" ││"Done" │ │   │
│  │  │        ││        ││        ││       │ │   │
│  │  │┌──────┐││┌──────┐││        ││       │ │   │
│  │  ││<task- ││││<task- │││        ││       │ │   │
│  │  ││ card> ││││ card> │││        ││       │ │   │
│  │  │└──────┘││└──────┘│││        ││       │ │   │
│  │  └────────┘└────────┘└────────┘└──────┘ │   │
│  └──────────────────────────────────────────┘   │
│                                                  │
│  ┌──────────────────────────────────────────┐   │
│  │         <task-detail-panel>               │   │
│  │  (slides in on task click)                │   │
│  │  - Title editor                           │   │
│  │  - Date picker                            │   │
│  │  - Estimate input                         │   │
│  │  - <markdown-editor> (description)        │   │
│  │  - Share button                           │   │
│  └──────────────────────────────────────────┘   │
│                                                  │
│  ┌──────────┐  ┌──────────────────┐             │
│  │<theme-   │  │ <app-toolbar>    │             │
│  │ toggle>  │  │ (settings, etc.) │             │
│  └──────────┘  └──────────────────┘             │
└─────────────────────────────────────────────────┘
```

**Each component has a single responsibility:**

| Component | Responsibility | Owns |
|---|---|---|
| `<kanban-board>` | Layout, column ordering, drop zones | Column arrangement |
| `<kanban-column>` | Renders tasks for one status, accepts drops | Nothing — reads from store |
| `<task-card>` | Displays task summary + urgency indicator | Urgency color computation |
| `<task-detail-panel>` | Full task editing surface | Form validation |
| `<markdown-editor>` | Edit/preview toggle for description | Markdown parsing + sanitization |
| `<theme-toggle>` | Switches theme, persists preference | Theme state |
| `<date-picker>` | Native date input wrapper | Date formatting |

### 2.3 State Management

A lightweight, observable store. Not Redux. Not MobX. A 40-line pub/sub pattern that any junior developer can read and debug.

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│  UI Actions   │────▶│    Store      │────▶│  UI Render    │
│ (user events) │     │ (single truth)│     │ (subscribers) │
└──────────────┘     └──────┬───────┘     └──────────────┘
                            │
                            ▼
                     ┌──────────────┐
                     │  Persistence  │
                     │  Adapter      │
                     │ (storage.local│
                     │  with debounce│
                     └──────────────┘
```

**Key design decisions:**

1. **Store is synchronous in-memory.** UI never waits for storage. Storage writes are debounced (300ms) and fire-and-forget with error logging.
2. **Persistence adapter is an interface, not a concrete dependency.** Today it's `chrome.storage.local`. Tomorrow it could be IndexedDB or sync storage. The store doesn't know or care.
3. **State mutations go through named actions only.** No direct state manipulation. `store.dispatch('MOVE_TASK', { taskId, newStatus })` — not `state.tasks[3].status = 'done'`.

**Why this matters at 3 AM:** When something goes wrong, you open the console, look at dispatched actions, and trace exactly what happened. No implicit state mutations hiding in component lifecycle methods.

### 2.4 Data Model (Refined)

The prior plan's model is close. Here's what I'd change:

```
StorageEnvelope {
  schemaVersion: 1,              // For migrations. Non-negotiable.
  tasks: {
    [taskId: string]: Task       // Map, not array. O(1) lookups.
  },
  columnOrder: {                 // Separate from task data.
    new: string[],               // Array of taskIds in display order.
    todo: string[],
    in_progress: string[],
    done: string[]
  },
  preferences: {
    theme: 'light' | 'dark',
    markdownEnabled: boolean
  }
}

Task {
  id: string,                    // crypto.randomUUID()
  title: string,
  description: string,           // Raw markdown or plain text
  dueDate: string | null,        // 'YYYY-MM-DD' or null. Not undefined.
  estimateMinutes: number | null, // Store in MINUTES. Display in whatever.
  createdAt: string,             // ISO 8601
  updatedAt: string              // ISO 8601
}
```

**Critical difference from the prior plan: estimates stored as minutes, always.**

The prior plan suggested `estimate + estimateUnit`. That means every piece of code that compares, sums, or sorts estimates needs unit conversion logic. Store canonical. Display flexible. This is the kind of "clever flexibility" that creates bugs.

**Why a task map + column order arrays:**
- Moving a task between columns is a splice from one array and push to another. No re-sorting the entire task list.
- Reordering within a column is an array splice. Clean.
- Looking up a task by ID is O(1), not a `.find()` on an array.

### 2.5 Schema Migration Strategy

```
function migrateStorage(data) {
  if (!data || !data.schemaVersion) {
    return createFreshStore();
  }

  let current = data;

  if (current.schemaVersion === 1) {
    // current = migrateV1toV2(current);
    // current.schemaVersion = 2;
  }

  // Future migrations chain here.
  return current;
}
```

This runs on every extension load. It's 5 lines per migration. It costs nothing. And it means we never have to write a Stack Overflow answer titled "users lost all their tasks after update."

---

## 3. Urgency Color Logic — Getting It Right

The prior plan identified this correctly, but let me be precise about the implementation because **this is the #1 source of user-reported bugs in date-based systems.**

```
function computeUrgency(task) {
  if (task.dueDate === null) return 'neutral';
  if (task.status === 'done')  return 'neutral';   // ← This check MUST be here

  const now = new Date();
  const endOfDueDate = new Date(task.dueDate + 'T23:59:59');

  const hoursRemaining = (endOfDueDate - now) / (1000 * 60 * 60);

  if (hoursRemaining < 0)  return 'overdue';       // Red
  if (hoursRemaining <= 24) return 'due-soon';      // Yellow
  return 'future';                                   // Green
}
```

**Why `T23:59:59`?** Because if I set a due date of February 9th, I mean I have until the end of February 9th to finish it. Not midnight at the start. The prior plan caught this — I'm reinforcing it because I've seen this bug in production more times than I've had hot meals.

**When does this recompute?** On render. Not on a timer. Not on a background interval. The board re-renders when you open the popup. Tasks that were green yesterday might be yellow now. That's correct behavior. If someone leaves the popup open for 6 hours (they won't), we add a single `visibilitychange` listener to re-render when the tab becomes visible again.

---

## 4. Markdown: Architecture Decision, Not a "Bonus"

The prior plan calls markdown a "bonus." I disagree. It's an architecture-level decision that needs to be made upfront because it determines:

1. **What goes in the description field** — raw markdown or plain text?
2. **Security surface** — rendered HTML from user input requires sanitization.
3. **Dependencies** — we need a markdown parser and a sanitizer.

**My recommendation: include markdown from the start, behind a simple toggle.**

- Parser: `marked` (12KB minified, well-maintained, zero dependencies)
- Sanitizer: `DOMPurify` (7KB minified, the gold standard)
- Total cost: ~19KB. For a browser extension, this is nothing.

**Why not defer it?** Because retrofitting markdown into a plain-text description field means either a migration or a dual-mode rendering path. Both are more expensive than just doing it right from the start.

**The editor approach:**
- A `<textarea>` with a "Preview" toggle button. Not a WYSIWYG editor. Not CodeMirror. A textarea and a preview div.
- Edit mode: `<textarea>` with the raw markdown.
- Preview mode: Rendered HTML via `marked` → `DOMPurify` → `innerHTML`.

Simple. Debuggable. No one needs to learn a rich-text framework.

---

## 5. The Share Feature

```
function taskToMarkdown(task) {
  const status = task.status.replace('_', ' ');
  const due = task.dueDate ?? 'No due date';
  const estimate = task.estimateMinutes
    ? formatEstimate(task.estimateMinutes)
    : 'No estimate';

  return [
    `## ${task.title}`,
    '',
    `**Status:** ${status}`,
    `**Due:** ${due}`,
    `**Estimate:** ${estimate}`,
    '',
    '---',
    '',
    task.description || '*No description*'
  ].join('\n');
}
```

On click: `navigator.clipboard.writeText(markdown)`. Show a brief "Copied!" toast. Done.

No sharing APIs. No external services. No permissions. Copy to clipboard and let the user paste it wherever they want. The simplest feature in the whole extension.

---

## 6. Theme System

CSS custom properties on `:root`, toggled via `document.documentElement.dataset.theme`:

```
:root[data-theme="light"] {
  --bg-primary: #ffffff;
  --bg-secondary: #f5f5f5;
  --text-primary: #1a1a1a;
  --text-secondary: #666666;
  --border: #e0e0e0;
  --urgency-green: #2e7d32;
  --urgency-yellow: #f9a825;
  --urgency-red: #c62828;
  --accent: #1976d2;
}

:root[data-theme="dark"] {
  --bg-primary: #1e1e1e;
  --bg-secondary: #2d2d2d;
  --text-primary: #e0e0e0;
  --text-secondary: #999999;
  --border: #444444;
  --urgency-green: #66bb6a;
  --urgency-yellow: #fdd835;
  --urgency-red: #ef5350;
  --accent: #42a5f5;
}
```

**Note the adjusted urgency colors for dark theme.** The light-theme greens and reds are too dark on a dark background. This is the kind of detail that separates "we have dark mode" from "dark mode actually works." Accessibility isn't optional.

---

## 7. File Structure

```
extension/
├── manifest.json
├── popup.html
├── popup.js                    // Entry point, initializes store + renders
├── css/
│   ├── variables.css           // Theme tokens
│   ├── base.css                // Reset, typography
│   └── components.css          // Component-specific styles
├── js/
│   ├── store.js                // State management (pub/sub store)
│   ├── storage-adapter.js      // chrome.storage.local persistence
│   ├── migration.js            // Schema migration logic
│   ├── urgency.js              // Date comparison + urgency computation
│   ├── markdown.js             // marked + DOMPurify wrapper
│   ├── share.js                // Task-to-markdown + clipboard
│   └── components/
│       ├── kanban-board.js
│       ├── kanban-column.js
│       ├── task-card.js
│       ├── task-detail-panel.js
│       ├── markdown-editor.js
│       ├── theme-toggle.js
│       └── date-picker.js
├── lib/
│   ├── marked.min.js           // Vendored, version-pinned
│   └── dompurify.min.js        // Vendored, version-pinned
└── icons/
    ├── icon-16.png
    ├── icon-48.png
    └── icon-128.png
```

Every file has one job. No file exceeds 200 lines. The next developer finds what they need by name alone.

---

## 8. Failure Modes

Because if you don't design for failure, failure designs itself.

| Failure | Impact | Mitigation |
|---|---|---|
| `chrome.storage.local` write fails | Data loss on next session | Retry once, then show non-blocking warning. Keep in-memory state valid. |
| Storage quota exceeded (~5MB) | Can't save new tasks | Monitor used bytes on save. Warn user at 80% capacity. Suggest archiving done tasks. |
| Popup closes mid-edit | Unsaved description changes | Debounced auto-save on every keystroke (300ms). No explicit "save" button needed. |
| Markdown renders malicious HTML | XSS in extension context | DOMPurify. No exceptions. No "we'll add it later." |
| User has 500 tasks in one column | UI becomes sluggish | Virtual scrolling in Phase 2. For Phase 1, cap visible tasks with a "Show more" pattern. |
| Extension update changes schema | Data incompatible | Migration function runs on every load. Schema version is checked first. |
| Date picker returns unexpected format | Urgency logic breaks | Validate date format at the storage boundary. Reject, don't coerce. |

---

## 9. Implementation Plan (VTS Format)

### TASK 1: Extension Scaffold + Manifest

**Description:** Create the Manifest V3 extension skeleton with popup entry point, CSP headers, and icon placeholders. Establish the file structure. This is the foundation — get it right, everything builds cleanly on top.
**Acceptance Criteria:**
- Extension loads in Chrome without errors
- Popup opens with a placeholder page
- Manifest declares only `storage` permission
- File structure matches the architecture spec
**Complexity:** S
**Dependencies:** None
**Files:** `manifest.json`, `popup.html`, `popup.js`, all directory scaffolding

### TASK 2: Store + Persistence Adapter

**Description:** Implement the pub/sub state store with named action dispatch, and the `chrome.storage.local` persistence adapter with debounced writes. Include schema migration runner (v1 only, but the framework is in place).
**Acceptance Criteria:**
- Store dispatches actions and notifies subscribers
- State persists across popup open/close cycles
- Migration function handles fresh install (no data) and v1 data
- Storage write errors are logged, not swallowed
**Complexity:** M
**Dependencies:** Task 1
**Files:** `js/store.js`, `js/storage-adapter.js`, `js/migration.js`

### TASK 3: Theme System

**Description:** Implement CSS custom properties for light/dark themes. Build the `<theme-toggle>` component that persists preference via the store.
**Acceptance Criteria:**
- Light and dark themes render correctly
- Theme preference persists across sessions
- Urgency colors are legible in both themes
- Toggle is accessible (keyboard-operable, has ARIA label)
**Complexity:** S
**Dependencies:** Task 2
**Files:** `css/variables.css`, `css/base.css`, `js/components/theme-toggle.js`

### TASK 4: Kanban Board + Columns

**Description:** Build the `<kanban-board>` and `<kanban-column>` components. Four columns rendered from store state. "Add task" button per column (or just in "New").
**Acceptance Criteria:**
- Four columns render with correct headers
- Columns subscribe to store and re-render on state change
- Horizontal layout fits within popup width (scrollable if needed)
- "Add task" creates a new task in the correct column
**Complexity:** M
**Dependencies:** Task 2, Task 3
**Files:** `js/components/kanban-board.js`, `js/components/kanban-column.js`, `css/components.css`

### TASK 5: Task Card + Urgency Indicators

**Description:** Build the `<task-card>` component that displays title, due date, estimate, and urgency color indicator. Implement the urgency computation module.
**Acceptance Criteria:**
- Cards display title, due date (if set), and estimate (if set)
- Color indicator follows the spec: green (>24h), yellow (<=24h), red (overdue), neutral (no date or done)
- Due date treats end-of-day as the deadline, not midnight
- Cards are clickable to open detail panel
**Complexity:** M
**Dependencies:** Task 4
**Files:** `js/components/task-card.js`, `js/urgency.js`, `css/components.css`

### TASK 6: Task Detail Panel

**Description:** Build the `<task-detail-panel>` slide-in component with editable title, date picker, estimate input, and description textarea. Auto-saves on change via debounced store dispatch.
**Acceptance Criteria:**
- Panel opens when task card is clicked
- All fields are editable and persist on change
- Date picker uses native `<input type="date">`
- Estimate input accepts numeric value with unit selector (minutes/hours/days), stores as minutes
- Panel can be closed (back button or click outside)
- Delete task action with confirmation
**Complexity:** L
**Dependencies:** Task 5
**Files:** `js/components/task-detail-panel.js`, `js/components/date-picker.js`, `css/components.css`

### TASK 7: Drag and Drop Between Columns

**Description:** Implement drag-and-drop for moving task cards between columns. Use the native HTML Drag and Drop API — no library needed for this scope.
**Acceptance Criteria:**
- Tasks can be dragged between columns
- Drop target highlights during drag
- Column order array updates in store on drop
- Fallback: column move buttons on each card for accessibility
**Complexity:** M
**Dependencies:** Task 5
**Files:** `js/components/kanban-column.js`, `js/components/task-card.js`

### TASK 8: Markdown Editor + Sanitized Rendering

**Description:** Build the `<markdown-editor>` component with edit/preview toggle. Integrate `marked` for parsing and `DOMPurify` for sanitization. Wrap behind the `markdownEnabled` preference.
**Acceptance Criteria:**
- Edit mode shows raw markdown in textarea
- Preview mode renders sanitized HTML
- Links in markdown open in new tabs (`target="_blank"`, `rel="noopener"`)
- XSS payloads in markdown are neutralized (test with common vectors)
- Falls back to plain text display if markdown preference is off
**Complexity:** M
**Dependencies:** Task 6
**Files:** `js/components/markdown-editor.js`, `js/markdown.js`, `lib/marked.min.js`, `lib/dompurify.min.js`

### TASK 9: Share as Markdown

**Description:** Implement the share button that serializes a task to a well-formatted markdown string and copies it to the clipboard. Show a brief confirmation toast.
**Acceptance Criteria:**
- Share button appears in task detail panel
- Click copies formatted markdown to clipboard
- Toast confirms "Copied to clipboard" for 2 seconds
- Markdown format matches the spec from the architecture doc
- Works when description contains markdown formatting
**Complexity:** S
**Dependencies:** Task 6
**Files:** `js/share.js`, `js/components/task-detail-panel.js`

### TASK 10: Polish, Edge Cases, and QA Pass

**Description:** Final sweep. Handle edge cases: empty states, long titles, storage warnings, popup resize behavior. Verify urgency colors in both themes. Test with 50+ tasks. Verify keyboard navigation.
**Acceptance Criteria:**
- Empty column states show helpful message ("No tasks yet")
- Long titles truncate with ellipsis on cards, full title in detail
- All interactions work with keyboard only
- Performance acceptable with 100 tasks loaded
- No console errors or warnings
**Complexity:** M
**Dependencies:** Tasks 1–9
**Files:** Various — all component and CSS files

---

## 10. What I'd Push Back On

1. **Four columns might be three too many for a popup.** Consider whether "New" and "Todo" are meaningfully different in a personal task manager. I'd argue for three: `Todo | In Progress | Done`. If the user wants "New" as a staging area, that's a Phase 2 conversation. But this is the user's call — I'll build four if they want four. I'll just note it.

2. **Drag and drop in a popup is tricky.** Popup dimensions are constrained. Drag targets are small. On laptops with trackpads, the UX can be frustrating. The architecture supports it, but I'd put status-change buttons (small arrow icons) on each card as the primary interaction, and DnD as the progressive enhancement.

3. **"Open in Tab" should be Phase 1, not Phase 2.** The popup constraint is the single biggest UX risk. A simple "Open in Tab" link in the popup header that opens the same app in a full tab eliminates the constraint entirely. It's 10 lines of code in the manifest and one `<a>` tag.

---

## 11. Trade-offs Summary

| Decision | Trade-off | Why |
|---|---|---|
| Vanilla JS over framework | More manual wiring, less abstraction | Extension is small enough that framework overhead isn't justified. Revisit if scope doubles. |
| Web Components | Slight learning curve for some devs | Native encapsulation, no build step, browser-native |
| Task map + column arrays | Slightly more complex storage | O(1) lookups, clean DnD, worth it from day one |
| Estimates in minutes | User must convert mentally (or we convert on display) | Single canonical unit prevents comparison bugs across the entire codebase |
| Markdown from day one | 19KB added bundle size | Avoids retrofitting later; security surface addressed once, correctly |
| No external sync | No multi-device | Privacy first; sync is a Phase 3 concern at earliest |

---

## 12. Final Assessment

This is a well-scoped project. The prior discovery plan captured requirements accurately. My contribution is structural: defining component boundaries, state flow, data model decisions that prevent bugs, and a file organization that a new developer can navigate in under 60 seconds.

The biggest risk isn't technical — it's scope creep. A Kanban board *wants* to become Trello. Hold the line on the phased plan. Ship Phase 1, use it for a week, then decide what Phase 2 actually needs to be based on real usage, not imagination.

Build the simple thing. Build it well. Build it so the developer maintaining it six months from now — possibly you, possibly someone else, definitely on a bad day — can understand every decision by reading the code.

---

Why did the microservices architect refuse to build a monolithic browser extension? Because they kept trying to decompose the popup into 47 independently deployable containers with their own CI pipelines. Sometimes a single `popup.html` is the right architecture... and that's the hardest thing for an architect to admit.

-- Architect Vern (measure twice, deploy once)
