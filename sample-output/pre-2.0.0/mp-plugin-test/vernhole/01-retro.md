

# Retro Vern's Take: CSV Upload & Display Website

*[leans back in chair, adjusts reading glasses, takes a long sip of coffee]*

Alright, let me tell you something. I've been building things that read CSVs since before half of you were born. And this plan? It's fine. It's *fine*. But it's got that unmistakable smell of a 2024 engineer who's never had to ship something on a deadline with nothing but Notepad and an FTP client.

Let me walk through this.

---

## What the plan gets right

Credit where it's due:

- **No backend.** Thank god. Someone in the room has sense. We don't need a Node server, a Docker container, and a Kubernetes cluster to read a local file. The browser can do this. We knew this in 2005 with ActiveX (okay, bad example), but the `FileReader` API has been solid for over a decade now.
- **Papa Parse.** Yeah, use it. I hand-rolled CSV parsers three times in my career. Every single time I regretted it by row 50,000 when some vendor decided to put a newline *inside* a quoted field. Papa Parse exists because people like me suffered. Use it.
- **No sorting, no filtering, no export.** Correct. The user said "display the content in a table." They didn't say "build me a spreadsheet." Scope creep is how simple projects die.

---

## What the plan overcomplicates

Here's where I start shaking my head.

**"Streaming parse with step mode, render first 500 rows, load-more button, chunk rendering..."**

Hold on. The user said "simple." You know what a simple CSV viewer looks like? It looks like this:

1. User picks a file.
2. You parse it.
3. You put it in a `<table>`.
4. Done.

If the file is 50 rows, which — let's be honest — covers 90% of the CSVs normal humans work with, all that streaming/chunking infrastructure is dead code. You've added complexity for a problem the user hasn't told you they have.

Now, do I think you should handle large files? Sure. But here's the thing: **you don't need streaming to handle large files. You need a row limit and a warning.** Parse the whole thing (Papa Parse is fast — it'll chew through a 10MB file in under a second on any modern machine), then render the first 500 rows. No step mode. No chunk rendering. Just... `slice(0, 500)`. That's it. We were doing this with Perl CGI scripts in 2001.

**"Settings panel: delimiter override, encoding override, header toggle..."**

A settings panel. For a CSV viewer. That the user asked to be "simple."

Look, the header toggle — fine, that's useful. But delimiter override? Encoding override? You're building UI for edge cases that Papa Parse already auto-detects. You're adding controls so the user can manually do what the library does automatically. That's not a feature, that's a confession that you don't trust your tools.

If auto-detection fails, *then* you add overrides. In v2. After someone actually complains.

---

## What I'd actually build

Here's my design. One HTML file. That's it. One file.

```
Structure:
- A heading that says what it does
- A file input (the native <input type="file"> — it's ugly, it works, 
  every human on earth knows what it is)
- A checkbox: "First row is header" (checked by default)
- A <div> where the table goes
- A status line: "Showing X of Y rows" or an error message

Behavior:
- User picks .csv file
- FileReader reads it as text
- Papa Parse parses it (auto-detect delimiter, auto-detect headers 
  based on checkbox)
- If more than 1000 rows, only render 1000, show the count
- Render a plain <table> with <thead> and <tbody>
- If parsing fails, show the error Papa Parse gives you
- That's it. Go home.

Styling:
- A little bit of CSS so the table doesn't look like 1997
- Sticky header (one CSS property: position: sticky)
- Alternating row colors (one CSS selector: tr:nth-child(even))
- max-height on the table container with overflow-y: scroll
- Done.
```

No build step. No npm. No framework. No bundler. No components. Open the HTML file in a browser, it works. Email it to someone, it works. Put it on any web server, it works. Put it on a file share, it works.

The only external dependency is Papa Parse via CDN. One `<script>` tag.

---

## The grizzled veteran's wisdom

Here's what 20+ years of building things has taught me about projects like this:

1. **The hardest part isn't the code, it's the CSV.** CSVs are a lie. There is no standard. Every tool exports them differently. Papa Parse handles this. Trust it and move on.

2. **"Simple" means simple.** When a user says "simple design," they mean they want to understand the whole thing at a glance. They don't want a settings panel. They don't want chunked rendering. They want to see their data in a table.

3. **One file is a feature.** The fact that you can send someone a single HTML file and they can use it — that's not a limitation, that's a deployment strategy. We used to call these "tools." Before everything became a "platform."

4. **You will be tempted to add features.** Resist. The moment you add sorting, someone will want filtering. The moment you add filtering, someone will want export. The moment you add export, congratulations — you've built a worse version of Excel, and Excel already exists.

5. **The "Load More" button is fine but keep it dead simple.** Don't over-engineer it. Parse all the data upfront. Store it in a variable. Render a slice. Button appends more rows. This is a `for` loop, not an architecture.

---

## My verdict on the discovery plan

The plan is 70% good and 30% over-engineered. Strip out the settings panel (except the header checkbox), drop the streaming parse, drop the encoding detection, and you've got a solid little tool that you can build in an hour and maintain forever.

The best code is the code you don't write. And the best features are the ones you ship, not the ones you planned.

*[finishes coffee, mutters something about how we used to do this with a Perl one-liner and a `<pre>` tag]*

---

**Bottom line:** One HTML file. One CDN script tag. A file input, a checkbox, a table, and some CSS. Ship it, move on, build the next thing. That's the design.
