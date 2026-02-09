---
id: VTS-002
title: "Define Product Data Schema and Stripe Products Catalog"
complexity: M
status: pending
owner: ""
source: discovery
source_ref: "07-architect-architect-breakdown.md"
dependencies:
  - VTS-001
files:
  - "`types/product.ts`"
  - "`data/products-fallback.ts`"
  - "`lib/products.ts`"
---

# Define Product Data Schema and Stripe Products Catalog

Design the product data model backed by Stripe Products. Stripe is the single source of truth for products and prices. Local TypeScript interfaces define the shape, but `getAllProducts()` reads from Stripe's API (cached aggressively since products change rarely). A static fallback file exists for development and build-time static generation, but production reads from Stripe. This means adding a product is done in the Stripe Dashboard — no code deploy needed. The static JSON approach creates a data sync liability from Day 1; using Stripe Products eliminates price sync bugs entirely — the price charged is always the price displayed because they come from the same source.

## Criteria

- TypeScript interface `Product` defined with: `id`, `name`, `slug`, `description`, `priceInCents`, `currency`, `images[]`, `variants[]` (size/color), `stripePriceId`, `available` boolean, `category`
- TypeScript interface `ProductVariant` with: `id`, `label`, `stripePriceId`, `available`
- `lib/products.ts` fetches from Stripe Products API (`stripe.products.list()` with `expand: ['data.default_price']`)
- Response mapped to local `Product` interface
- Cache layer: in-memory cache with 5-minute TTL (or ISR revalidation) — products don't change mid-session
- Static fallback: `data/products-fallback.ts` used during `npm run build` if Stripe is unavailable (for static generation)
- `getProductBySlug()` uses Stripe product metadata field `slug` for URL-friendly lookups
- Helper functions: `getAllProducts()`, `getProductBySlug()`, `getProductsByCategory()`
- All helper functions retain the same signatures — consumers don't know the data source changed
- All helpers are typed and return `Product | undefined` where appropriate (no silent nulls)
- Data access functions are pure — no side effects, easy to test, easy to swap out later
- Development mode: works with Stripe test mode products
