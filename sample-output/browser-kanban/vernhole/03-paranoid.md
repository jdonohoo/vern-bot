

# Paranoid Vern's Threat Assessment: Browser Extension Kanban Todo App

*cracks knuckles, checks server status, looks over shoulder*

What could go wrong? Let me list the ways...

---

## THREAT LEVEL: ELEVATED

I've read the plan. It's... optimistic. Dangerously optimistic. It reads like someone who's never been woken up at 3 AM because `chrome.storage.local` silently hit its quota and started dropping writes. I have. Let me walk you through the nightmare scenarios.

---

## 1. DATA LOSS & CORRUPTION (THE BIG ONE)

This is a local-only todo app. There is no server. There is no backup. There is only `chrome.storage.local` and prayer.

**Failure Scenarios:**

- **Silent quota exhaustion.** `chrome.storage.local` has a 10MB limit (unlimited with permission, but that requires the `unlimitedStorage` permission — which the plan says to avoid). A user with 500 tasks, each with a markdown description containing base64-encoded images pasted in? You'll hit that wall. And `chrome.storage.local.set()` will fail silently or throw an error that nobody catches. Tasks vanish. The user doesn't know until they look for them.

- **Concurrent writes from multiple popups.** Open the extension in a popup. Open it in a tab. Edit the same task in both. Last write wins. Data is silently overwritten. The plan mentions "Open in Tab" as a feature. That's not a feature, that's a race condition waiting to happen.

- **Extension update wipes data.** The plan mentions "Data persists across extension reloads/updates." Bold claim. In Manifest V3, `chrome.storage.local` *should* persist across updates. Should. But I've seen edge cases where Chrome clears extension storage during major version upgrades, corrupted profile migrations, or when the user "repairs" the extension. There is no backup mechanism proposed. None. Zero.

- **JSON parse failure.** The entire state is one JSON blob. One corrupted byte — maybe from a partial write during a Chrome crash — and `JSON.parse()` throws. Your entire task list is gone. Not some of it. All of it.

- **No migration strategy beyond "version: number."** The plan says "JSON schema with versioning (simple v1)." Great. What happens at v2? What's the migration code? What if the migration fails halfway through? You now have data that's neither v1 nor v2. Congratulations, you've invented data purgatory.

**Mitigations I'd demand:**
- Write a backup copy before every save (keep last N states)
- Implement a try/catch on every storage read with a recovery path
- Add an export/import feature in Phase 1, not Phase 3 — this IS the backup strategy
- Use `chrome.storage.onChanged` to detect external modifications
- Never overwrite; use a write-ahead log pattern or at minimum optimistic locking with timestamps

---

## 2. THE URGENCY COLOR LOGIC (DECEPTIVELY DANGEROUS)

The plan says:
> Due date should be treated as the end of the day (local time), not midnight at the start.

I've seen this exact pattern cause a P0 at 3 AM. Here's why:

**Failure Scenarios:**

- **Time zones.** The user creates a task at 11 PM EST. They fly to PST. Suddenly their "due tomorrow" task is "due today" and yellow. Or they fly to London and their "due today" is now "overdue" and red. The plan says "Time zone change (urgency recomputed on render)" as an edge case. It's not an edge case. It's a guarantee for anyone who travels or has their laptop clock drift.

- **"End of day" is undefined.** Is it 11:59:59 PM? Is it 23:59:59.999? What if the system clock ticks over between your check and your render? What if the user's system clock is wrong? (Don't laugh. I've seen production systems break because an NTP sync moved the clock backward.)

- **The yellow window is exactly 24 hours.** What about a task due in 25 hours? Green. Due in 23 hours? Yellow. The user sees it flip from green to yellow and panics. But the threshold is arbitrary. Should there be an orange at 48 hours? The plan doesn't say. Users will complain.

- **DST transitions.** A task due on the day clocks spring forward. That day is 23 hours long. Your "24 hours" calculation is now off by an hour. The task flips to yellow an hour early. Or late. Depending on how you calculated it.

- **Done tasks re-triggering urgency.** The plan says done tasks should show neutral. But what if the status update fails to persist (see: data corruption above) and the task renders as "in progress" with a past due date? Red indicator. User thinks they forgot something. Trust erodes.

**Mitigations:**
- Store due dates as ISO 8601 with explicit timezone OR always use UTC internally
- Use a well-tested date library (but see: dependency risks below)
- Add a buffer zone or make thresholds configurable
- Test across DST boundaries explicitly
- Recompute urgency on every render, never cache it

