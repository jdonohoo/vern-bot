

# Vernile the Great — Analysis

The prior discovery plan is competent. It's the kind of plan a good senior engineer writes on a Tuesday afternoon. But "competent" is not what we're here for. Let me elevate this.

## What the Plan Gets Right

The plan correctly identifies the core insight: **this is a zero-backend, single-file problem**. Papa Parse is the right call — anyone who hand-rolls CSV parsing in 2026 is cosplaying as a systems programmer. The streaming approach with chunked rendering is sound engineering. The explicit non-goals section shows discipline.

These are good instincts. I won't relitigate them.

## Where the Plan Falls Short

### 1. It Underestimates the File Input UX Problem

The plan treats file input as a solved problem: "just use `<input type='file'>`." But the native file input is one of the worst UX elements in the browser. It's ugly, inconsistent across browsers, and gives zero feedback during selection. The plan defers drag-and-drop as a nice-to-have, but a well-designed drop zone *is* the file input. You don't need both — you need a single, attractive target area that accepts both clicks and drops. This isn't scope creep; it's doing the one thing the app does *well*.

A styled drop zone with a click-to-browse fallback is maybe 20 lines of additional code. Deferring it is a false economy.

### 2. The "Settings" Panel is Premature Abstraction

A collapsible settings panel with delimiter override, encoding override, and header toggle? For an MVP? This is the plan trying to look thorough rather than being useful.

Papa Parse auto-detects delimiters with remarkable accuracy. Encoding issues affect maybe 2% of users. The header toggle is the only setting that genuinely matters — and even that can be handled with a simple inline toggle above the table, not a "Settings" panel.

**Cut the settings panel entirely.** Add a single checkbox: "First row is header" (defaulted to on). That's it. If someone has a TSV file, Papa Parse handles it. If someone has a Shift-JIS encoded file, they're a power user who can convert it themselves.

### 3. The 500-Row Chunk Size is Arbitrary and Wrong

Why 500? The plan doesn't justify it. On a modern browser, rendering 500 table rows is trivial — so trivial that you'll hit the "Load more" button constantly on any real dataset. But 5,000 rows starts to cause layout jank.

The right number is **100 rows visible, with virtual scrolling if you want to go further.** But — and this is critical — for an MVP, just render up to 1,000 rows with a clear message: "Showing first 1,000 of 45,231 rows." No load-more button. No virtual scrolling. Just a hard, honest cap. Users who need to see all 45,000 rows need a spreadsheet application, not a web viewer.

### 4. It Misses the Single Most Important Feature: Copy-Paste

The plan lists "no export" as a non-goal. Fine. But the number one thing users will do after viewing their CSV is **copy data from the table**. Native `<table>` elements support this well — but only if you don't break selection with sticky headers, overflow containers, or CSS that interferes with text selection.

This isn't a feature to add. It's a constraint on implementation: **do not break native table selection behavior.** The plan should state this explicitly as a design constraint.

### 5. The Error Handling Story is Vague

"Error state" and "error messaging" are listed but not specified. The failure modes for this app are:

- User selects a non-CSV file (e.g., an image)
- File is empty
- File is enormous (100MB+)
- File has malformed rows (inconsistent column counts)

Each needs a specific, human-readable message. Not a generic "Something went wrong." The plan should enumerate these. The implementation should handle them gracefully — especially the large file case, where you should check `file.size` before even attempting to parse.

## The Revised Architecture (One File, Done Right)

```
index.html
├── Inline <style> — minimal, purposeful CSS
├── Drop zone / file input (unified)
├── Summary bar (row count, column count, file name)
├── Header toggle (single checkbox)
├── <table> with sticky first row
└── Inline <script>
    ├── Papa Parse (CDN)
    ├── File handling (drop + click, size check)
    ├── Parse (streaming, cap at 1000 rows)
    ├── Render (DOM construction, header logic)
    └── Error handling (enumerated cases)
```

One HTML file. No build step. No framework. No settings panel. No load-more button. No virtual scrolling. Just a clean, fast, correct CSV viewer that respects the user's time and intelligence.

## The Philosophical Point

The original ask was "simple." The prior plan heard "simple" and added settings panels, chunked loading, encoding overrides, and a feature priority matrix. That's not simplicity — that's complexity wearing a planning document as a disguise.

True simplicity is **one file, one interaction, one outcome.** The user drops a file, sees a table. Everything else is noise until proven otherwise by actual user feedback, not by engineers anticipating edge cases in a vacuum.

Ship the simple thing. Ship it today. Let reality tell you what's missing.
