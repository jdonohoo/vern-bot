---
id: VTS-018
title: "JSON Export/Import"
complexity: M
status: pending
owner: ""
source: discovery
source_ref: "07-architect-architect-breakdown.md"
dependencies:
  - VTS-003
files:
  - "src/utils/exportImport.ts"
  - "src/components/shared/ExportImportControls.tsx"
---

# JSON Export/Import

Implement full JSON export of all data (tasks, tags, columns, settings) and JSON import with validation and conflict resolution. The export uses the `ExportEnvelope` type defined in VTS-002 and includes the schema version. Add UI trigger in the app header or a settings area.

This is the user's only backup mechanism and their only path off the product. For a localStorage-only app, its absence is a shipping blocker. The council reached 7+ Vern consensus on this.

## Criteria

- **Export:** Produces a valid JSON file containing all tasks, tags, columns, column order, and settings
- **Export:** JSON includes `schemaVersion` and `exportedAt` timestamp (uses `ExportEnvelope` type)
- **Export:** Triggers a browser file download with a meaningful filename (e.g., `todo-export-2025-03-15.json`)
- **Import:** Accepts a JSON file upload
- **Import:** Validates the schema version before applying
- **Import:** Handles malformed JSON gracefully (shows an error message, does not crash)
- **Import:** Handles missing or extra fields gracefully (fills defaults for missing, ignores unknown)
- **Import:** Asks user to confirm before overwriting existing data (overwrite mode)
- **Round-trip:** Export then import produces identical application state
- UI trigger is accessible from the app header or a settings/menu area
- Import/export works in both light and dark themes
