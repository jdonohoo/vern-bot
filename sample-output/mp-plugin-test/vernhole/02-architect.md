

# Architect Vern — Blueprint Review

Let me lay out how I'd actually structure this thing. The discovery plan is solid directionally, but it's making some choices I want to challenge and a few it's glossing over.

---

## Where the Plan Gets It Right

- **No backend.** Correct. A static single-page app is the right architecture here. Zero deployment complexity, zero server cost, zero auth surface.
- **Papa Parse over hand-rolled parsing.** Non-negotiable. CSV is a deceptively nasty format (quoted fields, embedded newlines, BOM markers). Don't reinvent this.
- **Chunked rendering.** The DOM is your bottleneck, not parsing. Capping visible rows is the right call.

---

## Where I'd Push Back

### 1. "Single HTML file" is underspecified

The plan says "single-page" but doesn't commit to a structure. Here's the decision:

**Option A — Single `index.html` with inline `<script>` and `<style>`**
- One file. CDN-load Papa Parse.
- Dead simple to host (GitHub Pages, S3, literally anywhere).
- No build step.

**Option B — Minimal project with separate files**
- `index.html`, `style.css`, `app.js`
- Still no build step, still static hosting.
- Easier to maintain if you ever touch it again.

**My call: Option A for this scope.** The entire app logic fits in ~150 lines of JS. Splitting it into files adds structure without adding clarity. One file means one thing to deploy, one thing to debug, one thing to share.

### 2. The 500-row default is arbitrary and the "Load More" pattern is wrong

500 rows into a `<table>` is fine for the DOM. But "Load 500 more" is a pagination UX bolted onto what should be a viewport problem. If someone uploads 50,000 rows, clicking "load more" 100 times is not a real interaction.

**Better approach:**
- Render first **100 rows** as the initial view (fast, scannable).
- Display total row count: `"Showing 100 of 48,291 rows."`
- One button: `"Show all"` — with a warning if row count exceeds 10,000.
- No incremental loading. Either you're previewing or you're looking at everything. Two states, not N states.

This is simpler to build, simpler to reason about, and honest about what it can handle.

### 3. The "Settings" panel is scope creep disguised as a collapsible div

The plan lists: header toggle, delimiter override, encoding override. For a "simple CSV viewer," delimiter and encoding overrides are solving problems that don't exist yet.

**What to actually ship:**
- **Header toggle:** Yes. This is the one setting that matters. Default to "first row is header" = true. One checkbox above the table.
- **Delimiter override:** No. Papa Parse auto-detects. If auto-detection fails on someone's weird file, that's an edge case you handle when a user reports it, not before.
- **Encoding override:** No. UTF-8 covers >99% of real-world CSVs. `FileReader.readAsText()` defaults to UTF-8. Done.

One checkbox. Not a settings panel.

### 4. Error handling needs a concrete contract

The plan says "error state" but doesn't define what errors are possible. Here's the actual list:

| Error | Cause | Response |
|---|---|---|
| No file selected | User clicks away from picker | No-op, keep current state |
| Wrong file type | User picks a `.png` | Show: "Please select a .csv file" |
| Empty file | 0 bytes | Show: "This file is empty" |
| Parse failure | Malformed CSV | Show Papa Parse error message verbatim |
| Enormous file | >50MB | Show warning before parsing, let user proceed |

That's it. Five cases. Handle them explicitly.

---

## The Actual Blueprint

```
┌─────────────────────────────────────┐
│  CSV Viewer                         │
│  Upload a CSV file to view its      │
│  contents.                          │
│                                     │
│  [Choose File]                      │
│                                     │
│  ☐ First row is header              │
│                                     │
│  ┌─ Results ──────────────────────┐ │
│  │ Showing 100 of 12,408 rows    │ │
│  │ [Show all]                     │ │
│  │                                │ │
│  │ ┌────────────────────────────┐ │ │
│  │ │ Name │ Email │ Status │... │ │ │  ← sticky header
│  │ ├────────────────────────────┤ │ │
│  │ │ ...  │ ...   │ ...    │   │ │ │  ← zebra-striped
│  │ └────────────────────────────┘ │ │
│  └────────────────────────────────┘ │
└─────────────────────────────────────┘
```

### Component Breakdown

```
index.html (single file)
├── <style>        — ~60 lines: table, layout, states
├── <div#app>      — static shell: title, file input, checkbox, results area
└── <script>
    ├── State      — { file, rows[], headers[], showAll, error }
    ├── handleFile()    — FileReader → Papa Parse (streaming) → state
    ├── render()        — state → DOM (clear + rebuild results area)
    └── Error handling  — 5 explicit cases per the table above
```

### Data Flow

```
[File Input] → FileReader.readAsText()
            → Papa.parse(text, { header: false })
            → rows[] (all parsed in memory)
            → render(rows.slice(0, 100))  or  render(rows)
```

Note: I'd parse the full file into memory, not stream. For the file sizes this tool will realistically handle (<50MB), full parse is simpler and Papa Parse handles it in under a second. Streaming adds complexity for a performance gain nobody will notice.

### Key Implementation Decisions

1. **Papa Parse via CDN** — `<script src="https://cdnjs.cloudflare.com/ajax/libs/PapaParse/5.4.1/papaparse.min.js">` — no install, no build.
2. **Vanilla JS, no framework** — This is a 150-line app. React/Vue/Svelte would triple the complexity for zero benefit.
3. **`<table>` element** — Semantically correct, accessible, copy-pasteable. The plan got this right.
4. **`table-layout: fixed` + `overflow-x: auto`** — Handles wide CSVs without breaking the layout.
5. **Sticky `<thead>`** via `position: sticky; top: 0` — The one CSS enhancement that actually matters.

---

## What I'd Cut From the Discovery Plan

- ~~Drag-and-drop (deferred)~~ — Don't even list it. It's not on the roadmap. It's noise.
- ~~Delimiter override~~ — Papa Parse handles it.
- ~~Encoding override~~ — UTF-8 or nothing at v1.
- ~~Streaming parse~~ — Unnecessary complexity for this scale.
- ~~"Load 500 more" incremental pattern~~ — Replaced with preview/show-all toggle.

## What I'd Add

- **File size check before parsing** — `if (file.size > 50 * 1024 * 1024)` warn the user. Costs 2 lines, prevents browser lockups.
- **`<table>` wrapper with `max-height` and `overflow-y: auto`** — Keeps the page scrollable even with 10k+ rows visible.
- **Column count in the summary** — `"12,408 rows × 7 columns"` — useful context, free to compute.

---

## Final Verdict

The discovery plan is 80% right but over-engineers the settings and under-specifies the error handling and rendering strategy. The architecture is one static HTML file, ~200 lines total, zero dependencies beyond a CDN-loaded Papa Parse. Ship it in an afternoon, host it anywhere.

That's the blueprint. Clean, minimal, no ambiguity about what to build or how it fits together.
