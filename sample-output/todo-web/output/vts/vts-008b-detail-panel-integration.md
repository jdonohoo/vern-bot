---
id: VTS-008b
title: "Detail Panel Integration"
complexity: S
status: pending
owner: ""
source: discovery
source_ref: "07-architect-architect-breakdown.md"
dependencies:
  - VTS-008a
  - VTS-009
  - VTS-010
  - VTS-011
files:
  - "src/components/detail/TaskDetail.tsx"
---

# Detail Panel Integration

Integrate the tag selector (VTS-009), date picker (VTS-010), and markdown editor/preview (VTS-011) into the detail panel shell (VTS-008a). Wire all field updates to the store. This is the final assembly step that turns the shell into the complete detail editing surface.

## Criteria

- Tag selector from VTS-009 is embedded in the detail panel and functional
- Date picker from VTS-010 is embedded in the detail panel and functional
- Markdown editor/preview from VTS-011 replaces the plain textarea for the body field
- All field changes (tags, dates, body) are persisted to the store
- Panel layout accommodates all integrated components without overflow or scroll issues
- Tab order through the panel fields is logical (title -> status -> tags -> dates -> body)
- All integrated components respect the theme tokens from VTS-016
