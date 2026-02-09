---
id: VTS-007
title: "Implement Client-Side Cart State Management"
complexity: M
status: pending
phase: 2
owner: ""
source: discovery
source_ref: "07-architect-architect-breakdown.md"
dependencies:
  - VTS-001
files:
  - "`lib/cart/CartProvider.tsx`"
  - "`lib/cart/cartReducer.ts`"
  - "`types/cart.ts`"
  - "`app/layout.tsx`"
---

# Implement Client-Side Cart State Management

**Phase 2 — Non-blocking for MVP.** The council agreed: for fewer than 15 SKUs, direct-to-checkout per item is the MVP. Cart adds complexity with marginal benefit until customers demonstrate multi-item purchasing behavior. This task remains valid but is no longer a blocker for the checkout flow. VTS-009 supports both direct-to-checkout (single item) and cart-based checkout.

Cart state that persists across page navigations and survives browser refresh. Use React Context + `localStorage` — no external state library needed for this scale. The cart holds items with product ID, variant ID, quantity, and price. This is client-side only; Stripe is the server-side source of truth for actual purchases. How will this fail? `localStorage` isn't available during SSR. Guard against it.

## Criteria

- `CartProvider` context wrapping the app in the root layout
- Cart state: array of `CartItem` objects (`productId`, `variantId`, `name`, `variantLabel`, `priceInCents`, `quantity`, `image`)
- Actions: `addToCart()`, `removeFromCart()`, `updateQuantity()`, `clearCart()`
- `getCartTotal()` and `getCartItemCount()` derived values
- Persistence to `localStorage` with hydration guard (no SSR crash)
- Cart state initializes empty on server, hydrates from `localStorage` on client
- Max quantity per item enforced (e.g., 10) to prevent accidental bulk orders
- Cart provider must gracefully handle the case where it's not yet integrated — components should work with and without cart context
- Add a feature flag or conditional: if cart is disabled, "Add to Cart" becomes "Buy Now" and routes directly to checkout
