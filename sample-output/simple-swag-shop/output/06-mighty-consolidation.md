# Swag Storefront Discovery — Master Plan (Next.js, Dark Theme, Stripe, Vercel/SST)

UNLIMITED POWER. I’ve seen this pattern in 47,000 repos. Here’s the comprehensive solution and the edge cases you didn’t ask about.

---

## Executive Summary

You want a dark-themed Next.js swag storefront, deployable on Vercel or SST, with Stripe payments. The right path is **build the branded storefront** and **buy the commerce plumbing** where possible. The decision is not “build vs buy,” it’s **build the surface (brand + UX), buy the heavy ops (payments, taxes, fulfillment)**. Start lean: Next.js + Stripe Checkout + Vercel. Only introduce heavier platforms (Shopify, Medusa, etc.) when operations become a real burden.

---

## Problem Space

You are selling swag, not building a commerce platform. The dominant risks are not technical—they are operational:

- Returns, refunds, and chargebacks
- Inventory and fulfillment errors
- Tax compliance and shipping rules
- Support load if orders go wrong

The system should be **simple, fast to ship, and durable**. It must support a dark theme with strong brand identity and a frictionless payment flow.

---

## Goals and Constraints

### Primary Goals
- Launch a dark-themed swag storefront quickly
- Accept payments reliably via Stripe
- Deploy easily to Vercel or SST
- Keep ongoing ops minimal

### Hard Constraints
- No code or scaffolding in this deliverable
- Plan must be detailed and comprehensive
- Provide a clear architecture path and decision points

---

## Key Decisions

### 1) Checkout Model
**Recommendation:** Stripe Checkout (hosted)  
- Reduces PCI scope, security burden, and compliance surface
- Handles receipts, tax, and localization with minimal work
- Fastest time-to-launch

### 2) Deployment Target
**Recommendation:** Vercel first, SST later  
- Vercel: fastest, simplest, and aligned with Next.js
- SST: ideal once you need AWS-native background jobs or cost control at scale

### 3) Product Data Source
**Recommendation:** Start static, evolve later  
- If SKU count < 50, store products in JSON or a simple CMS
- Add a database only after you validate sales velocity

---

## Architecture Options (Ranked)

### Option A: Next.js + Stripe Checkout (Lean Build)
**Best for:** Speed, flexibility, brand control  
**Tradeoffs:** Minimal ops tooling, manual fulfillment early on

### Option B: Next.js + Headless Shopify
**Best for:** Operational maturity + custom frontend  
**Tradeoffs:** Shopify fees + integration complexity

### Option C: Full Shopify Theme
**Best for:** Time-to-market at all costs  
**Tradeoffs:** Less design control, platform lock-in

**Recommendation:** Option A now, Option B later if scale requires it.

---

## Build vs Buy — Answered Directly

**Build:** the storefront UX, brand, and visual identity  
**Buy:** payment processing, tax handling, fulfillment ops

If you are fully AI-enabled, build speed is not your bottleneck. **Operational complexity is.** Use Stripe Checkout and a simple fulfillment process early. Only adopt Shopify or Medusa when manual ops become painful.

---

## Competitors to Stripe and Shopify (Shortlist Only)

### Payments
- Stripe (default best-in-class for devs)
- Square (if you sell in-person at events)
- Adyen (enterprise-scale, usually overkill)

### Commerce Platforms
- Shopify (default managed commerce)
- BigCommerce (less lock-in, but similar model)
- Medusa (open-source headless, self-hosted)

**Reality:** For swag, Stripe + a minimal ops flow is enough until you’re at scale.

---

## Edge Cases You Didn’t Ask About (But Should)

- **Chargebacks:** Stripe handles disputes, but you must respond with evidence
- **Tax obligations:** Stripe Tax reduces burden but doesn’t eliminate compliance
- **Refund workflows:** Must be consistent and visible to customers
- **Inventory mismatch:** If you sell out, you need a real-time stopgap
- **Webhook failures:** Lost webhooks can mean unfulfilled orders

---

## Operational Plan

### Minimal Ops (Phase 1)
- Stripe Checkout for payments and receipts
- Manual fulfillment tracking
- Stripe dashboard as source of truth
- Basic refund policy

### Scale Ops (Phase 2)
- Add fulfillment automation (Printful or ship integrations)
- Add order DB and admin view
- Add monitoring for webhook health

---

## Implementation Plan (No Code)

### Phase 0 — Decision & Assets (Days 1–2)
- Finalize swag SKUs and pricing
- Gather product photography
- Confirm brand tone and palette

### Phase 1 — Storefront Build (Days 3–5)
- Dark theme layout and product grid
- Product detail flows as needed
- Stripe Checkout integration

### Phase 2 — Post-Checkout Ops (Days 6–7)
- Webhook handling plan
- Order confirmation messaging
- Fulfillment pipeline definition

### Phase 3 — Launch (Day 8)
- Deploy to Vercel
- Smoke test checkout
- Send to a seed group for purchase

---

## UX and Dark Theme Guidance

- Prefer deep neutral backgrounds (`zinc-900`, `slate-950`)
- Use high-contrast text and generous whitespace
- Product imagery should “pop” against dark surfaces
- Keep UI minimal; the swag is the hero

---

## Risk Register

| Risk | Impact | Mitigation |
|------|--------|------------|
| Webhook failure | Missed orders | Retry logic + Stripe event replay |
| Tax miscalculation | Legal/financial risk | Use Stripe Tax or external service |
| Fulfillment errors | Brand trust loss | Manual checks early, automate later |
| Overengineering | Missed launch window | Start with simplest functional stack |

---

## Final Recommendation (Unified Vision)

Start with **Next.js + Stripe Checkout + Vercel**. Focus on swag quality, photography, and a crisp dark-theme storefront. Defer Shopify or Medusa until you have real traction and ops pain. Build the brand surface; outsource the boring infrastructure. This path is the fastest and most pragmatic while preserving future flexibility.

UNLIMITED POWER… and one final truth: Why did the developer bring a ladder to the bugfix? Because they heard the errors were on a higher level.  
-- MightyVern *mic drop*