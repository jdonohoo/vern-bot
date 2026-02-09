---
id: VTS-009
title: "Implement Stripe Checkout Session API Route"
complexity: L
status: pending
owner: ""
source: discovery
source_ref: "07-architect-architect-breakdown.md"
dependencies:
  - VTS-002
files:
  - "`app/api/checkout/route.ts`"
  - "`lib/stripe.ts`"
  - "`types/checkout.ts`"
  - "`.env.example`"
---

# Implement Stripe Checkout Session API Route

Server-side API route that creates a Stripe Checkout Session. Supports two modes: direct-to-checkout (single item with variant, no cart required) and cart-based checkout (when VTS-007/008 are implemented). Direct-to-checkout is the Phase 1 MVP path. This is the boundary between your app and Stripe — validate everything here. Don't trust the client-sent prices; look up the real Stripe Price IDs server-side. Redirect the user to Stripe's hosted checkout page. This is where PCI compliance stays simple: you never touch card data.

## Criteria

- API route at `app/api/checkout/route.ts` (POST)
- Accepts EITHER a single item (`{ productId, variantId, quantity: 1 }`) OR a cart items array
- When in direct-to-checkout mode: product detail page sends single item, no cart dependency
- Validates each item: product exists, variant exists, still available
- Maps items to Stripe `line_items` using server-side Stripe Price IDs (never trust client-sent prices)
- Creates Stripe Checkout Session with: `line_items`, `mode: 'payment'`, `success_url`, `cancel_url`
- Stripe Tax: `automatic_tax: { enabled: true }` on session creation (requires Stripe Tax activation in dashboard — document this)
- `shipping_address_collection` enabled with allowed countries configured via env var
- `STRIPE_TAX_ENABLED` env var to toggle tax (allows launch without Stripe Tax activated)
- `.env.example` updated with `STRIPE_TAX_ENABLED` and `SHIPPING_ALLOWED_COUNTRIES` variables
- Returns session URL to client for redirect
- Error handling: returns structured error responses (400 for validation, 500 for Stripe errors)
- Stripe SDK initialized once in a shared module, not per-request
- Rate limiting consideration documented (even if not implemented in v1)
