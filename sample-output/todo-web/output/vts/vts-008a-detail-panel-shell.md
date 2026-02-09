---
id: VTS-008a
title: "Detail Panel Shell"
complexity: S
status: pending
owner: ""
source: discovery
source_ref: "07-architect-architect-breakdown.md"
dependencies:
  - VTS-003
  - VTS-004
files:
  - "src/components/detail/TaskDetailShell.tsx"
  - "src/components/detail/TaskDetail.tsx"
---

# Detail Panel Shell

Build the minimal detail panel drawer that opens when a task is selected. This is the shell only -- it handles layout, open/close behavior, URL state, and basic field editing (title, body textarea, status dropdown, timestamps). Tags, dates, and markdown preview are integrated later in VTS-008b.

This split resolves the circular dependency where the detail panel previously depended on tags and dates, but tags and dates needed a detail panel to test in context. The shell ships first; integrations plug in after.

## Criteria

- Detail panel opens as a drawer/slide-over when a task is selected
- Title is editable inline (click to edit, blur to save)
- Body field is a plain textarea (markdown preview comes in VTS-008b via VTS-011)
- Status dropdown changes the task's column (uses `isCompletionColumn` to set/clear `completedAt`)
- Created and updated timestamps are displayed (read-only)
- Changes are saved to the store immediately (or on blur)
- **URL state:** Selected task ID is stored in URL search params (`?task=<id>`)
- **URL state:** Opening a URL with `?task=<id>` opens the detail panel for that task
- **URL state:** Browser back button closes the panel
- Closing the panel deselects the active task and clears the URL param
- Focus is trapped within the panel when open
- Escape key closes the panel
- Clicking outside the panel closes it
