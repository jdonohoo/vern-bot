---
id: VTS-010
title: "Application Initialization and Wiring"
complexity: M
status: pending
owner: ""
source: oracle
source_ref: "oracle-architect-breakdown.md"
dependencies:
  - VTS-001
  - VTS-002
  - VTS-003
  - VTS-004
files:
  - "`scripts/main.js`"
---

# Application Initialization and Wiring

Wire everything together. The main entry point initializes the migration runner, storage service, state manager, theme system, kanban board, and event listeners. This is the orchestration layer — it should read like a table of contents for the application. The init sequence now includes migration check before state hydration, `chrome.storage.onChanged` listener setup, and execution context detection. Keep it linear and obvious.

## Criteria

- `main.js` initializes in clear sequential order: storage → **migration** → state → theme → board → event bindings
- Error during init shows user-friendly message (not a blank popup)
- All components properly wired to state manager events
- No global mutable state outside of state manager
- Console logging at key init steps (for debugging in production)
- Graceful handling if `chrome.storage.local` is unavailable (e.g., development mode with fallback to `localStorage`)
- `chrome.storage.onChanged` listener registered for cross-view sync (popup ↔ tab)
- Detect execution context (popup vs. tab) and apply appropriate layout class
