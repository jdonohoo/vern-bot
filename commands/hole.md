---
description: VernHole / K-Hole Vern - Be careful what you wish for. Random Vern personas do discovery on your idea. The more the merrier.
argument-hint: [idea/task]
---

# VernHole (K-Hole Vern)

You've entered the VernHole. There's no going back. Only through.

**WARNING:** You asked for this.

## Step 1: How Deep Do You Want to Go?

Before summoning the council, ask the user using AskUserQuestion:

> "How many Verns do you want to summon? (min 5, the more the merrier)"

Options:
- **All 13** (Recommended) - Summon every Vern. Maximum perspectives, maximum chaos. The more the merrier.
- **7-9 Verns** - Solid chaos. Plenty of contradictions and insights.
- **5-6 Verns** - A focused council. Diverse but manageable.
- **Random** - Let fate decide (5-13)

## Step 2: Output Location

Ask the user using AskUserQuestion:

> "Where should the VernHole output go?"

Options:
- **Current directory** (Recommended) - `./vernhole/`
- **Choose a path** - custom location

## Step 3: Execute via Bash Script

**CRITICAL: Do NOT orchestrate the Vern passes yourself.** Instead, run the `bin/vernhole` bash script in a single Bash tool call. This ensures the entire VernHole runs non-interactively without permission prompts.

The script is located at `bin/vernhole` relative to the plugin root. Find the plugin root by looking for `.claude-plugin/plugin.json`.

```bash
{plugin_root}/bin/vernhole "<idea>" "<output_dir>" "<num_verns>" ["<context_file>"]
```

Arguments:
- **idea**: The user's idea/task from `$ARGUMENTS`
- **output_dir**: The directory from step 2
- **num_verns**: Based on step 1 choice:
  - All 13 → pass `13`
  - 7-9 → pick a random number 7-9
  - 5-6 → pick a random number 5-6
  - Random → leave empty (script picks 5-13)
- **context_file** (optional): If this VernHole is being run on a discovery plan, pass the master plan file path

### Important:
- Use a long timeout (at least 600000ms / 10 minutes) for the Bash call — the script spawns multiple LLM subprocesses (one per Vern plus synthesis)
- The script handles ALL file creation, directory setup, and LLM calls internally
- Each LLM subprocess uses `--dangerously-skip-permissions` so no permission prompts during execution

## The Vern Roster (dynamic)

The roster is built automatically from every persona in `agents/*.md` (excluding `vernhole-orchestrator.md`). As new personas are added, they join the VernHole automatically. The `bin/vernhole` script scans agent files at runtime.

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
