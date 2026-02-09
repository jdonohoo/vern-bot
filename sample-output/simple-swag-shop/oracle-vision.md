

# Oracle Vision

*The council has spoken. Now let me tell you what they actually said.*

---

## Summary

The architect's 15-task breakdown is structurally sound — the bones are good. But the VernHole synthesis reveals three critical gaps and several misalignments between what the council agreed on and what the tasks actually build. The biggest finding: **the task plan builds a cart-first storefront, but the council unanimously recommended direct-to-checkout as the MVP path.** That's not a minor detail — it changes the critical path, reduces complexity, and accelerates launch by 2-3 tasks.

Additionally, the council surfaced concerns about product photography on dark backgrounds, mobile-first design, monitoring/observability, and the Stripe Products-as-catalog pattern — none of which appear in the current task breakdown. The tasks also over-index on "Phase 2" work (webhooks, error handling polish) while under-indexing on "Phase 1" essentials (notifications, fulfillment workflow, trust signals).

**Key changes:**
- **3 new tasks** added (monitoring, trust signals/notifications, product photography guidelines)
- **5 tasks modified** (scope adjustments, acceptance criteria gaps, complexity reassessments)
- **0 tasks removed** (every task earns its place, but 2 are re-sequenced as optional/Phase 2)
- **Dependency chain restructured** to reflect direct-to-checkout-first approach
- **VTS-007 and VTS-008 demoted** from critical path to "add when customers ask for multi-item purchasing"

---

## New Tasks

### TASK 16: Implement Health Check Endpoint and Basic Monitoring

**Description:**
Add a `/api/health` endpoint that verifies connectivity to Stripe and returns service status. Configure Vercel Analytics for traffic monitoring. This is the "3 AM insurance" that Architect Vern and Retro Vern both circled: when Stripe has an incident or your deploy breaks, you need to know before your customers tell you. UptimeRobot (free tier) or Vercel's built-in monitoring should ping this endpoint. Without this, your first indication of downtime is an angry email.

**Acceptance Criteria:**
- API route at `app/api/health/route.ts` (GET)
- Returns JSON: `{ status: "ok" | "degraded", stripe: "connected" | "error", timestamp: ISO string }`
- Stripe check: lightweight API call (e.g., `stripe.products.list({ limit: 1 })`) to verify key validity
- Response time under 500ms
- Vercel Analytics enabled in `next.config.ts`
- Documentation: recommended UptimeRobot or equivalent setup for the health endpoint

**Complexity:** S
**Dependencies:** VTS-001, VTS-009 (needs Stripe SDK initialized)
**Files:** `app/api/health/route.ts`, `next.config.ts`

---

### TASK 17: Purchase Notification Pipeline and Fulfillment Trigger

**Description:**
When a purchase completes (via webhook), send a notification to the team. The council unanimously agreed: Stripe Dashboard is your admin panel for Phase 1, but you still need a push notification — nobody checks dashboards proactively. This is the bridge between "payment received" and "swag shipped." Start with email or Slack webhook; the mechanism matters less than the existence of the alert. Every Vern who discussed operations flagged this gap implicitly — the webhook handler (VTS-011) logs events, but logging without alerting is just writing a diary nobody reads.

**Acceptance Criteria:**
- On `checkout.session.completed`, send notification with: customer email, items purchased, total amount, session ID
- Notification channel: Slack webhook (preferred) or email via Resend/SendGrid
- Notification includes a direct link to the Stripe Dashboard session
- Failure to notify does NOT block the webhook response (fire-and-forget with error logging)
- Environment variable for notification target (`SLACK_WEBHOOK_URL` or `NOTIFICATION_EMAIL`)
- `.env.example` updated with new variable

**Complexity:** S
**Dependencies:** VTS-011
**Files:** `lib/notifications.ts`, `app/api/webhooks/stripe/route.ts` (modification), `.env.example`

---

### TASK 18: Product Photography Guidelines and Image Asset Pipeline

