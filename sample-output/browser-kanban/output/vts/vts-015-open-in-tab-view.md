---
id: VTS-015
title: "Open in Tab View"
complexity: M
status: pending
owner: ""
source: oracle
source_ref: "oracle-architect-breakdown.md"
dependencies:
  - VTS-010
  - VTS-004
files:
  - "`tab.html`"
  - "`scripts/main.js` (extend with tab detection)"
  - "`styles/tab-layout.css`"
---

# Open in Tab View

Add an "Open in Tab" button that opens the Kanban board in a full browser tab. The tab view uses the same codebase but with a wider layout. State is shared via `chrome.storage.local`. Cross-view sync is handled by the `StateManager`'s `chrome.storage.onChanged` listener (established in Task 4). Includes duplicate tab detection — the UX cost of duplicates is higher than the implementation cost of a 3-line check.

## Criteria

- "Open in Tab" button/icon in the popup toolbar
- Opens `tab.html` in a new browser tab via `chrome.tabs.create()`
- Tab view uses the same JS and CSS but with a wider layout (CSS media query or layout class)
- State is shared: changes in popup reflect in tab and vice versa
- `chrome.storage.onChanged` listener keeps both views in sync (via StateManager)
- Tab view has more generous column widths and card sizes
- Tab detection: `chrome.tabs.query` to check if tab view is already open; if so, focus it instead of opening a duplicate
