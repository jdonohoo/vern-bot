<p align="center">
  <img src="images/vern-bot-logo.png" alt="Vern-Bot Logo" width="300" />
</p>

# Vern-Bot

[![CI](https://github.com/jdonohoo/vern-bot/actions/workflows/ci.yml/badge.svg)](https://github.com/jdonohoo/vern-bot/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/jdonohoo/vern-bot)](https://github.com/jdonohoo/vern-bot/releases)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

https://jdonohoo.github.io/vern-bot/

AI discovery pipelines with personality disorders. Multi-LLM orchestration, competing AI personas, VernHole council tiers, Oracle vision, Historian indexing, and chaos-to-clarity workflows.

Runs as a **Claude Code plugin** or a **standalone terminal app** that drives the LLM CLIs you already have: Claude Code, Codex, Antigravity and Grok.

**One idea in, battle-tested plan out.** Give it an idea; several models and personas analyze it, argue about it, and hand you a consolidated plan broken into tasks. See real runs in [`sample-output/`](sample-output/) — for example [todo-web](sample-output/todo-web/): the [step-by-step analysis](sample-output/todo-web/output/), the [VernHole council](sample-output/todo-web/vernhole/), and the [Oracle's vision](sample-output/todo-web/oracle-vision.md).

Upgrading from 2.x? See the [changelog](CHANGELOG.md).

## The 17 Personas

The **Tier** column is each persona's weight class: `deep`, `balanced` or `fast`. Inside Claude Code it is the model in brackets. In the standalone CLI, discovery pipeline steps use the LLM set for each step in the config (MightyVern on Codex, YOLO Vern on Antigravity by default), VernHole council members run on Claude except the Historian, which prefers Antigravity, and the tier picks the model on whichever LLM runs (see [Models, effort and tiers](#models-effort-and-tiers)).

### Core Personas
| Command | Tier | Best For |
|---------|-------|----------|
| `/vern:mediocre` | Balanced (Sonnet) | Quick fixes, scripts, "just make it work" |
| `/vern:great` | Deep (Opus) | Architecture, elegant solutions, quality code |
| `/vern:nyquil` | Fast (Haiku) | Quick answers, minimal output, brevity |
| `/vern:ketamine` | Deep (Opus) | Deep exploration, multi-pass planning, patterns |
| `/vern:yolo` | Balanced (Sonnet) | Fast execution, no guardrails, chaos |
| `/vern:mighty` | Deep (Opus) | Comprehensive code gen, thorough analysis |

### Specialist Personas
| Command | Tier | Best For |
|---------|-------|----------|
| `/vern:architect` | Deep (Opus) | Systems design, scalable architecture, production-grade thinking |
| `/vern:inverse` | Balanced (Sonnet) | Devil's advocate, stress-testing assumptions |
| `/vern:paranoid` | Balanced (Sonnet) | Risk assessment, finding failure modes |
| `/vern:optimist` | Fast (Haiku) | Encouragement, positive framing, can-do energy |
| `/vern:academic` | Deep (Opus) | Evidence-based analysis, citing prior art |
| `/vern:startup` | Balanced (Sonnet) | MVP thinking, lean approach, cut scope |
| `/vern:enterprise` | Deep (Opus) | Governance, compliance, process rigor |
| `/vern:ux` | Deep (Opus) | User experience, empathy-driven design, journey mapping |
| `/vern:retro` | Balanced (Sonnet) | Historical perspective, proven tech, cutting through hype |
| `/vern:oracle` | Deep (Opus) | Reads the council's wisdom, recommends VTS task changes |
| `/vern:historian` | none (prefers Antigravity) | Indexes massive input folders into structured concept maps |

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
| "What did the council actually say?" | `/vern:oracle` |
| "Index these 50 input files" | `/vern:historian` |
| "Create a custom persona" | `/vern:generate` |
| "I want chaos/creativity" | `/vern:hole` |
| "Full project discovery" | `/vern:discovery` |
| "What commands are there?" | `/vern:help` |
| "Configure LLMs/pipeline" | `/vern:setup` |

## Install

### As a Claude Code Plugin

```bash
/plugin marketplace add https://github.com/jdonohoo/vern-bot
/plugin install vern
```

### Standalone (No Claude Code Required)

#### Homebrew (macOS / Linux)

```bash
brew tap jdonohoo/vern
brew install vern
```

To update later: `brew upgrade vern`

#### Scoop (Windows)

```powershell
scoop bucket add vern https://github.com/jdonohoo/homebrew-vern
scoop install vern
```

To update later: `scoop update vern`

#### Winget (Windows)

```powershell
winget install jdonohoo.vern
```

To update later: `winget upgrade jdonohoo.vern`

#### Manual Download

Download the latest binary from [GitHub Releases](https://github.com/jdonohoo/vern-bot/releases):

```bash
# macOS Apple Silicon
curl -Lo vern https://github.com/jdonohoo/vern-bot/releases/latest/download/vern-darwin-arm64
chmod +x vern
sudo mv vern /usr/local/bin/
```

#### Build from Source

```bash
git clone https://github.com/jdonohoo/vern-bot.git
cd vern-bot/go
go build -o bin/vern ./cmd/vern
sudo cp bin/vern /usr/local/bin/
```

#### First Run

Run the setup wizard to detect your LLMs and configure defaults:

```bash
vern setup
```

Then launch the TUI:

```bash
vern tui
```

**Available binaries:** macOS (Intel + Apple Silicon), Linux (x64 + ARM64), Windows (x64 + ARM64).

**Required:** At least one LLM CLI installed — `claude`, `codex`, `agy` ([Antigravity CLI](https://antigravity.google/cli), which replaced the Gemini CLI), or `grok` (xAI Grok Build).

## Standalone TUI

Launch the interactive terminal UI — no Claude Code needed:

```bash
vern tui
```

```
┌──────────────────────────────────┐
│         VERN-BOT v3.0            │
│                                  │
│  [1] Discovery Pipeline          │
│  [2] VernHole Council            │
│  [3] Single LLM Run              │
│  [4] Post-Processing             │
│  [5] Generate Persona            │
│  [6] Historian                   │
│  [7] Settings                    │
│  [q] Quit                        │
│                                  │
│  LLM Mode: Mixed + Claude FB     │
│  LLMs: claude codex antigravity  │
└──────────────────────────────────┘
```

The TUI runs multi-phase discovery with a split-panel layout:

- **Pipeline phase** — step outline (left) with live progress, activity log (right)
- **VernHole phase** — per-Vern roster with async status indicators (green/red/gray), activity log
- **Oracle phase** — consult → apply transitions with Architect working status
- **Results view** — scrollable synthesis + full activity log, press `c` to copy

Settings persist to `~/.config/vern/config.json` — LLM mode, pipeline preferences, default discovery folder, and LLM availability.

Or run commands directly from the CLI:

```bash
vern discovery "build a SaaS for freelancers"
vern hole --council conflict "monolith vs microservices"
vern run codex "analyze this codebase"
vern tobeads --apply ./discovery/my-project/output/vts/
```

### Supported LLM CLIs

| Name | CLI | Aliases | Notes |
|------|-----|---------|-------|
| `claude` | `claude` | `c` | Claude Code. Always the last-resort fallback. |
| `codex` | `codex` | `x` | OpenAI Codex CLI. |
| `antigravity` | `agy` | `a`, `agy`, `gemini`, `g` | Google's Antigravity CLI, which replaced the Gemini CLI. `gemini` in old configs and flags runs `agy`. |
| `grok` | `grok` | `k`, `xai` | xAI's [Grok Build CLI](https://docs.x.ai/build/cli/headless-scripting). Sign in with `grok login` or set `XAI_API_KEY`. |

Copilot support was removed in v3; configs that still name it fall back to claude.

### LLM Modes

Control which LLMs handle your pipeline steps and where failures fall back to:

| Mode | Description |
|------|-------------|
| `mixed_claude_fallback` | Default. Mixed LLMs per step, claude catches failures |
| `mixed_codex_fallback` | Mixed LLMs, codex as safety net |
| `mixed_antigravity_fallback` | Mixed LLMs, antigravity as safety net |
| `mixed_grok_fallback` | Mixed LLMs, grok as safety net |
| `single_llm` | One LLM for everything |

```bash
# Override LLM mode for a single run
vern discovery --llm-mode mixed_codex_fallback "my idea"

# Use a single LLM for everything
vern discovery --single-llm antigravity "my idea"

# Same flags work for VernHole
vern hole --single-llm codex "my idea"
```

### Read-only by default

Every LLM CLI runs read-only by default: vern writes the output files itself, so the models only need to read. Claude runs with `--permission-mode dontAsk` (reads allowed, edits and shell denied), Codex with `--sandbox read-only`, and Antigravity in plan mode with its terminal sandbox (it may run read commands, and writes the plan files it keeps under `~/.gemini/`). Grok has no documented read-only headless mode yet, so it keeps write access and vern warns once per run.

To let the CLIs edit files and run commands, pass `--allow-writes` to any command, set `VERN_ALLOW_WRITES=1`, or set `"allow_writes": true` in your config.

### Models, effort and tiers

Personas declare a tier (`deep`, `balanced`, `fast`; from `tier:` in the agent file, or derived from its Claude model). Your config maps each tier to a model per LLM and an effort; a provider with no model uses its CLI's default. The defaults only set Claude's models:

```json
"tiers": {
  "deep":     {"models": {"claude": "opus", "codex": "gpt-6-astra", "antigravity": "gemini-4-argon-high"}, "effort": "high"},
  "balanced": {"models": {"claude": "sonnet"}},
  "fast":     {"models": {"claude": "haiku"}, "effort": "low"}
}
```

Pipeline steps can override all three. A step's `model` goes only to that step's own LLM, so a fallback to another LLM uses its tier mapping instead:

```json
{"step": 1, "name": "Initial Analysis", "persona": "mighty", "llm": "codex",
 "model": "gpt-6-astra", "effort": "xhigh", "context_mode": "prompt_only", "prompt_prefix": "..."}
```

Effort is one of `low`, `medium`, `high`, `xhigh`, `max`. For single runs:

```bash
vern run claude --model sonnet --effort low "summarize this"
vern run codex --tier deep "review this design"
vern models            # tiers, personas by tier, and the models each CLI offers
vern doctor            # which CLIs are installed, signed in, and still take vern's flags
```

## Usage (Claude Code Plugin)

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

# Oracle vision (after VernHole)
/vern:oracle review the council's synthesis against the VTS tasks

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
| `/vern:discovery` | Full pipeline: Default (5-step) or Expanded (7-step) multi-LLM discovery + council + Oracle |
| `/vern:hole` | Summon a VernHole council tier to brainstorm your idea |
| `/vern:oracle` | Consult Oracle Vern — synthesize council wisdom into VTS modifications |
| `/vern:historian` | Index a directory of input files into a structured concept map |
| `/vern:generate` | Generate a new Vern persona using AI |
| `/vern:oracle-consult` | Run Oracle consult/apply on existing VernHole output |
| `/vern:rerun-discovery` | Rerun a discovery pipeline on an existing project with fresh config |
| `/vern:vernhole-existing` | Point a VernHole council at existing discovery output |
| `/vern:setup` | Configure LLMs, pipeline personas, preferences |
| `/vern:help` | Show all available commands and personas |

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
[Input Files?] → Historian (Antigravity) → input-history.md
                          │
                          ▼
Codex (Analysis) → Claude (Refinement) → Antigravity (Chaos Check)
                          │
                          ▼
                  Codex (Consolidation)
                          │
                          ▼
                  Architect Vern (Breakdown) → VTS Tasks
```

#### Expanded Pipeline (7 steps) — use `--expanded` flag

```
[Input Files?] → Historian (Antigravity) → input-history.md
                          │
                          ▼
Codex (Analysis) → Claude (Refinement) → Claude (Reality Check)
                                                  │
                                                  ▼
                  Antigravity (Chaos Check) → Claude (MVP Lens)
                          │
                          ▼
                  Codex (Consolidation)
                          │
                          ▼
                  Architect Vern (Breakdown) → VTS Tasks
```

The **Historian pre-step** runs automatically when input files exist. It uses Antigravity (large context window) to index all input materials into `input-history.md`, which downstream steps then reference. If Antigravity is unavailable, it falls back to the configured fallback LLM.

The expanded pipeline inserts **Vern the Mediocre** (Reality Check) and **Startup Vern** (MVP Lens) before consolidation for more thorough analysis.

Every pass receives the **original prompt + all input context** alongside the chain outputs, so nothing gets lost.

### Error Handling & Recovery

Pipelines are **failure-tolerant**. A single LLM step failure won't kill your entire run.

| Feature | Description |
|---------|-------------|
| **20-min timeout** | Each step has a 20-minute watchdog. Configurable via `timeout_seconds` in config or `VERN_TIMEOUT` env var. |
| **`--resume-from N`** | Resume a pipeline from step N after a failure. Skips completed steps, preserves context chaining. |
| **`--max-retries N`** | Retry failed steps (default: 2 retries). Failed LLMs automatically fall back based on your LLM mode (e.g. codex/antigravity fall back to claude in `mixed_claude_fallback` mode). |
| **Pipeline log** | `output/pipeline.log` tracks per-step status (OK/FAILED/SKIPPED), timestamps, exit codes, and retry counts. |
| **Pipeline status** | `output/pipeline-status.md` provides a human-readable progress summary with step results table, durations, output sizes, and resume hints. |
| **Failure markers** | Failed steps write a `# STEP FAILED` marker instead of halting. Downstream steps continue. |
| **Downstream guards** | VTS post-processing, VernHole, and Oracle automatically skip when upstream steps fail. |

```bash
# Resume from step 3 after fixing an issue
/vern:discovery --batch --resume-from 3 "my idea" ./discovery/my-project

# Set max retries to 3
/vern:discovery --batch --max-retries 3 "my idea"

# Override timeout (in seconds) via environment
VERN_TIMEOUT=600 bin/vern-run claude "say hello"
```

VernHole is also failure-tolerant — if a Vern fails or times out, it's excluded from synthesis and the remaining Verns carry on.

6. After the pipeline, choose a **VernHole council tier** to brainstorm the plan:

| Tier | Count | Name |
|------|-------|------|
| 1 | 3 | **Council of the Three Hammers** — great, mediocre, ketamine |
| 2 | 6 | **Max Conflict** — startup, enterprise, yolo, paranoid, optimist, inverse |
| 3 | 3–5 | **The Inner Circle** — architect, inverse, paranoid + random fill |
| 4 | 6–9 | **The Round Table** — mighty, yolo, startup, academic, enterprise + random fill |
| 5 | 10–13 | **The War Room** — round table core + ux, retro, optimist, nyquil + random fill |
| 6 | 15 | **The Full Vern Experience** — every summonable persona |
| 7 | ? | **Fate's Hand** — random count, random selection |

7. Optionally consult **Oracle Vern** to review VernHole synthesis against VTS tasks, producing `oracle-vision.md` with recommended modifications
8. Choose to **auto-apply** the Oracle's vision (Architect Vern rewrites VTS) or **manually review**

```
Pipeline → VTS Tasks → VernHole Council → Oracle Vision → Auto-Apply (optional)
```

### Standardized Folder Structure

```
discovery/{name}/
├── input/
│   ├── prompt.md              # Your idea description
│   ├── input-history.md       # Historian's index (auto-generated)
│   └── {reference files}      # Specs, designs, code, anything
├── output/
│   ├── 01-mighty-initial-analysis.md
│   ├── 02-great-refinement.md
│   ├── 03-yolo-chaos-check.md       # (or 03-mediocre-reality-check.md in expanded)
│   ├── 04-mighty-consolidation.md   # (or 06-... in expanded)
│   ├── 05-architect-architect-breakdown.md  # (or 07-... in expanded)
│   ├── vts/                   # Vern Task Spec files
│   ├── pipeline.log           # Per-step status, timestamps, exit codes
│   └── pipeline-status.md     # Human-readable progress summary
├── vernhole/                  # Only if you opted in
└── oracle-vision.md           # Only if Oracle ran
```

## VTS → Beads

Import your VTS tasks into [Beads](https://github.com/steveyegge/beads) by Steve Yegge. The `vern tobeads` command reads VTS files and creates Beads issues via the `br` CLI.

```bash
# Dry-run (default) — see what would be created
vern tobeads ./discovery/my-project/output/vts/

# Actually create issues in Beads
vern tobeads --apply ./discovery/my-project/output/vts/

# Target a different Beads repo + sync after
vern tobeads --apply --sync --beads-dir ~/my-beads-repo ./vts/
```

| Flag | Description |
|------|-------------|
| `--apply` | Create issues (default is dry-run) |
| `--beads-dir` | Target a specific Beads repo |
| `--sync` | Run `br sync --flush-only` after import |

The importer runs preflight validation (duplicate IDs, unknown statuses, missing dependency targets, cycle detection) before touching Beads. Re-run safe — existing issues are detected and skipped.

## The VernHole

<p align="center">
  <img src="images/vernhole.png" alt="The VernHole" width="400" />
</p>

```bash
/vern:hole "should we use microservices or monolith"
```

Choose a council tier — from the focused **Council of the Three Hammers** (3 Verns) to **The Full Vern Experience** (all 15). Each Vern gives their unique take. Chaos synthesizes into clarity. Maybe.

### Council Tiers

| Tier | Name | Count | Personas |
|------|------|-------|----------|
| **Council of the Three Hammers** | Fixed | 3 | great, mediocre, ketamine |
| **Max Conflict** | Fixed | 6 | startup, enterprise, yolo, paranoid, optimist, inverse |
| **The Inner Circle** | Core + random | 3–5 | architect, inverse, paranoid + random fill |
| **The Round Table** | Core + random | 6–9 | mighty, yolo, startup, academic, enterprise + random fill |
| **The War Room** | Core + random | 10–13 | round table + ux, retro, optimist, nyquil + random fill |
| **The Full Vern Experience** | All | 15 | every summonable persona |
| **Fate's Hand** | Random | 3–15 | random count, random selection |

Oracle and Historian are excluded from all council rosters (15 summonable personas).

## Requirements

**As a plugin:** No additional dependencies. The CLI binary auto-downloads on first use.

**Standalone:** At least one LLM CLI installed — `claude`, `codex`, `agy` (Antigravity) or `grok` (Grok Build). Supports 4 LLMs with configurable fallback chains.

Cross-platform: macOS (Intel + Apple Silicon), Linux (x64 + ARM64), Windows (x64 + ARM64).

## Plugin Structure

```
vern-bot/
├── .claude-plugin/
│   ├── plugin.json          # Plugin metadata
│   └── marketplace.json     # Marketplace manifest
├── agents/                   # 18 agent definitions (17 personas + the VernHole orchestrator)
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
│   ├── oracle.md
│   ├── historian.md
│   └── vernhole-orchestrator.md
├── commands/                  # 24 command definitions (17 personas + 7 workflows/utility)
│   ├── v.md                  # /vern:v shorthand router
│   ├── help.md               # /vern:help command reference
│   ├── setup.md              # /vern:setup configuration
│   └── {persona}.md          # Per-persona command files
├── skills/                    # 24 skill definitions
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
│   ├── oracle/SKILL.md
│   ├── oracle-consult/SKILL.md
│   ├── historian/SKILL.md
│   ├── generate/SKILL.md
│   ├── hole/SKILL.md
│   ├── discovery/SKILL.md
│   ├── rerun-discovery/SKILL.md
│   ├── vernhole-existing/SKILL.md
│   └── new-idea/SKILL.md
├── go/                        # Compiled CLI (vern run, discovery, hole, tobeads, historian, generate, oracle, tui, setup)
│   ├── cmd/vern/             # Cobra CLI entry points
│   ├── internal/             # Config, LLM runner, VTS, pipeline, council, TUI, generate
│   ├── go.mod
│   └── go.sum
├── bin/                       # Shell wrappers (delegate to Go CLI binary)
│   ├── vern-run              # Single LLM runner wrapper (bash)
│   ├── vern-run.cmd          # Single LLM runner wrapper (Windows)
│   ├── vern-discovery        # Full discovery pipeline wrapper (bash)
│   ├── vern-discovery.cmd    # Full discovery pipeline wrapper (Windows)
│   ├── vernhole              # VernHole chaos mode wrapper (bash)
│   ├── vernhole.cmd          # VernHole chaos mode wrapper (Windows)
│   ├── vern-historian        # Historian indexing wrapper (bash)
│   ├── install-vern-cli      # Auto-download CLI binary from GitHub releases (bash)
│   └── install-vern-cli.cmd  # Auto-download CLI binary from GitHub releases (Windows)
├── sample-output/             # Real discovery runs, kept as examples
├── discovery/                 # Discovery pipeline output (your runs)
└── README.md
```

## Development

### Prerequisites
- Go 1.25+ (`brew install go`)

### Build
```bash
cd go && go build -o bin/vern ./cmd/vern
```

### Test
```bash
cd go && go test ./...
```

### Subcommands
```bash
vern run <llm> <prompt>              # Single LLM run (--model, --effort, --tier)
vern doctor                           # Check CLIs, sign-ins and the flags vern uses
vern models [llm]                     # Persona tiers and the models each CLI offers
vern discovery <prompt>               # Full discovery pipeline
vern hole <idea>                      # VernHole council
vern tobeads <vts-dir>               # Import VTS tasks into Beads
vern historian <directory>            # Index a directory into a concept map
vern generate <name> <description>   # Generate a new Vern persona using AI
vern oracle consult <idea>            # Generate Oracle vision from VernHole output
vern oracle apply                     # Apply Oracle vision to rewrite VTS tasks
vern tui                              # Interactive terminal UI
vern setup                            # First-run configuration wizard
```

### Release
```bash
git tag v3.x.y && git push --tags
# GitHub Actions builds binaries automatically via GoReleaser
```

## Contributing

Issues and pull requests are welcome. Run `cd go && go vet ./... && go test ./...` before opening a PR; CI runs the same checks on Linux, macOS and Windows. If you change anything in `agents/` or `config.default.json`, run `go generate ./internal/embedded/` so the compiled-in copies stay in sync (CI checks this). Note user-facing changes in [CHANGELOG.md](CHANGELOG.md).

## License

[MIT](LICENSE)

## Dad Jokes

Every Vern ends with one. It's the law.

---

*From chaos, clarity. From the VernHole, wisdom. And always, dad jokes.*
