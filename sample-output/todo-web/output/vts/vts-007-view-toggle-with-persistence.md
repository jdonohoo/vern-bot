---
id: VTS-007
title: "View Toggle with Persistence"
complexity: S
status: pending
owner: ""
source: discovery
source_ref: "07-architect-architect-breakdown.md"
dependencies:
  - VTS-003
  - VTS-005
files:
  - "src/components/shared/ViewToggle.tsx"
  - "src/App.tsx"
---

# View Toggle with Persistence

Build ViewToggle component that switches between Board and List views. Persist preference to localStorage via the store. Switching views does not discard the currently open detail panel. Show a disabled/coming-soon state for views that are not yet implemented (e.g., board view if VTS-006 is not yet complete).

Note: VTS-006 (Board View) is no longer a dependency. This task can be built with only the list view available, showing a disabled toggle option for the board view until it ships.

## Criteria

- Toggle switches between board and list views
- Active view preference persists across page reloads
- Switching views does not close the detail panel if open
- Toggle is keyboard-accessible
- Visual indicator shows which view is active
- **Disabled state:** If a view's component is not yet available, its toggle option shows as disabled with a "Coming soon" tooltip or visual indicator
- **Graceful degradation:** App defaults to list view if the persisted preference references an unavailable view
