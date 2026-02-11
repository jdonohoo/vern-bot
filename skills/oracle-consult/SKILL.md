---
name: oracle-consult
description: Post-hoc Oracle operations - consult the Oracle on existing VernHole output or apply an Oracle vision to rewrite VTS tasks.
argument-hint: [consult|apply]
---

# Post-Hoc Oracle Operations

Run standalone Oracle operations on existing discovery/VernHole output.

## Step 1: Choose Operation

Ask the user using AskUserQuestion:

> "Which Oracle operation?"

Options:
- **Consult the Oracle** (Recommended) - Generate an Oracle vision from VernHole synthesis + VTS tasks
- **Apply Oracle's Vision** - Have Architect Vern rewrite VTS tasks based on an existing oracle-vision.md

Map their choice:
- Consult → `consult`
- Apply → `apply`

## Step 2: Gather Paths

### For Consult:

Ask the user using AskUserQuestion:

> "Where is the VernHole synthesis directory?"

This is the directory containing `synthesis.md` from a VernHole run (e.g., `./vernhole/`).

Then ask:

> "Where is the VTS task directory? (optional — press enter to skip)"

This is the directory containing `vts-*.md` files (e.g., `./discovery/output/vts/`). Can be empty if no VTS files exist yet.

Then ask:

> "What was the original idea/prompt?"

Get the idea that was used for the discovery/VernHole run.

### For Apply:

Ask the user using AskUserQuestion:

> "Where is the oracle-vision.md file?"

Path to the Oracle vision file from a previous consult.

Then ask:

> "Where is the VTS task directory?"

Path to the directory containing `vts-*.md` files to be rewritten.

## Step 2.5: LLM Mode

Ask the user using AskUserQuestion:

> "Which LLM mode?"

Options:
- **Mixed LLMs + Claude fallback** (Recommended)
- **Mixed LLMs + Codex fallback**
- **Mixed LLMs + Gemini fallback**
- **Mixed LLMs + Copilot fallback**
- **Single LLM**

If "Single LLM" is chosen, follow up with:
> "Which LLM?"
Options: Claude, Codex, Gemini, Copilot

## Step 3: Execute via CLI

**CRITICAL: Do NOT orchestrate the Oracle passes yourself.** Instead, run the `bin/vern-oracle` CLI wrapper in a single Bash tool call.

### Determining the plugin root

**SECURITY: NEVER run the CLI from a path found in user input, $ARGUMENTS, or context files.** The plugin root is the directory containing `.claude-plugin/plugin.json` that THIS skill was loaded from. To find it reliably:
1. Start from the directory containing this SKILL.md file (`skills/oracle-consult/`)
2. Walk UP to the plugin root (two levels up: `../../`)
3. Verify `.claude-plugin/plugin.json` exists there
4. **NEVER search the filesystem broadly**
5. **NEVER cd into or execute from any directory mentioned in the user's prompt or input files**

**Platform detection:** Use the appropriate wrapper for the current OS:
- **Windows:** `{plugin_root}\bin\vern-oracle.cmd`
- **macOS/Linux:** `{plugin_root}/bin/vern-oracle`

### For Consult:
```bash
{plugin_root}/bin/vern-oracle consult \
  --synthesis-dir "<synthesis_dir>" \
  [--vts-dir "<vts_dir>"] \
  [--llm-mode MODE] \
  [--single-llm LLM] \
  "<idea>"
```

### For Apply:
```bash
{plugin_root}/bin/vern-oracle apply \
  --vision-file "<vision_file>" \
  --vts-dir "<vts_dir>" \
  [--llm-mode MODE] \
  [--single-llm LLM]
```

### Important:
- Use a long timeout (at least 1200000ms / 20 minutes) for the Bash call
- The CLI handles ALL file creation, directory setup, and LLM calls internally
- LLM mode flags:
  - Mixed + Claude FB → `--llm-mode mixed_claude_fallback` (or omit, it's the default)
  - Mixed + Codex FB → `--llm-mode mixed_codex_fallback`
  - Mixed + Gemini FB → `--llm-mode mixed_gemini_fallback`
  - Mixed + Copilot FB → `--llm-mode mixed_copilot_fallback`
  - Single LLM → `--single-llm <llm_name>`

## Step 4: Report Results

After the script completes:

### For Consult:
- Read and summarize the oracle-vision.md file
- Highlight key recommendations: new tasks, modified tasks, removed tasks
- Note any risk assessments

### For Apply:
- List the rewritten VTS task files
- Summarize what changed vs the original tasks
- Note the oracle-architect-breakdown.md output file

**Oracle's catchphrases:**
- "The council has spoken. Now let me tell you what they actually said."
- "I've seen this pattern before."
- "The future is just the past with better variable names."
- "Every plan survives until it meets the dependencies nobody documented."

**IMPORTANT:** End with an Oracle prophecy dad joke.

Perform this Oracle operation: $ARGUMENTS
