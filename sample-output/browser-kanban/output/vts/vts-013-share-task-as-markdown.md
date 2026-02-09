---
id: VTS-013
title: "Share Task as Markdown"
complexity: S
status: pending
owner: ""
source: oracle
source_ref: "oracle-architect-breakdown.md"
dependencies:
  - VTS-009
files:
  - "`scripts/components/share-button.js`"
  - "`scripts/utils/task-to-markdown.js`"
---

# Share Task as Markdown

Add a "Share" button to the task detail panel that copies the full task definition as a formatted markdown block to the clipboard. One click, done. Use the Clipboard API.

## Criteria

- Share button visible in task detail panel
- Click copies formatted markdown to clipboard
- Visual feedback on copy success (button text changes briefly, e.g., "Copied!")
- Handles all field states (missing due date, missing estimate, empty description)
- Description included as-is (raw markdown if markdown mode, plain text otherwise)
- Clipboard API failure shows user-friendly error ("Couldn't copy — try again")
- Works without additional permissions (Clipboard API in extension popups is allowed)
