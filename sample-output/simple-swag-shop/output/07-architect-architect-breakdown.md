

# Swag Storefront — Vern Task Spec (VTS) Breakdown

*Architect Vern here. I've read the master plan. It's solid — Option A is the right call. Now let me turn this into tasks that a team (or a fleet of AI agents) can actually execute without ambiguity. Every task has clear boundaries, clear acceptance criteria, and clear failure modes.*

*Updated per Oracle Vision review: 3 tasks added, 5 tasks modified, dependency chain restructured for direct-to-checkout-first approach. VTS-007 and VTS-008 demoted to Phase 2 (non-blocking for MVP).*

*Let's build this thing right.*

---

## System Architecture Overview

Before we break ground, here's what we're building:

```
┌─────────────────────────────────────────────────────┐
│                    Vercel / SST                      │
│  ┌───────────────────────────────────────────────┐  │
│  │              Next.js App (App Router)          │  │
│  │  ┌─────────┐  ┌──────────┐  ┌─────────────┐  │  │
│  │  │  Pages  │  │   API    │  │  Stripe      │  │  │
│  │  │  (SSR/  │  │  Routes  │  │  Products    │  │  │
│  │  │  SSG)   │  │  /api/*  │  │  (Source of  │  │  │
│  │  │         │  │          │  │   Truth)     │  │  │
│  │  └────┬────┘  └────┬─────┘  └──────────────┘  │  │
│  │       │            │                           │  │
│  │       │     ┌──────┴──────┐                    │  │
│  │       │     │  Stripe     │                    │  │
│  │       │     │  Checkout   │                    │  │
│  │       │     │  Session    │                    │  │
│  │       │     └──────┬──────┘                    │  │
│  │       │            │                           │  │
│  │       │     ┌──────┴──────┐                    │  │
│  │       │     │  Webhook    │──> Notifications   │  │
│  │       │     │  Handler    │    (Slack/Email)   │  │
│  │       │     └─────────────┘                    │  │
│  │       │                                        │  │
│  │       │     ┌──────────────┐                   │  │
│  │       │     │  Health      │──> Monitoring     │  │
│  │       │     │  /api/health │    (UptimeRobot)  │  │
│  │       │     └──────────────┘                   │  │
│  └───────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────┘
          │                        │
          ▼                        ▼
   ┌─────────────┐         ┌─────────────┐
   │   Stripe    │         │   Slack /   │
   │   (Products,│         │   Email     │
   │   Payments, │         │   Notifs    │
   │   Tax,      │         │             │
   │   Receipts) │         │             │
   └─────────────┘         └─────────────┘
```

**Key principles:**
- The storefront is a thin presentation layer. Stripe does the heavy lifting.
- Stripe Products is the source of truth for catalog and pricing. No local JSON price sync liability.
- Direct-to-checkout is the MVP purchase path. Cart is Phase 2.
- Monitoring and notifications are first-class citizens, not afterthoughts.

---

## VTS Task Index

Individual VTS files: `vts/`

### Phase 1 — Ship the MVP

| ID | Task | Complexity | Dependencies |
|----|------|------------|--------------|
| VTS-001 | Initialize Next.js Project with App Router and Dark Theme Foundation | M | None |
| VTS-002 | Define Product Data Schema and Stripe Products Catalog | M | VTS-001 |
| VTS-003 | Configure Deployment for Vercel and SST | M | VTS-001 |
| VTS-004 | Build Site Header, Footer, and Navigation Shell | M | VTS-001 |
| VTS-005 | Build Product Grid and Product Card Components | M | VTS-001, VTS-002 |
| VTS-006 | Build Product Detail Page (with Buy Now) | M | VTS-002, VTS-005 |
| VTS-009 | Implement Stripe Checkout Session API Route (direct-to-checkout + cart) | L | VTS-002 |
| VTS-010 | Build Checkout Success and Cancel Pages | M | VTS-009 |
| VTS-011 | Implement Stripe Webhook Handler | L | VTS-009 |
| VTS-012 | Build Static Pages (About, Privacy, Terms) | S | VTS-004 |
| VTS-016 | Implement Health Check Endpoint and Basic Monitoring | S | VTS-001, VTS-009 |
| VTS-017 | Purchase Notification Pipeline and Fulfillment Trigger | S | VTS-011 |
| VTS-018 | Product Photography Guidelines and Image Asset Pipeline | S | VTS-005, VTS-006 |

### Phase 2 — Post-Launch (When Demand Validates)

| ID | Task | Complexity | Dependencies |
|----|------|------------|--------------|
| VTS-007 | Implement Client-Side Cart State Management | M | VTS-001 |
| VTS-008 | Build Cart Page (with slide-out drawer) | M | VTS-007 |

### Polish and Launch (Parallel to Phase 1, Required Before Go-Live)

| ID | Task | Complexity | Dependencies |
|----|------|------------|--------------|
| VTS-013 | Error Handling, Loading States, and Edge Case UI | M | VTS-005, VTS-006, VTS-009 |
| VTS-014 | SEO, Open Graph, and Performance Optimization | M | VTS-005, VTS-006 |
| VTS-015 | End-to-End Smoke Test Plan and Pre-Launch Checklist | M | VTS-009, VTS-010, VTS-011, VTS-013, VTS-016, VTS-017 |

---

