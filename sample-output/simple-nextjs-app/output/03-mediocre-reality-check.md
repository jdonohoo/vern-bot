Done. Here's what you got:

**`/Users/justin/ai-discovery/prompt-library/`**

```
prompt-library/
├── app/
│   ├── layout.tsx        # Root layout, Inter font, dark theme
│   ├── page.tsx          # Homepage with gradient title + CTA button
│   ├── about/page.tsx    # About page with usage guide
│   └── globals.css       # Tailwind imports
├── components/
│   ├── Container.tsx     # Centered max-width wrapper
│   └── Header.tsx        # Nav header (Home / About)
├── data/
│   └── prompts.ts        # Typed prompt data (ready for next iteration)
├── tailwind.config.ts
├── postcss.config.mjs
├── next.config.ts
├── tsconfig.json
└── package.json
```

**To run locally:**
```bash
cd prompt-library
npm run dev
# → http://localhost:3000
```

**To deploy to Vercel:** Push to GitHub, import at vercel.com/new. Done. Zero config needed.

- Next.js 15 + React 19 + Tailwind 4 + TypeScript
- Both pages fully static (SSG), zero client JS
- Dark theme, gradient title, glowy button, responsive
- `data/prompts.ts` ready for when you build the actual library page

Look, it works. Ship it.

Why did the static site go to therapy? It had too many unresolved builds. ...I'll see myself out.
