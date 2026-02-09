Here’s the comprehensive analysis and planning document for the requested browser plugin todo app. I’ve seen this pattern in 47,000 repos, and I’m going to cover every angle with UNLIMITED POWER.

**1) Problem Space Overview**
- **Core need**: A lightweight, in-browser todo/kanban tool that lives in a browser extension (plugin), supporting task flow across states: `new | todo | in progress | done`.
- **Key user value**: Fast capture, lightweight scheduling, quick status visibility, and minimal friction. The kanban view provides an at-a-glance workflow and makes state changes intuitive.
- **Differentiators**: 
  - Calendar-based due dates with visual urgency.
  - Task estimate field.
  - Rich description with links, and optional markdown editing/rendering.
  - Shareable task definition in markdown.

**2) Technical Requirements**

**Functional**
1. **Task lifecycle columns**: `New`, `Todo`, `In Progress`, `Done`.
2. **CRUD**:
   - Create, edit, delete tasks.
   - Update status by drag-and-drop or quick actions.
3. **Task fields**:
   - Title (required).
   - Description (plain text + links; optional markdown).
   - Due date via calendar picker.
   - Estimate (time; numeric + unit).
4. **Color indicators**:
   - Green: due date in future (more than 1 day).
   - Yellow: due within 24 hours.
   - Red: overdue.
   - Optional: Gray when no date set.
5. **Theme**:
   - Light/Dark preference.
   - Persisted per user.
6. **Markdown support (bonus)**:
   - Markdown editor with preview toggle.
   - Render markdown in display view.
7. **Share**:
   - Generate full task definition in markdown.
   - Copy to clipboard.
   - Optional: share modal or inline.

**Non-functional**
- **Performance**: instantaneous UI updates; small data sets but must remain responsive.
- **Storage**: local persistence (extension storage).
- **Reliability**: no data loss during extension updates.
- **Privacy**: all local; no network by default.
- **Portability**: Chrome + Edge; optional Firefox.

**3) Proposed Architecture**

**A. Extension Type**
- **Manifest V3** extension.
- **Popup UI** as primary interface.
- Optional **side panel** (Chrome) or full page for larger boards.

**B. Components**
1. **UI Layer**
   - Kanban board with columns.
   - Task cards with priority color bars or border.
   - Task detail drawer/modal for editing.
   - Theme toggle in settings.
   - Share modal.

2. **State Management**
   - In-memory state for current session.
   - Persisted to extension storage (Chrome storage API).
   - Optimistic updates for UI.

3. **Storage**
   - `chrome.storage.local` (or `browser.storage.local`).
   - Schema versioning to allow migrations.

4. **Rendering**
   - Markdown parser (e.g., `marked`, `markdown-it`).
   - Sanitizer (e.g., `DOMPurify`) to avoid XSS from markdown.

5. **Date & Time**
   - Due date stored in ISO date format (`YYYY-MM-DD`).
   - Color indicator logic computed at render-time.
   - Use local timezone for “within 24 hours.”

**C. Data Model (Logical)**
- `Task`
  - `id: string`
  - `title: string`
  - `description: string`
  - `markdownEnabled: boolean`
  - `status: enum` (`new|todo|in_progress|done`)
  - `dueDate?: string` (ISO date)
  - `estimate?: number`
  - `estimateUnit?: 'min'|'hour'|'day'`
  - `createdAt: string`
  - `updatedAt: string`

- `Preferences`
  - `theme: 'light'|'dark'`
  - `defaultEstimateUnit`
  - `markdownDefault`
  - `dateLocale`

**D. Board Behavior**
- Drag-and-drop between columns.
- If DnD not implemented, allow quick move button.
- Task cards show:
  - Title
  - Due date and estimated time
  - Color indicator strip

**4) UI/UX Considerations**

**Board Layout**
- Horizontal columns with scroll.
- Responsive layout in popup width/height constraints.
- If popup is too small, enable “Open in full tab” view.

**Task Detail**
- Clicking task opens panel or modal.
- Fields: title, description, links, due date, estimate, status.
- Markdown toggle:
  - Edit (textarea with markdown).
  - Preview (rendered markdown).

**Theme**
- Preference stored and applied at launch.
- Provide toggle in settings or header.

**Share Feature**
- “Share” button in task detail.
- Markdown output includes all fields:
  - Title
  - Status
  - Due date
  - Estimate
  - Description
- Copy to clipboard, with success feedback.

**5) Edge Cases You Didn’t Ask About**
- Tasks with no due date (no color or neutral).
- Invalid or cleared due date.
- Estimate field empty vs zero.
- Due date exactly now or within 24h (should be Yellow).
- Time zone edge:
  - User changes timezone: recompute on render.
- Markdown XSS risk: sanitize.
- Large number of tasks (pagination? scroll?).
- Drag-and-drop on touch devices (fallback).
- Storage quota (Chrome limit). Data remains small but add guard.
- Extension update migration for schema changes.

**6) Risks and Mitigations**
- **Data loss**: use atomic writes; store full data blob.
- **XSS**: sanitize markdown output.
- **Storage size**: encourage cleanup or add archive.
- **UI constraints**: popup size limits. Offer “open in tab”.

**7) Implementation Plan (No Code, Planning Only)**

1. **Specification Finalization**
   - Confirm columns, fields, and exact color logic.
   - Confirm markdown library and sanitization approach.
   - Confirm target browsers and extension type.

2. **Data Model + Storage Plan**
   - Define JSON schema.
   - Define storage key layout.
   - Versioning strategy.

3. **UX Flow Definition**
   - Board layout in popup and full tab.
   - Task detail interaction.
   - Settings and theme application.
   - Share modal UX.

4. **Validation & Testing Plan**
   - Unit: date logic, markdown rendering, storage.
   - Integration: create/edit/delete tasks.
   - Visual: theme toggling and color indicators.
   - Migration: old schema upgrade.

**8) Recommendations**
- **Priority 1**: Implement board + CRUD + storage first.
- **Priority 2**: Add due date color logic.
- **Priority 3**: Add markdown support with sanitization.
- **Priority 4**: Share as markdown with clipboard.
- **Priority 5**: Polish drag-and-drop UX and transitions.

**9) Unknowns**
- Target browser(s)?
- Preferred UI framework (if any)?
- Data portability: export/import to JSON or markdown?
- Should tasks sync across devices (requires cloud or Chrome sync storage)?

UNLIMITED POWER... and one final truth: Why don’t tasks ever win arguments? Because they can’t handle the due date. *mic drop*  
-- MightyVern *mic drop*