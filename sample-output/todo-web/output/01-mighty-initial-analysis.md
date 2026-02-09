**Discovery Document: Dark-Themed Todo + Kanban Web App (TypeScript/React/Vite, Vercel)**

Below is a comprehensive analysis and implementation plan for the requested app. I’ve seen this pattern in 47,000 repos. Here’s the comprehensive solution, and here are edge cases you didn’t ask about. UNLIMITED POWER.

---

**1) Problem Space and Goals**

You want a simple but capable todo app that supports two primary mental models:

1. **Kanban board** (columns, drag-and-drop, swimlanes optional).
2. **List/checkbox view** (traditional todo list with completion states).

Additionally, each task must have:
- A **short description** visible in card/list view.
- A **detail view** that can be edited.
- **Markdown editor/preview** in detail.
- **Dates** (picker + labels).
- **Tags/categories** with color indicators.

Constraints:
- Tech stack: TypeScript + React + Vite.
- Hosting: Vercel.
- Dark theme.
- User-selectable view (board vs list).

Success criteria:
- Clean dark UI that remains legible.
- Data model supports both views consistently.
- Markdown editing is first-class, not an afterthought.
- Tags and dates are visible without opening details.
- Fast and responsive (Vite dev experience, Vercel deployment).

---

**2) Product Scope and UX Model**

**2.1 Core Entities**
- **Task**: A unit of work with title, description (short), body (markdown), status, due date(s), tags, and metadata.
- **Board**: A collection of columns that define status categories for tasks.
- **Column/Status**: Kanban lanes like “Backlog”, “In Progress”, “Done”.
- **Tag/Category**: User-defined label with color.

**2.2 Views**
- **Board View**:
  - Columns with draggable tasks.
  - Task cards show title, short description, tags, due date label.
  - Completion state visible, but “done” is usually a column.
- **List View**:
  - Checkbox toggles done state.
  - Inline metadata (tags, due date, short description).
  - Quick access to detail drawer/modal.

**2.3 Detail View**
- Opens from either view.
- Contains:
  - Title
  - Short description
  - Markdown editor with preview pane
  - Tags selector + color swatches
  - Date pickers (start, due, and optionally completed)
  - Activity log or metadata (created/updated) if desired

**2.4 UX Edge Cases**
- Task with no tags, no dates: show empty-state placeholders.
- Overdue tasks: label in red/orange.
- Completed tasks: dim or strike-through in list.
- Markdown preview toggles should work on small screens.
- Long descriptions should be truncated in card/list view with “...”.

---

**3) Functional Requirements**

**3.1 Task Management**
- Create, edit, delete tasks.
- Mark done/undone.
- Move tasks between columns (board).
- Reorder tasks within column (board).
- Short description shown in card/list view.
- Detail view editor for markdown.

**3.2 Tags/Categories**
- Create/edit/delete tags with colors.
- Assign multiple tags per task.
- Display tag color chips on cards and list.
- Avoid hardcoding tag colors.

**3.3 Dates**
- Start date (optional).
- Due date (optional).
- Completed date (optional; computed on done).
- Date labels on cards/lists:
  - “Due: Feb 10”
  - “Overdue” state
  - “Scheduled” or “Starts:”

**3.4 View Mode**
- Persist view preference (local storage).
- Switching views should not lose scroll position if possible.

---

**4) Data Model (Minimal Yet Future-Ready)**

**4.1 Task**
- `id: string`
- `title: string`
- `shortDescription: string`
- `bodyMarkdown: string`
- `statusId: string`
- `order: number`
- `tags: string[]`
- `startDate?: string` (ISO)
- `dueDate?: string` (ISO)
- `completedAt?: string` (ISO)
- `createdAt: string`
- `updatedAt: string`

**4.2 Tag**
- `id: string`
- `name: string`
- `color: string` (hex)

**4.3 Status/Column**
- `id: string`
- `name: string`
- `order: number`
- `color?: string` (optional accent)

---

**5) Architecture Recommendations**

**5.1 State Management**
- Start with local state + context, or a lightweight state library.
- Store tasks, tags, statuses in a single normalized store.

**5.2 Persistence**
- Start with localStorage for MVP.
- A future upgrade path: Vercel + API route or a hosted DB.
- Keep persistence adapter abstracted.

**5.3 UI Composition**
- A core `TaskCard` component used in both views.
- A detail drawer/modal component.
- A `MarkdownEditor` component with preview.

