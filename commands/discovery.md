---
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

1. **LLM Mode**: How should LLMs be allocated?
   - **Mixed LLMs + Claude fallback** (Recommended) — Default pipeline uses codex/claude/gemini/copilot per step config, falls back to claude on failure
   - **Mixed LLMs + Codex fallback** — Same pipeline, falls back to codex instead of claude
   - **Mixed LLMs + Gemini fallback** — Same pipeline, falls back to gemini
   - **Mixed LLMs + Copilot fallback** — Same pipeline, falls back to copilot
   - **Single LLM** — All steps use one LLM (ask which one next)

   If "Single LLM" is chosen, follow up with:
   > "Which LLM should run all steps?"
   Options:
   - **Claude**
   - **Codex**
   - **Gemini**
   - **Copilot**

2. **Pipeline mode**: Which pipeline to run?
   - `Default (5-step)` — Analysis → Refinement → Chaos Check → Consolidation → Architect (Recommended)
   - `Expanded (7-step)` — Adds Reality Check (Vern the Mediocre) and MVP Lens (Startup Vern) before consolidation

3. **VernHole**: Run VernHole on the result after pipeline completes?
   - **No, just the pipeline** (Recommended)
   - **Fate's Hand (random)** — random count, random selection
   - **Council of the Three Hammers** (3) — great + mediocre + ketamine
   - **Max Conflict** (6) — startup, enterprise, yolo, paranoid, optimist, inverse
   - **The Full Vern Experience** (all 15) — every summonable persona speaks

   If the user wants finer control, they can also specify: `inner` (3-5, architect-led), `round` (6-9, round table), or `war` (10-13, war room).

4. **Oracle** (only ask if VernHole = yes): After VernHole, consult Oracle Vern?
   - **No, skip the Oracle** (Recommended)
   - **Yes, get the Oracle's vision**

5. **Oracle apply mode** (only if Oracle = yes):
   - **Manual** — review oracle-vision.md yourself (Recommended)
   - **Auto-apply** — Architect Vern executes the Oracle's vision

## Step 4: Execute Pipeline via CLI

**CRITICAL: Do NOT orchestrate the pipeline steps yourself.** Instead, build a single CLI command and run it via the Bash tool. This ensures the entire pipeline runs non-interactively without permission prompts.

The CLI wrapper is located relative to the plugin root. Find the plugin root by looking for `.claude-plugin/plugin.json`.

**Platform detection:** Use the appropriate wrapper for the current OS:
- **Windows:** `{plugin_root}\bin\vern-discovery.cmd`
- **macOS/Linux:** `{plugin_root}/bin/vern-discovery`

Build the command:

```bash
# macOS/Linux:
{plugin_root}/bin/vern-discovery --batch \
# Windows:
# {plugin_root}\bin\vern-discovery.cmd --batch ^
  [--llm-mode MODE]                       # LLM fallback mode
  [--single-llm LLM]                     # single LLM mode (overrides --llm-mode)
  [--expanded]                            # if user chose expanded pipeline mode
  [--skip-input]                          # if user said no to reading input files
  [--vernhole N]                          # if user wants VernHole with a specific count
  [--vernhole-council NAME]              # if user chose a named council
  [--oracle]                              # if user wants Oracle Vern
  [--oracle-apply]                        # if user wants auto-apply (implies --oracle)
  [--extra-context /path/to/file ...]     # for each extra context file the user provided
  "<idea prompt>" \
  "<discovery_dir>"
```

### Flag mapping:
- User chose **Mixed + Claude FB** → add `--llm-mode mixed_claude_fallback` (or omit, it's the default)
- User chose **Mixed + Codex FB** → add `--llm-mode mixed_codex_fallback`
- User chose **Mixed + Gemini FB** → add `--llm-mode mixed_gemini_fallback`
- User chose **Mixed + Copilot FB** → add `--llm-mode mixed_copilot_fallback`
- User chose **Single LLM** → add `--single-llm <chosen_llm>` (e.g. `--single-llm codex`)
- User chose **Expanded** pipeline → add `--expanded`
- User said **no** to reading input files → add `--skip-input`
- User said **yes** to VernHole:
  - Fate's Hand → add `--vernhole-council random`
  - Council of the Three Hammers → add `--vernhole-council hammers`
  - Max Conflict → add `--vernhole-council conflict`
  - The Full Vern Experience → add `--vernhole-council full`
  - The Inner Circle → add `--vernhole-council inner`
  - The Round Table → add `--vernhole-council round`
  - The War Room → add `--vernhole-council war`
- User said **yes** to Oracle → add `--oracle`
- User said **auto-apply** → add `--oracle-apply` (replaces `--oracle`)
- User provided extra files → add `--extra-context /path/to/file` for each one

### Important:
- Use a long timeout (at least 1200000ms / 20 minutes) for the Bash call — the pipeline spawns multiple LLM subprocesses
- The CLI handles ALL file creation, directory setup, and LLM calls internally
- Each LLM subprocess uses `--dangerously-skip-permissions` so no permission prompts during execution

## Step 5: Report Completion

After the CLI completes, read `{discovery_dir}/output/pipeline-status.md` for a structured overview. Then tell the user:
- The status table from pipeline-status.md (step results, durations, sizes)
- Where files were created
- Pipeline mode used (default or expanded) and number of steps
- Read and briefly summarize the master plan from the consolidation output file
- If any steps failed, show the resume command from the status file
- If VTS files were generated, how many
- If VernHole ran, which Verns were summoned
- If Oracle ran, summarize the oracle-vision.md
- If auto-apply ran, note that VTS files were updated

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
│   ├── pipeline.log                        # Per-step status, timestamps, exit codes
│   ├── pipeline-status.md                  # Structured status (read this first)
│   └── vts/                               # Vern Task Spec files
│       ├── vts-001-{slug}.md
│       └── ...
├── vernhole/                              # Only if user opted in
│   ├── 01-{persona}.md
│   └── synthesis.md
└── oracle-vision.md                       # Only if Oracle ran
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
├── vernhole/
└── oracle-vision.md
```

Begin discovery on: $ARGUMENTS
