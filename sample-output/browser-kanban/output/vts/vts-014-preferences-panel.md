---
id: VTS-014
title: "Preferences Panel"
complexity: S
status: pending
owner: ""
source: oracle
source_ref: "oracle-architect-breakdown.md"
dependencies:
  - VTS-005
  - VTS-012
files:
  - "`scripts/components/preferences-panel.js`"
  - "`styles/preferences.css`"
---

# Preferences Panel

Build a preferences/settings panel accessible from the toolbar. Houses the theme toggle, markdown toggle, and default estimate unit selector. Keep it simple — a slide-out panel or modal. This centralizes user configuration and keeps the main board clean.

## Criteria

- Settings icon/button in the toolbar
- Panel includes: theme toggle, markdown on/off, default estimate unit dropdown
- All changes persist immediately via state manager
- Panel closable via close button, escape key, or click-outside
- Current preference values shown on open (not defaults)
- Changes to markdown setting don't destroy existing description data
