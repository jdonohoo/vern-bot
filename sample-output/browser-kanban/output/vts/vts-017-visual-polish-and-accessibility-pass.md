---
id: VTS-017
title: "Visual Polish and Accessibility Pass"
complexity: M
status: pending
owner: ""
source: oracle
source_ref: "oracle-architect-breakdown.md"
dependencies:
  - VTS-010
files:
  - "All CSS files"
  - "component files as needed for ARIA"
---

# Visual Polish and Accessibility Pass

Final visual and accessibility review. Ensure consistent spacing, typography, and color usage. Verify keyboard navigation works throughout. Check color contrast ratios for both themes. Add ARIA labels where needed. Must run after core integration is complete — you need something to review.

## Criteria

- Consistent spacing scale (4px/8px/12px/16px/24px)
- Typography hierarchy: headings, body, captions
- All interactive elements keyboard-accessible (tab order, focus rings)
- ARIA labels on icon buttons and dynamic content
- Color contrast meets WCAG AA in both themes
- Urgency colors distinguishable for common color blindness (verify secondary indicators from Task 7 are effective)
- Smooth transitions on panel open/close, theme switch, card move
- No layout shifts on interaction
- Tested at popup min/max dimensions
