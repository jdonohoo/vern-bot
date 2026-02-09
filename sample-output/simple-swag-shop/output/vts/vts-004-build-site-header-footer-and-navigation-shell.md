---
id: VTS-004
title: "Build Site Header, Footer, and Navigation Shell"
complexity: M
status: pending
owner: ""
source: discovery
source_ref: "07-architect-architect-breakdown.md"
dependencies:
  - VTS-001
files:
  - "`components/layout/Header.tsx`"
  - "`components/layout/Footer.tsx`"
  - "`components/layout/MobileNav.tsx`"
  - "`app/layout.tsx`"
---

# Build Site Header, Footer, and Navigation Shell

Implement the persistent layout shell — header with logo/brand, navigation links, and a cart indicator; footer with links and legal copy. This is the frame that every page lives inside. Dark theme, minimal, brand-forward. The cart indicator should show item count from client-side state (Task 7). Navigation should be responsive — mobile hamburger menu that doesn't require a PhD to use.

## Criteria

- `Header` component with: logo/brand name, nav links (Shop, About), cart icon with item count badge
- `Footer` component with: copyright, links (Privacy, Terms, Contact), brand tagline
- Responsive: mobile menu (hamburger → slide or dropdown), desktop horizontal nav
- Dark theme styling: header/footer use `zinc-950` or `zinc-900`, high-contrast text
- Cart icon is a link/button that navigates to `/cart`
- Layout is accessible: semantic HTML (`<header>`, `<nav>`, `<main>`, `<footer>`), keyboard navigable
- No JavaScript required for initial paint (SSR-friendly)