**Description:**
UX Vern's most urgent flag: dark hoodies disappear on dark backgrounds. This is not a code task — it's a design specification task that the codebase must support. Define the image requirements (dimensions, aspect ratio, background treatment, file format), create placeholder images that demonstrate the correct approach, and document the photography guidelines. The council agreed 60-80% of effort should go to product quality. This task ensures the codebase supports that investment rather than fighting it.

**Acceptance Criteria:**
- Written specification: image dimensions (e.g., 1200x1200 minimum), aspect ratio (1:1 for grid, 4:3 for detail), file format (WebP primary, PNG fallback)
- Dark theme product image guidelines: slight gradient background (`zinc-800` to `zinc-850`), subtle shadow/glow to separate product from page background, avoid pure-black product shots
- Placeholder images updated to demonstrate correct treatment (not generic colored rectangles)
- `next.config.ts` image optimization settings configured for the specified dimensions
- Component-level `sizes` attributes in `ProductCard` and `ProductDetail` match the spec

**Complexity:** S
**Dependencies:** VTS-005, VTS-006
**Files:** `public/images/` (placeholder assets), `next.config.ts`, documentation in code comments or a `PRODUCT_IMAGES.md`

---

## Modified Tasks

### VTS-002: Define Product Data Schema and Stripe Products Catalog (was: Define Product Data Schema and Static Product Catalog)

**Changes:** The council (Retro, NyQuil, Startup) converged on using **Stripe Products as the source of truth** rather than a local JSON file. This eliminates price sync bugs entirely — the price charged is always the price displayed because they come from the same source. The static JSON approach creates a data sync liability from Day 1. Modified to use Stripe Products API with local cache/fallback, not pure static data.

**Description:**
Design the product data model backed by Stripe Products. Stripe is the single source of truth for products and prices. Local TypeScript interfaces define the shape, but `getAllProducts()` reads from Stripe's API (cached aggressively since products change rarely). A static fallback file exists for development and build-time static generation, but production reads from Stripe. This means adding a product is done in the Stripe Dashboard — no code deploy needed.

**Acceptance Criteria:**
- TypeScript interfaces `Product` and `ProductVariant` retained as specified
- `lib/products.ts` fetches from Stripe Products API (`stripe.products.list()` with `expand: ['data.default_price']`)
- Response mapped to local `Product` interface
- Cache layer: in-memory cache with 5-minute TTL (or ISR revalidation) — products don't change mid-session
- Static fallback: `data/products-fallback.ts` used during `npm run build` if Stripe is unavailable (for static generation)
- `getProductBySlug()` uses Stripe product metadata field `slug` for URL-friendly lookups
- All helper functions retain the same signatures — consumers don't know the data source changed
- Development mode: works with Stripe test mode products

**Complexity:** M (was S — Stripe API integration adds real complexity)
**Dependencies:** VTS-001

---

### VTS-007: Implement Client-Side Cart State Management (was: same title)

**Changes:** Demoted from critical path to Phase 2. The council (YOLO, Startup, Mediocre) agreed: for < 15 SKUs, direct-to-checkout per item is the MVP. Cart adds complexity with marginal benefit until customers demonstrate multi-item purchasing behavior. The task remains valid but is no longer a blocker for checkout flow. VTS-009 is modified to support both direct-to-checkout (single item) and cart-based checkout.

**Complexity:** M (unchanged, but now Phase 2)
**Dependencies:** VTS-001 (unchanged)
**Additional Acceptance Criteria:**
- Cart provider must gracefully handle the case where it's not yet integrated — components should work with and without cart context
- Add a feature flag or conditional: if cart is disabled, "Add to Cart" becomes "Buy Now" and routes directly to checkout

---

### VTS-008: Build Cart Page (was: same title)

**Changes:** Demoted to Phase 2 alongside VTS-007. Not needed for MVP launch if using direct-to-checkout. When implemented, should include the slide-out cart panel that UX Vern recommended for mobile rather than a dedicated page route only.

