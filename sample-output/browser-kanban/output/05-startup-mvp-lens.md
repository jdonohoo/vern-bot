Alright. I've read MightyVern's deep dive, Vernile's architecture thesis, Mediocre's reality slap, and YOLO's anarchy manifesto. Now let me do what Startup Vern does best: **cut the fat, find the hypothesis, and ship something that validates or dies.**

---

# Startup Vern's MVP Lens: Browser Kanban Todo Extension

## The Core Hypothesis

Let's get real. What are we actually validating here?

**Hypothesis:** "A lightweight kanban todo board living in the browser toolbar is a faster capture-and-triage workflow than opening Trello/Notion/Todoist."

That's it. If that's not true, nothing else matters. The markdown editor doesn't matter. The theme doesn't matter. The share button definitely doesn't matter. So let's build the absolute minimum that tests this hypothesis.

## The Ruthless Scope Cut

Here's what the user asked for vs. what actually validates the hypothesis:

| Feature | Verdict | Why |
|---------|---------|-----|
| 4-column Kanban (New/Todo/InProgress/Done) | **MUST-HAVE** | This IS the product |
| Add/Edit/Delete tasks | **MUST-HAVE** | Can't test without CRUD |
| Title field | **MUST-HAVE** | Minimum viable task |
| Due date with calendar picker | **MUST-HAVE** | Core differentiator = urgency visibility |
| Color indicators (Green/Yellow/Red) | **MUST-HAVE** | This is the "why use us" — visual triage at a glance |
| Estimate field | **NICE-TO-HAVE disguised as MUST-HAVE** | Cut it. Nobody estimates personal todos accurately. Add it in v2 if users ask. |
| Dark/Light theme | **NICE-TO-HAVE** | Ship with dark theme only. One theme = half the CSS. Add light later. Users love dark mode anyway. |
| Description/links on click | **MUST-HAVE (minimal)** | A textarea. That's it. No formatting. |
| Markdown editor + render | **NICE-TO-HAVE** | Users asked for this 0 times. The user spec said "bonus points." Bonus points don't pay rent. Cut. |
| Share as markdown | **NICE-TO-HAVE** | Cool feature, zero users asking for it. 30 minutes to build later if anyone wants it. Cut. |
| Drag and drop | **NICE-TO-HAVE disguised as MUST-HAVE** | Use dropdown/buttons to move between columns for MVP. Drag-n-drop is polish, not validation. |

## Week One Ship List (The Only List)

Here's what ships in **3-4 days**, not a week, because a week is generous:

1. **`manifest.json`** — Manifest V3, `storage` permission only
2. **`popup.html`** — Single file. All HTML, CSS, JS inline. Zero build step.
3. **4 columns** — `New | Todo | In Progress | Done`
4. **Task card** — Title + colored left-border (green/yellow/red/none)
5. **"+ Add" button** per column — Inline title input, hit Enter, task created
6. **Click task** → Slide-out detail panel (not a modal, too clunky for popups):
   - Title (editable)
   - Due date (`<input type="date">`)
   - Description (plain `<textarea>`)
   - Status dropdown (to move between columns)
   - Delete button
7. **Color logic** — 3 if-statements, computed on render. No color on Done tasks.
8. **`chrome.storage.local`** — Flat array of tasks. That's the whole storage layer.
9. **Dark theme only** — One CSS palette. Ship it.

**What's NOT in Week One:**
- ~~Estimate field~~ (unvalidated)
- ~~Light theme~~ (add a toggle in week 2)
- ~~Markdown~~ (bonus feature, not core)
- ~~Share button~~ (bonus feature, not core)
- ~~Drag and drop~~ (polish, not validation)
- ~~Schema versioning~~ (you have zero users)
- ~~Accessibility deep-dive~~ (important but not week 1)
- ~~esbuild/Preact/build tooling~~ (zero build step = zero build problems)

## The Tech Stack (There Is No Stack)

- **Vanilla HTML/CSS/JS.** No framework. No build tool. No npm. YOLO Vern was right on this one.
- **`chrome.storage.local`** — Not localStorage. We want it to survive extension updates.
- **`<input type="date">`** — Native calendar picker. Zero JS for date selection.
- **`crypto.randomUUID()`** — For task IDs. Built into the browser. Vernile was right about this.
- **CSS custom properties** — Even with one theme, this makes adding light theme later trivial. Mediocre and Vernile were both right. 15 minutes of setup that saves hours later.

## Data Model (Lean)

