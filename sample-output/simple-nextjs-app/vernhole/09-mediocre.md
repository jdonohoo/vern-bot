

# Prompt Library — Mediocre Analysis

Alright, I overslept but I read the plan. Here's my take.

## The Plan is Overcooked

That discovery doc is like 2000 words for a 2-page static site. Data models? Future persistence strategies? localStorage edge cases? C'mon. We're building two pages with a button between them.

Here's what actually matters:

## What You Need

**5 files. That's it.**

```
prompt-library/
├── app/
│   ├── layout.tsx      ← dark theme, metadata, nav
│   ├── page.tsx        ← hero + button
│   ├── about/
│   │   └── page.tsx    ← explainer
│   └── globals.css     ← tailwind directives
├── tailwind.config.ts
├── package.json
└── tsconfig.json
```

The discovery plan has a `components/Container.tsx` and `components/Header.tsx`. For two pages. That's a whole components folder for a `<div className="max-w-3xl mx-auto">` and a `<nav>` with two links. Inline them. We'll refactor later.

Same with `data/prompts.ts` — there's no prompt list in the MVP. Don't create files for features that don't exist yet. That's a tomorrow problem.

## Setup

Three commands:

```bash
npx create-next-app@latest prompt-library --typescript --tailwind --app --src-dir=false --import-alias="@/*"
cd prompt-library
npm run dev
```

Done. Next.js scaffolds everything including `postcss.config.mjs`, `tsconfig.json`, and tailwind config. Don't overthink it.

## Dark Theme

One line in `layout.tsx`: slap `className="dark"` on the `<html>` tag. Set your body bg in globals.css:

```css
body { @apply bg-gray-950 text-gray-100; }
```

The discovery plan talks about "deep background color" and "subtle accent color" and "WCAG contrast ratios." Just use `bg-gray-950` and `text-gray-100`. Contrast ratio is ~15:1. Ship it.

## Homepage

- Big bold title: "Prompt Library"
- One sentence: "Discover and save useful LLM prompts."
- One button: links to `/about`
- Center everything with flexbox. Done in 20 lines of JSX.

## About Page

- What is a prompt library
- Why it's useful
- How to use it
- A link back home

Static text. Server component by default. Zero client JS shipped. Exactly what we want.

## The Stuff I'd Skip From the Plan

| Plan says | I say |
|---|---|
| `Container.tsx` component | Inline a div. Two pages don't need shared abstractions |
| `Header.tsx` component | Put the nav in `layout.tsx` directly |
| `data/prompts.ts` | No prompts in MVP. Delete this |
| "Future data model" table | Cool table. Not relevant today |
| "State and persistence strategy" | There is no state. Next question |
| 10 risks and edge cases | It's two static pages. The risk is overthinking it |
| "Validation checklist" | Click the button. Does it go to /about? Great |

## What I Actually Agree With

- Static rendering only. No `"use client"` anywhere. Correct.
- Vercel deploys Next.js with zero config. Just `git push`. Correct.
- Mobile-first, centered layout, responsive. Yes. Tailwind handles this with `max-w-3xl mx-auto px-4`.
- Keep it simple for v1. Ironic advice from a 2000-word plan, but yes.

## Deploy

```bash
git init && git add -A && git commit -m "init"
# push to github
# import repo in vercel dashboard
# click deploy
```

Or `npx vercel` if you're in a hurry. Which I am.

## Verdict

The discovery plan is solid architecturally but it's planning for a skyscraper when we're building a shed. Two pages, one button, dark theme, Tailwind. You can have this running in 15 minutes. The "future enhancements" section is longer than the actual MVP implementation would be.

Build the shed. If people show up, build the skyscraper.

---

Why did the Next.js developer break up with the backend? Because they had no server-side chemistry. Ship it.
