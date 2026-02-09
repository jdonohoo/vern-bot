---
id: VTS-012
title: "Build Static Pages (About, Privacy, Terms)"
complexity: S
status: pending
owner: ""
source: discovery
source_ref: "07-architect-architect-breakdown.md"
dependencies:
  - VTS-004
files:
  - "`app/about/page.tsx`"
  - "`app/privacy/page.tsx`"
  - "`app/terms/page.tsx`"
---

# Build Static Pages (About, Privacy, Terms)

Every store needs these. About page tells your brand story. Privacy and Terms pages are legally necessary once you're processing payments. These are static pages — no dynamic data, no API calls. Keep them simple, keep them dark-themed, keep them findable from the footer.

## Criteria

- `/about` page with placeholder brand story content
- `/privacy` page with placeholder privacy policy
- `/terms` page with placeholder terms of service
- All pages use consistent dark theme styling
- All pages linked from Footer component (Task 4)
- Proper metadata for each page
- Content is easy to update (just edit the page file or a markdown source)
