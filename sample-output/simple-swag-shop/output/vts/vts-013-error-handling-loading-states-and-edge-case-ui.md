---
id: VTS-013
title: "Error Handling, Loading States, and Edge Case UI"
complexity: M
status: pending
owner: ""
source: discovery
source_ref: "07-architect-architect-breakdown.md"
dependencies:
  - VTS-005
  - VTS-006
  - VTS-009
files:
  - "`app/error.tsx`"
  - "`app/not-found.tsx`"
  - "`app/loading.tsx`"
  - "`app/products/[slug]/loading.tsx`"
  - "`components/ui/Skeleton.tsx`"
  - "`components/ui/ErrorMessage.tsx`"
---

# Error Handling, Loading States, and Edge Case UI

The unglamorous work that separates a prototype from a product. Global error boundary, loading skeletons, not-found page, and consistent error states throughout the app. Users should never see a blank white screen, an unhandled exception, or a cryptic error message. Especially on a dark theme — a flash of white from an error page is jarring and amateurish.

## Criteria

- Global `error.tsx` boundary in `/app` with dark-themed error UI and retry button
- Global `not-found.tsx` with dark-themed 404 page and link back to shop
- `loading.tsx` in key routes (`/`, `/products/[slug]`) with skeleton loaders
- Skeleton loaders match the layout of the content they replace (not generic spinners)
- API error responses follow a consistent shape: `{ error: string, details?: string }`
- Buy Now button disabled with clear messaging when item is unavailable
- Network error handling on client-side fetch calls (checkout creation, session retrieval)
