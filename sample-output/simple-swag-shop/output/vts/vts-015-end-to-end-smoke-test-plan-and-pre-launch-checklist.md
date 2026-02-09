---
id: VTS-015
title: "End-to-End Smoke Test Plan and Pre-Launch Checklist"
complexity: M
status: pending
owner: ""
source: discovery
source_ref: "07-architect-architect-breakdown.md"
dependencies:
  - VTS-009
  - VTS-010
  - VTS-011
  - VTS-013
  - VTS-016
  - VTS-017
files:
  - "`LAUNCH_CHECKLIST.md`"
---

# End-to-End Smoke Test Plan and Pre-Launch Checklist

Before you ship, you verify. This task produces a manual smoke test script and a pre-launch checklist. Not automated tests (that's a future task) — this is the "can a human buy swag without hitting a wall" verification. Cover the happy path, the sad paths, and the "what if Stripe is having a bad day" path. Also verify both deployment targets work. Must include mobile device testing and live mode preflight.

## Criteria

- Written smoke test script covering:
- Browse products on homepage
- Navigate to product detail, select variant, click "Buy Now" (direct-to-checkout path)
- Complete payment with Stripe test card
- Verify success page shows order confirmation
- Verify webhook fires and logs order
- Verify notification is sent (Slack/email) on purchase completion
- Test cancel flow (return to product page)
- Test 404 page (invalid product slug)
- Test single-variant product (no variant selector shown)
- Test health check endpoint returns `{ status: "ok" }`
- Mobile device testing: iOS Safari and Android Chrome at minimum
- Verify sticky bottom bar and Buy Now button on mobile
- Verify tap targets meet minimum 44x44px
- Live mode preflight section:
- Stripe webhook endpoints re-registered for production
- API keys rotated to live mode
- Stripe Tax registration verified (if enabled)
- Pre-launch checklist:
- All env vars set in production (including `STRIPE_TAX_ENABLED`, `SLACK_WEBHOOK_URL`)
- Stripe webhook endpoint registered and verified
- Stripe test mode vs live mode toggle documented
- Domain and SSL configured
- Vercel deployment successful with production build
- SST deployment tested (if applicable)
- Health check endpoint monitored (UptimeRobot or equivalent)
- Fulfillment checklist: who packs, who ships, how is tracking communicated
- Documented in a `LAUNCH_CHECKLIST.md` file
