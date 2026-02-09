---
id: VTS-002
title: "Data Model and Type Definitions"
complexity: S
status: pending
owner: ""
source: discovery
source_ref: "07-architect-architect-breakdown.md"
dependencies: []
files:
  - "src/types/index.ts"
---

# Data Model and Type Definitions

Define TypeScript interfaces for Task, Tag, Column, and the overall store shape. Define default columns (To Do, In Progress, Done). Include a schema version constant for persistence migration support. Include an export envelope type for JSON export/import (VTS-018).

Compared to earlier specs: `shortDescription` is dropped (derive a summary from `bodyMarkdown` if needed). `startDate` is dropped from v1. `isCompletionColumn` is added to Column to identify which column(s) represent "done" status.

## Data Model

### Task
```
id: string (nanoid)
title: string
bodyMarkdown: string
statusId: string (FK -> Column)
order: number
tagIds: string[] (FK -> Tag[])
dueDate: string | null (ISO date, no time)
completedAt: string | null (ISO datetime)
createdAt: string (ISO datetime)
updatedAt: string (ISO datetime)
```

### Tag
```
id: string (nanoid)
name: string
color: string (hex)
```

### Column (Kanban Status)
```
id: string (nanoid)
name: string
order: number
color: string | null (optional accent)
isCompletionColumn: boolean
```

### Schema Version
```
SCHEMA_VERSION: number (starts at 1)
```

### Export Envelope
```
ExportEnvelope {
  schemaVersion: number
  exportedAt: string (ISO datetime)
  tasks: Record<string, Task>
  tags: Record<string, Tag>
  columns: Record<string, Column>
  columnOrder: string[]
}
```

## Criteria

- All entity types exported from `src/types/index.ts`
- Default column definitions defined (To Do, In Progress, Done)
- Done column has `isCompletionColumn: true`, others have `false`
- `shortDescription` field is NOT present on Task (removed from v1)
- `startDate` field is NOT present on Task (removed from v1)
- `SCHEMA_VERSION` constant is exported
- `ExportEnvelope` type is exported for use by VTS-018
- Types include all specified fields with correct types
- No `any` types
