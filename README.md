# Vern-Bot

A Claude Code plugin. 12 AI personas with different personalities, models, and approaches to problem-solving. Plus VernHole chaos mode and a full multi-LLM discovery pipeline.

## Install as Plugin

```bash
claude plugin add /path/to/vern-bot
# or from GitHub:
claude plugin add jdonohoo/vern-bot
```

## Usage

```
/vern <persona> <task>
```

## The 12 Personas

### Core Personas
| Persona | Aliases | Model | Best For |
|---------|---------|-------|----------|
| `mediocre` | `med`, `m` | Sonnet | Quick fixes, scripts, "just make it work" |
| `great` | `vernile`, `g` | Opus | Architecture, elegant solutions, quality code |
| `nyquil` | `nq`, `n` | Haiku | Quick answers, minimal output, brevity |
| `ketamine` | `ket`, `k` | Opus | Deep exploration, multi-pass planning, patterns |
| `yolo` | `y` | Gemini | Fast execution, no guardrails, chaos |
| `mighty` | `codex`, `c` | Codex | Comprehensive code gen, thorough analysis |

### Specialist Personas
| Persona | Aliases | Model | Best For |
|---------|---------|-------|----------|
| `inverse` | `inv`, `i` | Sonnet | Devil's advocate, stress-testing assumptions |
| `paranoid` | `para`, `p` | Sonnet | Risk assessment, finding failure modes |
| `optimist` | `opt`, `o` | Haiku | Encouragement, positive framing, can-do energy |
| `academic` | `acad`, `a` | Opus | Evidence-based analysis, citing prior art |
| `startup` | `su`, `s` | Sonnet | MVP thinking, lean approach, cut scope |
| `enterprise` | `ent`, `e` | Opus | Governance, compliance, process rigor |

### Workflows & Pipelines
| Command | Aliases | Description |
|---------|---------|-------------|
| `new-idea` | `new`, `ni` | Create discovery folder with input/output structure |
| `discovery` | `disco`, `d` | Full pipeline: Codex->Claude->Gemini->Consolidate->Architect |
| `hole` | `khole`, `vh` | 5-12 random Verns brainstorm your idea |

## Examples

```bash
# Quick and scrappy
/vern m write a bash script to rename all .txt files to .md

# Quality architecture
/vern great design a REST API for user authentication

# Deep exploration
/vern ket explore different approaches to state management

# Contrarian check
/vern inverse challenge this API design

# Risk assessment
/vern paranoid review this deployment plan

# MVP scope
/vern startup what's the minimum viable version of this feature

# Full chaos
/vern hole should we use GraphQL or REST

# Prepared discovery (recommended for complex ideas)
/vern new-idea my-saas-app
# -> drop reference files into discovery/my-saas-app/input/
# -> edit input/prompt.md with your idea
/vern discovery my-saas-app

# Quick discovery (skip the prep)
/vern discovery build a SaaS for freelancers
# -> prompts for a name, creates the folder, runs the pipeline
```

## When to Use Each

| Situation | Persona |
|-----------|---------|
| "I need this NOW" | `mediocre` |
| "I need this RIGHT" | `great` |
| "Just answer the question" | `nyquil` |
| "Help me think through this" | `ketamine` |
| "YOLO" | `yolo` |
| "Give me EVERYTHING" | `mighty` |
| "Play devil's advocate" | `inverse` |
| "What could go wrong?" | `paranoid` |
| "Hype me up" | `optimist` |
| "What does the research say?" | `academic` |
| "What's the MVP?" | `startup` |
| "Enterprise-grade analysis" | `enterprise` |
| "I want chaos/creativity" | `hole` |
| "Full project discovery" | `discovery` |

## The Discovery Pipeline

Two ways to use it:

### Prepared Discovery (for complex ideas)
```bash
/vern new-idea my-saas-app                    # Creates the folder structure
# Drop specs, designs, existing code into discovery/my-saas-app/input/
# Edit input/prompt.md with your idea description
/vern discovery my-saas-app "focus on API design"  # Runs pipeline with all that context
```

### Quick Discovery (skip the prep)
```bash
/vern discovery "AI-powered code review tool"
# Prompts for name + location, creates folders, runs pipeline
```

### What happens when you run discovery

1. Checks for input files, asks if you want them read as context
2. Asks if you want to add any other files (one at a time, easy to stop)
3. Builds the full prompt from your idea + all gathered context
4. Runs the multi-LLM gauntlet:

```
Codex (Analysis) -> Claude (Refinement) -> Gemini (Chaos Check)
                          |
                          v
                  Codex (Consolidation)
                          |
                          v
                  Systems Architect (Breakdown)
                          |
                          v
                  [Optional] VernHole on the plan
```

Every pass receives the **original prompt + all input context** alongside the chain outputs, so nothing gets lost.

After discovery completes, you're prompted: *"Do you want to run a VernHole on the synthesised plan?"*

### Standardized Folder Structure

```
discovery/{name}/
├── input/
│   ├── prompt.md              # Your idea description
│   └── {reference files}      # Specs, designs, code, anything
├── output/
│   ├── 01-codex-analysis.md
│   ├── 02-claude-refinement.md
│   ├── 03-gemini-chaos.md
│   ├── 04-master-plan.md
│   ├── 05-architect-breakdown.md
│   └── tasks/                 # (or issues/ tickets/ beads/)
└── vernhole/                  # Only if you opted in
```

## The VernHole

```bash
/vern hole "should we use microservices or monolith"
```

You choose how many Verns to summon (5-12). Each gives their unique take. Chaos synthesizes into clarity. Maybe.

## Plugin Structure

```
vern-bot/
├── .claude-plugin/
│   └── plugin.json          # Plugin metadata
├── agents/                   # 13 agent definitions
│   ├── vern-mediocre.md
│   ├── vernile-great.md
│   ├── nyquil-vern.md
│   ├── ketamine-vern.md
│   ├── yolo-vern.md
│   ├── mighty-vern.md
│   ├── inverse-vern.md
│   ├── paranoid-vern.md
│   ├── optimist-vern.md
│   ├── academic-vern.md
│   ├── startup-vern.md
│   ├── enterprise-vern.md
│   └── vernhole-orchestrator.md
├── commands/
│   └── vern.md               # /vern slash command (router)
├── skills/                    # 14 skill definitions
│   ├── vern-mediocre/SKILL.md
│   ├── vernile-great/SKILL.md
│   ├── nyquil-vern/SKILL.md
│   ├── ketamine-vern/SKILL.md
│   ├── yolo-vern/SKILL.md
│   ├── mighty-vern/SKILL.md
│   ├── inverse-vern/SKILL.md
│   ├── paranoid-vern/SKILL.md
│   ├── optimist-vern/SKILL.md
│   ├── academic-vern/SKILL.md
│   ├── startup-vern/SKILL.md
│   ├── enterprise-vern/SKILL.md
│   ├── vernhole/SKILL.md
│   └── vern-discovery/SKILL.md
├── bin/                       # Pipeline orchestration scripts
│   ├── vern-run              # Single LLM runner
│   ├── vern-discovery        # Full discovery pipeline
│   └── vernhole              # VernHole chaos mode
├── discovery/                 # Discovery pipeline output (plugin marketplace plans)
└── README.md
```

## Dad Jokes

Every Vern ends with one. It's the law.

---

*From chaos, clarity. From the VernHole, wisdom. And always, dad jokes.*
