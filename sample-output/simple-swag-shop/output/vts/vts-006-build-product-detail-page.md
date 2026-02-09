---
id: VTS-006
title: "Build Product Detail Page"
complexity: M
status: pending
owner: ""
source: discovery
source_ref: "07-architect-architect-breakdown.md"
dependencies:
  - VTS-002
  - VTS-005
files:
  - "`app/products/[slug]/page.tsx`"
  - "`components/product/ProductDetail.tsx`"
  - "`components/product/VariantSelector.tsx`"
  - "`components/product/ImageGallery.tsx`"
  - "`components/product/BuyNowButton.tsx`"
---

# Build Product Detail Page

Dynamic route page for individual products (`/products/[slug]`). Shows full product info: large image(s), name, description, price, variant selector (size/color), and a "Buy Now" button as the primary purchase path. This is where the buying decision happens — clarity and confidence are everything. Variants must be selectable without page reload. If a variant is unavailable, it should be visually disabled, not hidden. The "Buy Now" button is the MVP purchase path (direct-to-checkout); "Add to Cart" is conditionally rendered only when the cart system (VTS-007) is active.

## Criteria

- Dynamic route at `app/products/[slug]/page.tsx`
- Static generation via `generateStaticParams()` from product data
- Large product image with optional thumbnail gallery (if multiple images)
- Product name, description, and formatted price prominently displayed
- Variant selector (buttons or dropdown) — disabled state for unavailable variants
- Variant selector uses large tap targets (minimum 44x44px per WCAG)
- If product has only one variant (or no variants), skip the selector entirely — don't show a dropdown with one option
- "Buy Now" button that creates a Stripe Checkout session for the single selected item+variant (calls VTS-009 API directly)
- "Add to Cart" button conditionally rendered only when cart system (VTS-007) is active
- Mobile layout: product image full-width, sticky bottom bar with price + "Buy Now" button
- Proper 404 handling via `notFound()` if slug doesn't match any product
- SEO metadata via `generateMetadata()` with product name and description
- Breadcrumb navigation back to shop
