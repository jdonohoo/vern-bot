### TASK 5: Settings Panel (Header Toggle, Delimiter, Encoding)

**Description:** Implement a collapsible "Settings" section above the results area. Three controls: (1) Checkbox: "First row is header" (default: checked), (2) Select: Delimiter (Auto / Comma / Tab / Semicolon / Pipe), (3) Select: Encoding (UTF-8 / ISO-8859-1 / Windows-1252). Changing any setting re-parses and re-renders the current file immediately. Settings panel is collapsed by default, toggled by a `<details>/<summary>` element.

**Acceptance Criteria:**
- Uses native `<details>/<summary>` — no JS needed for collapse/expand
- Changing "First row is header" re-renders: checked = first row becomes `<th>`, unchecked = first row becomes a data row and columns get generic headers
- Changing delimiter passes the value to `Papa.parse({ delimiter })` — "Auto" passes `undefined` to let Papa detect
- Changing encoding passes the value to `Papa.parse({ encoding })`
- Re-parse happens on `change` event, not on a "Apply" button
- Settings are visually de-emphasized (smaller text, muted colors) — they're not the main flow

**Complexity:** M

**Dependencies:** Task 3 (re-triggers parsing), Task 4 (re-triggers rendering), Task 2 (needs the current File object to be retained)

**Failure Modes to Consider:**
- User changes settings before uploading a file — settings change is a no-op, no error shown
- User changes encoding and the file becomes garbled — this is expected behavior, user can switch back. No magic.