**Complexity:** M (unchanged, but now Phase 2)
**Dependencies:** VTS-007
**Additional Acceptance Criteria:**
- Consider slide-out cart panel (drawer) as primary interaction on mobile, with `/cart` page as fallback/desktop view
- "Buy Now" direct-to-checkout option remains available even after cart is implemented (for single-item impulse purchases)

---

### VTS-009: Implement Stripe Checkout Session API Route (was: same title)

**Changes:** Must support two modes: direct-to-checkout (single item with variant, no cart required) and cart-based checkout (when VTS-007/008 are implemented). The current spec assumes cart-only. Additionally, the council flagged that Stripe Tax should be considered even in v1 — it's a single config line and avoids the post-Wayfair landmine Academic Vern raised.

**Additional Acceptance Criteria:**
- Accepts EITHER a single item (`{ productId, variantId, quantity: 1 }`) OR a cart items array
- When in direct-to-checkout mode: product detail page sends single item, no cart dependency
- Stripe Tax: `automatic_tax: { enabled: true }` on session creation (requires Stripe Tax activation in dashboard — document this)
- `shipping_address_collection` enabled with allowed countries configured via env var
- `STRIPE_TAX_ENABLED` env var to toggle tax (allows launch without Stripe Tax activated)
- `.env.example` updated with new variables

**Complexity:** L (unchanged — the dual-mode adds scope but the Stripe SDK makes it straightforward)
**Dependencies:** VTS-002 (removed VTS-007 as hard dependency; cart is optional)

---

### VTS-006: Build Product Detail Page (was: same title)

**Changes:** Must include a "Buy Now" direct-to-checkout button in addition to (or instead of) "Add to Cart." This is the primary purchase path in Phase 1. UX Vern's mobile-first mandate also applies here — the variant selector and buy button must be thumb-reachable on mobile. Added acceptance criteria for mobile layout.

**Additional Acceptance Criteria:**
- "Buy Now" button that creates a Stripe Checkout session for the single selected item+variant (calls VTS-009 API directly)
- "Add to Cart" button conditionally rendered only when cart system (VTS-007) is active
- Mobile layout: product image full-width, sticky bottom bar with price + "Buy Now" button
- Variant selector uses large tap targets (minimum 44x44px per WCAG)
- If product has only one variant (or no variants), skip the selector entirely — don't show a dropdown with one option

---

## Removed Tasks

No tasks are removed. Every task in the current breakdown earns its place in the final product. However, VTS-007 and VTS-008 are **re-sequenced as Phase 2 / non-blocking** — they should be completed, but not before the core checkout flow is live.

---

## Dependency Changes

### Critical Path (Phase 1 — Ship the MVP):
```
VTS-001 (Foundation)
  ├── VTS-002 (Product Schema + Stripe Products) [was S, now M]
  │     ├── VTS-005 (Product Grid)
  │     │     └── VTS-006 (Product Detail + Buy Now) [modified]
  │     │           └── VTS-009 (Checkout API — direct-to-checkout mode) [modified]
  │     │                 └── VTS-010 (Success/Cancel Pages)
  │     │                 └── VTS-011 (Webhook Handler)
  │     │                       └── TASK 17 (Notifications) [NEW]
  │     └── TASK 18 (Image Guidelines) [NEW]
  ├── VTS-003 (Deployment)
  ├── VTS-004 (Header/Footer/Nav)
  │     └── VTS-012 (Static Pages)
  └── TASK 16 (Health Check + Monitoring) [NEW]
```

### Phase 2 (Post-Launch, When Demand Validates):
```
VTS-007 (Cart State) ← VTS-001
  └── VTS-008 (Cart Page) ← VTS-007
        └── VTS-009 (Checkout API — cart mode already supported)
```

### Polish & Launch (Parallel to Phase 1, Required Before Go-Live):
```
VTS-013 (Error Handling) ← VTS-005, VTS-006, VTS-009
VTS-014 (SEO/OG/Perf) ← VTS-005, VTS-006
VTS-015 (Smoke Test) ← VTS-009, VTS-010, VTS-011, VTS-013
```

