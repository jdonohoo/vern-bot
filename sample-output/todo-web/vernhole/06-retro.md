# Retro Vern's Take: The Dark-Themed Todo App

*pulls up chair, adjusts bifocals, sips coffee from a mug that says "I survived CORBA"*

Alright. Let me look at this.

## What Is This, Really?

Strip away the buzzwords — TypeScript, React, Vite, Zustand, Vercel, @dnd-kit — and what do you have?

**A todo list.**

A todo list with columns. We called those "kanban boards" after Toyota popularized them in the *1950s*. Before Trello existed, before Jira consumed our souls, people were moving index cards between columns on a corkboard. The software version of this has been built approximately fourteen million times since 2005.

And you know what? That's fine. Todo apps are the "Hello World" of real applications. Everyone builds one. Most of them teach you something. Just... let's be honest about what we're building.

## The Prior Plan: A Reasonable Document With Too Many Meetings

MightyVern's plan up there is *thorough*. I'll give it that. It reads like someone who's built this before and wants to make sure you don't miss anything. Good instincts.

But let me point out a few things that made my eye twitch.

### The Data Model Is Overthought

```
Task: id, title, shortDescription, bodyMarkdown, statusId, order, 
      tagIds, startDate, dueDate, completedAt, createdAt, updatedAt
```

You know what a todo item looked like in 2004?

```
id, title, done, category, due_date, notes
```

And it worked. For years. For millions of users.

Now, I'm not saying the proposed model is *wrong* — `shortDescription` separate from `bodyMarkdown` is actually a reasonable call if you're showing card previews. But `startDate` AND `dueDate` AND `completedAt` AND `createdAt` AND `updatedAt`? That's five date fields on a todo item. You're building a todo app, not Microsoft Project.

**My advice:** Start with `title`, `body` (markdown), `status`, `order`, `tagIds`, `dueDate`, `done`, `createdAt`. Add fields when you actually need them, not when you might theoretically need them. YAGNI has been right more often than it's been wrong.

### The Technology Stack

Let me walk through this:

- **TypeScript/React/Vite** — Fine. React is the new jQuery. It's everywhere, everyone knows it, it'll be maintainable. Vite is genuinely better than webpack. I'll allow it.

- **Zustand** — A state management library for... a single-user todo app that stores everything in localStorage. You know what else manages state? `useState` and `useReducer`. They ship with React. They're free. They work. Zustand is a fine library, but you're adding a dependency to manage maybe 50 tasks in memory. We used to call this "a variable."

- **@dnd-kit for drag-and-drop** — *sighs* Look, HTML5 drag and drop has been in browsers since 2011. It's ugly to implement, I'll grant you that. But @dnd-kit is a non-trivial dependency for one feature (reordering cards). If drag-and-drop is core to your app, fine. If it's a nice-to-have, maybe start without it and add it later. The list view doesn't even need it.

- **Vercel** — For a static site with localStorage persistence. You could host this on GitHub Pages. You could host this on literally any static file server. `python -m http.server` would technically work. Vercel is fine, it's free for this use case, but let's not pretend there's a deployment challenge here. This is an `index.html` and some JavaScript.

- **localStorage persistence** — Absolutely correct for v1. I've seen too many todo apps that require you to set up a Postgres database, a Redis cache, and a Kubernetes cluster before you can write "buy milk." localStorage has been reliable since IE8. It's the right call.

### The Plan Has 6 Phases

Six. Phases. For a todo app.

Back in my day, the phases were:

1. **Build it.** 
2. **Ship it.**

I'm being a little unfair — the phased breakdown makes sense if you're managing scope. But don't let the plan become bigger than the project. I've seen teams spend more time planning their todo app than it would take to build three of them.

## What's Actually New Here vs. Old Wine

Let me be fair about what's genuinely improved since the old days:

