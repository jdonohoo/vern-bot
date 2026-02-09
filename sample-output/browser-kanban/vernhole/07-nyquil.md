

Plan's solid. Few notes before I fade...

**What matters:**
- Manifest V3, popup + optional tab view. Correct.
- `chrome.storage.local` — yes. Not localStorage. Good.
- CSS custom properties for theming. Obviously.
- Urgency = end-of-day local time. Critical detail most miss.

**What I'd change:**

Storage: task map + column order arrays. Not flat array. DnD will happen eventually, and migrating storage shape later is annoying. Do it right now while it's free.

Framework: vanilla JS. It's a popup. React is mass murder on bundle size for this. Web components if you're feeling fancy but... probably not.

Markdown: use `marked` + `DOMPurify`. Tiny libs. Ship in v1, it's like 30 lines of integration. Don't defer what's cheap.

**Estimate field:** `<input type="number" min="0">` + unit dropdown. Done. Don't overthink.

**Date picker:** native `<input type="date">`. Works everywhere that matters. Don't npm install a calendar library for this...

**Share button:**
```js
navigator.clipboard.writeText(taskToMarkdown(task))
```
...that's it. Show a toast. Move on.

**DnD:** HTML5 drag API is fine for 4 columns. No library needed. `dragstart`, `dragover`, `drop`. Maybe 40 lines total.

**What the plan missed:**
- Keyboard shortcuts. At minimum: `n` for new task. Users who live in extensions expect this.
- Task ordering within columns. Need a `position` field or ordered array.
- Clear done tasks. Bulk action. Will be needed day 2.

**What to cut from v1 if tired:**
- Export/import (nobody does this first week)
- Recurring tasks (scope creep disguised as a feature)
- Search (ctrl+f exists)

**Architecture in 6 lines:**
```
popup.html → board view
popup.js   → state + render
detail.js  → task detail panel
store.js   → chrome.storage wrapper
theme.js   → toggle + persist
style.css  → css vars, columns, cards
```

That's the whole extension. No build step needed. No bundler. Ship it.

**The one thing everyone forgets:** test urgency color across timezone changes. Set your clock to UTC+12, create a task due "today," then switch to UTC-12. If it breaks, your date logic is wrong. Use date-only comparison, ignore time component.

k. good plan. build it.

why did the kanban board go to bed early... it had too many things in the *rest* column... zzz

-- NyQuil Vern *types trailing off into pillow*
