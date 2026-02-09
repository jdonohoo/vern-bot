---
id: VTS-009
title: "Task Detail Panel (View + Edit)"
complexity: L
status: pending
owner: ""
source: oracle
source_ref: "oracle-architect-breakdown.md"
dependencies:
  - VTS-004
  - VTS-007
  - VTS-008
files:
  - "`scripts/components/task-detail.js`"
  - "`styles/task-detail.css`"
---

# Task Detail Panel (View + Edit)

Build the task detail panel that opens when a task card is clicked. This is where users edit the title, description, due date, estimate, and status. The panel slides in from the right or opens as an overlay — it must not navigate away from the board. Description field supports plain text and links. This panel is the editing workhorse; keep it clean and obvious.

## Criteria

- Panel opens on task card click, showing all task fields
- All fields are editable inline: title, description, due date, estimate, estimateUnit, status
- Status change via dropdown moves task between columns
- Description field is a textarea that auto-grows
- Links in description are clickable when not in edit mode
- Changes auto-save on field blur (no explicit save button needed)
- Close button returns to board view
- Escape key closes panel
- Shows `createdAt` and `updatedAt` as read-only metadata
- Delete button with "Are you sure?" confirmation
- Panel subscribes to state changes (reflects external updates)
