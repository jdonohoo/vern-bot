---
id: VTS-006
title: "Kanban Board Layout (Four Columns)"
complexity: M
status: pending
owner: ""
source: oracle
source_ref: "oracle-architect-breakdown.md"
dependencies:
  - VTS-004
  - VTS-005
files:
  - "`scripts/components/kanban-board.js`"
  - "`styles/kanban-board.css`"
---

# Kanban Board Layout (Four Columns)

Build the four-column Kanban board UI: New, Todo, In Progress, Done. Each column renders task cards from state. The board must work within popup constraints. The "New" column has special behavior — an always-visible inline text input for zero-friction capture — while other columns get a "+" button. This architectural distinction is defined here so it doesn't get discovered during task creation implementation.

## Criteria

- Four columns rendered: New, Todo, In Progress, Done
- Each column shows its name and task count badge
- Columns scroll vertically independently when content overflows
- Board scrolls horizontally if popup width is insufficient
- Board subscribes to `tasks-changed` events and re-renders
- Empty state message per column ("No tasks")
- Column layout uses CSS Grid or Flexbox (no absolute positioning hacks)
- Responsive within popup size constraints (min 400px, max 800px width)
- "New" column renders an always-visible inline text input at the top for zero-friction capture
- Other columns (Todo, In Progress, Done) show a "+" button in the column header
- Popup minimum usable width tested at 400px with all four columns visible
- Column widths distribute evenly with CSS Grid, minimum 90px per column in popup
- Board detects popup vs. tab context and adjusts layout accordingly (ties to Task 15)
