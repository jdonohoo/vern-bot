---
id: VTS-001
title: "Project Scaffolding and Tooling"
complexity: S
status: pending
owner: ""
source: discovery
source_ref: "07-architect-architect-breakdown.md"
dependencies: []
blocked_by: []
blocks:
  - VTS-016
files:
  - "package.json"
  - "vite.config.ts"
  - "tailwind.config.ts"
  - "src/App.tsx"
  - "src/main.tsx"
  - "vercel.json"
---

# Project Scaffolding and Tooling

Initialize Vite + React + TypeScript project. Configure Tailwind CSS with a custom config stub (theme tokens are deferred to VTS-016). Configure ESLint and Prettier. Install core dependencies (Zustand, date-fns, nanoid, @dnd-kit/core). Create the project directory structure. Verify Vercel deployment config handles SPA routing.

Dark theme setup is NOT part of this task -- that responsibility belongs to VTS-016 (Theme Token System & Dark Mode Foundation).

## Criteria

- `npm run dev` serves a blank page with no errors
- Vite + React + TypeScript configured and building cleanly
- Tailwind installed with custom config stub (empty `extend` block ready for VTS-016 tokens)
- ESLint and Prettier configured and passing on all files
- Project structure directories created (`src/components/`, `src/store/`, `src/types/`, `src/utils/`, `src/theme/`)
- Vercel config handles SPA routing (rewrites to index.html)
- All core dependencies installed and importable
- No dark theme tokens or color system in this task (see VTS-016)
