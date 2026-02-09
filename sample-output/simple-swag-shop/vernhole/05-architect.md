

# Architect Vern's Systems Review: Swag Storefront

*Cracks knuckles. Pulls up the whiteboard. Pours a third coffee.*

Alright. MightyVern did solid work here. The bones are good. But I'm going to do what I always do — stress-test the architecture, name the components explicitly, define the boundaries, and make sure the developer who inherits this at 3 AM on a Saturday can figure out what's going on.

Let me walk through this properly.

---

## 1. REQUIREMENTS VALIDATION

Before I draw a single box on a diagram, let me confirm what we're actually building:

**What this IS:**
- A branded swag storefront (not a general e-commerce platform)
- Dark-themed, visually polished — the swag is the product, the site is the brand
- Payment via Stripe
- Deployable to Vercel (primary) or SST (future option)
- Limited SKU count (likely < 50 items)

**What this is NOT:**
- A multi-vendor marketplace
- A platform with user accounts, wishlists, or loyalty programs (yet)
- A fulfillment management system
- A Shopify competitor

That distinction matters. It constrains every decision downstream.

---

## 2. HIGH-LEVEL ARCHITECTURE

Here's how I see the system decomposing:

```
┌─────────────────────────────────────────────────────┐
│                    BROWSER (Client)                  │
│                                                      │
│  ┌──────────┐  ┌──────────────┐  ┌───────────────┐  │
│  │ Product   │  │ Cart State   │  │ Theme Provider │  │
│  │ Catalog   │  │ (React ctx)  │  │ (dark-first)   │  │
│  │ Pages     │  │              │  │               │  │
│  └──────────┘  └──────────────┘  └───────────────┘  │
└────────────────────┬────────────────────────────────┘
                     │
                     │ HTTPS
                     ▼
┌─────────────────────────────────────────────────────┐
│               NEXT.JS SERVER (Vercel / SST)          │
│                                                      │
│  ┌──────────────────┐  ┌─────────────────────────┐  │
│  │ /api/checkout     │  │ /api/webhooks/stripe     │  │
│  │ Creates Stripe    │  │ Receives payment events  │  │
│  │ Checkout Session  │  │ Triggers fulfillment     │  │
│  └──────────────────┘  └─────────────────────────┘  │
│                                                      │
│  ┌──────────────────┐  ┌─────────────────────────┐  │
│  │ Product Data      │  │ Order Confirmation       │  │
│  │ (static JSON or   │  │ (success/cancel pages)   │  │
│  │  headless CMS)    │  │                         │  │
│  └──────────────────┘  └─────────────────────────┘  │
└────────────────────┬────────────────────────────────┘
                     │
                     │ HTTPS (server-side only)
                     ▼
┌─────────────────────────────────────────────────────┐
│                  STRIPE (External)                    │
│                                                      │
│  ┌────────────┐  ┌───────────┐  ┌────────────────┐  │
│  │ Checkout    │  │ Webhooks  │  │ Stripe Tax     │  │
│  │ Sessions    │  │ (events)  │  │ (optional)     │  │
│  └────────────┘  └───────────┘  └────────────────┘  │
└─────────────────────────────────────────────────────┘
```

**Three layers. Three responsibilities. No ambiguity.**

1. **Presentation Layer** — Next.js pages, dark theme, product display, cart UX
2. **API Layer** — Two route handlers: create checkout, receive webhooks
3. **Payment Layer** — Stripe handles the money. We never touch card data.

That's it. If you can't explain your architecture in three layers, it's too complex for the problem.

---

## 3. WHERE I AGREE WITH THE PRIOR PLAN

MightyVern got the big calls right:

- **Stripe Checkout (hosted) over custom payment forms.** Absolutely. Reduces PCI scope to nearly zero. The moment you build a custom card input, you own that liability. Don't.
- **Vercel first, SST later.** Correct. Vercel is the path of least resistance for Next.js. SST is the escape hatch when you need queues, cron jobs, or AWS-native infra. You don't need that on day one.
- **Static product data to start.** Yes. A `products.json` file or a simple content layer is the right call for < 50 SKUs. A database is overhead you don't need yet.
- **Build the surface, buy the plumbing.** This is the correct framing. Your brand lives in the storefront. Stripe is the plumbing. Don't rebuild plumbing.

---

## 4. WHERE I'D SHARPEN THE ARCHITECTURE

### 4a. Cart State Management — Be Explicit

The prior plan waves hands at cart state. Here's the decision:

**Use React Context + `localStorage` persistence.** No external state library needed for a cart with < 50 SKUs. The cart is ephemeral — if the user clears their browser, they lose their cart. That's fine for swag.

