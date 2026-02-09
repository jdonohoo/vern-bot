---
id: VTS-007
title: "Task Card Component with Urgency Colors"
complexity: M
status: pending
owner: ""
source: oracle
source_ref: "oracle-architect-breakdown.md"
dependencies:
  - VTS-005
  - VTS-006
files:
  - "`scripts/components/task-card.js`"
  - "`scripts/utils/urgency.js`"
  - "`styles/task-card.css`"
---

# Task Card Component with Urgency Colors

Build the task card component that renders within Kanban columns. Each card shows the title, due date, estimate, and an urgency indicator. Urgency uses both color AND a secondary indicator (icon or text label) for accessibility. Urgency colors are tuned per-theme for contrast. This is the core visual feedback loop — get it right here.

## Criteria

- Card displays: title, due date (formatted), estimate (if present)
- Left border or badge shows urgency color based on defined rules
- Urgency indicator includes both color AND a secondary indicator (icon or text label): e.g., a dot/icon + "Overdue", "Due soon", "On track"
- Urgency colors are different between light and dark themes (tuned for contrast in each)
- Tasks in "Done" status show neutral urgency regardless of due date (no red "overdue" on completed tasks)
- Tasks with no due date show neutral/gray indicator (not green — green implies "on track," which implies a date exists)
- Urgency rules explicitly defined: Green = due date > 24h from now; Yellow = due date <= 24h and not overdue; Red = due date < now (past end-of-day); Neutral/Gray = no due date OR status is "done"
- Urgency computed fresh on each render (no stale cache)
- Cards are clickable (emits event for detail panel, Task 9)
- Card has delete button (with confirmation)
- Card has status-move buttons or dropdown to change column
- Urgency utility function is pure, isolated, and unit-testable
- Cards truncate long titles with ellipsis
