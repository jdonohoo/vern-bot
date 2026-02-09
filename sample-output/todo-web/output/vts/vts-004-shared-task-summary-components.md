---
id: VTS-004
title: "Shared Task Summary Components"
complexity: M
status: pending
owner: ""
source: discovery
source_ref: "07-architect-architect-breakdown.md"
dependencies:
  - VTS-002
  - VTS-003
files:
  - "src/components/shared/TaskSummary.tsx"
  - "src/components/shared/TagChip.tsx"
  - "src/components/shared/DateBadge.tsx"
---

# Shared Task Summary Components

Build TaskSummary, TagChip, and DateBadge components. TaskSummary displays title, tag chips, and date badges. TagChip renders with computed opacity for dark-mode readability. DateBadge shows neutral/warning/danger states based on due date proximity.

Note: `shortDescription` was removed from the data model (VTS-002). TaskSummary shows the title, tags, and dates only. A body preview can be derived from `bodyMarkdown` if needed in a future iteration.

## Criteria

- TaskSummary renders title, tags, and dates
- TagChip uses opacity-based rendering (15% bg, 60% border, 90% text)
- DateBadge shows correct state: neutral (>3 days), warning (≤3 days), danger (overdue)
- Completed tasks show "Done [date]" instead of overdue
- Components are reusable by both board and list views
- Tasks with no tags or dates render cleanly without empty space
