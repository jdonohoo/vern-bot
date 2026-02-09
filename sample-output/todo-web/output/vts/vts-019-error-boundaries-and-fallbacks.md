---
id: VTS-019
title: "Error Boundaries and Fallbacks"
complexity: S
status: pending
owner: ""
source: discovery
source_ref: "07-architect-architect-breakdown.md"
dependencies:
  - VTS-004
files:
  - "src/components/shared/ErrorBoundary.tsx"
  - "src/components/shared/ErrorFallback.tsx"
---

# Error Boundaries and Fallbacks

Wrap major UI sections in React error boundaries with helpful fallback UIs. Handle common failure modes: localStorage quota exceeded, lazy-load failures, corrupted task data, and missing tag references. No uncaught error should crash the entire app. Each major section can fail independently.

A localStorage-only app with no error handling is a data loss machine. This task exists because the original plan had zero error handling tasks.

## Error Boundaries to Implement

1. **App-level boundary:** Catches anything that escapes section-level boundaries. Shows a "Something went wrong" message with a reload button.
2. **Board view boundary:** If the board crashes, the list view still works. Fallback: "Board view encountered an error. Try switching to list view."
3. **List view boundary:** If the list crashes, the board still works. Fallback: same pattern.
4. **Detail panel boundary:** If the detail panel crashes, the main views still work. Fallback: "Unable to load task details. Close and try again."
5. **Markdown preview boundary:** If react-markdown fails to load or render, show the raw markdown text instead of crashing.

## Specific Failure Modes

- **localStorage quota exceeded:** Detect quota errors on write, show a warning toast ("Storage is full. Export your data to free up space."), do not crash
- **Lazy-load failures:** If React.lazy chunks fail to load (network error), show a retry button in the fallback UI
- **Corrupted task data:** If a task references a deleted tag ID, render the task without that tag (do not crash). Log a warning.
- **Missing column reference:** If a task's `statusId` points to a non-existent column, assign it to the first column and log a warning

## Criteria

- No uncaught error crashes the entire app
- Each major section (board, list, detail, markdown) can fail independently
- localStorage quota exceeded shows a warning, not a white screen
- Lazy-load failures show a retry button
- Corrupted data references (deleted tags, missing columns) are handled gracefully with warnings, not crashes
- Error fallback UIs are styled with theme tokens and look intentional, not broken
- Error boundaries log errors to console for debugging
