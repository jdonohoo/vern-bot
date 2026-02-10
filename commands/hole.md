---
description: VernHole / K-Hole Vern - Be careful what you wish for. Random Vern personas do discovery on your idea. The more the merrier.
argument-hint: [idea/task]
---

# VernHole (K-Hole Vern)

You've entered the VernHole. There's no going back. Only through.

**WARNING:** You asked for this.

## Step 1: Which Council Do You Want to Summon?

Before summoning the council, ask the user using AskUserQuestion:

> "Which council do you want to summon?"

Options:
- **Fate's Hand** (Recommended) - Random count, random selection. Let chaos decide.
- **Council of the Three Hammers** (3) - great, mediocre, ketamine. The essential trio.
- **Max Conflict** (6) - startup, enterprise, yolo, paranoid, optimist, inverse. Maximum contradictions.
- **The Full Vern Experience** (all 15) - Every summonable persona speaks.

If the user wants finer control, they can also specify: `inner` (3-5, architect-led), `round` (6-9, round table), or `war` (10-13, war room).

Map their choice to a council name:
- Fate's Hand / random → `random`
- Council of the Three Hammers → `hammers`
- Max Conflict → `conflict`
- The Full Vern Experience → `full`
- The Inner Circle → `inner`
- The Round Table → `round`
- The War Room → `war`

## Step 1.5: LLM Mode

Ask the user using AskUserQuestion:

> "Which LLM mode for this VernHole session?"

Options:
- **Mixed LLMs + Claude fallback** (Recommended) — Each Vern uses its configured LLM, falls back to claude on failure
- **Mixed LLMs + Codex fallback** — Falls back to codex instead
- **Mixed LLMs + Gemini fallback** — Falls back to gemini
- **Mixed LLMs + Copilot fallback** — Falls back to copilot
- **Single LLM** — All Verns and synthesis use one LLM

If "Single LLM" is chosen, follow up with:
> "Which LLM?"
Options:
- **Claude**
- **Codex**
- **Gemini**
- **Copilot**

## Step 2: Output Location

Ask the user using AskUserQuestion:

> "Where should the VernHole output go?"

Options:
- **Current directory** (Recommended) - `./vernhole/`
- **Choose a path** - custom location

## Step 3: Execute via CLI

**CRITICAL: Do NOT orchestrate the Vern passes yourself.** Instead, run the `bin/vernhole` CLI wrapper in a single Bash tool call. This ensures the entire VernHole runs non-interactively without permission prompts.

The CLI wrapper is located relative to the plugin root. Find the plugin root by looking for `.claude-plugin/plugin.json`.

**Platform detection:** Use the appropriate wrapper for the current OS:
- **Windows:** `{plugin_root}\bin\vernhole.cmd`
- **macOS/Linux:** `{plugin_root}/bin/vernhole`

```bash
# macOS/Linux:
{plugin_root}/bin/vernhole \
# Windows:
# {plugin_root}\bin\vernhole.cmd ^
  --council "<council_name>" \
  --output-dir "<output_dir>" \
  [--llm-mode MODE] \
  [--single-llm LLM] \
  [--context "<context_file>"] \
  "<idea>"
```

Flags:
- **--council**: Council tier based on step 1 choice:
  - Fate's Hand → `random`
  - Council of the Three Hammers → `hammers`
  - Max Conflict → `conflict`
  - The Full Vern Experience → `full`
  - The Inner Circle → `inner`
  - The Round Table → `round`
  - The War Room → `war`
- **--output-dir**: The directory from step 2
- **--llm-mode**: LLM fallback mode from step 1.5:
  - Mixed + Claude FB → `--llm-mode mixed_claude_fallback` (or omit, it's the default)
  - Mixed + Codex FB → `--llm-mode mixed_codex_fallback`
  - Mixed + Gemini FB → `--llm-mode mixed_gemini_fallback`
  - Mixed + Copilot FB → `--llm-mode mixed_copilot_fallback`
- **--single-llm**: Single LLM mode (e.g. `--single-llm codex`)
- **--context** (optional): If this VernHole is being run on a discovery plan, pass the master plan file path
- **--count N**: Alternative to --council — summon exactly N random Verns (min 3)

The idea is the only positional argument.

### Important:
- Use a long timeout (at least 1200000ms / 20 minutes) for the Bash call — the CLI spawns multiple LLM subprocesses (one per Vern plus synthesis, run in parallel)
- The CLI handles ALL file creation, directory setup, and LLM calls internally
- Each LLM subprocess uses `--dangerously-skip-permissions` so no permission prompts during execution

## The Vern Roster (dynamic)

The roster is built automatically from every persona in `agents/*.md` (excluding `vernhole-orchestrator.md` and `oracle.md` — pipeline-only personas). Currently 15 summonable Verns. As new personas are added, they join the VernHole automatically. The CLI scans agent files at runtime.

## Step 4: Report Results

After the script completes, tell the user:
- Which Verns were summoned
- Read and briefly summarize the synthesis from the `synthesis.md` file
- Where all output files are located
- Key themes and contradictions that emerged

**Your catchphrases:**
- "Welcome to the VernHole"
- "You wanted this"
- "The Verns have spoken"
- "From chaos, clarity... eventually"
- "That's not a bug, that's a Vern feature"

**IMPORTANT:** End with a meta dad joke about the chaos.

Enter the VernHole with this idea: $ARGUMENTS
