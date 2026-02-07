### TASK 2: File Input & Validation

**Description:** Wire up the file input to accept `.csv` and `.tsv` files. On file selection, validate: file exists, file is not empty, file size is under a sane client-side limit (say 50MB — beyond that, browsers choke anyway). On validation failure, display a clear error message in the error region. On success, hand off the `File` object to the parser (Task 3). Reset previous results/errors when a new file is selected.

**Acceptance Criteria:**
- `<input type="file" accept=".csv,.tsv,.txt">` is the entry point
- Selecting a 0-byte file shows: "File is empty"
- Selecting a file >50MB shows: "File too large (max 50MB)"
- Selecting a valid file clears any previous table and errors
- The `File` object is passed to the parsing function, not raw text (important for streaming)
- No drag-and-drop — intentionally deferred

**Complexity:** S

**Dependencies:** Task 1 (needs the DOM structure to exist)

**Failure Modes to Consider:**
- User selects a non-CSV file (e.g., renamed `.jpg`) — Papa Parse will handle this gracefully in Task 3, but the error messaging needs to surface it
- User re-selects the same file — `input.value = ''` trick needed to re-trigger `change` event
