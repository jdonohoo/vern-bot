---
id: VTS-011
title: "Implement Stripe Webhook Handler"
complexity: L
status: pending
owner: ""
source: discovery
source_ref: "07-architect-architect-breakdown.md"
dependencies:
  - VTS-009
files:
  - "`app/api/webhooks/stripe/route.ts`"
  - "`lib/stripe.ts`"
  - "`lib/orders.ts`"
---

# Implement Stripe Webhook Handler

The most critical piece of backend infrastructure. Stripe sends webhook events for completed checkouts, failed payments, and disputes. You must verify the webhook signature (never trust unverified payloads), handle `checkout.session.completed` to confirm orders, and log events for debugging. How will this fail at 3 AM? A webhook fails silently, an order goes unfulfilled, and a customer emails you angry. Add logging. Add retry awareness.

## Criteria

- API route at `app/api/webhooks/stripe/route.ts` (POST)
- Reads raw request body for signature verification (Next.js body parsing must be disabled for this route)
- Verifies Stripe webhook signature using `STRIPE_WEBHOOK_SECRET`
- Handles `checkout.session.completed` event: logs order details (session ID, customer email, amount)
- Handles `payment_intent.payment_failed` event: logs failure details
- Returns 200 to Stripe promptly (processing can be async if needed)
- Returns 400 for invalid signatures with no information leakage
- Idempotent: processing the same event twice doesn't create duplicate side effects
- Structured logging with event type, session ID, and timestamp
