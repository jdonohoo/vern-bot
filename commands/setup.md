---
description: Configure vern-bot - LLM availability, discovery pipeline personas, and preferences.
allowed-tools: [Read, Write, Glob, AskUserQuestion]
---

# Vern-Bot Setup

Interactive setup wizard for configuring vern-bot. Run this after installing the plugin, or anytime you want to change your configuration.

## Step 1: Check for Existing Config

Look for config at `~/.claude/vern-bot-config.json`. If it exists, read it and show current settings before asking if they want to reconfigure.

If it exists, ask:
> "You already have a vern-bot config. What would you like to do?"

Options:
- **Reconfigure everything** - start fresh
- **View current config** - just show what's set
- **Keep current config** - exit setup

## Step 2: LLM Availability

Ask using AskUserQuestion:

> "Which LLMs do you have installed? (Claude is assumed)"

Options (multiSelect: true):
- **Codex** (used by MightyVern for raw analysis power)
- **Gemini** (used by YOLO Vern for chaos checks)

Based on their answer, set `llms.codex` and `llms.gemini` to true/false. `llms.claude` is always true.

## Step 3: Configure Discovery Pipeline

If they have ALL three LLMs, ask:

> "Use the default discovery pipeline? (Codex -> Claude -> Gemini -> Codex -> Claude)"

Options:
- **Yes, use defaults** (Recommended) - skip to Step 4
- **No, let me customize** - proceed to custom config

If they're MISSING an LLM, or they chose to customize:

> "Let's configure each pipeline step. Every persona can run on any LLM you have - the persona is the personality, the LLM is the engine."

For each of the 5 pipeline steps, ask using AskUserQuestion:

### Step 1 - Initial Analysis
> "Who should do the initial analysis?"

Options (build dynamically based on available LLMs):
- **MightyVern on Codex** (Recommended if codex available) - raw power analysis
- **MightyVern on Claude** - thorough but different flavor
- **Vernile on Claude** - elegant initial analysis
- **Academic Vern on Claude** - research-heavy initial analysis

### Step 2 - Refinement
> "Who should refine the analysis?"

Options:
- **Vernile on Claude** (Recommended) - elegant refinement
- **Ketamine Vern on Claude** - multi-dimensional refinement
- **Academic Vern on Claude** - evidence-based refinement
- **Enterprise Vern on Claude** - governance-focused refinement

### Step 3 - Chaos Check / Stress Test
> "Who should stress-test the plan?"

Options (build dynamically):
- **YOLO Vern on Gemini** (Recommended if gemini available) - full chaos
- **YOLO Vern on Claude** - chaos energy, Claude engine
- **Inverse Vern on Claude** - contrarian stress test
- **Paranoid Vern on Claude** - risk-focused stress test
- **Startup Vern on Claude** - "is this actually an MVP?" test

### Step 4 - Consolidation
> "Who should consolidate everything into a master plan?"

Options:
- **MightyVern on Codex** (Recommended if codex available) - comprehensive synthesis
- **MightyVern on Claude** - synthesis without Codex
- **Vernile on Claude** - elegant consolidation
- **Ketamine Vern on Claude** - pattern-finding consolidation

### Step 5 - Architect Breakdown
> "Who should break down the plan into actionable items?"

Options:
- **Architect Vern on Claude** (Recommended) - structured breakdown
- **Enterprise Vern on Claude** - governance-heavy breakdown
- **Startup Vern on Claude** - MVP-focused breakdown
- **Vernile on Claude** - elegant task breakdown

## Step 4: VernHole Defaults

Ask using AskUserQuestion:

> "Default number of Verns for VernHole? (can always override per-run)"

Options:
- **Random (5-12)** (Recommended) - let fate decide
- **Small council (5-6)** - focused and manageable
- **Medium chaos (7-9)** - good balance
- **Full VernHole (10-12)** - maximum perspectives

## Step 5: Save Config

Write the config to `~/.claude/vern-bot-config.json`:

```json
{
  "version": "1.0.0",
  "llms": {
    "claude": true,
    "codex": true,
    "gemini": false
  },
  "discovery_pipeline": [
    {
      "step": 1,
      "name": "Initial Analysis",
      "persona": "mighty-vern",
      "llm": "codex",
      "description": "MightyVern on Codex"
    },
    {
      "step": 2,
      "name": "Refinement",
      "persona": "vernile-great",
      "llm": "claude",
      "description": "Vernile on Claude"
    },
    {
      "step": 3,
      "name": "Chaos Check",
      "persona": "yolo-vern",
      "llm": "claude",
      "description": "YOLO Vern on Claude"
    },
    {
      "step": 4,
      "name": "Consolidation",
      "persona": "mighty-vern",
      "llm": "codex",
      "description": "MightyVern on Codex"
    },
    {
      "step": 5,
      "name": "Architect Breakdown",
      "persona": "architect-vern",
      "llm": "claude",
      "description": "Architect Vern on Claude"
    }
  ],
  "vernhole": {
    "default_count": "random",
    "min": 5,
    "max": 12
  }
}
```

## Step 6: Confirm Setup

Show the user a summary:

```
=== VERN-BOT CONFIGURED ===

LLMs: Claude ✓  Codex ✓  Gemini ✗

Discovery Pipeline:
  1. Initial Analysis    → MightyVern on Codex
  2. Refinement          → Vernile on Claude
  3. Chaos Check         → YOLO Vern on Claude
  4. Consolidation       → MightyVern on Codex
  5. Architect Breakdown → Architect Vern on Claude

VernHole: Random (5-12 Verns)

Run /setup anytime to reconfigure.
```

## Notes

- Config is stored at `~/.claude/vern-bot-config.json`
- The bin scripts read this config to determine which LLM to invoke
- Any persona can run on any LLM - the persona prompt is the personality, the LLM is the engine
- If no config exists, the pipeline falls back to defaults (codex/claude/gemini)

Begin setup: $ARGUMENTS