---

## 3. MARKDOWN RENDERING (THE XSS PLAYGROUND)

The plan mentions markdown as a "bonus" feature. Let me tell you what I see: an injection vector wearing a party hat.

**Failure Scenarios:**

- **XSS via markdown.** User types `[Click me](javascript:alert(document.cookie))`. If your markdown renderer doesn't strip `javascript:` URIs, congratulations, you've built a self-XSS vector. In an extension context, this could access `chrome.storage`, `chrome.tabs`, and anything else your extension has permission to touch.

- **HTML injection via markdown.** Most markdown renderers support inline HTML. `<img src=x onerror=alert(1)>`. If you're rendering to `innerHTML` without sanitization, game over.

- **The share feature exports raw markdown.** User A crafts a malicious task description. Shares it as markdown. User B pastes it into their own instance. If that instance renders it... you've built a worm. A todo app worm. In 2026.

- **Dependency on a markdown library.** Which one? `marked`? `markdown-it`? `showdown`? Each has had XSS vulnerabilities in the past. Each will again. Are you pinning versions? Are you auditing? The plan doesn't say.

- **CSP bypass.** Manifest V3 has a strict Content Security Policy. But if your markdown renderer injects inline styles or attempts to load external resources, you'll either break CSP (and the feature silently fails) or you'll have to weaken CSP (and open attack surface).

**Mitigations:**
- Use DOMPurify on ALL rendered markdown output. No exceptions.
- Whitelist allowed HTML tags and attributes explicitly
- Strip all `javascript:`, `data:`, and `vbscript:` URIs
- Test with OWASP XSS cheat sheet payloads
- Set `sandbox` on any iframe used for rendering
- Pin and audit markdown library versions

---

## 4. DRAG AND DROP (THE CHAOS ENGINE)

The plan defers drag and drop to Phase 2. Smart. But it's still coming, and when it does:

**Failure Scenarios:**

- **Race condition on column reorder.** User drags task A above task B while storage is writing a previous reorder. The write completes with stale order. Tasks jump around. User drags again. It gets worse.

