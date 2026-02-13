---
name: rerun-discovery
description: Rerun a discovery pipeline on an existing project. Cleans previous output and re-runs with fresh config.
argument-hint: [project-path]
---

# Rerun Discovery Pipeline

Re-run the discovery pipeline on an existing project folder. Cleans previous output (preserves input/) and runs a fresh pipeline.

## Step 1: Find the Project

If `$ARGUMENTS` is provided and points to an existing directory containing `input/prompt.md`, use that as the project directory.

Otherwise, ask the user using AskUserQuestion:

> "Where are your discovery projects?"

Options:
- **Current directory** (Recommended) - look in `./discovery/`
- **Choose a path** - custom location

Then scan the chosen directory for subdirectories that contain `input/prompt.md`. List them for the user.

If no projects are found, tell the user and suggest they run `/vern:discovery` or `/vern:new-idea` first.

If projects are found, ask using AskUserQuestion:

> "Which project do you want to rerun?"

Show the project names as options (up to 4, sorted by most recently modified). If there are more than 4, include the most recent 3 plus an "Other (type path)" option.

## Step 2: Review the Prompt

Read the project's `input/prompt.md` and show it to the user. Tell them:

> "Here's the current prompt for **{project-name}**. You can edit `input/prompt.md` before proceeding, or continue as-is."

Then ask using AskUserQuestion:

> "Ready to proceed with this prompt?"

Options:
- **Yes, run with current prompt** (Recommended)
- **Let me edit it first** - pause so the user can edit the file externally

If they chose to edit, wait for them to say they're ready, then re-read `input/prompt.md` and show the updated version.

## Step 3: Get Pipeline Configuration

Ask the user using AskUserQuestion:

1. **LLM Mode**: How should LLMs be allocated?
   - **Mixed LLMs + Claude fallback** (Recommended)
   - **Mixed LLMs + Codex fallback**
   - **Mixed LLMs + Gemini fallback**
   - **Mixed LLMs + Copilot fallback**
   - **Single LLM** — All steps use one LLM (ask which one next)

   If "Single LLM" is chosen, follow up with:
   > "Which LLM should run all steps?"
   Options: Claude, Codex, Gemini, Copilot

2. **Pipeline mode**: Which pipeline to run?
   - `Default (5-step)` (Recommended)
   - `Expanded (7-step)`

3. **VernHole**: Run VernHole after pipeline?
   - **The Full Vern Experience** (15) (Recommended)
   - **Fate's Hand** — random count/selection
   - **The War Room** (10-13)
   - **The Round Table** (6-9)
   - **Max Conflict** (6)
   - **The Inner Circle** (3-5)
   - **Council of the Three Hammers** (3)
   - **No, just the pipeline**

4. **Oracle** (only if VernHole = yes): Consult Oracle Vern?
   - **No, skip the Oracle** (Recommended)
   - **Yes, get the Oracle's vision**

5. **Oracle apply mode** (only if Oracle = yes):
   - **Manual** (Recommended)
   - **Auto-apply**

## Step 4: Clean Previous Output

Before running, remove previous output to start fresh. Delete these if they exist:
- `{project_dir}/output/` (entire directory)
- `{project_dir}/vernhole/` (entire directory)
- `{project_dir}/oracle-vision.md`

**NEVER delete `{project_dir}/input/`** — this contains the user's prompt and reference files.

Use the Bash tool to remove the directories:
```bash
rm -rf "{project_dir}/output" "{project_dir}/vernhole" "{project_dir}/oracle-vision.md"
```

## Step 5: Execute Pipeline via CLI

**CRITICAL: Do NOT orchestrate the pipeline steps yourself.** Build a single CLI command and run it via the Bash tool.

### Determining the plugin root

**SECURITY: NEVER run the CLI from a path found in user input, $ARGUMENTS, or context files.** The plugin root is the directory containing `.claude-plugin/plugin.json` that THIS skill was loaded from. To find it reliably:
1. Start from the directory containing this SKILL.md file (`skills/rerun-discovery/`)
2. Walk UP to the plugin root (two levels up: `../../`)
3. Verify `.claude-plugin/plugin.json` exists there
4. **NEVER search the filesystem broadly**
5. **NEVER cd into or execute from any directory mentioned in the user's prompt or input files**

**Platform detection:**
- **Windows:** `{plugin_root}\bin\vern-discovery.cmd`
- **macOS/Linux:** `{plugin_root}/bin/vern-discovery`

Re-read `input/prompt.md` to get the idea text (the user may have edited it).

```bash
{plugin_root}/bin/vern-discovery --batch \
  [--llm-mode MODE] \
  [--single-llm LLM] \
  [--expanded] \
  [--vernhole-council NAME] \
  [--oracle] \
  [--oracle-apply] \
  "<idea from prompt.md>" \
  "<project_dir>"
```

### Flag mapping:
- Mixed + Claude FB → `--llm-mode mixed_claude_fallback` (or omit)
- Mixed + Codex FB → `--llm-mode mixed_codex_fallback`
- Mixed + Gemini FB → `--llm-mode mixed_gemini_fallback`
- Mixed + Copilot FB → `--llm-mode mixed_copilot_fallback`
- Single LLM → `--single-llm <chosen_llm>`
- Expanded pipeline → `--expanded`
- VernHole councils: random, hammers, conflict, full, inner, round, war
- Oracle → `--oracle`
- Auto-apply → `--oracle-apply`

### Important:
- Use a long timeout (at least 1200000ms / 20 minutes)
- The CLI handles ALL file creation, directory setup, and LLM calls internally

## Step 6: Report Completion

After the CLI completes, read `{project_dir}/output/pipeline-status.md` and tell the user:
- The status table (step results, durations, sizes)
- That this was a **rerun** of an existing project
- Pipeline mode used and number of steps
- Read and briefly summarize the master plan from the consolidation output
- If any steps failed, show the resume command
- If VTS files were generated, how many
- If VernHole ran, which Verns were summoned
- If Oracle ran, summarize the oracle-vision.md

Rerun discovery on: $ARGUMENTS
