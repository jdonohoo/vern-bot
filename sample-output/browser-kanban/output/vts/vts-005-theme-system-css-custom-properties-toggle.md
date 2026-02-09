---
id: VTS-005
title: "Theme System (CSS Custom Properties + Toggle)"
complexity: M
status: pending
owner: ""
source: oracle
source_ref: "oracle-architect-breakdown.md"
dependencies:
  - VTS-001
  - VTS-004
files:
  - "`styles/themes.css`"
  - "`styles/variables.css`"
  - "`scripts/components/theme-toggle.js`"
---

# Theme System (CSS Custom Properties + Toggle)

Implement the light/dark theme system using CSS custom properties toggled via a `data-theme` attribute on the document root. Define a complete set of design tokens. Build the toggle UI component. Persist preference through the state manager. Respect `prefers-color-scheme` as the initial default if no preference is saved.

## Criteria

- CSS variables defined for both `light` and `dark` themes covering: background, surface, card, text-primary, text-secondary, border, accent, urgency-green, urgency-yellow, urgency-red, urgency-neutral
- Theme applied via `data-theme="light|dark"` on `<html>` element
- Toggle component (button or switch) in the header/toolbar area
- Theme preference persisted via `StateManager.updatePreferences()`
- On load: reads saved preference, falls back to `prefers-color-scheme`, falls back to light
- Transition animation on theme switch (subtle, 200ms)
- All urgency colors maintain accessible contrast ratios in both themes
