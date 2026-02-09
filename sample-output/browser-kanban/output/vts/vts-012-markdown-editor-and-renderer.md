---
id: VTS-012
title: "Markdown Editor and Renderer"
complexity: L
status: pending
owner: ""
source: oracle
source_ref: "oracle-architect-breakdown.md"
dependencies:
  - VTS-009
files:
  - "`scripts/components/markdown-editor.js`"
  - "`scripts/utils/markdown-renderer.js`"
  - "`styles/markdown.css`"
  - "`lib/marked.min.js`"
  - "`lib/purify.min.js`"
---

# Markdown Editor and Renderer

Add optional markdown support for the task description field. When enabled in preferences, the description textarea becomes a markdown editor with a live preview toggle. Use a lightweight markdown library bundled with the extension. **Critical: sanitize all rendered HTML.** Use DOMPurify or equivalent. Markdown XSS in a browser extension is a particularly bad day.

## Criteria

- Markdown toggle in preferences (default: off)
- When enabled: description field shows "Edit" and "Preview" tabs
- Edit mode: plain textarea with markdown syntax
- Preview mode: rendered HTML from markdown
- All rendered HTML is sanitized (no script execution, no event handlers)
- Links in rendered markdown open in new tab with `rel="noopener noreferrer"`
- Basic markdown supported: headings, bold, italic, lists, links, code blocks, blockquotes
- Graceful fallback: if markdown parsing fails, show raw text
- Libraries bundled with extension (no CDN, no external fetches)
