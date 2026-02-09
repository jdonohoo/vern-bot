---
name: discovery
description: Full Vern Discovery flow - multi-LLM planning pipeline with Default (5-step) and Expanded (7-step) modes.
argument-hint: [name] [prompt] or just [prompt]
---

# Vern Discovery Pipeline

The ultimate multi-LLM discovery process. Your idea goes through the gauntlet and emerges battle-tested.

## Step 1: Determine Workflow

Parse `$ARGUMENTS` to figure out which path the user is on:

### Path A: Prepared Discovery (existing idea folder)
If the first argument matches an existing `discovery/{name}/` folder (check `./discovery/{name}/input/`):
- Use that folder as DISCOVERY_DIR
- The name is the first argument
- Any remaining arguments are additional prompt context

### Path B: Quick Discovery (no prep)
If no matching folder exists, treat all arguments as the prompt:
1. Ask the user for an **idea name** (for the folder) using AskUserQuestion
2. Ask for **output location** using AskUserQuestion:
   - **Current directory** (Recommended) - `./discovery/{name}/`
   - **Choose a path** - custom location
3. Write the user's prompt to `{location}/discovery/{name}/input/prompt.md`

## Step 2: Gather Input Context

Once the discovery folder is established (either path):

### Read input/ files
If there are files in `input/`, ask the user using AskUserQuestion:

> "I found files in `input/`. Should I read them as context for discovery?"

Options:
- **Yes, read all input files** (Recommended)
- **No, just use the prompt**

Track this as READ_INPUT (yes or no).

### Additional file paths
Then ask using AskUserQuestion:

> "Do you want to add any other files as context?"

Options:
- **No, that's everything** (Recommended)
- **Yes, I have more files to add**

If yes: Ask the user to provide file paths one at a time. Collect them as EXTRA_FILES list.

## Step 3: Get Pipeline Configuration

Ask the user using AskUserQuestion:

1. **Pipeline mode**: Which pipeline to run?
   - `Default (5-step)` — Analysis → Refinement → Chaos Check → Consolidation → Architect (Recommended)
   - `Expanded (7-step)` — Adds Reality Check (Vern the Mediocre) and MVP Lens (Startup Vern) before consolidation

2. **VernHole**: Run VernHole on the result after pipeline completes?
   - **No, just the pipeline** (Recommended)
   - **Yes, 5-6 Verns** (focused council)
   - **Yes, 7-9 Verns** (getting chaotic)
   - **Yes, 10-12 Verns** (full VernHole)

## Step 4: Execute Pipeline via Bash Script

**CRITICAL: Do NOT orchestrate the pipeline steps yourself.** Instead, build a single bash command and run it via the Bash tool. This ensures the entire pipeline runs non-interactively without permission prompts.

The script is located at `bin/vern-discovery` relative to the plugin root. Find the plugin root by looking for `.claude-plugin/plugin.json`.

Build the command:

```bash
{plugin_root}/bin/vern-discovery --batch \
  [--expanded]                            # if user chose expanded pipeline mode
  [--skip-input]                          # if user said no to reading input files
  [--vernhole N]                          # if user wants VernHole (pick random N in their range)
  [--extra-context /path/to/file ...]     # for each extra context file the user provided
  "<idea prompt>" \
  "<discovery_dir>"
```

### Flag mapping:
- User chose **Expanded** pipeline → add `--expanded`
- User said **no** to reading input files → add `--skip-input`
- User said **yes** to VernHole → add `--vernhole N` where N is a random number in their chosen range (5-6, 7-9, or 10-12)
- User provided extra files → add `--extra-context /path/to/file` for each one

### Important:
- Use a long timeout (at least 600000ms / 10 minutes) for the Bash call — the pipeline spawns multiple LLM subprocesses
- The script handles ALL file creation, directory setup, and LLM calls internally
- Each LLM subprocess uses `--dangerously-skip-permissions` so no permission prompts during execution

## Step 5: Report Completion

After the script completes, tell the user:
- Where files were created
- Pipeline mode used (default or expanded) and number of steps
- Summary of the pipeline output
- Read and briefly summarize the master plan from the consolidation output file
- The folder structure for reference

## The Pipelines

### Default (5-step)

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   CODEX     │────▶│   CLAUDE    │────▶│   GEMINI    │
│  (MightyVern)│     │(Vernile/Ket)│     │ (YOLO Vern) │
│  Analysis   │     │ Refinement  │     │ Chaos Check │
└─────────────┘     └─────────────┘     └─────────────┘
                            │
                            ▼
                    ┌─────────────┐
                    │   CODEX     │
                    │   Master    │
                    │ Consolidate │
                    └─────────────┘
                            │
                            ▼
                    ┌─────────────┐
                    │  ARCHITECT  │
                    │    VERN     │
                    │ Break Down  │
                    └─────────────┘
```

### Expanded (7-step)

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   CODEX     │────▶│   CLAUDE    │────▶│   CLAUDE    │
│  (MightyVern)│     │  (Vernile)  │     │ (Mediocre)  │
│  Analysis   │     │ Refinement  │     │Reality Check│
└─────────────┘     └─────────────┘     └─────────────┘
                                               │
        ┌──────────────────────────────────────┘
        ▼
┌─────────────┐     ┌─────────────┐
│   GEMINI    │────▶│   CLAUDE    │
│ (YOLO Vern) │     │(Startup Vern)│
│ Chaos Check │     │  MVP Lens   │
└─────────────┘     └─────────────┘
                            │
                            ▼
                    ┌─────────────┐
                    │   CODEX     │
                    │   Master    │
                    │ Consolidate │
                    └─────────────┘
                            │
                            ▼
                    ┌─────────────┐
                    │  ARCHITECT  │
                    │    VERN     │
                    │ Break Down  │
                    └─────────────┘
```

**IMPORTANT:** Every pass receives the original idea + all input context alongside the chain outputs, so nothing gets lost.

## Folder Structure (Final)

### Default mode
```
discovery/{name}/
├── input/
│   ├── prompt.md                          # The original prompt/idea
│   └── {any reference files}              # User-provided context
├── output/
│   ├── 01-mighty-initial-analysis.md
│   ├── 02-great-refinement.md
│   ├── 03-yolo-chaos-check.md
│   ├── 04-mighty-consolidation.md
│   ├── 05-architect-architect-breakdown.md
│   └── vts/                               # Vern Task Spec files
│       ├── vts-001-{slug}.md
│       └── ...
└── vernhole/                              # Only if user opted in
    ├── 01-{persona}.md
    └── synthesis.md
```

### Expanded mode
```
discovery/{name}/
├── input/
│   ├── prompt.md
│   └── {any reference files}
├── output/
│   ├── 01-mighty-initial-analysis.md
│   ├── 02-great-refinement.md
│   ├── 03-mediocre-reality-check.md
│   ├── 04-yolo-chaos-check.md
│   ├── 05-startup-mvp-lens.md
│   ├── 06-mighty-consolidation.md
│   ├── 07-architect-architect-breakdown.md
│   └── vts/
└── vernhole/
```

Begin discovery on: $ARGUMENTS
