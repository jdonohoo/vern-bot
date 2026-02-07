# Vern-Bot

A Claude Code plugin. 15 AI personas with different personalities, models, and approaches to problem-solving. Plus VernHole chaos mode and a full multi-LLM discovery pipeline.

## The 15 Personas

### Core Personas
| Command | Model | Best For |
|---------|-------|----------|
| `/vern-bot:vern-mediocre` | Sonnet | Quick fixes, scripts, "just make it work" |
| `/vern-bot:vernile-great` | Opus | Architecture, elegant solutions, quality code |
| `/vern-bot:nyquil-vern` | Haiku | Quick answers, minimal output, brevity |
| `/vern-bot:ketamine-vern` | Opus | Deep exploration, multi-pass planning, patterns |
| `/vern-bot:yolo-vern` | Gemini | Fast execution, no guardrails, chaos |
| `/vern-bot:mighty-vern` | Codex | Comprehensive code gen, thorough analysis |

### Specialist Personas
| Command | Model | Best For |
|---------|-------|----------|
| `/vern-bot:architect-vern` | Opus | Systems design, scalable architecture, production-grade thinking |
| `/vern-bot:inverse-vern` | Sonnet | Devil's advocate, stress-testing assumptions |
| `/vern-bot:paranoid-vern` | Sonnet | Risk assessment, finding failure modes |
| `/vern-bot:optimist-vern` | Haiku | Encouragement, positive framing, can-do energy |
| `/vern-bot:academic-vern` | Opus | Evidence-based analysis, citing prior art |
| `/vern-bot:startup-vern` | Sonnet | MVP thinking, lean approach, cut scope |
| `/vern-bot:enterprise-vern` | Opus | Governance, compliance, process rigor |
| `/vern-bot:ux-vern` | Opus | User experience, empathy-driven design, journey mapping |
| `/vern-bot:retro-vern` | Sonnet | Historical perspective, proven tech, cutting through hype |

### When to Use Each

| Situation | Command |
|-----------|---------|
| "I need this NOW" | `/vern-bot:vern-mediocre` |
| "I need this RIGHT" | `/vern-bot:vernile-great` |
| "Just answer the question" | `/vern-bot:nyquil-vern` |
| "Help me think through this" | `/vern-bot:ketamine-vern` |
| "YOLO" | `/vern-bot:yolo-vern` |
| "Give me EVERYTHING" | `/vern-bot:mighty-vern` |
| "Design the system" | `/vern-bot:architect-vern` |
| "Play devil's advocate" | `/vern-bot:inverse-vern` |
| "What could go wrong?" | `/vern-bot:paranoid-vern` |
| "Hype me up" | `/vern-bot:optimist-vern` |
| "What does the research say?" | `/vern-bot:academic-vern` |
| "What's the MVP?" | `/vern-bot:startup-vern` |
| "Enterprise-grade analysis" | `/vern-bot:enterprise-vern` |
| "Can a real person use this?" | `/vern-bot:ux-vern` |
| "Haven't we solved this before?" | `/vern-bot:retro-vern` |
| "I want chaos/creativity" | `/vern-bot:vernhole` |
| "Full project discovery" | `/vern-bot:vern-discovery` |

## Install

```bash
/plugin marketplace add https://github.com/jdonohoo/vern-bot
/plugin install vern-bot
```

## Usage

```
/vern-bot:<persona> <task>
```

Each persona is a separate skill invoked directly:

```bash
/vern-bot:vern-mediocre <task>
/vern-bot:vernile-great <task>
/vern-bot:vernhole <idea>
/vern-bot:vern-discovery <prompt>
```

Or use the shorthand router:

```bash
/vern-bot:vern med <task>
/vern-bot:vern great <task>
/vern-bot:vern hole <idea>
/vern-bot:vern disco <prompt>
```

## Examples

