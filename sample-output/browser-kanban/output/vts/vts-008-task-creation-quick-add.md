---
id: VTS-008
title: "Task Creation (Quick Add)"
complexity: M
status: pending
owner: ""
source: oracle
source_ref: "oracle-architect-breakdown.md"
dependencies:
  - VTS-004
  - VTS-006
  - VTS-007
files:
  - "`scripts/components/task-form.js`"
  - "`styles/task-form.css`"
---

# Task Creation (Quick Add)

Build the quick-add task flow. Bifurcated UX per the council's resolution: the "New" column gets an always-visible inline text input (type title, press Enter, done), while other columns get a "+" button that opens an inline form or lightweight modal with title, due date, and estimate fields. The user should go from intent to captured task in under 3 seconds.

## Criteria

- "New" column: always-visible inline text input. Type title, press Enter, task created in "New" status. No modal, no extra clicks.
- Other columns: "+" button opens inline form or lightweight modal with title, due date, and estimate fields
- Global keyboard shortcut: `N` key (when no input is focused) opens quick-add in "New" column
- Title field is required; form won't submit without it
- Due date field uses native `<input type="date">` (reliable, accessible, no library needed)
- Estimate field: numeric input + unit dropdown (minutes/hours/days)
- Submit creates task via `StateManager.addTask()` with auto-generated ID and timestamps
- Default status matches the column where "Add" was initiated
- Form resets after successful creation
- Escape key or click-outside closes form without saving
- Validation: empty title shows inline error message
- Focus management: title field auto-focused on open
