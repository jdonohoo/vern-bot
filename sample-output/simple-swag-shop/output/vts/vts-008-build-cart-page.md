---
id: VTS-008
title: "Build Cart Page"
complexity: M
status: pending
phase: 2
owner: ""
source: discovery
source_ref: "07-architect-architect-breakdown.md"
dependencies:
  - VTS-007
files:
  - "`app/cart/page.tsx`"
  - "`components/cart/CartItemRow.tsx`"
  - "`components/cart/CartSummary.tsx`"
  - "`components/cart/CartDrawer.tsx`"
---

# Build Cart Page

**Phase 2 — Non-blocking for MVP.** Not needed for MVP launch when using direct-to-checkout. Implement after cart state management (VTS-007) is in place and customer behavior validates multi-item purchasing demand.

Dedicated `/cart` page showing all items, quantities, line totals, and a grand total. Users can adjust quantities or remove items. The "Checkout" button triggers Stripe Checkout Session creation. Empty cart state should be friendly, not a dead end — show a link back to the shop. This is the last stop before money changes hands; make it clear and trustworthy.

## Criteria

- Route at `app/cart/page.tsx`
- List of cart items with: image thumbnail, name, variant, unit price, quantity selector, line total, remove button
- Grand total displayed prominently
- "Checkout" button triggers POST to `/api/checkout` (VTS-009)
- Loading state on checkout button while session is being created
- Empty cart state: friendly message + "Continue Shopping" link
- Quantity changes update cart state immediately (VTS-007)
- Responsive layout: stacked on mobile, table-like on desktop
- Consider slide-out cart panel (drawer) as primary interaction on mobile, with `/cart` page as fallback/desktop view
- "Buy Now" direct-to-checkout option remains available even after cart is implemented (for single-item impulse purchases)