### Key Dependency Changes:
- **VTS-009 no longer depends on VTS-007.** Checkout works with direct-to-checkout (single item). Cart is additive, not prerequisite.
- **VTS-006 now directly connects to VTS-009** via "Buy Now" — this is the MVP purchase path.
- **VTS-015 (Smoke Test) should also depend on TASK 16 and TASK 17** — can't verify launch readiness without monitoring and notifications.
- **TASK 18 (Image Guidelines) depends on VTS-005 and VTS-006** but can be done in parallel as a design spec task.

---

## Risk Assessment

### Risks Addressed by These Changes:
1. **Cart-before-checkout bottleneck** — Eliminated. Direct-to-checkout removes 2 tasks from the critical path.
2. **Price sync bugs** — Eliminated. Stripe Products as source of truth means prices can't drift.
3. **Silent order failures** — Addressed. TASK 17 ensures someone knows when money changes hands.
4. **Dark-on-dark product photography** — Addressed. TASK 18 provides guidelines before anyone uploads images.
5. **"Is it up?" uncertainty** — Addressed. TASK 16 gives you a heartbeat.

### Remaining Risks (The Oracle Sees These Shadows):

1. **Stripe Test Mode → Live Mode Transition.** The plan builds entirely in test mode. The moment you flip to live, you'll discover: webhook endpoints need re-registering, API keys need rotating, and Stripe Tax requires a real tax registration. **Mitigation:** VTS-015 smoke test should include a "live mode preflight" section documenting every Stripe Dashboard toggle.

2. **No Order Persistence.** The council agreed "no database for Phase 1" — but this means if a webhook fails and isn't retried, you lose the order record. Stripe retains it, but you have no local audit trail. **Mitigation:** At minimum, write webhook events to a log file or structured logging service. The gap between "we logged it" and "we have a database" is where orders disappear.

3. **Fulfillment Is Entirely Manual.** No task addresses how orders get from "payment received" to "swag shipped." TASK 17 notifies you, but the fulfillment workflow itself is undefined. For Phase 1 this is acceptable — but the gap between "someone got a Slack message" and "someone packed a box and called UPS" is where customer satisfaction dies. **Recommendation:** Add a fulfillment checklist to VTS-015's launch checklist.

4. **International Sales Tax.** Academic Vern flagged post-*Wayfair* complexity. The modified VTS-009 adds Stripe Tax as an option, but if you sell internationally without a Merchant of Record, you're potentially liable for VAT/GST in every country you ship to. **Mitigation:** For launch, restrict `shipping_address_collection` to US-only. Document the Paddle/Lemon Squeezy evaluation trigger: first international order request.

5. **Mobile UX Is Specified But Not Verified.** The modifications add mobile-first acceptance criteria, but no task explicitly includes mobile testing across devices. **Mitigation:** Add mobile device testing (iOS Safari, Android Chrome at minimum) to VTS-015 smoke test plan.

6. **The Swag Itself.** Every single Vern said the product matters more than the tech. No task in this breakdown addresses sourcing, quality control, or sample ordering. This is the biggest risk of all — you could build a flawless storefront selling mediocre hoodies. **Mitigation:** This is out of scope for the VTS (it's not a code task), but it should be a parallel workstream with higher priority than anything after VTS-009.

---

*I've seen this pattern a thousand times. The plan focuses on the code because the code is controllable. The product, the photography, the fulfillment — those are messy, human, uncertain. But they're where the margin lives. The council told you this in eight different voices. I'm telling you in one: ship the storefront fast, then spend your time on the swag.*

---

Why did the Oracle refuse to review the sprint backlog? Because every task said "estimated: 1 day" and the Oracle knows the only thing that ships in one day is the decision to add another planning meeting. ...The prophecy has been spoken.

-- Oracle Vern *(the signal was there all along)*
