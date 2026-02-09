---
id: VTS-013
title: "Task Deletion"
complexity: S
status: pending
owner: ""
source: discovery
source_ref: "07-architect-architect-breakdown.md"
dependencies:
  - VTS-003
  - VTS-004
files:
  - "src/components/shared/TaskSummary.tsx"
  - "src/components/detail/TaskDetail.tsx"
  - "src/store/taskSlice.ts"
---

# Task Deletion

Implement task deletion from the TaskSummary component (accessible from both board cards and list rows) and from the detail panel. Require confirmation before permanent deletion. Support undo within a short window via soft-delete.

Note: This task no longer depends on VTS-008 (Detail Panel). Deletion is primarily accessible from TaskSummary, which is available in both views without the detail panel.

## Criteria

- Delete action is available from TaskSummary (e.g., a menu or icon button on hover/focus)
- Delete action is also available from the detail panel (when open)
- **Confirmation dialog** before deletion (not a browser `confirm()` -- a styled modal/dialog)
- **Undo/soft-delete:** After confirmation, task is soft-deleted with a 5-second undo toast/snackbar
- After the undo window expires, task is permanently removed from the store (and localStorage)
- If the detail panel is open for the deleted task, it closes
- Task disappears from both views immediately on soft-delete
- Undo restores the task to its original position and column
