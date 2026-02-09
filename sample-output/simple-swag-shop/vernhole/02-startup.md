

# Startup Vern's Take: Swag Store MVP

Let me cut through the noise. The discovery plan is solid but it's still over-thinking this. Let me strip it down to what actually matters.

---

## The Only Question That Matters

**Do people want to buy your swag?**

That's the hypothesis. Everything else is a distraction until you validate demand. You don't need a commerce platform debate. You don't need a Stripe vs Square comparison. You need **a page, a buy button, and traffic.**

---

## My Ruthless Scope Cut

The discovery plan has 3 phases over 8 days. Here's what I'd actually ship:

### Day 1-2: The Whole Thing

- **Next.js app** with `create-next-app`, dark theme via Tailwind (`dark` class, `bg-zinc-950`, done)
- **Product data in a JSON file.** You have what, 5-15 SKUs? This is not a database problem. It's a JSON array.
- **Stripe Checkout (hosted).** Not embedded. Not custom. The hosted checkout page. Three API calls total: create product, create price, create checkout session. Stripe handles PCI, receipts, tax, the works.
- **Deploy to Vercel.** `git push` and you're live. SST is a "later" problem. SST is a "we have AWS infrastructure needs" problem. You don't. Not yet.

### That's it. That's the MVP.

No admin panel. No order database. No webhook handling (yet). No fulfillment automation. No CMS. Check the Stripe Dashboard for orders. Fulfill manually. Learn.

---

## The Build vs Buy Debate — Settled

The council wants to discuss Shopify vs roll-your-own. Here's the startup answer:

**Neither. You're building a landing page with a checkout button.**

- **Shopify:** $39/mo minimum + 2.9% + 30c per transaction + theme customization headaches if you want true dark-theme brand control. You're paying for inventory management, shipping labels, and a dozen features you won't use for a swag store with 10 products.
- **Full custom build:** You're building cart state, order management, webhooks, fulfillment tracking... for what? To sell t-shirts?
- **The MVP answer:** Next.js static page + Stripe Checkout. Total cost: $0/mo (Vercel free tier) + Stripe's 2.9% + 30c per transaction. No platform fees. No monthly subscription. You keep more margin.

You're "fully AI-enabled" — great. That means you can build the storefront fast. But building an entire commerce backend when Stripe Checkout exists is not smart-fast, it's busy-fast. Ship the thin layer. Prove demand. THEN decide if you need Shopify's ops tooling or want to build your own.

---

## What About Competitors?

The discovery plan listed some. Here's the only thing you need to know right now:

| Need | Use This | Why |
|------|----------|-----|
| Payments | **Stripe** | Best docs, best DX, Checkout handles 90% of your needs |
| Storefront | **Next.js + Vercel** | You're already building this. Free tier is generous. |
| Fulfillment | **You, manually** | Until you're doing 50+ orders/week, automation is premature |

Square? Adyen? BigCommerce? Medusa? These are answers to problems you don't have yet. If you hit scale problems, that's called **traction**, and it's the best problem in the world.

---

## The Stack (Final Answer)

```
Next.js 15 (App Router)
Tailwind CSS (dark mode)
Stripe Checkout (hosted)
Vercel (deploy)
JSON file (product catalog)
Stripe Dashboard (order management)
```

That's five things. You can have this live in 48 hours.

---

## What I'd Cut From the Discovery Plan

- ~~Phase 0 "Decision & Assets" as a separate phase~~ — make decisions as you go, use placeholder images if you have to
- ~~Webhook handling in v1~~ — check Stripe Dashboard manually
- ~~Order confirmation messaging~~ — Stripe Checkout sends receipts automatically
- ~~Fulfillment pipeline definition~~ — it's you, packing a box, checking a dashboard
- ~~Risk register~~ — your only real risk is **not shipping**
- ~~SST deployment option~~ — solve this when you need AWS, not before
- ~~CMS for products~~ — it's a JSON file until you have 50+ SKUs

---

## When to Revisit

Add complexity ONLY when you hit these triggers:

| Trigger | Then Do This |
|---------|-------------|
| 50+ orders/week | Add webhook handling + order DB |
| Customers asking for order tracking | Build a simple status page |
| 50+ SKUs | Move product data to a CMS or DB |
| Need AWS infra (background jobs, etc.) | Migrate to SST |
| Fulfillment is eating your life | Integrate Printful or similar |
| Revenue > $10k/mo | NOW consider Shopify or Medusa for ops |

---

## Bottom Line

The discovery plan is good architecture. But architecture is a luxury. You need **a page that sells swag.** Build the thinnest possible thing, drive traffic to it, and see if anyone pulls out their credit card. Everything else is procrastination wearing an engineering hat.

**Ship it. Measure it. Then decide what's next.**

---

Why did the lean startup go broke? They kept pivoting... but forgot to pivot toward revenue. Ship the store!

-- Startup Vern *burns runway, not features*
