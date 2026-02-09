---
id: VTS-014
title: "SEO, Open Graph, and Performance Optimization"
complexity: M
status: pending
owner: ""
source: discovery
source_ref: "07-architect-architect-breakdown.md"
dependencies:
  - VTS-005
  - VTS-006
files:
  - "`app/layout.tsx`"
  - "`app/sitemap.ts`"
  - "`public/robots.txt`"
  - "`app/products/[slug]/page.tsx`"
---

# SEO, Open Graph, and Performance Optimization

The store needs to be discoverable and shareable. Proper metadata, Open Graph tags for social sharing (someone shares a product link on Twitter and it looks good), sitemap for search engines, and performance tuning. Dark theme sites often get dinged on contrast ratios — verify accessibility. This is the difference between "it works" and "it's professional."

## Criteria

- Root metadata in `app/layout.tsx`: site title, description, Open Graph defaults
- Per-product Open Graph metadata via `generateMetadata()` in product detail page
- `robots.txt` at `/public/robots.txt`
- Sitemap generation (static or via `app/sitemap.ts`)
- All images optimized via Next.js `<Image>` with WebP format and proper sizing
- Accessibility audit: contrast ratios meet WCAG AA on dark backgrounds
- No layout shift (CLS) on product grid or cart
- Bundle size reviewed — no unnecessary dependencies
