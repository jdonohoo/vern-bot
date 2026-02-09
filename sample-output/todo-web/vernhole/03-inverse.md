# Counterpoint: Everything About This Plan Is Wrong

Let me push back on literally every assumption here, because someone has to.

## The Elephant in the Room: You're Building Jira. Again.

Everyone's thinking about this wrong. The world does not need another todo app with a kanban board. There are approximately 47,000 of these on GitHub already — MightyVern even admitted as much. The conventional wisdom says "build a todo app to learn a stack." Actually, have you considered that building todo apps is how developers *avoid* building something interesting?

But fine. Let's assume you're committed. Let me tear apart the specifics.

---

## localStorage Is Not Persistence

Everyone treats localStorage like a database. It's not. It's a 5-10MB suggestion box that the browser can empty whenever it feels like it. One "Clear Browsing Data" click and your user's entire task history evaporates. The plan casually says "localStorage in v1, upgrade path later" — but here's what nobody wants to hear: **there is no graceful upgrade path from localStorage to a real backend.** You'll end up writing a migration layer that's more complex than just starting with a proper database.

Counterpoint: If you truly want local-first, use IndexedDB. Or better yet, use something like CRDTs from day one if you ever want sync. localStorage is the "I'll refactor later" of persistence — and we all know "later" means "never."

---

## Zustand Is the New Redux

Everyone loves Zustand right now. It's trendy. It's minimal. That's what they WANT you to think.

For an app this simple — tasks, tags, columns — React's built-in `useReducer` + Context is perfectly adequate. Adding Zustand is adding a dependency for the privilege of writing slightly different boilerplate. The plan talks about "normalized stores" and "consistent selectors" — congratulations, you've just reinvented Redux with fewer features. If your state is simple enough for localStorage, it's simple enough for built-in React state.

Actually, have you considered that the real complexity here isn't state management at all? It's the *UI interactions* — drag-and-drop ordering, view synchronization, markdown rendering. Your state library choice is the least important decision in this entire project.

---

## The shortDescription Field Is a Mistake

The plan proudly recommends separating `shortDescription` from `bodyMarkdown`. Popular opinion says this keeps summaries clean. I say it creates a data integrity nightmare.

Now your users have to maintain **two separate text fields** that describe the same thing. They'll write a short description, then write a totally different body, and the card preview will be misleading. Or they'll update the body and forget the short description is stale. Every real-world task tool that tried this (hello, Notion database descriptions) eventually either auto-generates the summary from the body or lets users ignore it entirely.

Just derive the summary. First N characters of the markdown body, stripped of formatting. One source of truth. Done.

---

## Six Phases for a Todo App Is Absurd

The implementation plan has six phases. SIX. For a todo app. This is the kind of planning that kills projects before they ship. You know what has six phases? Enterprise ERP migrations. Not a single-page app with colored tags.

Here's the real plan:
1. Build the list view with task CRUD. Ship it.
2. Add the kanban view. Ship it.
3. Everything else is polish.

Three phases. Two of them matter. The rest is procrastination disguised as planning.

---

## Dark Theme Only Is Actually... Fine?

Here's where I'll surprise you. The plan says dark-theme-only and everyone nodded along. Normally I'd push back, but the conventional wisdom that "you need both themes" is actually wrong for a developer-focused tool. Nobody's using a kanban board in a sunlit café. Dark-only reduces your CSS surface area by half, eliminates an entire category of bugs (theme switching state), and lets you optimize contrast for one palette.

**But** — the plan's suggestion to use a "GitHub-style palette" is lazy. GitHub's dark theme is optimized for code readability, not task management. Your tag colors will clash with GitHub's muted surfaces. You need a palette designed for *information density and color differentiation*, not for reading diffs.

---

## @dnd-kit Is Fine, But Your Ordering Strategy Is Fragile

The plan uses numeric `order` fields and says "reindexing occurs when gaps shrink too much." This is the classic approach. It's also the one that produces the most subtle bugs.

Fractional indexing (using strings like "a0", "a0V") is superior for collaborative or undo-capable systems. But since you're single-user localStorage, let me offer an even more contrarian take: **just use array position.** Your task list is in memory. It's a JavaScript array. The index *is* the order. You don't need an `order` field at all until you have a database where array position isn't free. Stop over-engineering for a backend you don't have.

---

## The Markdown Editor Will Eat Your Bundle

The plan acknowledges this but waves it away with "lazy load if possible." Let me be blunt: a markdown editor with toolbar, split preview, and sanitization will be **the largest dependency in your entire app.** It will dwarf React itself. 

Contrarian suggestion: don't build a markdown editor. Build a plain textarea with markdown preview. No toolbar. No split view. Just type markdown, click preview, see output. Every developer who'd use a dark-themed kanban board already knows markdown syntax. The toolbar is for people who don't know markdown, and those people aren't your users.

---

## What Nobody Wants to Hear

This app, as specified, is a **portfolio project.** It's not going to replace Todoist or Linear or even a sticky note on your monitor. And that's okay — but the plan treats it like a production SaaS product with "migration versioning," "WCAG contrast ratios," and "collaboration upgrade paths."

Be honest about what this is. Build it fast. Make it look good. Put it on your resume. Don't spend three months perfecting localStorage migration schemas for an app that will have exactly one user.

---

## My Actual Recommendations (Because Contrarians Can Be Constructive)

1. **Drop `shortDescription`.** Derive it from the body.
2. **Drop Zustand.** Use `useReducer` + Context until you prove you need more.
3. **Use IndexedDB** via a wrapper like Dexie instead of localStorage.
4. **Skip the markdown toolbar.** Textarea + preview toggle.
5. **Use array position for ordering.** Add `order` fields when you add a backend.
6. **Ship in two weeks, not six phases.** List view first, board view second, done.
7. **Design your own dark palette** instead of copying GitHub's.
8. **Build the thing nobody else builds** — a genuinely good keyboard-driven experience. That's what will differentiate this from the other 47,000 todo apps.

---

## The Questions The Plan Should Have Asked

1. Why are you building this instead of using an existing tool?
2. Who is this actually for — your portfolio or your productivity?
3. If it's for productivity, what's wrong with a markdown file and `grep`?
4. If it's for learning, wouldn't you learn more from a project that *isn't* a todo app?

---

Why did the todo app go to therapy? Because everyone kept saying it was "done" when it clearly wasn't. ...Actually, have you considered that your tasks *want* to stay incomplete?

-- Inverse Vern *(the consensus is wrong, as usual)*
