---
id: VTS-018
title: "JSON Export and Import"
complexity: M
status: pending
owner: ""
source: oracle
source_ref: "oracle-architect-breakdown.md"
dependencies:
  - VTS-002
  - VTS-004
  - VTS-014
files:
  - "`scripts/utils/export-import.js`"
  - "`scripts/components/preferences-panel.js` (extend)"
---

# JSON Export and Import

Build a one-click JSON export of all application data (tasks, column order, preferences) as a downloadable `.json` file, and an import function that can restore from that file. This is disaster recovery for a local-only application. The difference between "I lost a task" and "I lost everything."

## Criteria

- Export button in preferences panel downloads `kanban-backup-YYYY-MM-DD.json`
- Export includes full storage envelope (version, tasks, columnOrder, preferences)
- Import button accepts `.json` file via native file picker
- Import validates schema version and data structure before overwriting
- Import shows confirmation dialog with task count before applying ("Import 23 tasks? This will replace your current data.")
- Invalid file shows clear error message
- Export works in both popup and tab views
- No additional permissions required (`chrome.downloads` not needed — use `Blob` + `URL.createObjectURL` + anchor click)
