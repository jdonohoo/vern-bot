---
id: VTS-003
title: "Configure Deployment for Vercel and SST"
complexity: M
status: pending
owner: ""
source: discovery
source_ref: "07-architect-architect-breakdown.md"
dependencies:
  - VTS-001
files:
  - "`vercel.json`"
  - "`sst.config.ts`"
  - "`.env.example`"
  - "`.gitignore`"
  - "`next.config.ts`"
---

# Configure Deployment for Vercel and SST

Set up dual deployment paths. Vercel is the primary target (fastest, simplest). SST configuration is the secondary target for when you need AWS-native capabilities. Both must work from the same codebase without forking. Environment variables must be managed cleanly — no secrets in code, ever. How will this fail at 3 AM? If your env vars are missing or misconfigured and nobody notices until a customer can't check out.

## Criteria

- `vercel.json` configured (if needed beyond defaults) with proper build settings
- SST config file (`sst.config.ts`) defining the Next.js site resource with proper permissions
- Environment variable schema documented: `STRIPE_SECRET_KEY`, `STRIPE_PUBLISHABLE_KEY`, `STRIPE_WEBHOOK_SECRET`, `NEXT_PUBLIC_BASE_URL`
- `.env.example` file with all required variables listed (no real values)
- `.env.local` in `.gitignore`
- Build succeeds in both Vercel and SST contexts
- README section on deployment for both paths
