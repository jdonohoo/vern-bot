### TASK 1: Project Scaffolding & HTML Shell

**Description:** Create the single `index.html` file that serves as the entire application. Set up the semantic HTML structure: page header with title/description, file input area, collapsible settings panel (hidden by default), results area (empty state), and error display region. Include the Papa Parse CDN link. All CSS will be in a `<style>` block, all JS in a `<script>` block — single file, zero build tools.

**Acceptance Criteria:**
- Single `index.html` loads in any modern browser with no server
- Page displays title, description, file input button, and an empty state message ("Upload a CSV to get started")
- Error display region exists but is hidden
- Results area exists but is hidden
- Papa Parse is loaded via CDN and `Papa` is available on `window`
- Page is responsive enough to not look broken on mobile

**Complexity:** S

**Dependencies:** None — this is the foundation.
