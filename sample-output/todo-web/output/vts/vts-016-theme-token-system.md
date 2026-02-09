---
id: VTS-016
title: "Theme Token System and Dark Mode Foundation"
complexity: S
status: pending
owner: ""
source: discovery
source_ref: "07-architect-architect-breakdown.md"
dependencies:
  - VTS-001
files:
  - "src/theme/tokens.css"
  - "src/theme/ThemeProvider.tsx"
  - "tailwind.config.ts"
---

# Theme Token System and Dark Mode Foundation

Define Tailwind theme tokens (CSS custom properties) for all semantic colors. Configure dark mode via Tailwind `darkMode: 'class'`. Create a ThemeProvider component that reads system preference and persists user choice. Every subsequent component uses these tokens exclusively -- no hardcoded color values anywhere.

This is a foundation task. Without it, every component will hardcode colors and require rework when dark mode ships. Building it early (right after scaffolding) saves 8+ hours of refactoring later.

## Tokens

```css
/* Backgrounds */
--bg-base           /* Page background */
--bg-surface         /* Cards, panels */
--bg-elevated        /* Dropdowns, modals */
--bg-hover           /* Interactive hover state */

/* Text */
--text-primary       /* High-emphasis text */
--text-secondary     /* Low-emphasis, metadata */
--text-muted         /* Disabled, placeholder */

/* Borders */
--border             /* Dividers, card borders */
--border-active      /* Focus rings, active states */

/* Semantic */
--accent             /* Primary actions */
--success            /* Completed, positive */
--warning            /* Due soon */
--danger             /* Overdue, destructive */
```

## Criteria

- Token CSS file exists at `src/theme/tokens.css` with all semantic color variables
- Light and dark values defined for every token
- Tailwind config uses `darkMode: 'class'` and references the CSS custom properties
- ThemeProvider component reads `prefers-color-scheme` system preference on first load
- ThemeProvider persists user's explicit choice (light/dark/system) to the store
- Dark/light toggle works and applies the correct class to the document root
- WCAG AA contrast ratios (4.5:1 for normal text) are met in both themes
- No component uses hardcoded color values -- all colors reference theme tokens