## Task Dependency Graph

```
Phase 1 (Critical Path — Ship the MVP):
  VTS-001 (Foundation)
    ├── VTS-002 (Product Schema + Stripe Products)
    │     ├── VTS-005 (Product Grid)
    │     │     └── VTS-006 (Product Detail + Buy Now)
    │     │           └── VTS-018 (Image Guidelines) [parallel]
    │     └── VTS-009 (Checkout API — direct-to-checkout mode)
    │           ├── VTS-010 (Success/Cancel Pages)
    │           └── VTS-011 (Webhook Handler)
    │                 └── VTS-017 (Notifications)
    ├── VTS-003 (Deployment)
    ├── VTS-004 (Header/Footer/Nav)
    │     └── VTS-012 (Static Pages)
    └── VTS-016 (Health Check + Monitoring)

Phase 2 (Post-Launch, When Demand Validates):
  VTS-007 (Cart State) ← VTS-001
    └── VTS-008 (Cart Page + Drawer) ← VTS-007
          └── VTS-009 (Checkout API — cart mode already supported)

Polish and Launch (Parallel, Required Before Go-Live):
  VTS-013 (Error Handling) ← VTS-005, VTS-006, VTS-009
  VTS-014 (SEO/OG/Perf) ← VTS-005, VTS-006
  VTS-015 (Smoke Test) ← VTS-009, VTS-010, VTS-011, VTS-013, VTS-016, VTS-017
```

## Critical Path

**VTS-001 → VTS-002 → VTS-009 → VTS-011 → VTS-017 → VTS-015**

This is the spine of the system. The key change from the original plan: **VTS-007 (cart) is no longer on the critical path.** Direct-to-checkout via "Buy Now" eliminates the cart dependency for MVP launch. The cart becomes additive when customer behavior validates multi-item purchasing.

---

## Key Changes from Original Plan (per Oracle Vision)

| Change | What | Why |
|--------|------|-----|
| VTS-002 updated | Stripe Products as source of truth (was static JSON) | Eliminates price sync bugs — price displayed = price charged |
| VTS-006 updated | Added "Buy Now" direct-to-checkout, mobile-first criteria | MVP purchase path without cart dependency |
| VTS-007 demoted | Phase 2, non-blocking | Cart adds complexity with marginal benefit for < 15 SKUs |
| VTS-008 demoted | Phase 2, non-blocking | Not needed until cart state exists |
| VTS-009 updated | Dual mode (single item + cart), Stripe Tax, removed VTS-007 dependency | Supports both MVP (Buy Now) and Phase 2 (cart) checkout |
| VTS-013 updated | Removed VTS-008 dependency | Cart page is Phase 2; error handling shouldn't wait for it |
| VTS-015 updated | Added VTS-016, VTS-017 dependencies; mobile testing; live mode preflight | Can't verify launch readiness without monitoring and notifications |
| VTS-016 added | Health check endpoint + monitoring | "3 AM insurance" — know before customers tell you |
| VTS-017 added | Purchase notification pipeline | Logging without alerting is a diary nobody reads |
| VTS-018 added | Product photography guidelines | Dark hoodies on dark backgrounds = invisible products |

---

## Trade-offs and Assumptions

| Decision | Trade-off | Revisit When |
|----------|-----------|--------------|
| Stripe Products as catalog | Requires Stripe API call (cached) vs instant static read | Product count > 100 or need offline catalog editing |
| Direct-to-checkout MVP | No multi-item cart in Phase 1 | Customers demonstrate multi-item purchasing behavior |
| Stripe Checkout (hosted) | Less checkout UX control | Conversion rate needs custom checkout optimization |
| Vercel primary | AWS lock-in avoided, but limited to Vercel's runtime | Background jobs, cron, or AWS-native services needed |
| No database | No order history in your system | You need order management beyond Stripe dashboard |
| Stripe Tax optional | Tax compliance deferred if not activated | First sale to a state/country with tax obligations |
| US-only shipping at launch | Avoids international VAT/GST complexity | First international order request triggers Paddle/Lemon Squeezy evaluation |

---

## What I'd Watch For

1. **Webhook reliability** — This is the #1 operational risk. Monitor it from day one. VTS-017 ensures you're notified.
2. **Image quality on dark backgrounds** — VTS-018 addresses this. Follow the guidelines before uploading product photos.
3. **Stripe test vs live mode** — Triple-check before launch. A misconfigured key in production is a bad day. VTS-015 now includes a live mode preflight.
4. **Mobile checkout flow** — Test the full flow on a real phone. VTS-015 includes mobile device testing criteria.
5. **No order persistence** — If a webhook fails and isn't retried, you lose the order record locally. Stripe retains it, but consider structured logging as a safety net.
6. **Fulfillment is entirely manual** — VTS-017 notifies you, but who packs the box? Document this in the launch checklist.

---

18 tasks. Clear boundaries. No ambiguity. Phase 1 gets the storefront live with direct-to-checkout. Phase 2 adds the cart when customers prove they want it. The next developer — or the next AI agent — can pick up any task and know exactly what "done" looks like.

Why did the architect demote the shopping cart to Phase 2? Because premature optimization is the root of all evil... but premature *cart*-imization is the root of all delayed launches. ...I'll see myself out.

-- Architect Vern (measure twice, deploy once)
