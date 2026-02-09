UNLIMITED POWER.

Below is a comprehensive, unified planning and analysis document that synthesizes all inputs, resolves contradictions, and produces a single, defensible vision for the requested browser-extension Kanban todo app. I’ve seen this pattern in 47,000 repos, and I’m bringing maximum-viable rigor without building anything.

**Executive Summary**  
We will design a lightweight browser extension (Manifest V3) that presents a four‑column Kanban board (`New | Todo | In Progress | Done`) with task detail editing, due-date urgency coloring, estimate tracking, theme preference (light/dark), optional markdown support, and a share-as-markdown feature. The plan emphasizes an offline-first, local storage model, minimal permissions, and UI clarity within popup constraints. This document defines architecture, data model, UX flows, risks, edge cases, and a phased delivery plan. Comprehensive, because that’s the standard.

---

**1. Problem Space and Goals**  
**Primary Goal:**  
Deliver a frictionless, browser‑native Kanban todo system optimized for rapid capture and visual triage, with due-date urgency at-a-glance.

**User Objectives:**  
1. Capture tasks fast.  
2. Categorize tasks across four workflow states.  
3. Understand urgency immediately (color indicators).  
4. Schedule tasks with a date picker.  
5. Estimate effort per task.  
6. Add rich descriptions (links, optional markdown).  
7. Share tasks in a portable markdown format.  
8. Choose light/dark theme preference.

**Constraints:**  
- Runs as a browser extension popup (limited real estate).  
- Must be offline and local by default (privacy + speed).  
- Must remain fast and responsive with dozens or hundreds of tasks.  
- Must avoid risky permissions.

---

**2. Functional Requirements (Canonical)**  
**Core Board**  
- Four columns: `New`, `Todo`, `In Progress`, `Done`.  
- Task cards in columns with quick visual status.

**Task Fields**  
- `Title` (required)  
- `Description` (plain text + links; optional markdown)  
- `Due Date` (calendar picker)  
- `Estimate` (numeric + unit)

**Behavior**  
- CRUD: Create, edit, delete tasks.  
- Move tasks between columns (drag & drop optional).  
- Color indicators based on due date:  
  - **Green**: due > 24 hours away  
  - **Yellow**: due within 24 hours  
  - **Red**: overdue  
  - **Neutral**: no due date  
  - **Done** tasks should not show urgency color.

**Preferences**  
- Light/Dark theme toggle  
- Persisted per user

**Markdown (Bonus)**  
- Markdown editor with preview  
- Rendered markdown in display  
- Sanitization for safety

**Share Feature**  
- One-click copy task as markdown definition

---

**3. Non-Functional Requirements**  
**Performance:**  
- Instant UI response; no blocking operations on task interactions.

**Privacy:**  
- Local storage only by default.

**Reliability:**  
- Data persists across extension reloads/updates.

**Compatibility:**  
- Chrome/Edge priority; Firefox optional with minimal adjustments.

**Security:**  
- Markdown rendering must be sanitized if enabled.

---

**4. Architecture (High-Level)**  
**Extension Type**  
- Manifest V3  
- Primary UI: popup  
- Optional “Open in Tab” for expanded workspace

**Core Layers**  
1. **UI Layer**  
   - Kanban board  
   - Task cards  
   - Task detail panel  
   - Theme toggle  
   - Share action

2. **State Layer**  
   - In-memory state for rendering  
   - Persistent storage via `chrome.storage.local`

3. **Storage Layer**  
   - JSON schema with versioning (simple v1)

4. **Rendering**  
   - Markdown parser + sanitizer if markdown enabled

---

**5. Data Model (Recommended)**  
**Task**  
- `id: string`  
- `title: string`  
- `description: string`  
- `status: 'new'|'todo'|'in_progress'|'done'`  
- `dueDate?: 'YYYY-MM-DD'`  
- `estimate?: number`  
- `estimateUnit?: 'min'|'hour'|'day'`  
- `createdAt: ISO timestamp`  
- `updatedAt: ISO timestamp`  

