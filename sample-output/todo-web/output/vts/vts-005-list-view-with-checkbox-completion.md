---
id: VTS-005
title: "List View with Checkbox Completion"
complexity: M
status: pending
owner: ""
source: discovery
source_ref: "07-architect-architect-breakdown.md"
dependencies:
  - VTS-003
  - VTS-004
files:
  - "src/components/list/ListView.tsx"
  - "src/components/list/TaskRow.tsx"
---

# List View with Checkbox Completion

Build the ListView and TaskRow components. Rows use TaskSummary for display. Checkboxes toggle completion (sets `completedAt`, moves to Done column). Tasks are sorted: incomplete first (by due date), completed last.

## Criteria

- List renders all tasks with checkboxes
- Checking a task sets `completedAt` and moves to Done column
- Unchecking a task clears `completedAt` and moves to first non-Done column
- Completed tasks are visually dimmed
- Sort order: incomplete (by due date asc) → completed (by completion date desc)
- Clicking a task row opens the detail view
- Empty state shows a helpful message
