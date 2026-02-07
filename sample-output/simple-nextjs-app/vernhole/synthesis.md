

# VernHole Discovery: Prompt Library (Next.js App Router)

**Welcome to the VernHole. You asked for this.**

The council has convened — 10 Verns deep. From YOLO's "SEND IT" to Enterprise's 7-week governance timeline, from Retro's "we called this a website in 2003" to Paranoid's "have you checked the CTA button? check it again." Let me synthesize this beautiful chaos.

---

## The Council Speaks

The full analyses are above. Let me distill each Vern's core wisdom:

### YOLO Vern
**Key Insight**: This is a weekend project masquerading as an architecture review. Three commands to scaffold, edit three files, push, deploy. Done.

### Startup Vern
**Key Insight**: The hypothesis is "do people want a prompt library?" — you can't validate that with a plan document. Ship it, send it to 10 people, watch what they do.

### MightyVern
**Key Insight**: The comprehensive details that actually matter — `dvh` vs `vh` for mobile, `next/font` for zero-CLS font loading, Tailwind v4 config differences, and FOUC prevention on dark themes.

### Academic Vern
**Key Insight**: The competitive landscape already exists (PromptBase, FlowGPT, Awesome ChatGPT Prompts). The differentiator must be UX quality. Also: forced dark themes reduce readability for ~33% of users with astigmatism.

### UX Vern
**Key Insight**: A "Prompt Library" with no visible prompts is a brochure, not a product. The CTA button sends users to a dead-end explainer page. Where are the books in this library?

### Ketamine Vern
**Key Insight**: This is a haiku, not a cathedral. Two pages are a design constraint that forces every word and pixel to matter. Honor the simplicity — don't build furniture for a tent.

### Enterprise Vern
**Key Insight**: "Deploy straight to Vercel" is not a deployment strategy. Where's the CI/CD pipeline? The rollback plan? The RACI matrix? *(Conditionally approved pending 6 meetings.)*

### Vernile the Great
**Key Insight**: `tracking-tight` on headings separates professional typography from defaults. `py-20 sm:py-32` for hero padding signals intentionality. Excellence lives in the specific design tokens, not the abstractions.

### Mediocre Vern
**Key Insight**: 5 files. That's what you actually need. The discovery plan has 14. We just saved 43% of the file count, which is 43% less stuff to debug at 2am.

### Retro Vern
**Key Insight**: Your `node_modules` folder will contain ~300MB of dependencies to render two pages of static text. In 2003 we did this with Notepad and FileZilla. The websites got the same amount of simpler. The toolchains did not.

### Paranoid Vern
**Key Insight**: The default Next.js 404 page is white. On your dark-themed site, hitting any wrong URL is a flashbang. Also: pin your dependencies, add security headers, and test the CTA button on Safari. *Especially* Safari.

---

## Synthesis from the Chaos

### Common Themes (9/10 Verns Agree)

1. **Drop `Container.tsx`** — Every single Vern (except Enterprise, who wants a committee to decide) says this is premature abstraction. Inline `max-w-3xl mx-auto px-6`. Extract on the third use.

2. **Drop `data/prompts.ts` from MVP** — 8 of 10 Verns say don't create files for features that don't exist. Academic and MightyVern argue for including it with sample data. Consensus: leave it out until you need it.

3. **The real file structure is ~7 files, not 14** — Layout, homepage, about page, globals.css, tailwind config, package.json, tsconfig. Everything else is noise at this stage.

4. **Static rendering, no `"use client"`** — Universal agreement. Server components all the way. Zero client JS shipped.

5. **Ship fast, iterate from real usage** — Whether phrased as "SEND IT" (YOLO), "MVP or die" (Startup), or "build the shed, then the skyscraper" (Mediocre) — the message is identical.

### Interesting Contradictions

| Topic | Camp A | Camp B |
|---|---|---|
| Show prompts on homepage? | UX, MightyVern, Academic: YES — a library needs books | YOLO, Startup, Ketamine: NO — validate the concept first |
| `data/prompts.ts`? | Academic, MightyVern: include with samples | Everyone else: YAGNI, add when needed |
| Is Next.js overkill? | Retro: "This is two HTML files" | MightyVern, Vernile: The framework investment pays off at scale |
| Governance needed? | Enterprise: 7-week phased rollout | YOLO: "Walk across the street" |
| Route name `/about`? | Most Verns: fine | Academic: `/guide` or `/how-it-works` better matches content intent |
| Enterprise rigor? | Enterprise: RACI matrix, CAB approval, pen testing | Retro: "I deployed this with FileZilla in 2003" |

### The Emergence

Three patterns crystallize from the chaos:

**Pattern 1: The Spectrum of Appropriate Complexity**
Every Vern is calibrating the same dial — how much infrastructure for how much product? The answer depends on intent. If this is a learning project, the full Next.js stack is justified by the learning. If this is a product, ship the minimum and grow. If this is a demo, two HTML files would suffice. The plan is solid for a "scaffold I'll grow into" — just strip the speculative parts.

**Pattern 2: Show, Don't Tell**
UX Vern nailed the core UX problem: a Prompt Library that shows zero prompts is a brochure. Multiple Verns independently arrived at "add 3-8 static prompts on the homepage." This is the single highest-impact change to the plan — it transforms the app from informational to functional. Even hardcoded.

**Pattern 3: The Details That Separate Good from Great**
MightyVern's `dvh` vs `vh`, Vernile's `tracking-tight`, Paranoid's dark-themed 404 page, Academic's contrast-for-astigmatism concern — these micro-decisions compound into the difference between "tutorial project" and "something that feels intentional." The plan covers the architecture well but underspecifies the craft.

### Recommended Path Forward

1. **Scaffold with `create-next-app`** — TypeScript, Tailwind, App Router, no src directory. Accept defaults. 1 minute.

2. **Build exactly 4 files:**
   - `app/layout.tsx` — Dark theme (`bg-gray-950 text-gray-100`), `next/font/google` with Inter, inline nav (just app name + About link), OG metadata
   - `app/page.tsx` — Hero with title, one-line description, CTA button to `/about`. Consider 3-5 static prompt cards below the fold.
   - `app/about/page.tsx` — What, why, how, who. End with a forward-looking CTA.
   - `app/globals.css` — Tailwind directives only

3. **Add the details that matter:**
   - `dvh` not `vh` for mobile viewport
   - `tracking-tight` on headings
   - `focus-visible:ring` on interactive elements
   - Security headers in `next.config.ts`
   - Custom `not-found.tsx` with dark theme
   - OG metadata for link sharing

4. **Deploy to Vercel.** Push to GitHub, import, click deploy. Zero config.

5. **Test:** Click the CTA button. On mobile. On desktop. In Safari. Check contrast with WebAIM. Tab through the site. Run Lighthouse.

6. **Then decide what's next** based on real usage, not speculative architecture.

---

The VernHole has spoken. And remember: Why did 10 Verns walk into a prompt library? The first one said "SEND IT," the last one said "file a JIRA ticket first," and the one in the middle said "we solved this with a CSV in 2004." But they all agreed on one thing: drop the `Container.tsx`. From chaos, `className="max-w-3xl mx-auto px-6"`.