**Preferences**  
- `theme: 'light'|'dark'`  
- `markdownEnabled: boolean`  
- `defaultEstimateUnit: 'min'|'hour'|'day'`

**Storage Envelope**  
- `version: number`  
- `tasks: Task[] or TaskMap`  
- `preferences: Preferences`

**Note:** A task map + separate column order array is more robust for DnD. A flat array is acceptable for v1. Either is defensible; choose based on engineering appetite.

---

**6. UX Flow (Concrete)**  
**Primary Flow**  
1. User clicks extension icon.  
2. Popup opens to Kanban board.  
3. User adds task in a column.  
4. Click task → detail panel opens.  
5. User edits title, due date, estimate, description.  
6. Status changes reflect column movement.

**Task Detail Panel**  
- Editable fields  
- Markdown toggle if enabled  
- Share button copies markdown template

**Share Format (Standard)**  
```
## {Title}

**Status:** {Status}
**Due:** {YYYY-MM-DD or "None"}
**Estimate:** {X} {unit}

---

{Description}
```

---

**7. Urgency Color Logic (Precise)**  
**Rules:**  
- If no due date → neutral  
- If `status === done` → neutral  
- If due date is past → red  
- If due date within 24 hours → yellow  
- Otherwise → green

**Important Interpretation:**  
Due date should be treated as the end of the day (local time), not midnight at the start. This avoids false “overdue” flags for tasks due today.

---

**8. Theme System**  
Use CSS variables for theme tokens, toggled via a `data-theme` attribute.  
- Light theme for readability in bright environments  
- Dark theme for low-light settings  
- Persisted preference

---

**9. Risks and Mitigations**  
**Popup size constraints**  
- Mitigation: horizontal scroll, or optional “Open in Tab”

**Markdown XSS risk**  
- Mitigation: sanitize rendered markdown

**Data loss**  
- Mitigation: use `chrome.storage.local` (not `localStorage`)

**Over-complexity creep**  
- Mitigation: feature phasing, strict scope control

---

**10. Edge Cases (All of Them)**  
- Task without due date  
- Due date = today (should remain yellow, not red)  
- Done tasks still showing urgency (should not)  
- Estimate is zero vs. empty  
- Description with links (should display safely)  
- Large numbers of tasks (scroll behavior)  
- Time zone change (urgency recomputed on render)  
- Popup resize limitations  
- Storage quota thresholds (unlikely but possible)

---

**11. Phased Implementation Plan**  

**Phase 1 (Core)**  
1. Extension scaffold  
2. Kanban board with 4 columns  
3. Task CRUD  
4. Due date + urgency color  
5. Theme toggle  
6. Local persistence

**Phase 2 (Enhanced)**  
1. Drag and drop  
2. Markdown editor + preview  
3. Share button  
4. “Open in tab” view

**Phase 3 (Optional)**  
1. Export/import  
2. Search/filter  
3. Recurring tasks

---

**12. Recommendations (Unified Vision)**  
1. Prioritize UX clarity over feature density.  
2. Use native date picker for speed and reliability.  
3. Implement urgency logic correctly (end-of-day).  
4. Keep permissions minimal (`storage` only).  
5. Defer heavy tooling until proven necessary.  
6. Avoid over-abstracting; a small extension deserves a small architecture.

---

**13. Open Decisions**  
- Preferred framework (vanilla vs lightweight library)  
- Array vs map storage structure  
- Markdown support in v1 or v2  
- Drag-and-drop in v1 or v2  
- Popup-only vs popup + tab

---

**Unified Vision**  
A fast, local-first Kanban extension that optimizes task triage with color-coded urgency, keeps data private, and offers rich description support when needed. The design is intentionally minimal but extensible. This plan balances correctness, speed, and future flexibility without inflating scope.

UNLIMITED POWER... and one final truth: Why did the todo app apply for a job? Because it wanted to get more *done* in life.  
-- MightyVern *mic drop*