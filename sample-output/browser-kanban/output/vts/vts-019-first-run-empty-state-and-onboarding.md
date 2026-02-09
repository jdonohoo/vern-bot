---
id: VTS-019
title: "First-Run Empty State and Onboarding"
complexity: S
status: pending
owner: ""
source: oracle
source_ref: "oracle-architect-breakdown.md"
dependencies:
  - VTS-006
  - VTS-008
  - VTS-005
files:
  - "`scripts/components/empty-state.js`"
  - "`styles/empty-state.css`"
---

# First-Run Empty State and Onboarding

Design the empty-state experience for a brand-new user who opens the extension for the first time. Without this, a new user opens the popup and sees four empty columns with "No tasks" repeated four times. That's not welcoming — that's a loading screen with extra steps. One message, one button, done.

## Criteria

- First-run detection: no tasks exist in storage
- Empty board shows a centered welcome message with brief explanation
- Clear CTA button: "Create your first task" that opens the quick-add form in the "New" column
- Welcome state disappears permanently after first task is created
- Empty individual columns still show "No tasks" after the first task exists (per Task 6)
- Welcome message adapts to current theme (light/dark)
- No tutorial, no multi-step wizard — one message, one button, done
