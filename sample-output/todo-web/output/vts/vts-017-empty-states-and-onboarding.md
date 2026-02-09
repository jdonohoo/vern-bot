---
id: VTS-017
title: "Empty States and Onboarding"
complexity: S
status: pending
owner: ""
source: discovery
source_ref: "07-architect-architect-breakdown.md"
dependencies:
  - VTS-005
  - VTS-006
files:
  - "src/components/shared/EmptyState.tsx"
  - "src/components/board/BoardView.tsx"
  - "src/components/list/ListView.tsx"
  - "src/components/tags/TagManager.tsx"
---

# Empty States and Onboarding

Design and implement empty states for every view in the app. Empty states are the user's first experience -- they are not polish, they are onboarding. Each empty state guides the user toward the next action.

## Empty States to Implement

1. **Empty board** (no tasks at all): Welcome message + prominent "Add your first task" call to action
2. **Empty column** (board view, specific column has no tasks): Subtle placeholder text, accepts drag drops
3. **Empty list view** (no tasks): Welcome message + quick-add input is visually highlighted
4. **Empty tag list** (tag manager, no tags created): "Create your first tag" prompt with example
5. **Empty detail body** (task with no markdown content): Placeholder text ("Add details...") that invites editing

## Criteria

- Every view has a meaningful empty state (no blank white or black voids)
- Empty states include a clear call to action guiding the user to the next step
- Empty board state specifically highlights the quick-add input or provides an inline creation flow
- Empty column state is visually distinct from a column with tasks but not distracting
- Empty states use theme tokens and look correct in both light and dark modes
- Empty states are not just text -- they include visual affordances (icons, highlighted inputs, etc.)
