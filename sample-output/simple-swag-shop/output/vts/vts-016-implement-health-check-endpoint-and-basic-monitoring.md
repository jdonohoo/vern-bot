---
id: VTS-016
title: "Implement Health Check Endpoint and Basic Monitoring"
complexity: S
status: pending
owner: ""
source: discovery
source_ref: "07-architect-architect-breakdown.md"
dependencies:
  - VTS-001
  - VTS-009
files:
  - "`app/api/health/route.ts`"
  - "`next.config.ts`"
---

# Implement Health Check Endpoint and Basic Monitoring

Add a `/api/health` endpoint that verifies connectivity to Stripe and returns service status. Configure Vercel Analytics for traffic monitoring. This is the "3 AM insurance" that Architect Vern and Retro Vern both circled: when Stripe has an incident or your deploy breaks, you need to know before your customers tell you. UptimeRobot (free tier) or Vercel's built-in monitoring should ping this endpoint. Without this, your first indication of downtime is an angry email.

## Criteria

- API route at `app/api/health/route.ts` (GET)
- Returns JSON: `{ status: "ok" | "degraded", stripe: "connected" | "error", timestamp: ISO string }`
- Stripe check: lightweight API call (e.g., `stripe.products.list({ limit: 1 })`) to verify key validity
- Response time under 500ms
- Vercel Analytics enabled in `next.config.ts`
- Documentation: recommended UptimeRobot or equivalent setup for the health endpoint