```bash
# Quick and scrappy
/vern-bot:vern-mediocre write a bash script to rename all .txt files to .md

# Quality architecture
/vern-bot:vernile-great design a REST API for user authentication

# Systems design
/vern-bot:architect-vern design a scalable notification system

# Deep exploration
/vern-bot:ketamine-vern explore different approaches to state management

# Contrarian check
/vern-bot:inverse-vern challenge this API design

# Risk assessment
/vern-bot:paranoid-vern review this deployment plan

# MVP scope
/vern-bot:startup-vern what's the minimum viable version of this feature

# Full chaos
/vern-bot:vernhole should we use GraphQL or REST

# Prepared discovery (recommended for complex ideas)
/vern-bot:vern-new-idea my-saas-app
# -> drop reference files into discovery/my-saas-app/input/
# -> edit input/prompt.md with your idea
/vern-bot:vern-discovery my-saas-app

# Quick discovery (skip the prep)
/vern-bot:vern-discovery build a SaaS for freelancers
# -> prompts for a name, creates the folder, runs the pipeline
```

## Workflows & Pipelines

| Command | Description |
|---------|-------------|
| `/vern-bot:vern-new-idea` | Create discovery folder with input/output structure |
| `/vern-bot:vern-discovery` | Full pipeline: Default (5-step) or Expanded (7-step) multi-LLM discovery |
| `/vern-bot:vernhole` | Random Verns brainstorm your idea |

## The Discovery Pipeline

Two ways to use it:

### Prepared Discovery (for complex ideas)
```bash
/vern-bot:vern-new-idea my-saas-app                         # Creates the folder structure
# Drop specs, designs, existing code into discovery/my-saas-app/input/
# Edit input/prompt.md with your idea description
/vern-bot:vern-discovery my-saas-app "focus on API design"   # Runs pipeline with all that context
```

### Quick Discovery (skip the prep)
```bash
/vern-bot:vern-discovery "AI-powered code review tool"
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
│   ├── 01-mighty-vern-initial-analysis.md
│   ├── 02-vernile-great-refinement.md
│   ├── 03-yolo-vern-chaos-check.md       # (or 03-vern-mediocre-reality-check.md in expanded)
│   ├── 04-mighty-vern-consolidation.md   # (or 06-... in expanded)
│   ├── 05-architect-vern-architect-breakdown.md  # (or 07-... in expanded)
│   └── tasks/                 # (or issues/ tickets/ beads/)
└── vernhole/                  # Only if you opted in
```

## The VernHole

```bash
/vern-bot:vernhole "should we use microservices or monolith"
```

You choose how many Verns to summon. Each gives their unique take. Chaos synthesizes into clarity. Maybe.

## Plugin Structure

```
vern-bot/
├── .claude-plugin/
│   ├── plugin.json          # Plugin metadata
│   └── marketplace.json     # Marketplace manifest
├── agents/                   # 16 agent definitions (15 personas + orchestrator)
│   ├── vern-mediocre.md
│   ├── vernile-great.md
│   ├── nyquil-vern.md
│   ├── ketamine-vern.md
│   ├── yolo-vern.md
│   ├── mighty-vern.md
│   ├── architect-vern.md
│   ├── inverse-vern.md
│   ├── paranoid-vern.md
│   ├── optimist-vern.md
│   ├── academic-vern.md
│   ├── startup-vern.md
│   ├── enterprise-vern.md
│   ├── ux-vern.md
│   ├── retro-vern.md
│   └── vernhole-orchestrator.md
├── commands/
│   └── vern.md               # /vern-bot:vern shorthand router
├── skills/                    # 18 skill definitions
│   ├── vern-mediocre/SKILL.md
│   ├── vernile-great/SKILL.md
│   ├── nyquil-vern/SKILL.md
│   ├── ketamine-vern/SKILL.md
│   ├── yolo-vern/SKILL.md
│   ├── mighty-vern/SKILL.md
│   ├── architect-vern/SKILL.md
│   ├── inverse-vern/SKILL.md
│   ├── paranoid-vern/SKILL.md
│   ├── optimist-vern/SKILL.md
│   ├── academic-vern/SKILL.md
│   ├── startup-vern/SKILL.md
│   ├── enterprise-vern/SKILL.md
│   ├── ux-vern/SKILL.md
│   ├── retro-vern/SKILL.md
│   ├── vernhole/SKILL.md
│   ├── vern-discovery/SKILL.md
│   └── vern-new-idea/SKILL.md
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
