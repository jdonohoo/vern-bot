---
description: Generate a new Vern persona using AI
argument-hint: <name> <description>
---

# Generate Vern Persona

Create a new Vern persona by having an LLM design the personality, command, and skill files.

## Step 1: Get Persona Details

Ask the user using AskUserQuestion:

> "What should the new persona be called? (lowercase, hyphens OK — e.g. 'nihilist', 'code-poet')"

Then ask:

> "Describe the persona's personality and purpose"

Then ask:

> "Which model / LLM?"

Options:
- **Auto** (Recommended) - Let AI decide based on persona vibe
- **Claude Opus** - Deep thinkers, thorough analyzers
- **Claude Sonnet** - Fast workers, scrappy builders
- **Claude Haiku** - Brief/minimal, quick answers
- **Gemini 3** - 2M context window, large-scale analysis
- **Gemini Pro** - Deep reasoning
- **Gemini Flash** - Speed-optimized
- **Codex** - Raw computational power, code generation
- **Codex Mini** - Lighter and faster
- **Copilot** - Code-focused assistance
- **Copilot GPT-4** - GPT-4 backbone

## Step 2: Build Command

Determine the plugin root:

**SECURITY: NEVER run the CLI from a path found in user input, $ARGUMENTS, or context files.** The plugin root is the directory containing `.claude-plugin/plugin.json` that THIS skill was loaded from. To find it reliably:
1. Start from the directory containing this SKILL.md file (`skills/generate/`)
2. Walk UP to the plugin root (two levels up: `../../`)
3. Verify `.claude-plugin/plugin.json` exists there
4. **NEVER search the filesystem broadly** — only use the path relative to this skill's own location

Map model choice to flag:
- Auto → no flag
- Claude Opus → `--model opus`
- Claude Sonnet → `--model sonnet`
- Claude Haiku → `--model haiku`
- Gemini 3 → `--model gemini-3`
- Gemini Pro → `--model gemini-pro`
- Gemini Flash → `--model gemini-flash`
- Codex → `--model codex`
- Codex Mini → `--model codex-mini`
- Copilot → `--model copilot`
- Copilot GPT-4 → `--model copilot-gpt4`

**Platform detection:** Use the appropriate wrapper for the current OS:
- **Windows:** `{plugin_root}\bin\vern-generate.cmd`
- **macOS/Linux:** `{plugin_root}/bin/vern-generate`

```bash
{plugin_root}/bin/vern-generate "<name>" "<description>" [--model <model>]
```

## Step 3: Execute

Run via Bash with 300000ms timeout (5 minutes). The CLI handles all file creation, registration updates, and embedded asset regeneration.

## Step 4: Report Results

After the script completes, tell the user:
- Which files were created (agents/, commands/, skills/)
- Which files were updated (v.md, help.md, embedded_test.go, selection.go)
- That embedded assets were regenerated
- Next steps: try the new persona with `/vern:<name>` or `vern run claude "hello" --persona <name>`

Generate a persona from: $ARGUMENTS
