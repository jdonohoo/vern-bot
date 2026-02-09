---
id: VTS-018
title: "Product Photography Guidelines and Image Asset Pipeline"
complexity: S
status: pending
owner: ""
source: discovery
source_ref: "07-architect-architect-breakdown.md"
dependencies:
  - VTS-005
  - VTS-006
files:
  - "`public/images/`"
  - "`next.config.ts`"
  - "`PRODUCT_IMAGES.md`"
---

# Product Photography Guidelines and Image Asset Pipeline

UX Vern's most urgent flag: dark hoodies disappear on dark backgrounds. This is not a code task — it's a design specification task that the codebase must support. Define the image requirements (dimensions, aspect ratio, background treatment, file format), create placeholder images that demonstrate the correct approach, and document the photography guidelines. The council agreed 60-80% of effort should go to product quality. This task ensures the codebase supports that investment rather than fighting it.

## Criteria

- Written specification: image dimensions (e.g., 1200x1200 minimum), aspect ratio (1:1 for grid, 4:3 for detail), file format (WebP primary, PNG fallback)
- Dark theme product image guidelines: slight gradient background (`zinc-800` to `zinc-850`), subtle shadow/glow to separate product from page background, avoid pure-black product shots
- Placeholder images updated to demonstrate correct treatment (not generic colored rectangles)
- `next.config.ts` image optimization settings configured for the specified dimensions
- Component-level `sizes` attributes in `ProductCard` and `ProductDetail` match the spec
