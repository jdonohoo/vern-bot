### TASK 7: Error Handling & Edge Cases

**Description:** Audit and harden all error paths. This is not a new feature — it's a pass over Tasks 2–6 to ensure every failure mode surfaces a clear, non-technical message. Specific cases: binary file uploaded, file read interrupted (e.g., user navigates away), Papa Parse returns errors, completely empty parse result, file with only whitespace.

**Acceptance Criteria:**
- Binary/non-text file: "This doesn't appear to be a text file. Please upload a CSV."
- Parse returns zero rows: "No data found in this file."
- Parse returns errors: Show count and first 3 errors below the table (not blocking render)
- All errors are displayed in the error region, styled distinctly (red border, not aggressive red background)
- Errors are cleared when a new file is selected
- No `console.error` calls reach production — all caught and surfaced to the UI
- No uncaught promise rejections

**Complexity:** S

**Dependencies:** Tasks 2, 3, 4, 5, 6 (this is a hardening pass over everything)
