---
id: VTS-003
title: "State Store with Persistence"
complexity: M
status: pending
owner: ""
source: discovery
source_ref: "07-architect-architect-breakdown.md"
dependencies:
  - VTS-002
  - VTS-016
files:
  - "src/store/index.ts"
  - "src/store/taskSlice.ts"
  - "src/store/tagSlice.ts"
  - "src/store/columnSlice.ts"
  - "src/store/persistence.ts"
  - "src/store/migrations.ts"
---

# State Store with Persistence

Implement Zustand store with normalized state for tasks, tags, and columns. Include CRUD actions for all entities. Add localStorage persistence middleware with schema versioning, migration support, and a migration stub function. Include debounced writes, cross-tab synchronization via the `storage` event listener, and quota monitoring that warns at 80% capacity. Persist theme preference (set by VTS-016).

## Criteria

- Store initializes with default columns
- Task CRUD operations work (create, read, update, delete)
- Tag CRUD operations work
- Column reordering works
- State persists to localStorage on changes (debounced, 500ms)
- State loads from localStorage on app start
- Schema version is stored and checked on load
- Migration stub function exists and is called when stored version < current version
- `migrations.ts` exports a `migrate(oldState, fromVersion, toVersion)` function with a clear pattern for adding version-specific transforms
- Malformed localStorage data falls back to defaults without crashing
- Cross-tab synchronization: listens to `window.storage` events and updates store when another tab writes
- Quota monitoring: warns user when localStorage usage exceeds 80% of estimated 5MB limit
- Theme preference (light/dark/system) is persisted alongside app state