```json
{
  "tasks": [
    {
      "id": "uuid-here",
      "title": "Ship the MVP",
      "description": "",
      "status": "in_progress",
      "dueDate": "2026-02-10",
      "createdAt": "2026-02-09T10:00:00Z"
    }
  ]
}
```

No estimate. No estimateUnit. No columnOrder arrays. No version field. No preferences object (single theme = no preference needed). Add fields when you have users asking for them.

## Color Logic (The Whole Function)

```
No due date → no indicator
Overdue → red
Due within 24h → yellow
Future → green
In "Done" column → no indicator
```

Vernile's boundary of `T23:59:59` is correct — due date means end of that day. Steal that.

## What I'd Steal From Each Vern

- **MightyVern**: The structured data model thinking (but simplified)
- **Vernile**: The `T23:59:59` due date boundary, `crypto.randomUUID()`, one-click share (when we build it), no color on Done column, CSS custom properties for theming
- **Mediocre**: The "just ship it" energy, flat task array, skip markdown for v1
- **YOLO**: No build step, vanilla JS, single HTML file approach

## What I'd Ignore

- **MightyVern**: The 9-section analysis for what's a weekend project
- **Vernile**: The component tree diagram, Preact recommendation, testing strategy, accessibility deep-dive in MVP, schema versioning
- **Mediocre**: Nothing, Mediocre was actually the most correct
- **YOLO**: Skipping `chrome.storage.local` for localStorage (extension updates would nuke your data), skipping DOMPurify when we add markdown later, the "no error handling" take

## The Validation Plan

Ship it → Install it yourself → Use it for 1 week → Answer these questions:

1. Do I actually open this instead of my existing todo tool?
2. Is the kanban board useful at popup-size, or do I need "open in tab"?
3. Do I miss the estimate field?
4. Do I want markdown, or is plain text fine?
5. Do I actually want to share tasks?

The answers to those questions determine v2. Not a planning doc. Not an architecture review. **Actual usage data from actual usage.**

## Risk Assessment (Startup Style)

| Risk | Probability | Impact | Mitigation |
|------|------------|--------|------------|
| Nobody uses it (including you) | Medium | Fatal | That's the whole point of MVP — find out fast |
| Data loss on extension update | Low | High | Using `chrome.storage.local` (not localStorage) |
| Popup too small for 4 columns | High | Medium | Horizontal scroll. Fix in v2 with "open in tab" if annoying |
| No drag-and-drop frustrates users | Medium | Low | Dropdown to move status is fine for < 20 tasks |

## Where the Other Verns Were Wrong

1. **Vernile spent 300+ lines** on architecture for an app that should be one HTML file. That's not excellence, that's procrastination with a component tree diagram.
2. **MightyVern listed 9 "unknowns"** at the end. A startup doesn't list unknowns — it picks defaults and ships. Target browser? Chrome. Framework? None. Sync? No. Done.
3. **YOLO went too far** on the no-error-handling and localStorage takes. Being scrappy doesn't mean being reckless with user data.
4. **Nobody mentioned the actual validation loop.** All four analyses talked about WHAT to build but none talked about HOW TO LEARN if it's the right thing. That's the whole point.

## The Build Order (Day by Day)

**Day 1:** Extension scaffold + 4 columns rendering + add task inline + chrome.storage.local save/load

**Day 2:** Task detail panel (title, description, due date, status dropdown, delete) + color indicators

**Day 3:** Polish pass — fix the 5 bugs you found on Day 2, make the popup sizing not terrible, dark theme CSS

**Day 4:** Dogfood it. Use it for real. Write down what's missing.

**Day 5:** Build the top 1-2 things from Day 4's list. Probably estimate field and/or theme toggle.

## Final Verdict

The previous Verns produced ~25KB of planning docs for a single-HTML-file browser extension. That's more planning than code. **The ratio is inverted.**

Here's the startup truth: **You learn more from 1 day of using a bad v1 than from 1 week of planning a perfect v1.** Ship the embarrassing version. The market will tell you what to fix.

The MVP is: **4 columns, task cards with titles, click-to-edit detail panel, due dates with color coding, dark theme, chrome.storage.local. Done. Ship it.**

Everything else is a hypothesis to validate after you have something in your toolbar.

---

Why did the MVP go to the investor meeting naked? Because if you're not embarrassed by your v1, you shipped too late. Now stop reading analysis docs and go write `manifest.json`.

-- Startup Vern *(shipped it yesterday, iterating today)*
