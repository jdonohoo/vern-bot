---
id: VTS-001
title: "Initialize Next.js Project with App Router and Dark Theme Foundation"
complexity: M
status: pending
owner: ""
source: discovery
source_ref: "07-architect-architect-breakdown.md"
dependencies: []
files:
  - "`package.json`"
  - "`tsconfig.json`"
  - "`tailwind.config.ts`"
  - "`app/layout.tsx`"
  - "`app/globals.css`"
  - "`app/page.tsx`"
  - "`.eslintrc.json`"
  - "`prettier.config.js`"
---

# Initialize Next.js Project with App Router and Dark Theme Foundation

Scaffold a new Next.js project using the App Router (`/app` directory). Configure TypeScript, Tailwind CSS v4, and establish the dark theme design token system. This is the foundation everything else sits on — get it right. Set up the project structure with clear separation of concerns: `/app` for routes, `/components` for UI, `/lib` for business logic, `/data` for static product data, `/types` for shared TypeScript interfaces.

## Criteria

- Next.js 15+ project bootstrapped with App Router
- TypeScript strict mode enabled
- Tailwind CSS configured with dark theme as default (not toggled — dark is the only mode)
- Design tokens defined: background (`zinc-950`, `zinc-900`), surface (`zinc-800`, `zinc-700`), text (`zinc-50`, `zinc-300`), accent color for CTAs
- Base layout component with dark background, proper font stack, and metadata
- ESLint and Prettier configured
- `npm run dev`, `npm run build`, and `npm run lint` all pass cleanly
- Directory structure documented in a top-level comment in the layout file
