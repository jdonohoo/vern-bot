---
id: VTS-006
title: "Board View with Drag-and-Drop"
complexity: L
status: pending
owner: ""
source: discovery
source_ref: "07-architect-architect-breakdown.md"
dependencies:
  - VTS-003
  - VTS-004
files:
  - "src/components/board/BoardView.tsx"
  - "src/components/board/Column.tsx"
  - "src/components/board/TaskCard.tsx"
  - "src/utils/ordering.ts"
---

# Board View with Drag-and-Drop

Build the BoardView and Column components using @dnd-kit. Cards use TaskSummary. Implement drag-and-drop within and across columns using float-based ordering. Include reindexing when order gaps get too small. Empty columns accept drops. Include performance optimizations to prevent render thrashing during drag operations. Support keyboard-accessible DnD and mobile touch with appropriate fallbacks.

## Criteria

- Board renders columns in order with task cards
- Tasks are draggable within a column (reorder)
- Tasks are draggable across columns (status change + reorder)
- Float-based ordering works correctly (midpoint insertion)
- Reindexing triggers when gaps < 0.001
- Empty columns show a drop zone and accept drops
- Drop placeholder is visible during drag
- "Add Task" affordance at the bottom of each column
- **Performance:** TaskCard components are memoized (React.memo) to prevent cross-column re-renders during drag
- **Performance:** Zustand selectors are granular -- dragging a card in Column A does not re-render cards in Column B
- **Performance:** Drag preview does not cause layout thrashing
- **Keyboard DnD:** @dnd-kit keyboard sensor is configured with accessible announcements
- **Keyboard DnD:** Users can pick up, move, and drop tasks using only the keyboard
- **Mobile touch:** Touch sensor configured with appropriate activation constraints (delay or distance)
- **Mobile touch:** Drag handles are minimum 44x44px
- **Mobile fallback:** A status dropdown is visible on each card on mobile/touch devices as an alternative to drag-and-drop
