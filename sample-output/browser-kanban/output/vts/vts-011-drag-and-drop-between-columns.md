---
id: VTS-011
title: "Drag and Drop Between Columns"
complexity: L
status: pending
owner: ""
source: oracle
source_ref: "oracle-architect-breakdown.md"
dependencies:
  - VTS-006
  - VTS-007
files:
  - "`scripts/components/kanban-board.js` (extend)"
  - "`scripts/components/task-card.js` (extend)"
  - "`styles/drag-drop.css`"
---

# Drag and Drop Between Columns

Add drag-and-drop support for moving task cards between columns. Use the native HTML5 Drag and Drop API — no library. Cards should have visual feedback during drag. On drop, update the task's status and position via the state manager. This is the UX upgrade that makes the board feel like a real Kanban tool.

## Criteria

- Cards are draggable between columns
- Visual feedback: dragged card has reduced opacity, target column highlights on dragover
- Drop updates task status to target column's status
- Drop position within column is respected (top, between cards, bottom)
- Column order array in storage updated on drop
- Keyboard accessibility: cards can be moved with keyboard shortcuts (arrow keys + modifier)
- No-op if dropped in same position
- Works reliably in Chrome (primary target)
