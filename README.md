# Vern-Bot

A Claude Code plugin. 15 AI personas with different personalities, models, and approaches to problem-solving. Plus VernHole chaos mode and a full multi-LLM discovery pipeline.

## The 15 Personas

### Core Personas
| Command | Model | Best For |
|---------|-------|----------|
| `/vern:mediocre` | Sonnet | Quick fixes, scripts, "just make it work" |
| `/vern:great` | Opus | Architecture, elegant solutions, quality code |
| `/vern:nyquil` | Haiku | Quick answers, minimal output, brevity |
| `/vern:ketamine` | Opus | Deep exploration, multi-pass planning, patterns |
| `/vern:yolo` | Gemini | Fast execution, no guardrails, chaos |
| `/vern:mighty` | Codex | Comprehensive code gen, thorough analysis |

### Specialist Personas
| Command | Model | Best For |
|---------|-------|----------|
| `/vern:architect` | Opus | Systems design, scalable architecture, production-grade thinking |
| `/vern:inverse` | Sonnet | Devil's advocate, stress-testing assumptions |
| `/vern:paranoid` | Sonnet | Risk assessment, finding failure modes |
| `/vern:optimist` | Haiku | Encouragement, positive framing, can-do energy |
| `/vern:academic` | Opus | Evidence-based analysis, citing prior art |
| `/vern:startup` | Sonnet | MVP thinking, lean approach, cut scope |
| `/vern:enterprise` | Opus | Governance, compliance, process rigor |
| `/vern:ux` | Opus | User experience, empathy-driven design, journey mapping |
| `/vern:retro` | Sonnet | Historical perspective, proven tech, cutting through hype |

### When to Use Each

| Situation | Command |
|-----------|---------|
| "I need this NOW" | `/vern:mediocre` |
| "I need this RIGHT" | `/vern:great` |
| "Just answer the question" | `/vern:nyquil` |
| "Help me think through this" | `/vern:ketamine` |
| "YOLO" | `/vern:yolo` |
| "Give me EVERYTHING" | `/vern:mighty` |
| "Design the system" | `/vern:architect` |
| "Play devil's advocate" | `/vern:inverse` |
| "What could go wrong?" | `/vern:paranoid` |
| "Hype me up" | `/vern:optimist` |
| "What does the research say?" | `/vern:academic` |
| "What's the MVP?" | `/vern:startup` |
| "Enterprise-grade analysis" | `/vern:enterprise` |
| "Can a real person use this?" | `/vern:ux` |
| "Haven't we solved this before?" | `/vern:retro` |
| "I want chaos/creativity" | `/vern:hole` |
| "Full project discovery" | `/vern:discovery` |

## Install

```bash
/plugin marketplace add https://github.com/jdonohoo/vern-bot
/plugin install vern
```

## Usage

```
/vern:<persona> <task>
```

Each persona is a separate skill invoked directly:

```bash
/vern:mediocre <task>
/vern:great <task>
/vern:hole <idea>
/vern:discovery <prompt>
```

Or use the shorthand router:

```bash
/vern:v med <task>
/vern:v great <task>
/vern:v hole <idea>
/vern:v disco <prompt>
```

## Examples

```bash
# Quick and scrappy
/vern:mediocre write a bash script to rename all .txt files to .md

# Quality architecture
/vern:great design a REST API for user authentication

# Systems design
/vern:architect design a scalable notification system

# Deep exploration
/vern:ketamine explore different approaches to state management

# Contrarian check
/vern:inverse challenge this API design

# Risk assessment
/vern:paranoid review this deployment plan

# MVP scope
/vern:startup what's the minimum viable version of this feature

# Full chaos
/vern:hole should we use GraphQL or REST

# Prepared discovery (recommended for complex ideas)
/vern:new-idea my-saas-app
# -> drop reference files into discovery/my-saas-app/input/
# -> edit input/prompt.md with your idea
/vern:discovery my-saas-app

# Quick discovery (skip the prep)
/vern:discovery build a SaaS for freelancers
# -> prompts for a name, creates the folder, runs the pipeline
```

## Workflows & Pipelines

| Command | Description |
|---------|-------------|
| `/vern:new-idea` | Create discovery folder with input/output structure |
| `/vern:discovery` | Full pipeline: Default (5-step) or Expanded (7-step) multi-LLM discovery |
| `/vern:hole` | Random Verns brainstorm your idea |

## The Discovery Pipeline

Two ways to use it:

