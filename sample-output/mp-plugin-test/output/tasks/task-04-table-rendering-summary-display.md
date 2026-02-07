### TASK 4: Table Rendering & Summary Display

**Description:** Take the parsed data (header + rows) and render an HTML `<table>` in the results area. Render only the first 500 rows initially. Display a summary line: `Showing {rendered} of {total} rows · {columns} columns`. Apply basic styling: sticky `<thead>`, zebra striping on rows, horizontal scroll wrapper for wide tables, monospace font for cell content. Keep it semantic — `<thead>`, `<tbody>`, `<th>`, `<td>`.

**Acceptance Criteria:**
- Table renders inside a `<div>` with `overflow-x: auto` for horizontal scrolling
- `<thead>` uses `position: sticky; top: 0` so headers stay visible on scroll
- Zebra striping via `tr:nth-child(even)` — no JS needed
- Summary line is always visible above the table
- Empty cells render as empty `<td>` (not "null" or "undefined")
- Column headers show column index if no header row (e.g., "Column 1", "Column 2")
- DOM contains at most 500 `<tr>` elements after initial render
- Table is cleared and rebuilt (not appended to) when a new file is uploaded

**Complexity:** M

**Dependencies:** Task 3 (needs parsed data), Task 1 (needs results area DOM)

**Failure Modes to Consider:**
- Table with 100+ columns — horizontal scroll handles it, but cells should have `white-space: nowrap` with a `max-width` and `text-overflow: ellipsis` to prevent layout blowout
- Cell containing HTML-like content — must be set via `textContent`, never `innerHTML`
- Very long cell values (e.g., 10,000 chars) — truncate display to ~200 chars with a `title` attribute for full value
