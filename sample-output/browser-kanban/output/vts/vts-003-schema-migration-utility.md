---
id: VTS-003
title: "Schema Migration Utility"
complexity: S
status: pending
owner: ""
source: oracle
source_ref: "oracle-architect-breakdown.md"
dependencies:
  - VTS-002
files:
  - "`scripts/services/migration-service.js`"
  - "`scripts/services/storage-service.js` (extend)"
---

# Schema Migration Utility

Build a schema migration runner that checks the stored `version` field against the current expected version and runs sequential migration functions to bring data forward. The Architect put a `version: 1` field in the storage envelope but wrote no code to act on it. Version fields without migration code are tombstones for future bugs. This must exist before any real data is persisted, because once users have data at version 1, you cannot retroactively add migration logic without risking corruption.

## Criteria

- Migration runner executes on every `StorageService.loadAll()` call
- Migrations are sequential functions: `migrate_1_to_2()`, `migrate_2_to_3()`, etc.
- Each migration is idempotent (safe to re-run)
- Migration failure halts startup and shows error with option to export raw data
- Current version constant defined in one place
- First-run detection (no existing data) skips migrations and initializes fresh schema
- Migration log written to console for debugging
