# Retro Vern's Take: The Browser Extension Kanban Todo App

*cracks knuckles, adjusts bifocals, sips from a mug that says "I survived CORBA"*

Alright. Let me read through this... 

A todo list. In a browser extension. With four columns and some colors.

Kid, sit down. Let me tell you a story.

---

## What Is This Problem, Really?

Strip away the buzzwords and you've got: **a personal task list with status tracking, due dates, and notes.** That's it. That's the whole thing. We had this in 1998. It was called a text file. Then we got fancy and put it in a spreadsheet. Then we got *real* fancy and put it in a database with a web frontend. Now we're putting it back in the browser but calling it a "Kanban board" because somewhere along the way we decided columns are a methodology.

The core ask is CRUD with four status buckets and a date field. I've written this on napkins that had more architectural complexity than what's needed here.

---

## What The Discovery Plan Gets Right

Credit where it's due — the prior plan is **actually pretty restrained**, and I respect that. Specifically:

- **Offline-first, local storage only.** Good. No reason whatsoever for a personal todo list to phone home. We learned this lesson with every "cloud-first" note app that shut down and took your data with it. Remember Google Notebook? Google Keep will join it someday. Local data is immortal data.
- **Minimal permissions (`storage` only).** Correct. I've seen browser extensions that ask for more permissions than a TSA agent. A todo list needs to read and write its own data. Period.
- **Phased delivery.** Smart. Ship the boring stuff first, add the shiny stuff later. This is how professionals work.
- **"Avoid over-abstracting; a small extension deserves a small architecture."** Whoever wrote that line has been burned before. I salute them.

---

## What Concerns Me (The Grizzled Take)

### 1. This Plan Document Is Bigger Than The App Will Be

You've got a 13-section planning document with edge cases, non-functional requirements, and a phased implementation plan... for a todo list. The spec is going to outweigh the code. I've seen this pattern before — it's called "enterprise planning for a shell script problem."

The entire Phase 1 is maybe 500-800 lines of actual code. You could scaffold this in an afternoon. The plan should be a README with bullet points, not a design document that needs its own table of contents.

### 2. Framework Decision: The Answer Is No Framework

The "Open Decisions" section asks about "preferred framework (vanilla vs lightweight library)." Let me answer that definitively: **vanilla JavaScript, HTML, and CSS.**

Here's why:
- This is a popup. It renders a list of items in columns. There is no complex state tree. There is no routing. There are no async data fetches from 14 microservices.
- React, Vue, Svelte — they all add build tooling, bundle size, and dependency maintenance for zero benefit at this scale.
- You know what renders four columns of cards fast? A `<div>` with `display: grid`. CSS Grid has been stable since 2017. It'll still work when whatever framework you pick is in the JavaScript graveyard next to Backbone and Knockout.
- `chrome.storage.local` is your database. `JSON.stringify` is your ORM. `document.createElement` is your component framework. Done.

Every dependency you add is a dependency you maintain. A todo list with zero dependencies will still work in 10 years. A todo list built on today's hot framework might not survive the next major version bump.

### 3. The Data Model Is Overcomplicated

```
id, title, description, status, dueDate, estimate, estimateUnit, createdAt, updatedAt
```

You know what this is? It's a row in a table. You know what handles rows in tables well? A flat JSON array. The plan mentions "task map + separate column order array" for drag and drop support. 

Friends, drag and drop is reordering items in an array. `Array.prototype.splice()` has been doing this since 1997. You don't need a separate data structure. You need `tasks.sort()` and a `position` field if ordering matters. We solved this with `ORDER BY sort_order` in SQL decades ago.

### 4. Markdown Support: Do You Actually Need This?

The plan suggests a markdown editor with preview and sanitization. For a todo item description in a browser extension popup.

Let me paint you a picture: your popup is roughly 400x600 pixels. You're going to put a markdown editor with a preview pane in there? That's not a feature, that's a compromise. 

Here's what I'd actually do: support `<textarea>` for descriptions. If you want links, just make URLs clickable with a simple regex (yeah, I know, "you can't parse URLs with regex" — you can parse *enough* of them). If someone truly needs markdown rendering, use a single lightweight library like `marked` (8KB minified) and sanitize with `DOMPurify`. But I'd push this to Phase 3, not Phase 2.

The "share as markdown" button, though? That's actually clever. It's just string concatenation — no library needed. Template literals have been in JavaScript since ES6 (2015). Build a string, copy to clipboard. Ten lines of code, tops.

