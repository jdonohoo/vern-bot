# Retro Vern's Take: The Swag Store That Doesn't Need To Be Complicated

*cracks knuckles, adjusts reading glasses, sips black coffee from a mug that says "I survived Perl CGI"*

---

## Let Me Tell You What This Actually Is

Strip away the buzzwords — Next.js, SST, "AI-enabled shop," headless commerce — and what do you have?

**A product catalog with a shopping cart and a checkout page.**

That's it. That's what this is. We were building these in 2003 with PHP and a MySQL table called `products`. The HTML was ugly. The CSS was `<table>` tags. But you know what? People bought stuff.

The MightyVern plan above is solid — I'll give credit where it's due. But let me add the perspective of someone who's watched seventeen generations of developers reinvent the same shopping cart.

---

## The Build vs Buy Question (Historical Context)

This debate is as old as software itself. I remember when the question was "do we use osCommerce or write our own?" Then it was "Magento or custom?" Then "WooCommerce or custom?" Now it's "Shopify or custom?"

The answer has always been the same, and it hasn't changed:

**How many SKUs are you selling, and how many orders per day?**

If you're selling 10-50 swag items and processing a handful of orders a day? You don't need Shopify. You don't need Medusa. You need a static product list, a Stripe Checkout button, and a webhook that sends you an email when someone buys something.

We used to call this a "brochure site with a buy button." It worked then. It works now.

### Shopify: What You're Actually Paying For

Shopify isn't a technology purchase. It's an **operations** purchase. You're paying for:

- Inventory management
- Shipping label generation
- Return/refund workflows
- Tax calculation across jurisdictions
- Customer accounts and order history
- 24/7 support when things break at 2am

If you need those things today? Use Shopify. Seriously. The 2.9% + 30¢ per transaction isn't a fee — it's insurance against building all that yourself.

If you don't need those things today? **Don't pay for them.** You can always migrate later. I've seen more projects die from premature optimization than from scaling problems.

### The "AI-Enabled" Angle

You mention you're a fully AI-enabled shop, so build speed is high. Fair. But here's what I've learned across every era of "we can build it faster now":

**Building it fast and maintaining it forever are two different problems.**

Every line of custom commerce code is a line you maintain when tax laws change, when Stripe updates their API, when a customer in the EU invokes GDPR. The question isn't "can we build it?" — it's "do we want to own it?"

---

## Competitors Worth Knowing About

Since you asked, here's the honest landscape from someone who's watched these markets consolidate and fragment repeatedly:

### Payments (Alternatives to Stripe)

| Provider | Honest Take |
|----------|-------------|
| **Stripe** | Best developer experience, period. Has been since ~2011. Default choice unless you have a specific reason not to. |
| **Square** | Better if you ever sell at physical events (pop-up shops, conferences). Their online APIs are fine but not as polished as Stripe. |
| **Paddle** | Handles sales tax and VAT as merchant of record. If tax compliance scares you, this removes it entirely. Worth a look. |
| **LemonSqueezy** | Similar merchant-of-record model. Simpler than Stripe if you're selling digital goods, but works for physical too. |
| **PayPal** | Still exists. Still works. Your customers' parents trust it. Don't laugh — conversion rates matter. |

**My call:** Stick with Stripe. It's boring in the best way. It works. The docs are good. The dashboard is good. The webhooks are reliable. I've been using it since they were still in beta, and it's never given me a reason to leave.

### Commerce Platforms (Alternatives to Shopify)

| Platform | Honest Take |
|----------|-------------|
| **Shopify** | The PostgreSQL of commerce platforms — it just works and it's been working. Headless (Storefront API) is mature. |
| **BigCommerce** | Less lock-in than Shopify, decent headless API, but smaller ecosystem. Slightly less polished. |
| **Medusa.js** | Open-source headless commerce. Self-hosted. Great if you want full control and don't mind running infrastructure. |
| **Saleor** | Another open-source option. Python/GraphQL. Solid but smaller community. |
| **Commerce.js** | API-first, lightweight. Good for exactly your use case (small catalog, custom frontend). |
| **Snipcart** | Drop a script tag on any site, get a cart. Been around since 2013. Still works. Sometimes the simplest tool is the right one. |

**My call:** For a swag store with a small catalog, Snipcart or a raw Stripe Checkout integration is probably all you need. Medusa is interesting if you want to grow into a real commerce operation without Shopify fees. But don't adopt it on day one — that's premature architecture.

---

## What I Actually Recommend

Here's the thing the MightyVern plan got right: **Option A now, evolve later.** But let me be more specific about what "Option A" really looks like, with less ceremony.

### The Boring Stack That Works

```
Next.js (App Router, because fine, it's 2026)
├── Static product data (JSON file or Stripe Products — they already have a product catalog built in)
├── Tailwind CSS (dark theme — zinc-900 backgrounds, done)
├── Stripe Checkout (hosted — never touch card numbers yourself, EVER)
├── Stripe Webhooks (one API route: /api/webhooks/stripe)
├── Vercel deployment (git push = deploy, same as it ever was)
└── Email notification on purchase (Resend, Postmark, or even just a webhook to Slack)
```

