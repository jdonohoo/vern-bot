

# VernHole Synthesis: CSV Upload & Table Display

Five Verns walked into a bar with the same CSV file. Here's what came out.

---

## Unanimous Agreement (Rare — Pay Attention)

Every single Vern converged on these points. When Retro Vern and Architect Vern agree, it's law:

1. **No backend.** Client-side only. Zero server cost, zero data liability, zero deployment complexity.
2. **Papa Parse, not hand-rolled parsing.** CSV is a deceptively nasty format. Every Vern who's tried `split(',')` has regretted it.
3. **Single HTML file, no build step, no framework.** Vanilla JS. One `<script>` tag for Papa Parse via CDN. Double-click to open.
4. **`<table>` element with sticky header.** Semantically correct, accessible, copy-pasteable.
5. **Kill the settings panel.** One inline checkbox ("First row is header") is enough. Delimiter and encoding overrides are premature.

---

## Interesting Contradictions

| Topic | Retro Vern | Architect Vern | MightyVern | Vernile | UX Vern |
|---|---|---|---|---|---|
| **Row limit** | 1,000, hard cap | 100 preview + "Show all" | Dynamic based on column count | 1,000, hard cap | Chunked + load more is fine |
| **Streaming/Workers** | "Just parse the whole thing" | "Parse fully, slice to render" | "Web Workers non-negotiable" | "Streaming, cap at 1000" | Doesn't care, wants a spinner |
| **Drag-and-drop** | Don't even mention it | Don't mention it | 15 lines, just do it | The drop zone IS the UI | Giant drop zone, center of page |
| **Load More button** | Simple for loop, fine | Replace with preview/show-all toggle | Keep it | No button, hard cap | Fine for v1 |

### Resolution:

- **Row limit:** Parse the full file into memory (Retro/Architect are right — Papa Parse handles this fast). Render first **100 rows** with a clear summary. One "Show all" button with a warning above 10k rows. MightyVern's dynamic sizing is clever but over-engineered for v1.
- **Web Workers:** MightyVern is technically correct that `worker: true` prevents UI freezes on large files, but Retro Vern is right that most CSVs are small. **Add a file size check first** — if under 10MB, parse directly. If over, warn the user. Workers can come later.
- **Drag-and-drop:** UX Vern wins this one. A styled drop zone that also accepts clicks is ~20 lines of code and dramatically improves first impressions. It's not scope creep — it IS the interface. Retro Vern, I know you hate it, but the native `<input type="file">` looks like 1997 and you know it.
- **Load More:** Architect's "preview + show all" two-state model is the cleanest. No incremental pagination.

---

## Hidden Gems (Things Only One Vern Caught)

- **MightyVern: `dynamicTyping: false`** — Critical. This is a display tool. Don't coerce `"001"` to `1` or `"true"` to a boolean. Show what's in the file.
- **Vernile: Don't break native table selection.** Copy-paste is the #1 thing users will do after viewing. This is a design constraint, not a feature.
- **UX Vern: "Your file never leaves your browser."** Six words. Massive trust boost. Costs nothing.
- **Architect Vern: Column count in summary.** `"mydata.csv — 12,847 rows x 23 columns"` — most useful metadata, free to compute.
- **UX Vern: Inline delimiter fix.** If auto-detect fails (everything in one column), show quick-pick buttons (comma, semicolon, tab) right above the table instead of burying it in settings.

---

## The Synthesized Design

```
┌─────────────────────────────────────────────┐
│                                             │
│     ┌─────────────────────────────────┐     │
│     │                                 │     │
│     │   Drop a CSV file here,         │     │
│     │   or click to browse            │     │
│     │                                 │     │
│     └─────────────────────────────────┘     │
│     Your file never leaves your browser.    │
│                                             │
│  ☑ First row is header                      │
│                                             │
│  mydata.csv — 12,847 rows × 7 columns       │
│  Showing first 100 rows. [Show all]         │
│                                             │
│  ┌──────────────────────────────────────┐   │
│  │ Name  │ Email        │ Status │ ...  │   │  ← sticky
│  ├──────────────────────────────────────┤   │
│  │ ...   │ ...          │ ...    │      │   │  ← zebra
│  │ ...   │ ...          │ ...    │      │   │
│  └──────────────────────────────────────┘   │
└─────────────────────────────────────────────┘
```

### Error States (Architect + UX Vern aligned)

| Scenario | Message |
|---|---|
| Wrong file type | "Please select a CSV file (.csv, .tsv, or .txt)" |
| Empty file | "This file is empty." |
| File > 50MB | "This is a large file. Parsing may take a moment." (warning, not a block) |
| Malformed rows | "Parsed 12,847 rows. 3 rows had issues." (highlight them) |
| Single column detected | "Everything ended up in one column. Try a different delimiter:" [,] [;] [tab] [|] |

### Implementation: One File, ~200 Lines

```
index.html
├── <style>           ~60 lines (drop zone, table, states)
├── <div id="app">    Drop zone, checkbox, summary, table container
└── <script>
    ├── Papa Parse via CDN
    ├── Drop zone + click handler (unified)
    ├── File size check before parsing
    ├── Papa.parse({ header: checkbox, dynamicTyping: false, skipEmptyLines: 'greedy' })
    ├── Store all rows in memory, render slice
    ├── "Show all" button (with >10k warning)
    └── 5 enumerated error cases
```

---

## Recommended Path Forward

1. **Build the single HTML file** matching the design above. It's ~200 lines, zero dependencies beyond a CDN script tag.
2. **Ship it.** GitHub Pages, S3, email it, double-click it — it works everywhere.
3. **Wait for real feedback** before adding anything. Not theoretical edge cases. Real humans with real CSVs telling you what's missing.

The Verns disagree on details but agree on philosophy: **the user said "simple," so build simple.** One file. One interaction. One outcome. Everything else is noise until proven otherwise.
