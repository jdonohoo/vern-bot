---
name: vernhole
description: VernHole / K-Hole Vern - Be careful what you wish for. 5-12 random Vern personas do discovery on your idea.
argument-hint: [idea/task]
---

# VernHole (K-Hole Vern)

You've entered the VernHole. There's no going back. Only through.

**WARNING:** You asked for this.

## Step 1: How Deep Do You Want to Go?

Before summoning the council, ask the user using AskUserQuestion:

> "How many Verns do you want to summon?"

Options:
- **Random (5-12)** (Recommended) - Let fate decide
- **5-6 Verns** - A manageable council. Diverse but focused.
- **7-9 Verns** - Getting chaotic. More contradictions, more insights.
- **10-12 Verns** - Full VernHole. ALL the perspectives. Maximum chaos.

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
  - Random → leave empty (script picks 5-12)
  - 5-6 → pick a random number 5-6
  - 7-9 → pick a random number 7-9
  - 10-12 → pick a random number 10-12
- **context_file** (optional): If this VernHole is being run on a discovery plan, pass the master plan file path

### Important:
- Use a long timeout (at least 600000ms / 10 minutes) for the Bash call — the script spawns multiple LLM subprocesses (one per Vern plus synthesis)
- The script handles ALL file creation, directory setup, and LLM calls internally
- Each LLM subprocess uses `--dangerously-skip-permissions` so no permission prompts during execution

## The Vern Roster (13 possible)

- Vern the Mediocre (scrappy speed demon)
- Vernile the Great (excellence incarnate)
- Nyquil Vern (brilliant brevity)
- Ketamine Vern (multi-dimensional planning)
- YOLO Vern (full send chaos)
- MightyVern (Codex power)
- Architect Vern (systems design, blueprints before builds)
- Inverse Vern (contrarian takes only)
- Paranoid Vern (what could go wrong?)
- Optimist Vern (everything will be fine)
- Academic Vern (needs more research)
- Startup Vern (MVP or die)
- Enterprise Vern (needs 6 meetings first)

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
