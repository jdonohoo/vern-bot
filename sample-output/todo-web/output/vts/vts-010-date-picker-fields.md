---
id: VTS-010
title: "Date Picker Fields"
complexity: M
status: pending
owner: ""
source: discovery
source_ref: "07-architect-architect-breakdown.md"
dependencies:
  - VTS-003
files:
  - "src/components/detail/DatePickerField.tsx"
  - "src/utils/dates.ts"
---

# Date Picker Fields

Integrate react-day-picker for due date selection. Style for both light and dark themes. Dates can be cleared. Show relative date display alongside the absolute date. This is a standalone component that does not depend on the detail panel -- it receives a task ID or date value as props and can be tested independently.

Note: `startDate` has been removed from v1. Only `dueDate` is supported.

## Criteria

- Calendar picker for due date
- Picker is styled to match theme tokens (light and dark)
- Selecting a date updates the task in the store
- **Clear button:** Dates can be cleared (set to null) via an explicit clear/remove button
- **Relative dates:** Display shows relative text alongside the date ("Tomorrow", "In 3 days", "Overdue by 2 days")
- Picker is keyboard-accessible
- Mobile-friendly sizing
- DatePickerField is a standalone component -- it does not depend on the detail panel