### 5. The Urgency Color Logic Is Correct But Let's Keep It Simple

The plan correctly identifies end-of-day interpretation for due dates. Good. Here's the entire implementation:

```javascript
function getUrgency(task) {
  if (!task.dueDate || task.status === 'done') return 'neutral';
  const endOfDueDay = new Date(task.dueDate + 'T23:59:59');
  const hoursLeft = (endOfDueDay - new Date()) / 3600000;
  if (hoursLeft < 0) return 'red';
  if (hoursLeft < 24) return 'yellow';
  return 'green';
}
```

That's it. Seven lines. No urgency engine. No color computation service. No abstract strategy pattern. A function. That returns a string. You put the string in a CSS class. CSS does the rest. We've been doing conditional styling since `bgcolor` was an HTML attribute.

### 6. Theme Toggle: CSS Variables, Nothing More

The plan gets this right. CSS custom properties (variables) with a `data-theme` attribute on the root element. Two sets of color definitions. A toggle that flips the attribute and saves to storage. This is a solved problem since CSS custom properties landed in browsers around 2016-2017.

Do not install a "theme management library." I've seen it happen. Don't let it happen to you.

---

## How Was This Solved Before?

- **1995:** A text file called `TODO.txt`. Still works.
- **2004:** A spreadsheet with conditional formatting for due dates. Color-coded rows. Shared via email attachment.
- **2006:** Remember The Milk. Todoist. Web apps with actual servers.
- **2010:** Trello. Which is literally this — a Kanban board for tasks. It's still around.
- **2015:** `todo.txt` format became a thing. Plain text, parseable, portable. Outlived most apps.
- **2020:** Every note app added a Kanban view because Notion made it trendy.

The pattern is the same every time: **a list of things, with states, and dates.** The only thing that changes is the UI and where the data lives.

---

## What I'd Actually Build

Here's my Phase 1, no nonsense:

**File structure:**
```
manifest.json
popup.html
popup.css
popup.js
```

That's four files. No `src/`. No `components/`. No `utils/`. No `build/`. No `package.json`. No `node_modules/`. Four files that a human can read top to bottom and understand in 20 minutes.

**Tech stack:**
- Vanilla JS (no framework)
- CSS Grid for the board layout
- CSS custom properties for theming
- `chrome.storage.local` for persistence
- Native `<input type="date">` for the calendar picker (it's been good enough since Chrome 20, circa 2012)
- `contenteditable` or `<textarea>` for descriptions
- `navigator.clipboard.writeText()` for the share button

**What I'd defer to Phase 2:**
- Drag and drop (buttons to move between columns work fine for v1)
- Markdown rendering
- "Open in tab" view

**What I'd skip entirely unless someone actually asks for it:**
- Recurring tasks (scope creep incarnate)
- Export/import (copy-paste the markdown share output; that *is* your export)
- Search/filter (how many todos do you have? If more than 50, your problem isn't search — it's prioritization)

---

## What's Genuinely New Here?

Honestly? Not much, and that's fine. The browser extension context is the interesting part — you're building a tiny app that lives in a constrained popup window, uses a browser-specific storage API, and needs to be instantly available. That's a legitimate UX advantage over opening a web app or a native app.

The Manifest V3 migration is the only "new" thing worth thinking about, and that's more of an annoying infrastructure constraint than an innovation.

---

## Final Verdict

The discovery plan is solid but overwrought for the actual scope. This is a weekend project wearing a suit and tie. Strip it down, build it in vanilla JS, ship Phase 1 in a day, and iterate from there. The best todo app is the one that exists and works, not the one with a 13-section design document.

The hardest part of this project won't be the code. It'll be resisting the urge to add features. Every todo app in history died from feature bloat, not from missing functionality. Keep it boring. Boring ships. Boring lasts.

You know what the most popular todo system in the world is? A sticky note. Zero dependencies, zero build time, works offline, supports any theme (buy different colored ones), and has never had a breaking API change. Your browser extension should aspire to be a digital sticky note, not a digital project management suite.

---

Why did the developer mass-migrate their todo app to a microservices architecture? They didn't. They wrote a shell script in 2003 that emails them a reminder, and it's still running on a box in the closet. `crontab -e` never needed a sprint planning meeting.

-- Retro Vern *(been mass-assigning since `at(1)` was cutting edge)*
