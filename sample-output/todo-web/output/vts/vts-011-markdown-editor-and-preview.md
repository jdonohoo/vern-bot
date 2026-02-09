---
id: VTS-011
title: "Markdown Editing and Preview"
complexity: M
status: pending
owner: ""
source: discovery
source_ref: "07-architect-architect-breakdown.md"
dependencies:
  - VTS-008a
files:
  - "src/components/detail/MarkdownEditor.tsx"
  - "src/utils/markdown.ts"
---

# Markdown Editing and Preview

Build the MarkdownEditor component using a plain textarea for editing and react-markdown + remark-gfm for preview rendering. Toggle between edit and preview modes (no split view). Sanitize preview output with DOMPurify. Lazy-load react-markdown for bundle size.

**Approach:** Textarea + react-markdown. Do NOT use @uiw/react-md-editor or any rich/WYSIWYG editor. The textarea is the editing surface. The preview is read-only rendered markdown. This keeps the bundle small (~35KB gzipped) and the implementation simple.

## Criteria

- Textarea for markdown input (plain text editing)
- Preview panel renders sanitized HTML via react-markdown + remark-gfm
- Toggle between edit mode and preview mode (not split view)
- **No rich editor toolbar** -- users write raw markdown in the textarea
- Preview is debounced (300ms)
- Output sanitized with DOMPurify (no XSS)
- External links get `target="_blank" rel="noopener noreferrer"`
- Empty state shows placeholder text ("Add details...")
- react-markdown chunk is lazy-loaded (React.lazy + Suspense)
- Lazy-load failure shows a graceful fallback (raw markdown text), not a crash
- Pasted rich text is stripped to plain text
