---
id: VTS-016
title: "Error Handling and Edge Case Hardening"
complexity: M
status: pending
owner: ""
source: oracle
source_ref: "oracle-architect-breakdown.md"
dependencies:
  - VTS-010
files:
  - "All component files"
  - "`scripts/utils/error-handler.js`"
---

# Error Handling and Edge Case Hardening

Systematic pass through all components to ensure errors are handled explicitly, edge cases are covered, and the user never sees a blank screen or silent failure. This is the task that separates "demo" from "daily driver." Must run after core integration is complete — you can't harden error handling on code that doesn't exist yet.

## Criteria

- Storage read/write failures show user-visible error messages
- Empty states handled in all views (no tasks, no due date, no description)
- Long titles truncated with ellipsis, not layout-breaking
- Zero and negative estimates handled gracefully
- Due date in past doesn't crash urgency calc
- Timezone changes handled (urgency recomputed on render, not cached)
- Storage quota approaching limit shows warning (warn if usage exceeds 80% of `chrome.storage.local` quota)
- All forms validate inputs before submission
- No unhandled promise rejections
- Console errors in production: zero tolerance target
- Backup integrity check: verify `_backups` array isn't corrupted on load
- Graceful degradation if `crypto.randomUUID()` is unavailable (fallback to timestamp-based ID)
- Rate-limit protection on debounced writes (if state changes faster than debounce clears, coalesce — don't queue)
