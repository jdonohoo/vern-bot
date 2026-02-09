### TASK 3: CSV Parsing with Papa Parse (Streaming)

**Description:** Parse the uploaded file using Papa Parse's streaming mode (`Papa.parse(file, { step, complete })`). Accumulate rows into an array. Track total row count separately from rendered rows. Detect whether the first row looks like a header (default: yes, overridable in Task 5). On parse error, surface Papa Parse's error array to the user. On success, pass the parsed data structure to the renderer (Task 4).

**Acceptance Criteria:**
- Uses `Papa.parse()` with `step` callback for streaming — does NOT read entire file into a string first
- Collects first 500 rows into a render buffer
- Continues counting total rows beyond 500 (for the summary)
- Stores all parsed rows in memory up to a reasonable cap (10,000 rows in memory, discard beyond that but keep the count)
- Parse errors (malformed rows, inconsistent column counts) are collected and shown after the table, not as blocking errors
- Header row is separated from data rows based on the "first row is header" setting

**Complexity:** M

**Dependencies:** Task 2 (needs the File object), Task 1 (needs error display region)

**Failure Modes to Consider:**
- File with inconsistent column counts — render anyway, pad short rows with empty cells
- File with only a header and no data — show header, display "No data rows found"
- File with 1M+ rows — streaming prevents memory blowup, but the in-memory cap (10k) prevents browser tab crashes. Summary shows "Showing 500 of 1,000,000+ rows"
