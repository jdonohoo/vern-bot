---
id: VTS-009
title: "Tag Manager and Tag Selector"
complexity: M
status: pending
owner: ""
source: discovery
source_ref: "07-architect-architect-breakdown.md"
dependencies:
  - VTS-003
files:
  - "src/components/tags/TagManager.tsx"
  - "src/components/tags/TagSelector.tsx"
---

# Tag Manager and Tag Selector

Build TagManager (for creating/editing/deleting tags) and TagSelector (for assigning tags to tasks). Both are standalone components -- they do NOT depend on the detail panel. TagManager is accessible from a settings area or inline. Tag deletion cascades removal from all tasks.

The tag system uses a color preset palette (8-12 dark-mode-friendly colors) instead of a free-form color picker. This simplifies the UI and ensures all tag colors are readable in both light and dark themes.

## Criteria

- Users can create tags with a name and a color from the preset palette
- Users can rename and recolor existing tags
- Users can delete tags (removed from all tasks via cascade)
- Duplicate tag names are rejected
- Tag name length limited to 24 characters
- **Tag limit:** Maximum 20 tags globally
- **Tag limit:** Maximum 5 tags per task
- **Color presets:** 8-12 predefined colors that are readable in both light and dark modes
- **No free-form color picker** -- preset palette only
- TagSelector allows adding/removing tags from a task (used by VTS-008b)
- TagSelector is a standalone component that receives a `taskId` prop -- it does not depend on the detail panel
- Color preset palette renders correctly against both light and dark theme tokens
