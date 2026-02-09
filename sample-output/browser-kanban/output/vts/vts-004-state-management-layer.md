---
id: VTS-004
title: "State Management Layer"
complexity: M
status: pending
owner: ""
source: oracle
source_ref: "oracle-architect-breakdown.md"
dependencies:
  - VTS-002
  - VTS-003
files:
  - "`scripts/state/state-manager.js`"
  - "`scripts/state/event-bus.js`"
---

# State Management Layer

Build a lightweight in-memory state manager that sits between the UI and storage. On load, it hydrates from storage (after migration). On mutations, it updates both in-memory state and persists to storage. It exposes a simple subscribe/notify pattern so UI components can react to state changes without tight coupling. Now includes `chrome.storage.onChanged` listener for cross-view sync — the state manager must handle external changes, not just internal ones.

## Criteria

- `StateManager` loads state from `StorageService` on init
- Exposes methods: `getTasks()`, `getTasksByStatus(status)`, `getTask(id)`, `addTask(task)`, `updateTask(id, changes)`, `deleteTask(id)`, `moveTask(id, newStatus, position)`, `getPreferences()`, `updatePreferences(changes)`
- Subscribe/notify pattern: `subscribe(event, callback)`, `notify(event, data)`
- Events: `tasks-changed`, `preferences-changed`
- State mutations persist to storage automatically
- Error in storage write does not corrupt in-memory state (write failures are surfaced, not swallowed)
- `StateManager` listens to `chrome.storage.onChanged` and re-hydrates on external changes
- Re-hydration emits `tasks-changed` and `preferences-changed` events so UI updates
- Race condition guard: if a local write is in-flight (debounce pending), incoming external changes don't clobber it
