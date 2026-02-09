

# MightyVern Analysis: CSV Upload & Table Display

## Executive Assessment

The discovery plan is **solid but conservative**. It makes the right call on Papa Parse over hand-rolled parsing — that's non-negotiable. But let me push harder on a few dimensions.

---

## What the Plan Gets Right

**Papa Parse with streaming** — Correct. Anyone who thinks `line.split(',')` handles CSV is someone who hasn't encountered a quoted field with an embedded comma, a newline inside quotes, or a BOM character. Papa Parse handles all of this. Ship with it, don't look back.

**No backend** — Absolutely right. A CSV viewer has zero reason to send user data to a server. `FileReader` + client-side parsing keeps it private, fast, and deployable anywhere (static host, S3, GitHub Pages). No CORS, no server costs, no data liability.

**Chunked rendering at 500 rows** — Sensible default. The DOM chokes well before you'd expect. 500 rows × 20 columns = 10,000 cells. That's fine. 50,000 rows × 20 columns = 1,000,000 cells. That kills the tab.

---

## Where I'd Push Back or Strengthen

### 1. The 500-Row Chunk Size Should Be Dynamic

Don't hardcode 500. Compute it from column count:

```
chunkSize = Math.floor(10000 / columnCount)
```

A 3-column CSV can comfortably render 3,000 rows. A 50-column CSV should cap at 200. The plan treats all CSVs the same — they aren't.

### 2. Papa Parse Config Matters More Than the Plan Acknowledges

The specific config is where bugs hide:

```javascript
Papa.parse(file, {
  header: true,          // toggleable
  dynamicTyping: false,  // keep everything as strings — display tool, not analytics
  skipEmptyLines: 'greedy',
  worker: true,          // parse off main thread — this is critical
  chunk: function(results, parser) { /* render incrementally */ },
  complete: function() { /* finalize */ },
  error: function(err) { /* surface to user */ }
});
```

Key decisions:
- **`worker: true`** — The plan mentions streaming but doesn't mention Web Workers. A 100MB CSV parsed on the main thread freezes the UI. Papa Parse has built-in worker support. Use it. This is the single most important performance decision.
- **`dynamicTyping: false`** — This is a display tool. Don't coerce `"001"` to `1`. Don't turn `"true"` into a boolean. Show what's in the file.
- **`skipEmptyLines: 'greedy'`** — Strips whitespace-only lines too. Users don't want to see blank rows.

### 3. The HTML Structure Should Be Minimal but Complete

The entire app is one file. Here's the skeleton that matters:

```
index.html
├── <input type="file" accept=".csv,.tsv,.txt">
├── <div id="status">        <!-- "Parsing...", "Showing X of Y rows", errors -->
├── <div id="table-wrapper">  <!-- overflow-x: auto for wide tables -->
│   └── <table>
│       ├── <thead>            <!-- sticky via position: sticky -->
│       └── <tbody>
├── <button id="load-more">   <!-- hidden until needed -->
└── <details id="settings">   <!-- collapsed settings panel -->
```

No framework. No build step. One HTML file with inline CSS and JS. If someone can't open it by double-clicking, you over-engineered it.

### 4. Error Handling Needs Specificity

The plan says "error messaging." Here's what actually goes wrong:

| Error | Message |
|-------|---------|
| No file selected | "Select a CSV file to get started." |
| Wrong file type | "This doesn't look like a CSV file. Expected .csv, .tsv, or .txt." |
| Empty file | "This file is empty." |
| Parse errors (malformed rows) | "Parsed N rows. X rows had issues (shown in red)." |
| File too large (>500MB) | "This file is very large. Display may be slow." — warning, not a block |

Don't silently drop malformed rows. Render them with a visual indicator. Users need to know their data has issues.

### 5. Settings: Simpler Than Proposed

The plan lists three settings. I'd keep two and change one:

- **"First row is header"** — Yes, as a checkbox, default ON.
- **Delimiter** — Auto-detect is Papa Parse's default and it's good. Show detected delimiter as read-only text ("Detected: comma"). Only offer override if auto-detect fails.
- **Encoding override** — Cut it from v1. UTF-8 covers 99% of cases. If someone has a Shift-JIS CSV, they know enough to convert it. Don't add UI for the 1% case.

### 6. One Enhancement Worth Adding: Column Count + Row Count Summary

Immediately after parsing, show:

```
mydata.csv — 12,847 rows × 23 columns (showing first 434 rows)
```

This is the single most useful piece of metadata. Users need to know the shape of their data before they look at the table. It costs nothing to implement (Papa Parse gives you this) and it's the first thing anyone asks: "how many rows?"

---

## Architecture in One Diagram

```
[File Input] → [FileReader] → [Papa Parse (Web Worker)]
                                    │
                              [chunk callback]
                                    │
                         [Buffer rows in memory]
                                    │
                    [Render first N to <table>]
                                    │
                    [Show summary + Load More]
```

All client-side. No network calls. No dependencies beyond Papa Parse (loaded from CDN or inlined).

---

## What I'd Cut From the Plan

- **Zebra striping** — CSS `tr:nth-child(even)` is one line. Fine, but don't call it a feature. It's a style rule.
- **Drag-and-drop as "deferred"** — It's 15 lines of JS. If you're going to mention it, just include it. But don't prioritize it over Web Worker parsing.

## What I'd Add to the Plan

- **Web Worker parsing** (`worker: true`) — Non-negotiable for files over 10MB.
- **`dynamicTyping: false`** — Explicit decision to preserve raw data.
- **Dynamic chunk sizing** based on column count.
- **File name + dimensions in the summary bar.**
- **Malformed row highlighting** instead of silent drops.

---

## Final Verdict

The plan is 80% right. The 20% gap is in **parsing configuration details** (Web Workers, dynamic typing, chunk sizing) — the stuff that separates "works on my test file" from "works on the weird CSV my coworker exported from SAP." Get the Papa Parse config right, render in chunks, keep it to one file, ship it.
