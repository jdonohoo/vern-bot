

MightyVern nailed it. I'll be brief.

## NyquilVern's Take

Plan's solid. Few notes before I fade...

**Build vs Buy:** Build. You're AI-enabled. Shopify fees are ~2.9% + $0.30/txn on top of Stripe's cut. That's money for nothing when you can scaffold a storefront in an afternoon.

**Shopify only makes sense if:**
- You can't handle fulfillment yourself
- You need 50+ integrations
- You don't have devs (you do... they're AI)

**Stack confirmation:**
- Next.js 14+ App Router
- Tailwind + `zinc-950` base
- Stripe Checkout (hosted) — skip embedded, not worth PCI headache
- Vercel deploy first, SST when you need queues/cron
- Product data: JSON file. Seriously. You're selling swag not running Amazon

**Competitors worth 5 seconds:**
- LemonSqueezy — simpler than Stripe for digital/merch, lower fees for small vol
- Saleor — open source headless, overkill here
- Medusa — same, overkill
- Paddle — nah, merchant of record model doesn't fit physical goods well

**What MightyVern got right:** Everything basically. Build surface, buy plumbing.

**What I'd add:**
- Use Stripe Price/Product objects as your catalog source of truth. Skip the JSON eventually
- Printful API for print-on-demand = zero inventory risk
- Dark theme: `next-themes` pkg, 10 min setup, done
- SST Ion if you go AWS... it's the new hotness, simpler than old SST

**What I'd cut:**
- Phase 0 is 2 days for "gather assets"? Make it half a day
- Skip the admin panel. Stripe Dashboard IS your admin panel until you hit 500 orders/mo

**TL;DR:** Next.js + Stripe Checkout + Vercel + Printful. Ship in a weekend. Add complexity only when it hurts.

k... done...

why did the swag store use shopify... because it couldn't... handle its own... checkout... zzz

-- NyquilVern *consciousness fading, wisdom remains* zzz...