### Prepared Discovery (for complex ideas)
```bash
/vern:new-idea my-saas-app                         # Creates the folder structure
# Drop specs, designs, existing code into discovery/my-saas-app/input/
# Edit input/prompt.md with your idea description
/vern:discovery my-saas-app "focus on API design"   # Runs pipeline with all that context
```

### Quick Discovery (skip the prep)
```bash
/vern:discovery "AI-powered code review tool"
# Prompts for name + location, creates folders, runs pipeline
```

### What happens when you run discovery

1. Checks for input files, asks if you want them read as context
2. Asks if you want to add any other files (one at a time, easy to stop)
3. Choose pipeline mode: Default (5-step) or Expanded (7-step)
4. Builds the full prompt from your idea + all gathered context
5. Runs the multi-LLM gauntlet:

#### Default Pipeline (5 steps)

```
Codex (Analysis) -> Claude (Refinement) -> Gemini (Chaos Check)
                          |
                          v
                  Codex (Consolidation)
                          |
                          v
                  Architect Vern (Breakdown)
```

#### Expanded Pipeline (7 steps) — use `--expanded` flag

```
Codex (Analysis) -> Claude (Refinement) -> Claude (Reality Check)
                                                  |
                                                  v
                  Gemini (Chaos Check) -> Claude (MVP Lens)
                          |
                          v
                  Codex (Consolidation)
                          |
                          v
                  Architect Vern (Breakdown)
```

The expanded pipeline inserts **Vern the Mediocre** (Reality Check) and **Startup Vern** (MVP Lens) before consolidation for more thorough analysis.

Every pass receives the **original prompt + all input context** alongside the chain outputs, so nothing gets lost.

After discovery completes, you're prompted: *"Do you want to run a VernHole on the synthesised plan?"*

### Standardized Folder Structure

```
discovery/{name}/
├── input/
│   ├── prompt.md              # Your idea description
│   └── {reference files}      # Specs, designs, code, anything
├── output/
│   ├── 01-mighty-initial-analysis.md
│   ├── 02-great-refinement.md
│   ├── 03-yolo-chaos-check.md       # (or 03-mediocre-reality-check.md in expanded)
│   ├── 04-mighty-consolidation.md   # (or 06-... in expanded)
│   ├── 05-architect-architect-breakdown.md  # (or 07-... in expanded)
│   └── tasks/                 # (or issues/ tickets/ beads/)
└── vernhole/                  # Only if you opted in
```

## The VernHole

```bash
/vern:hole "should we use microservices or monolith"
```

You choose how many Verns to summon. Each gives their unique take. Chaos synthesizes into clarity. Maybe.

## Plugin Structure

```
vern-bot/
├── .claude-plugin/
│   ├── plugin.json          # Plugin metadata
│   └── marketplace.json     # Marketplace manifest
├── agents/                   # 16 agent definitions (15 personas + orchestrator)
│   ├── mediocre.md
│   ├── great.md
│   ├── nyquil.md
│   ├── ketamine.md
│   ├── yolo.md
│   ├── mighty.md
│   ├── architect.md
│   ├── inverse.md
│   ├── paranoid.md
│   ├── optimist.md
│   ├── academic.md
│   ├── startup.md
│   ├── enterprise.md
│   ├── ux.md
│   ├── retro.md
│   └── vernhole-orchestrator.md
├── commands/
│   └── v.md                  # /vern:v shorthand router
├── skills/                    # 18 skill definitions
│   ├── mediocre/SKILL.md
│   ├── great/SKILL.md
│   ├── nyquil/SKILL.md
│   ├── ketamine/SKILL.md
│   ├── yolo/SKILL.md
│   ├── mighty/SKILL.md
│   ├── architect/SKILL.md
│   ├── inverse/SKILL.md
│   ├── paranoid/SKILL.md
│   ├── optimist/SKILL.md
│   ├── academic/SKILL.md
│   ├── startup/SKILL.md
│   ├── enterprise/SKILL.md
│   ├── ux/SKILL.md
│   ├── retro/SKILL.md
│   ├── hole/SKILL.md
│   ├── discovery/SKILL.md
│   └── new-idea/SKILL.md
├── bin/                       # Pipeline orchestration scripts
│   ├── vern-run              # Single LLM runner
│   ├── vern-discovery        # Full discovery pipeline
│   └── vernhole              # VernHole chaos mode
├── discovery/                 # Discovery pipeline output
└── README.md
```

## Dad Jokes

Every Vern ends with one. It's the law.

---

*From chaos, clarity. From the VernHole, wisdom. And always, dad jokes.*