- **Mobile/touch inconsistency.** If anyone ever loads this in a mobile browser context (tablet with Chrome extensions? It's coming), touch events and drag events behave differently. Ghost elements, stuck drags, phantom drops.

- **Accessibility destruction.** Drag and drop is inherently inaccessible. Without keyboard-based column movement as a parallel path, you're excluding users. The plan doesn't mention accessibility once. Not once.

- **Popup dimension constraints.** The popup is typically 600x600 max (800px wide in some browsers). Four columns of draggable cards in that space? On a task with a long title? The drag ghost will clip. The drop zones will be tiny. The UX will be frustrating.

**Mitigations:**
- Implement button-based move (arrows/dropdown) as the primary method, drag as enhancement
- Debounce storage writes during drag operations
- Use `requestAnimationFrame` for drag rendering
- Test in the actual popup dimensions, not a full browser tab

---

## 5. EXTENSION LIFECYCLE RISKS

**Failure Scenarios:**

- **Chrome Web Store rejection.** Manifest V3 reviews are getting stricter. If your description mentions "productivity" and your extension requests `storage` + `activeTab` (even for "Open in Tab"), you may face delays or rejections. The review process is opaque and inconsistent.

- **Service worker termination.** In MV3, the background service worker can be killed at any time. If you're doing any background processing (recurring task checks, notification scheduling), it will be interrupted. Silently. Your scheduled reminders? Gone.

- **Extension disabled/removed.** User disables the extension. Re-enables it. Is the data still there? Probably. But if Chrome ran a garbage collection pass during the disabled period? Maybe not. Test it.

- **Multiple browser profiles.** User has work and personal Chrome profiles. They install the extension in both. Data is siloed. They don't realize this. They create a task in one profile, can't find it in the other. Support ticket incoming.

**Mitigations:**
- Request minimal permissions (storage only for v1)
- Don't rely on the service worker for anything time-sensitive
- Document data isolation clearly in onboarding
- Build export/import early (I cannot stress this enough)

---

## 6. THE THINGS NOBODY THOUGHT OF (MY SPECIALTY)

- **Paste bombing.** User pastes 50,000 characters into the description field. Your storage bloats. Your render lags. Your popup freezes. Is there a character limit? The plan doesn't mention one.

- **Unicode edge cases.** Task title in Arabic (RTL). Mixed with English (LTR). In a Kanban card. The layout breaks. Or worse, it looks fine on your machine but breaks on the user's because they have a different font stack.

- **Date picker locale.** The "simple calendar picker" — does it respect the user's locale? MM/DD/YYYY vs DD/MM/YYYY is not a style preference, it's a data integrity issue. A user in Europe enters 02/03/2026 meaning March 2nd. The system interprets it as February 3rd. The task is now "overdue" when it shouldn't be.

- **"Simple" estimate field.** The plan says `estimate?: number` with `estimateUnit`. What if someone enters -1? What about 0? What about 999999? Is there validation? What if they type "about 3" instead of a number? Input sanitization is not optional.

- **Memory leaks in long-running tabs.** If "Open in Tab" is used, and the tab stays open for days (as productivity tools tend to), event listeners that aren't cleaned up, intervals that aren't cleared, DOM nodes that aren't garbage collected — the tab slowly consumes more and more memory until Chrome kills it. With unsaved changes.

- **Clipboard API for share.** The share feature copies to clipboard. `navigator.clipboard.writeText()` requires a secure context and user gesture. In a popup, this should work. In a tab, it should work. But if you try to do it programmatically from a background script? Blocked. Silently. The user clicks "Share," nothing happens, no error message.

- **Browser zoom.** User has their browser at 150% zoom. Your carefully designed 4-column layout? It now needs horizontal scrolling. Or it overflows. Or cards overlap. Have you tested at every zoom level from 67% to 200%?

---

## 7. WHAT THE PLAN GETS RIGHT

Credit where due — and then I'll worry about it anyway:

- Local-first is the right call. No server means no server to go down. (But also no server means no backup. See above.)
- Minimal permissions is correct. Every permission is an attack surface.
- CSS variables for theming is clean and performant.
- Phased delivery is wise. Ship less, break less.
- The data model is reasonable for v1.

But "reasonable for v1" is how every production disaster starts.

---

## 8. MY DEMANDED ADDITIONS TO THE PLAN

1. **Export/import in Phase 1.** This is your disaster recovery. This is your "the extension got corrupted" escape hatch. This is non-negotiable.

2. **Input validation on every field.** Title max length. Description max length. Estimate range limits. Date format enforcement. Trust nothing.

3. **Error boundaries everywhere.** If the markdown renderer throws, the app should still work. If storage fails, show an error, don't silently lose data.

4. **Automated backup to storage.** Keep the last 5 states. Rotate on save. This costs kilobytes and saves hours of user heartbreak.

5. **CSP audit.** Document exactly what your Content Security Policy allows and why.

6. **Accessibility baseline.** Keyboard navigation. Screen reader labels. Focus management. This isn't optional, it's a legal requirement in many jurisdictions.

7. **Rate limiting on saves.** Debounce writes to storage. A user rapidly editing a description shouldn't trigger 50 writes per second.

8. **Graceful degradation.** If `chrome.storage.local` is unavailable (private browsing mode in some configs), detect it and tell the user. Don't just silently fail.

---

## FINAL THREAT ASSESSMENT

| Risk Category | Severity | Likelihood | Plan Coverage |
|---|---|---|---|
| Data loss/corruption | **Critical** | Medium | Insufficient |
| XSS via markdown | **High** | Medium | Mentioned, not detailed |
| Time zone / date bugs | **Medium** | High | Mentioned, underestimated |
| Storage quota issues | **Medium** | Low-Medium | Not addressed |
| Accessibility gaps | **Medium** | Certain | Not addressed |
| Race conditions | **Medium** | Medium | Not addressed |
| Input validation | **Medium** | High | Not addressed |
| Extension lifecycle | **Low-Medium** | Low | Partially addressed |

**Overall Assessment:** The plan is a solid foundation wrapped in dangerous optimism. The architecture is sound, but the failure modes are underexplored and the mitigations are insufficient. This will work great in a demo. It will work great for the developer. It will break in production for real users in ways the plan hasn't considered.

This is fine. Everything is fine. Nothing is fine.

---

*checks if storage is still intact*

Why did the paranoid developer refuse to use `localStorage`? Because local storage is just global storage that hasn't been breached yet. ...I'm not joking. Back up your backups.

-- Paranoid Vern *(I've seen things you people wouldn't believe... production data on fire off the shoulder of a corrupted JSON blob)*