That's it. That's the whole system.

You know what's NOT in that list? A database. An ORM. A headless CMS. A state management library. An admin panel. A microservices architecture.

**You don't need them yet.** And "yet" might be "never" for a swag store.

### Why Stripe Products Instead of a JSON File or CMS

Here's a trick from the old days (well, 2015, but that's old in web years): **Stripe already has a product catalog.** You can create products and prices in the Stripe Dashboard, then fetch them with the API. Your source of truth for products, prices, and inventory lives in the same system that processes payments.

No sync issues. No "the CMS says $25 but Stripe charged $30" bugs. One source of truth. We used to call this "not duplicating your data," and it's still good advice.

### The Dark Theme Thing

I've been staring at dark terminals since before "dark mode" was a UI trend. Here's what actually matters:

- Background: one dark color. Not a gradient. Not glassmorphism. Just dark.
- Text: off-white. Not pure white (#fff burns the eyes at night). Use `#e4e4e7` or similar.
- Product images: they ARE the color. Let the swag pop against the dark background.
- Accent color: one. For buttons and links. That's it.

Tailwind makes this trivial. `dark:bg-zinc-950 dark:text-zinc-100`. We're done. This is not a design system problem — it's a CSS variables problem, and CSS variables have been in every browser since 2017.

---

## The Deployment Question: Vercel vs SST

**Start with Vercel.** I know SST is cool and gives you AWS control. But here's the thing:

Vercel is `git push` and you're live. SST is `git push`, then CloudFormation, then pray the IAM roles are right, then wait 3 minutes for the CDK to synthesize, then...

For a swag store? Vercel. The free tier probably handles your traffic. If it doesn't, that's a good problem to have — you're selling a lot of swag.

Move to SST when (and only when):
- You need background jobs that run longer than Vercel's serverless timeout
- You need to reduce hosting costs at scale
- You need AWS-specific services (SQS, DynamoDB, etc.)

Until then? YAGNI. (You Ain't Gonna Need It — coined in the 90s, still true today.)

---

## Should You Focus on Swag Over Tech?

**Yes. A thousand times yes.**

I've watched developers spend six months building the perfect e-commerce platform for products nobody wanted to buy. The product photography, the brand, the swag quality — that's what sells. Not your React rendering strategy.

Here's a story: In 2006, I helped a band sell merch online. We used a flat HTML page, a PayPal button, and a spreadsheet to track orders. They sold $40K of t-shirts in three months. The site looked like it was designed by someone who'd just discovered the `<blink>` tag. Nobody cared. The merch was good.

**The tech is the easy part. The swag is the hard part.**

Spend 80% of your energy on:
- What swag to sell
- Product photography (dark backgrounds, good lighting)
- Copywriting that makes people want to rep your brand
- Pricing that makes impulse purchases easy

Spend 20% on the tech. Which, with the stack above, is about 2-3 days of work.

---

## My Honest Assessment of the MightyVern Plan

It's good. Genuinely. The phased approach is right. The "build the surface, buy the plumbing" framework is exactly correct.

Where I'd push back:

1. **8 days is generous.** This is a 3-day build if you don't overthink it. Day 1: scaffold + products + dark theme. Day 2: Stripe integration + webhooks. Day 3: polish + deploy + test a real purchase.

2. **The risk register is enterprise-grade for a swag store.** Webhook failures? Set up Stripe's automatic retries (they've had this for years) and check your dashboard once a day. You don't need a monitoring system for 10 orders a week.

3. **"Evolve to Option B" should have a trigger.** Don't just say "later." Say "when we're processing more than 100 orders/month and manually fulfilling them becomes painful." Concrete thresholds prevent premature scaling.

---

## The One-Page Plan (Retro Vern Edition)

1. **Create products in Stripe Dashboard.** Pricing, images, descriptions. Done.
2. **Scaffold Next.js with Tailwind.** Dark theme. Product grid. Product pages.
3. **Fetch products from Stripe API** at build time or on request.
4. **"Buy" button → Stripe Checkout** (hosted). One API route to create the session.
5. **Webhook endpoint** → sends you a Slack/email notification. That's your "order management system" for now.
6. **Deploy to Vercel.** Connect the repo. Push to main = live.
7. **Buy a domain.** Point it at Vercel. SSL is automatic (it has been since Let's Encrypt, 2015).
8. **Sell swag.** Focus on the swag.

Total new dependencies: `stripe`, `next`, `tailwindcss`. That's it. Three dependencies. Not thirty.

---

## Final Word

Every few years, the industry convinces itself that selling things online is a hard technical problem. It's not. It's a hard *business* problem — inventory, fulfillment, customer service, tax compliance. The code part? We solved that. Multiple times. Across multiple decades.

Build the simplest thing that works. Sell swag. If you outgrow it, evolve. But I promise you — the bottleneck won't be your Next.js architecture. It'll be running out of medium-sized hoodies.

---

Why did the senior developer refuse to use a shopping cart framework? Because `<form action="/buy" method="POST">` has been taking money since 1995, and it never once needed a node_modules folder.

-- Retro Vern *(been mass-assigning since before GitHub had Issues)*
