

# Startup Vern's Take: Prompt Library MVP

Alright, let's cut through this. I just read a *very thorough* plan that's about 3x more document than we need for what is essentially a two-page static site. Let me strip this down to what actually matters.

---

## The Hypothesis

**People who use LLMs want a clean, bookmarkable place to discover and save useful prompts.**

That's it. That's the bet. Everything else is noise until we validate that.

## What's the MVP Here?

Two pages. Dark theme. Deploys to Vercel. Done.

The prior plan is solid on the *what*, but it's over-engineered on the *thinking about the what*. A data model? A `data/prompts.ts` file for "future expansion"? A `Container.tsx` component for two pages? Cut it. Cut it. Cut it.

Here's what actually ships:

### The Ruthlessly Scoped File Structure

```
prompt-library/
├── app/
│   ├── layout.tsx       ← dark theme, metadata, nav
│   ├── page.tsx         ← hero + CTA
│   ├── about/
│   │   └── page.tsx     ← what this is & how to use it
│   └── globals.css      ← tailwind imports + dark base
├── tailwind.config.ts
├── postcss.config.mjs
├── tsconfig.json
├── next.config.ts
└── package.json
```

That's **8 files**. The prior plan has 14. We just saved 43% of our file count, which means 43% less stuff to debug when something breaks at 2am.

### What I Cut and Why

| Cut | Why |
|---|---|
| `components/Container.tsx` | A `<div className="max-w-3xl mx-auto px-6">` is not a component. It's a div. Inline it. |
| `components/Header.tsx` | Two nav links. Put them in `layout.tsx`. Extract when you have 5+ pages. |
| `data/prompts.ts` | You have zero prompts to serve dynamically. This is speculative architecture. YAGNI. |
| `public/favicon.ico` | Next.js handles this. Don't overthink it. |
| `README.md` | You asked me not to write files, but also — ship first, document later. |

### What Each File Actually Does

**`app/layout.tsx`** — The shell. Sets `<html className="dark">`, adds a `<body>` with your dark background (`bg-gray-950 text-gray-100`), drops in two nav links (Home, About) in a simple header. Wraps `{children}`. That's it. This is a server component by default — don't add `"use client"` unless you need interactivity, which you don't.

**`app/page.tsx`** — The landing page. Three things, vertically centered:
1. Big bold title: "Prompt Library"
2. One-liner description: "Discover and save useful prompts for LLMs."
3. A fat CTA button linking to `/about`

No prompt cards. No search. No filtering. Just the value prop and a single action. If nobody clicks that button, the idea doesn't work and everything else is moot.

**`app/about/page.tsx`** — The explainer. Answer three questions:
1. What is a prompt library?
2. Why would you use one?
3. How do you get started?

Keep it to ~200 words. If you can't explain it in 200 words, the product is too complicated.

**`app/globals.css`** — Three Tailwind directives (`@tailwind base/components/utilities`) and maybe 5 lines of base styles for the dark theme. That's it.

**`tailwind.config.ts`** — Default config. Enable `darkMode: "class"`. Maybe add one accent color. Don't build a design system for two pages.

### Setup: 4 Commands to Shipped

```bash
npx create-next-app@latest prompt-library --typescript --tailwind --app --src-dir=false --import-alias="@/*"
cd prompt-library
npm run dev          # local validation
# push to GitHub → connect to Vercel → deployed
```

Vercel auto-detects Next.js. Zero config. Zero custom `vercel.json`. If you're writing deployment config, you're doing it wrong.

---

## What I Actually Agree With From the Prior Plan

- **Static rendering only.** Correct. No client components, no `"use client"`, no state. Server components all the way. This is a brochure site that happens to be built with React.
- **Dark theme with high contrast.** Yes. `bg-gray-950` background, `text-gray-100` body text, a bright accent (blue-500 or violet-500) for the CTA button. Don't get cute with gray-on-gray.
- **Mobile-first, centered layout.** Absolutely. `max-w-3xl mx-auto px-6` on the content wrapper. Responsive by default. Done.
- **Vercel-ready with zero backend.** This is the whole point. Static output, edge-deployed, fast everywhere.

## What I Disagree With

**"Define a data model to future-proof."** No. Future-proofing is the enemy of shipping. You know what future-proofs your app? Having users. Build the data model when you have prompts to put in it.

**"Add a static prompt list on the homepage for library feel."** Tempting, but scope creep. The hypothesis isn't "do people want to browse prompts on a static page." The hypothesis is "do people understand and want a prompt library." Validate the concept first, then add content.

**"Container.tsx and Header.tsx components."** Premature abstraction. You have two pages. The shared layout goes in `layout.tsx`. When you hit 4-5 pages and the nav gets complex, extract. Not before.

**The 10-item risk list.** Half of those risks are for features that don't exist yet. "Long prompt text can overflow cards" — what cards? We don't have cards. We have two pages of static text. Focus on the risks that exist *right now*:
1. Does the CTA button contrast enough to get clicks?
2. Does the hero fit above the fold on mobile?
3. Does Vercel deploy work on first push?

That's three risks. Handle those.

---

## The Real Prioritized Enhancement Roadmap

Don't build any of this until the two-page MVP is live and you've seen real traffic:

1. **Add 5-10 static prompt cards on the homepage** — only after you confirm people actually visit and engage
2. **Copy-to-clipboard button** — this is the first *real* feature, the moment the site becomes useful instead of just informational
3. **localStorage favorites** — lets people "save" without auth, zero backend cost
4. **Tag filtering** — only when you have 20+ prompts and users complain about finding things
5. **Search** — only when filtering isn't enough

Everything else (auth, databases, APIs, accounts) is a Series A problem. You're pre-seed. Act like it.

---

## The Metrics That Matter

Track exactly two things at launch:

1. **Homepage → About click-through rate.** If people don't click, your value prop is wrong. Rewrite the copy. Don't add features.
2. **Time on About page.** If they click but bounce in 3 seconds, your explanation isn't compelling. Rewrite it. Still don't add features.

Use Vercel Analytics (free tier). Don't integrate Google Analytics, Mixpanel, or anything else. One tool. Two metrics. Ship it.

---

## Bottom Line

The prior plan is architecturally sound but strategically bloated. It's treating a hypothesis like a product. This is a landing page with a second page. It should take one afternoon to build, not one afternoon to *plan*.

**Build it. Deploy it. Send the link to 10 people who use LLMs. Watch what they do. Then decide what's next.**

The market doesn't care about your `Container.tsx`. The market cares about whether your prompts are useful. And you can't find that out until you ship.

---

*Why did the lean startup cross the road? They didn't — they A/B tested which side had better conversion first. Then they pivoted to a bridge. Ship it!*
