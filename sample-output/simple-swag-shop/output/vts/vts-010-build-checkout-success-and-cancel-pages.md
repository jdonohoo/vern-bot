---
id: VTS-010
title: "Build Checkout Success and Cancel Pages"
complexity: M
status: pending
owner: ""
source: discovery
source_ref: "07-architect-architect-breakdown.md"
dependencies:
  - VTS-009
files:
  - "`app/checkout/success/page.tsx`"
  - "`app/checkout/cancel/page.tsx`"
  - "`app/api/checkout/session/route.ts`"
---

# Build Checkout Success and Cancel Pages

Post-checkout landing pages. Success page confirms the order and clears the cart. Cancel page reassures the user and sends them back to the cart. The success page should retrieve session details from Stripe to show a confirmation — don't just say "thanks" with no proof. These pages matter more than people think; they're the last impression before the customer waits for their swag.

## Criteria

- Success page at `app/checkout/success/page.tsx`
- Reads `session_id` from URL search params
- Fetches session details from Stripe via API route (`/api/checkout/session`)
- Displays: order confirmation message, customer email, line items summary
- Clears cart state on success page load
- Cancel page at `app/checkout/cancel/page.tsx`
- Cancel page: empathetic message ("No worries"), link back to cart (items preserved)
- Both pages handle missing/invalid session IDs gracefully (redirect to home)
- No sensitive data exposed on the client
