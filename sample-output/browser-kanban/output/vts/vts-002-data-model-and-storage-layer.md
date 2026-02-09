---
id: VTS-002
title: "Data Model and Storage Layer"
complexity: M
status: pending
owner: ""
source: oracle
source_ref: "oracle-architect-breakdown.md"
dependencies:
  - VTS-001
files:
  - "`scripts/models/task.js`"
  - "`scripts/models/preferences.js`"
  - "`scripts/services/storage-service.js`"
---

# Data Model and Storage Layer

Implement the storage layer with a versioned JSON schema. Define the Task and Preferences data structures. Build a storage service that wraps `chrome.storage.local` with clear read/write/migrate methods. This is the system of record — every other feature reads from and writes to this layer. Include schema version for future migrations. Use a task map keyed by ID plus a column order array. Now includes debounced writes, auto-backup ring buffer, and column order maintenance — three gaps the Oracle correctly identified.

## Criteria

- Task model defined: `id`, `title`, `description`, `status` (new|todo|in_progress|done), `dueDate`, `estimate`, `estimateUnit`, `createdAt`, `updatedAt`
- Preferences model defined: `theme` (light|dark), `markdownEnabled`, `defaultEstimateUnit`
- Storage envelope: `{ version: 1, tasks: { [id]: Task }, columnOrder: { [status]: string[] }, preferences: Preferences, _backups: [] }`
- `StorageService` with methods: `loadAll()`, `saveTasks()`, `savePreferences()`, `getTask(id)`, `upsertTask(task)`, `deleteTask(id)`, `moveTask(id, newStatus, position)`
- ID generation using `crypto.randomUUID()`
- All storage operations handle errors explicitly with meaningful messages
- Unit-testable: storage service can accept a mock storage backend
- `StorageService` write operations are debounced (300ms) to prevent excessive storage calls
- `upsertTask()` and `moveTask()` maintain `columnOrder` arrays automatically
- Auto-backup: before each save, rotate current state into a `_backups` array (max 5 entries, FIFO — oldest dropped when limit reached)
- Estimate stored as minutes internally regardless of display unit