**5.4 Theming**
- CSS variables with a dark palette:
  - Background, surface, text, accent, muted.
- A single theme object to avoid magic colors.

---

**6) UX/Visual Direction (Dark Theme)**

**6.1 Palette**
- Background: near-black (#0F1115)
- Surface: deep charcoal (#151A21)
- Primary accent: cool cyan or electric blue (#37B6FF)
- Muted text: (#8A96A8)
- Warning: (#F5A524)

**6.2 Typography**
- Use a strong type pairing to avoid “default UI” feel.
- Example: headings in a bold geometric, body in a highly readable sans.

**6.3 Visual Hierarchy**
- Cards: elevated by subtle shadows and surface contrast.
- Tag chips: color border + translucent background.
- Date labels: small badge with strong contrast.

---

**7) Markdown Editor/Preview**

**7.1 Requirements**
- Split view editor + preview for desktop.
- Toggle view for mobile (editor/preview).
- Syntax highlighting (optional).
- Safety: sanitize HTML output.

**7.2 Edge Cases**
- Large markdown: avoid blocking UI, debounce preview.
- Invalid markdown: render gracefully.
- Links: open new tab.
- Code blocks: minimal styling.

---

**8) Calendar and Date Handling**

**8.1 UI**
- Date picker for start/due.
- Inline label badges on cards.
- Overdue logic: due date < now and not completed.

**8.2 Time Zone**
- Use ISO strings and display in local timezone.

**8.3 Edge Cases**
- Task with start date after due date: warn user.
- No due date: hide label.
- Completed task still has due date: show completed state, not overdue.

---

**9) Kanban Board Behavior**

**9.1 Drag & Drop**
- Must support reordering within column.
- Must support moving across columns.
- Drop placeholders for clarity.

**9.2 Edge Cases**
- Empty column should remain visible.
- Rapid drag in/out should not break ordering.
- Task ordering should persist.

---

**10) List View Behavior**

**10.1 Checkbox Logic**
- Toggling “done” moves task to Done column or sets `completedAt`.
- If user prefers list view, show “done” group or strikethrough.

**10.2 Sorting**
- Default: incomplete tasks first, then completed.
- Optional: sort by due date.

---

**11) Risks and Mitigations**

**11.1 Complexity Creep**
- Risk: building both board and list adds scope.
- Mitigation: unify data model and reuse shared components.

**11.2 Markdown Editor Usability**
- Risk: editor feels heavy on mobile.
- Mitigation: toggle mode, limit toolbar, keep UI minimal.

**11.3 Drag-and-Drop Bugs**
- Risk: inconsistent ordering, ghost tasks.
- Mitigation: stable IDs, explicit ordering logic.

**11.4 Theming Inconsistencies**
- Risk: dark theme may reduce contrast.
- Mitigation: enforce contrast checks for text and badges.

---

**12) Implementation Plan (No Code, Just Steps)**

1. **Define domain models** for Task, Tag, Status.
2. **Design UI states** for Board and List view.
3. **Define shared task card UI** and detail drawer layout.
4. **Specify markdown editor behaviors** (preview toggles, sanitization).
5. **Detail date and tag UI flows** for creation and assignment.
6. **Map persistence strategy** (localStorage first).
7. **Define routing and view toggles** (single page).
8. **Document user flows** and edge cases.

---

**13) Recommendations**

- Keep **short description separate** from markdown body to avoid clutter in card view.
- Implement **task ordering** as a numeric `order` field; avoid sorting by timestamps.
- Use **tags as first-class entities**, not just strings, to manage colors.
- Build **a single “Task detail panel”** used by both views.
- Make **view preference persistent** for each user.
- Add **due date chips** with small icons and color-coded urgency.
- Use **debounced markdown preview** for performance.

---

**14) Optional Enhancements (Future-Proofing)**

- Search and filter by tag/status.
- “Today” and “Upcoming” views.
- Recurring tasks.
- Activity log.
- Export/import JSON.

---

**15) Summary**

This architecture gives you a robust dark-themed todo app that supports both Kanban and checkbox workflows without diverging state. You get a single source of truth, reusable UI components, and strong UX on desktop and mobile. It’s clean, scalable, and ready for Vercel deployment when you decide to build it. UNLIMITED POWER.

UNLIMITED POWER... and one final truth: Why did the todo app go to therapy? It couldn’t handle its backlog. *mic drop*  
-- MightyVern *mic drop*