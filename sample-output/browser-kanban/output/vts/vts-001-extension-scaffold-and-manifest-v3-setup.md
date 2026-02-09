---
id: VTS-001
title: "Extension Scaffold and Manifest V3 Setup"
complexity: S
status: pending
owner: ""
source: oracle
source_ref: "oracle-architect-breakdown.md"
dependencies: []
files:
  - "`manifest.json`"
  - "`popup.html`"
  - "`styles/reset.css`"
  - "`styles/main.css`"
  - "`scripts/main.js`"
  - "`assets/icon-16.png`"
  - "`assets/icon-48.png`"
  - "`assets/icon-128.png`"
---

# Extension Scaffold and Manifest V3 Setup

Create the foundational browser extension structure with Manifest V3. This includes the manifest file, popup HTML entry point, basic CSS reset, and the JS entry point. No functionality yet — just the skeleton that loads, opens a popup, and proves the extension lifecycle works. This is the foundation everything else builds on, so get it right.

## Criteria

- `manifest.json` with Manifest V3 format, `storage` permission only
- Popup HTML loads when clicking the extension icon
- Basic CSS reset and root font sizing applied
- JS entry point loads without errors
- Extension installs in Chrome without warnings
- Folder structure follows separation of concerns (`styles/`, `scripts/`, `assets/`)
