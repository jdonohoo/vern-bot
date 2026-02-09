---
id: VTS-012
title: "Task Creation Flow"
complexity: S
status: pending
owner: ""
source: discovery
source_ref: "07-architect-architect-breakdown.md"
dependencies:
  - VTS-003
  - VTS-005
files:
  - "src/components/shared/QuickAdd.tsx"
  - "src/components/list/ListView.tsx"
  - "src/store/taskSlice.ts"
---

# Task Creation Flow

Implement quick-add task creation. A visible input field allows the user to type a title and press Enter to create a task immediately. The task is created in the default (first) column. Body, tags, dates, and other metadata are added later via the detail panel -- quick-add is title-only for speed.

Note: This task no longer depends on VTS-006 (Board View) or VTS-008 (Detail Panel). It works with just the store and the list view. Board-specific "Add Task" at column bottom is part of VTS-006.

## Criteria

- **Quick-add input:** A visible, always-present text input for task creation
- **Quick-add flow:** Type title + press Enter = task created immediately
- New task is created in the first (default) column with an empty body
- New task is assigned the next order value in its column
- Task appears immediately in the list view
- **Focus behavior:** After creating a task, focus returns to the quick-add input (ready for next task)
- Input is cleared after successful creation
- Empty title (whitespace only) does not create a task
- Quick-add input has clear placeholder text ("Add a task...")
