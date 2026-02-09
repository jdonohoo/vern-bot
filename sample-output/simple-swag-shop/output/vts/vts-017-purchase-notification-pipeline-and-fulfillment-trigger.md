---
id: VTS-017
title: "Purchase Notification Pipeline and Fulfillment Trigger"
complexity: S
status: pending
owner: ""
source: discovery
source_ref: "07-architect-architect-breakdown.md"
dependencies:
  - VTS-011
files:
  - "`lib/notifications.ts`"
  - "`app/api/webhooks/stripe/route.ts`"
  - "`.env.example`"
---

# Purchase Notification Pipeline and Fulfillment Trigger

When a purchase completes (via webhook), send a notification to the team. The council unanimously agreed: Stripe Dashboard is your admin panel for Phase 1, but you still need a push notification — nobody checks dashboards proactively. This is the bridge between "payment received" and "swag shipped." Start with email or Slack webhook; the mechanism matters less than the existence of the alert. Every Vern who discussed operations flagged this gap implicitly — the webhook handler (VTS-011) logs events, but logging without alerting is just writing a diary nobody reads.

## Criteria

- On `checkout.session.completed`, send notification with: customer email, items purchased, total amount, session ID
- Notification channel: Slack webhook (preferred) or email via Resend/SendGrid
- Notification includes a direct link to the Stripe Dashboard session
- Failure to notify does NOT block the webhook response (fire-and-forget with error logging)
- Environment variable for notification target (`SLACK_WEBHOOK_URL` or `NOTIFICATION_EMAIL`)
- `.env.example` updated with new variable
