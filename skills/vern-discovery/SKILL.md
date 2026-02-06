---
name: vern-discovery
description: Full Vern Discovery flow - Codex→Claude→Gemini→Consolidate→Architect. Multi-LLM planning pipeline.
argument-hint: [name] [prompt] or just [prompt]
---

# Vern Discovery Pipeline

The ultimate multi-LLM discovery process. Your idea goes through the gauntlet and emerges battle-tested.

## Step 1: Determine Workflow

Parse `$ARGUMENTS` to figure out which path the user is on:

### Path A: Prepared Discovery (existing idea folder)
If the first argument matches an existing `discovery/{name}/` folder (check `./discovery/{name}/input/`):
- Use that folder
- The name is the first argument
- Any remaining arguments are additional prompt context

### Path B: Quick Discovery (no prep)
If no matching folder exists, treat all arguments as the prompt:
1. Ask the user for an **idea name** (for the folder) using AskUserQuestion
2. Ask for **output location** using AskUserQuestion:
   - **Current directory** (Recommended) - `./discovery/{name}/`
   - **Choose a path** - custom location
3. Create the standardized folder structure:
   ```
   {location}/discovery/{name}/
   ├── input/
   │   └── prompt.md    # Save the user's prompt here
   └── output/          # Pipeline writes here
   ```
4. Write the user's prompt to `input/prompt.md`

## Step 2: Gather Input Context

Once the discovery folder is established (either path):

### Read input/ files
Ask the user using AskUserQuestion:

> "I found files in `input/`. Should I read them as context for discovery?"

Options:
- **Yes, read all input files** (Recommended)
- **No, just use the prompt**

If yes: Read every file in `input/` and concatenate their contents as `INPUT_CONTEXT`.

### Additional file paths
Then ask using AskUserQuestion:

> "Do you want to add any other files as context?"

Options:
- **No, that's everything** (Recommended)
- **Yes, I have more files to add**

If yes: Ask the user to provide file paths. Read each one and append to `INPUT_CONTEXT`. After each file, ask again:

> "Add another file?"

Options:
- **No, done adding files**
- **Yes, add another**

This loop prevents clunkiness - one question per file, easy to stop.

### Build the full prompt
Combine into `FULL_PROMPT`:
```
IDEA: {prompt from prompt.md or user arguments}

{If additional prompt context from arguments:}
ADDITIONAL CONTEXT: {extra args}

{If input files were read:}
=== INPUT MATERIALS ===
{contents of each file, labeled with filename}
=== END INPUT MATERIALS ===
```

## Step 3: Get Pipeline Configuration

Ask the user using AskUserQuestion:

1. **Output format**: What deliverable format?
   - `tasks` - Simple task list (Recommended)
   - `issues` - GitHub Issues
   - `tickets` - Jira-style tickets
   - `beads` - Beads format

2. **Claude mode**: Standard or deep?
   - `vernile` - Standard Vernile (Recommended)
   - `ketamine` - Ketamine multi-pass

## The Pipeline

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   CODEX     │────▶│   CLAUDE    │────▶│   GEMINI    │
│  (MightyVern)│     │(Vernile/Ket)│     │ (YOLO Vern) │
│  First Pass │     │ Refinement  │     │ Chaos Check │
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
                    │  SYSTEMS    │
                    │  ARCHITECT  │
                    │ Break Down  │
                    └─────────────┘
```

## Step 4: Execute Pipeline

All output goes to `{discovery_dir}/output/`.

### Pass 1: MightyVern (Codex) - Initial Analysis
```bash
~/.claude/skills/vern-bot/bin/vern-run codex "{FULL_PROMPT}" output/01-codex-analysis.md
```

### Pass 2: Claude - Refinement
Feed Pass 1 output + `FULL_PROMPT` to Claude.
Save to `output/02-claude-refinement.md`

### Pass 3: YOLO Vern (Gemini) - Chaos Check
Feed Pass 2 output + `FULL_PROMPT` to Gemini.
Save to `output/03-gemini-chaos.md`

### Pass 4: Codex - Master Consolidation
Feed ALL previous outputs + `FULL_PROMPT` to Codex.
Save to `output/04-master-plan.md`

### Pass 5: Systems Architect - Breakdown
Feed master plan + `FULL_PROMPT` to Claude.
Save to `output/05-architect-breakdown.md`
Generate individual items in `output/{format}/` folder.

**IMPORTANT:** Every pass receives the `FULL_PROMPT` (original idea + all input context) alongside the chain outputs. This ensures no context is lost as the pipeline progresses.

## Step 5: Report Completion

Tell the user:
- Where files were created
- Summary of what was discovered
- The folder structure for reference

## Step 6 (Optional): VernHole on the Synthesised Plan

Ask the user using AskUserQuestion:

> "Discovery complete! Want to run a VernHole on the synthesised plan?"

Options:
- **No, I'm good** (Recommended)
- **Yes, send it to the VernHole**

If yes:
- Run the `vernhole` script with the **original idea** AND the **master plan** (`output/04-master-plan.md`) as context
- VernHole output goes to a `vernhole/` subfolder alongside `output/`

## Folder Structure (Final)

```
discovery/{name}/
├── input/
│   ├── prompt.md              # The original prompt/idea
│   └── {any reference files}  # User-provided context
├── output/
│   ├── 01-codex-analysis.md
│   ├── 02-claude-refinement.md
│   ├── 03-gemini-chaos.md
│   ├── 04-master-plan.md
│   ├── 05-architect-breakdown.md
│   └── {format}/              # tasks/ or issues/ or tickets/
│       ├── 001-item.md
│       └── ...
└── vernhole/                  # Only if user opted in
    ├── 01-{persona}.md
    └── synthesis.md
```

Begin discovery on: $ARGUMENTS
