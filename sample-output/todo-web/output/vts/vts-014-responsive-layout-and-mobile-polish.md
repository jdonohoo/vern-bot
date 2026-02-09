---
id: VTS-014
title: "Responsive Layout and Mobile Polish"
complexity: M
status: pending
owner: ""
source: discovery
source_ref: "07-architect-architect-breakdown.md"
dependencies:
  - VTS-005
  - VTS-006
  - VTS-007
  - VTS-008a
files:
  - "Multiple component files"
  - "src/theme/tokens.css"
---

# Responsive Layout and Mobile Polish

Ensure the app works on mobile screens. Board view shows one column at a time on mobile (with swipe or tab navigation between columns). List view is full-width. Detail panel is full-screen on mobile. All touch targets are appropriately sized.

Note: VTS-011 (Markdown Editor) is no longer a dependency. The markdown editor handles its own mobile behavior (toggle mode) independently.

## Criteria

- **Mobile board:** Board shows one column at a time on mobile (not horizontal scroll)
- **Mobile board:** User can swipe or tap tabs to switch between columns
- List view is full-width and touch-friendly on mobile
- **Mobile detail:** Detail panel is a full-screen overlay/drawer on mobile (not a side panel)
- Markdown editor uses toggle mode (not split) on mobile (handled by VTS-011, verified here)
- **Touch targets:** All interactive elements are at least 44x44px on touch devices
- No horizontal overflow on list view
- Tag manager is accessible on mobile
- View toggle is accessible and functional on mobile
