---
id: VTS-005
title: "Build Product Grid and Product Card Components"
complexity: M
status: pending
owner: ""
source: discovery
source_ref: "07-architect-architect-breakdown.md"
dependencies:
  - VTS-001
  - VTS-002
files:
  - "`app/page.tsx`"
  - "`components/product/ProductGrid.tsx`"
  - "`components/product/ProductCard.tsx`"
  - "`lib/formatPrice.ts`"
---

# Build Product Grid and Product Card Components

The main storefront page — a responsive grid of product cards. Each card shows the product image, name, price, and a quick visual cue for availability. The grid is the first thing visitors see. It must load fast, look sharp on dark backgrounds, and make the swag the hero. Product images should "pop" against the dark surface. No carousels, no infinite scroll — keep it honest.

## Criteria

- `ProductGrid` component: responsive CSS grid (1 col mobile, 2 col tablet, 3-4 col desktop)
- `ProductCard` component: product image, name, formatted price, availability badge, link to product detail
- Price formatting via a shared `formatPrice()` utility (handles cents → dollars, currency symbol)
- "Sold Out" visual treatment when `available === false` (desaturated image, overlay badge)
- Images use Next.js `<Image>` with proper `sizes` and lazy loading
- Cards have hover state: subtle scale or border glow (dark theme appropriate)
- Grid renders from static product data (Task 2)
- Lighthouse performance score > 90
