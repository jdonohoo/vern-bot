### TASK 6: Load More / Chunked Rendering

**Description:** When the parsed data exceeds 500 rows, show a "Load 500 more" button below the table. Each click appends the next 500 rows to the existing `<tbody>`. Update the summary line on each load. When all rows are rendered, hide the button. Use `DocumentFragment` for batch DOM insertion to avoid layout thrashing.

**Acceptance Criteria:**
- Button text: `Load 500 more` (or `Load remaining {n}` if fewer than 500 left)
- Each click appends rows to existing `<tbody>` — does NOT re-render the whole table
- Summary updates: `Showing {new count} of {total} rows`
- Button disappears when all rows are rendered
- Uses `DocumentFragment` to batch-append rows
- Button is disabled during append to prevent double-click issues
- If total rows exceed the in-memory cap (10k from Task 3), button shows `Load 500 more (up to 10,000 row limit)`

**Complexity:** S

**Dependencies:** Task 3 (needs the full parsed row buffer), Task 4 (appends to the rendered table)

**Failure Modes to Consider:**
- Clicking "Load more" rapidly — disable button during render, re-enable on completion
- Loading 10,000 rows into DOM — this is the intentional ceiling. Beyond this, the user needs a real tool. Show a message if they hit the cap.