- **Vite's dev experience** is legitimately great. Hot module replacement that actually works. We used to refresh the browser manually and lose all our state. Dark times.
- **CSS custom properties for theming** — This is better than the SASS variable nightmare of 2015. One set of tokens, everything updates. Good.
- **Markdown in a browser** — We used to need server-side rendering for this. Now `marked` or `markdown-it` runs client-side in milliseconds. Genuine progress.
- **The dark theme trend** — Fine, my eyes are grateful. GitHub's dark palette is well-tested. Don't reinvent it.

But kanban boards? Tag systems? Date pickers? Due date badges? CRUD operations on local data? We solved all of this. Repeatedly. For decades.

## What I'd Actually Do

If I were building this — and I've built variations of this more times than I care to admit:

1. **Skip Zustand.** Use React's built-in `useReducer` with a context provider. Your state is a list of tasks, a list of tags, and a list of columns. That's three arrays. You don't need a library for three arrays.

2. **Skip the markdown library initially.** Start with a `<textarea>`. Ship it. Add markdown preview in week two when you've confirmed people actually write markdown in their todos. Most people write "buy eggs" and "fix login bug," not dissertations.

3. **Build the list view first.** It's simpler, it validates your data model, and honestly it's what people use 80% of the time. The kanban board is the flashy demo — the list is what gets used daily.

4. **Use a `<input type="date">` for dates.** The browser's native date picker works on every platform, requires zero dependencies, and has been perfectly adequate since 2017. You do not need a date picker library for a todo app.

5. **Keep your CSS simple.** CSS custom properties for your dark theme tokens. No CSS-in-JS library needed. A single CSS file with well-named variables has outlasted every styling framework I've ever seen adopted and abandoned.

6. **localStorage is your database.** Treat it like one. Write a `save()` and `load()` function. Add a version number to your schema. That's your "persistence layer." It doesn't need an "adapter pattern" or "migration versioning strategy" — it needs an `if (version < 2) { migrateV1toV2(data) }`.

## The Edge Cases That Actually Matter

The plan lists 10 edge cases. Most of them are real but let me rank them by "will a user actually hit this":

**Will happen constantly:**
- Tasks with no tags or dates (most tasks)
- Empty columns (every new board)
- Mobile usage (everyone's on their phone)

**Will happen sometimes:**
- Task deleted while detail view is open (close the panel, move on)
- Tag deleted while assigned to tasks (filter it out, done)

**Will probably never happen in a personal todo app:**
- "Large markdown bodies" — in a *todo app*?
- "Multiple tasks with identical titles" — who cares? Let them.
- "Start date after due date" — you cut `startDate`, remember?

## The Real Risk Nobody Mentioned

The biggest risk with this project isn't technical. It's **completion**.

The graveyard of side projects is full of todo apps that got 70% done and were abandoned when the developer got bored after implementing drag-and-drop and before adding the boring stuff like "edit a task." Ship the simplest possible version first. Then iterate.

A shipped todo app with just a list view and checkboxes is infinitely more valuable than an unfinished todo app with a half-working kanban board.

## My Honest Assessment

The MightyVern plan is solid engineering for a well-scoped project. I just think it's *more* engineering than this problem requires. The dependency count can be halved. The data model can be simpler. The phases can be compressed.

Build the list view. Make it work. Make it dark. Ship it. Then add the kanban board. Then add markdown. Then add the fancy date badges. Each step should be a working, deployable application.

You know what survived every JavaScript framework war? The todo list that somebody actually finished and put on the internet.

## Bottom Line

This is a fine project. Good scope for learning React patterns, good portfolio piece, legitimate utility if you actually use it. Just don't let the architecture astronauts convince you that a todo app needs six phases, five date fields, and a state management library. 

The best todo app is the one you finish building.

---

Why did the senior developer refuse to use a kanban board library? Because `mv todo.txt doing.txt && mv doing.txt done.txt` has worked since 1973, and mass adoption never improved on `mv`.

-- Retro Vern *(mass adoption never improved on `mv`)*