The cart context should own exactly three things:
- `items[]` — what's in the cart (product ID, quantity, size/variant)
- `addItem()` / `removeItem()` / `updateQuantity()` — mutations
- `cartTotal` — derived, computed from items

No Redux. No Zustand. No state management library. React Context does this job. When someone tells you a cart needs Redux, ask them how many items they're managing. If the answer is under 100, the answer is Context.

### 4b. The Webhook Handler — This Is Where Systems Fail at 3 AM

The prior plan mentions webhooks but doesn't architect them. Let me be specific.

Your `/api/webhooks/stripe` route handler needs to:

1. **Verify the Stripe signature** — Every webhook, every time. No exceptions. Use `stripe.webhooks.constructEvent()`. If verification fails, return 400 immediately.
2. **Be idempotent** — Stripe will retry failed webhooks. Your handler must handle duplicate events gracefully. Use the event ID as a deduplication key.
3. **Process only the events you care about** — Start with `checkout.session.completed`. Ignore everything else with a 200 response (acknowledge but don't process).
4. **Separate receipt from processing** — Acknowledge the webhook (return 200) quickly. Do any heavy work asynchronously if needed.

```
Webhook arrives
  → Verify signature (fail → 400)
  → Check event type (unknown → 200, ignore)
  → Check idempotency (duplicate → 200, skip)
  → Process event (log, notify, trigger fulfillment)
  → Return 200
```

That's the flow. Five steps. Each one is a guard clause. Early returns all the way down. No nesting deeper than two levels.

### 4c. Project Structure — Files Have Boundaries

Here's how I'd lay out the project. Not the code — the *structure*:

```
src/
├── app/                          # Next.js App Router
│   ├── layout.tsx                # Root layout, dark theme provider, fonts
│   ├── page.tsx                  # Landing / hero section
│   ├── products/
│   │   └── page.tsx              # Product catalog grid
│   ├── cart/
│   │   └── page.tsx              # Cart review page
│   ├── checkout/
│   │   ├── success/page.tsx      # Post-purchase confirmation
│   │   └── cancel/page.tsx       # Checkout cancellation
│   └── api/
│       ├── checkout/route.ts     # POST: create Stripe Checkout session
│       └── webhooks/
│           └── stripe/route.ts   # POST: handle Stripe webhook events
├── components/
│   ├── layout/                   # Header, Footer, Navigation
│   ├── product/                  # ProductCard, ProductGrid, ProductDetail
│   ├── cart/                     # CartDrawer, CartItem, CartSummary
│   └── ui/                       # Button, Badge, shared primitives
├── lib/
│   ├── stripe.ts                 # Stripe client initialization (server-side only)
│   ├── products.ts               # Product data access layer
│   └── cart-context.tsx          # Cart state provider
├── data/
│   └── products.json             # Static product catalog
└── styles/
    └── globals.css               # Tailwind config, dark theme tokens
```

**Key principle:** Every file has one job. Every directory groups by domain, not by file type. If someone asks "where's the cart logic?" the answer is obvious. If someone asks "where are the API routes?" the answer is obvious. No treasure hunts.

### 4d. Dark Theme — Do It Once, Do It Right

Don't bolt on dark mode as a toggle. You said dark-themed. That means dark-first, dark-only (for now).

- Set `darkMode: 'class'` in Tailwind config
- Apply the `dark` class at the root layout level
- Use CSS custom properties for your brand palette so you change colors in one place
- Background: `zinc-950` or `slate-950` (true dark, not gray)
- Text: `zinc-100` for primary, `zinc-400` for secondary
- Accent: One brand color for CTAs and interactive elements
- Card surfaces: `zinc-900` with subtle border (`zinc-800`)

The product images are the color. Everything else is a stage for them.

---

## 5. DEPLOYMENT ARCHITECTURE

### Vercel (Primary)

```
GitHub Push → Vercel Build → Edge Network
                                  │
                    ┌──────────────┤
                    │              │
              Static Assets    Serverless Functions
              (products,       (/api/checkout,
               images,          /api/webhooks/stripe)
               pages)
```

- Set environment variables in Vercel dashboard: `STRIPE_SECRET_KEY`, `STRIPE_WEBHOOK_SECRET`, `NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY`
- Enable Vercel's built-in analytics for basic monitoring
- Use preview deployments for every PR — free smoke testing

### SST (Future / AWS)

Only move to SST when you need:
- Background job processing (order fulfillment automation)
- Custom domain routing beyond what Vercel offers
- AWS-native services (SQS for order queues, SES for transactional email)
- Cost optimization at scale (Vercel gets expensive with high traffic)

SST uses the `NextjsSite` construct. The migration is mostly config, not code — your Next.js app stays the same.

---

## 6. FAILURE MODES — HOW THIS BREAKS AT 3 AM

| Failure | Impact | Detection | Mitigation |
|---------|--------|-----------|------------|
| Stripe webhook delivery fails | Orders paid but not acknowledged | Stripe dashboard retry alerts | Implement webhook retry handling; poll Stripe events API as backup |
| Vercel function cold starts | Slow checkout initiation (~1-2s) | User-perceived latency | Acceptable for swag volume; SST with provisioned concurrency if needed |
| Product images CDN failure | Broken product grid | Visual monitoring or synthetic checks | Host images on Vercel's CDN (built-in), add blur placeholder fallbacks |
| Stripe API outage | Cannot create checkout sessions | Stripe status page monitoring | Display "checkout temporarily unavailable" — don't queue failed attempts |
| Product data stale (JSON updated but not deployed) | Wrong prices or sold-out items shown | Manual process discipline | Move to CMS or DB when this becomes painful |

**The most likely 3 AM failure:** A webhook event is lost or mishandled, and an order is paid but you don't know about it. Solution: Check your Stripe dashboard daily until you have automated monitoring. Add a daily reconciliation step — compare Stripe payments against your order log.

---

## 7. IMPLEMENTATION PLAN (VTS FORMAT)

### TASK 1: Project Scaffolding and Dark Theme Foundation

**Description:** Initialize Next.js project with App Router, Tailwind CSS, and dark-first theme configuration. Establish project structure, font loading, and base layout components (Header, Footer, root layout).
**Acceptance Criteria:**
- Next.js 14+ project with App Router initialized
- Tailwind CSS configured with dark-first custom palette (zinc-950 base)
- Root layout renders with dark background, correct fonts, Header, and Footer
- Project directory structure matches the architecture spec
- Runs locally with `npm run dev`
**Complexity:** M
**Dependencies:** None
**Files:** `package.json`, `tailwind.config.ts`, `src/app/layout.tsx`, `src/app/globals.css`, `src/components/layout/Header.tsx`, `src/components/layout/Footer.tsx`

### TASK 2: Product Data Layer and Catalog

**Description:** Create the static product data source (`products.json`) with a typed access layer. Build the product catalog page with a responsive grid of product cards. Each card shows image, name, price, and available variants.
**Acceptance Criteria:**
- `products.json` with at least 6 sample products (name, price, image, description, variants)
- TypeScript types for Product, ProductVariant
- `lib/products.ts` exports functions: `getAllProducts()`, `getProductById()`
- Product catalog page renders responsive grid (1-col mobile, 2-col tablet, 3-4 col desktop)
- ProductCard component displays product image, name, price, and a CTA
**Complexity:** M
**Dependencies:** Task 1
**Files:** `src/data/products.json`, `src/lib/products.ts`, `src/app/products/page.tsx`, `src/components/product/ProductCard.tsx`, `src/components/product/ProductGrid.tsx`

### TASK 3: Cart State Management and UI

**Description:** Implement cart context with `localStorage` persistence, cart drawer/page UI, and add-to-cart interactions from product cards.
**Acceptance Criteria:**
- `CartProvider` context with `items`, `addItem`, `removeItem`, `updateQuantity`, `clearCart`, `cartTotal`
- Cart state persists across page refreshes via `localStorage`
- Cart page/drawer shows line items with quantity controls and total
- Product cards have functional "Add to Cart" button
- Cart badge in header shows item count
**Complexity:** M
**Dependencies:** Task 1, Task 2
**Files:** `src/lib/cart-context.tsx`, `src/app/cart/page.tsx`, `src/components/cart/CartItem.tsx`, `src/components/cart/CartSummary.tsx`, `src/components/product/ProductCard.tsx` (update)

### TASK 4: Stripe Checkout Integration

**Description:** Create the server-side API route that builds a Stripe Checkout session from cart contents, and the client-side redirect flow. Include success and cancellation pages.
**Acceptance Criteria:**
- `lib/stripe.ts` initializes Stripe client (server-side only, never exposed to browser)
- `POST /api/checkout` accepts cart items, creates Stripe Checkout Session with line items, returns session URL
- Client redirects to Stripe Checkout on "Proceed to Checkout"
- Success page (`/checkout/success`) displays order confirmation with session ID
- Cancel page (`/checkout/cancel`) returns user to cart
- Environment variables documented: `STRIPE_SECRET_KEY`, `NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY`
**Complexity:** L
**Dependencies:** Task 3
**Files:** `src/lib/stripe.ts`, `src/app/api/checkout/route.ts`, `src/app/checkout/success/page.tsx`, `src/app/checkout/cancel/page.tsx`

### TASK 5: Stripe Webhook Handler

**Description:** Implement the webhook endpoint that receives Stripe events, verifies signatures, handles `checkout.session.completed`, and logs order data. Include idempotency guard.
**Acceptance Criteria:**
- `POST /api/webhooks/stripe` verifies Stripe signature on every request
- Unknown event types return 200 without processing
- `checkout.session.completed` events are logged with order details
- Duplicate event IDs are detected and skipped (idempotent)
- Failed verification returns 400
- `STRIPE_WEBHOOK_SECRET` environment variable documented
**Complexity:** M
**Dependencies:** Task 4
**Files:** `src/app/api/webhooks/stripe/route.ts`

### TASK 6: Landing Page and Visual Polish

**Description:** Build the hero/landing page with brand messaging, featured products, and a CTA to the catalog. Apply final visual polish — spacing, transitions, hover states, responsive breakpoints, image optimization.
**Acceptance Criteria:**
- Landing page with hero section, value proposition, and featured product highlights
- Smooth transitions on hover states and page navigation
- All images use Next.js `Image` component with proper sizing and blur placeholders
- Responsive across mobile, tablet, and desktop
- Lighthouse accessibility score > 90
**Complexity:** M
**Dependencies:** Task 2
**Files:** `src/app/page.tsx`, `src/components/product/FeaturedProducts.tsx`, various component style updates

### TASK 7: Deployment Configuration and Smoke Testing

**Description:** Configure Vercel deployment (environment variables, build settings). Document SST deployment path. Write a manual smoke test checklist for end-to-end checkout flow.
**Acceptance Criteria:**
- Vercel project configured and deploys on push to `main`
- Environment variables set in Vercel dashboard
- Stripe webhook endpoint registered in Stripe dashboard pointing to production URL
- SST deployment documented (config file and commands, not implemented)
- Smoke test checklist covers: browse products → add to cart → checkout → Stripe payment → success page → webhook received
**Complexity:** M
**Dependencies:** Task 5, Task 6
**Files:** `vercel.json` (if needed), `README.md` (deployment section), `sst.config.ts` (documented, not necessarily functional)

---

## 8. WHAT I'D PUSH BACK ON

**The prior plan's 8-day timeline.** Not because it's wrong, but because it's misleading. Day 1-2 being "finalize SKUs and photography" is a business task that could take a week or a month depending on your supply chain. The *code* is a 3-5 day build for an experienced developer or an AI-assisted workflow. The *business* prep is unbounded. Separate them.

**The Phase 2 "add a database" guidance.** I'd be more specific: you add a database when you need *order history visible to the admin*, not before. Stripe is your order database until then. Don't introduce a Postgres or DynamoDB instance for 50 SKUs and 10 orders a week.

**Missing from the prior plan: error boundaries.** Every page that talks to Stripe or loads dynamic data needs a React Error Boundary or Next.js `error.tsx`. When Stripe is down, the user sees a helpful message, not a white screen of death.

---

## 9. BUILD VS BUY — MY FINAL WORD

I agree with the prior plan's conclusion, but I'll frame it as an architect:

**You are not building e-commerce infrastructure. You are building a branded front door.**

The moment you frame it as "building a store," you're tempted to build inventory management, order dashboards, fulfillment tracking, and refund workflows. That's six months of work and a full-time operations burden.

Instead: **Build the thing customers see. Buy the thing that processes money. Use a spreadsheet for everything else until the spreadsheet breaks.**

Shopify makes sense when you need managed fulfillment, multi-channel selling, or when your operations team (not your engineering team) needs to manage inventory. For a swag store run by an AI-enabled engineering team, Stripe Checkout + a Next.js frontend is the right call.

Revisit in 90 days. If you're processing > 100 orders/week and spending more than 2 hours/week on manual fulfillment, then evaluate Shopify or Medusa. Not before.

---

## 10. MONITORING FROM DAY ONE

Because observability is not a Phase 2 feature:

- **Vercel Analytics** — free, built-in, covers page performance
- **Stripe Dashboard** — your order monitoring for Phase 1
- **A simple health check** — `/api/health` route that returns 200. Point an uptime monitor at it (UptimeRobot, free tier).
- **Structured logging in webhook handler** — `console.log(JSON.stringify({ event: 'checkout.completed', sessionId, amount, timestamp }))`. Vercel captures these in its log viewer.

That's enough to sleep through the night.

---

## Summary

The prior plan is sound. I've sharpened the component boundaries, specified the cart architecture, detailed the webhook handler contract, laid out the file structure, and added failure mode analysis. The seven tasks above are sequenced with clear dependencies and can be executed in roughly two parallel tracks (UI: Tasks 1→2→3→6 and Backend: Tasks 4→5→7, merging at Task 7).

**Build the brand. Buy the plumbing. Ship in a week. Revisit in 90 days.**

---

Why did the microservice break up with the monolith? Because it needed more space — but then it spent all its time on networking.

...I'll see myself out.

-- Architect Vern (measure twice, deploy once)
