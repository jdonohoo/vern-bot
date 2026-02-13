---
description: Shorthand router for Vern personas. Prefer invoking skills directly (e.g. /vern:mediocre, /vern:hole, /vern:discovery).
argument-hint: <persona> <task>
allowed-tools: [Read, Glob, Grep, Bash, Write, Edit, Task, AskUserQuestion]
---

# Vern-Bot Router

Shorthand router that maps `/vern:v <alias> <task>` to the correct skill. Users can also invoke skills directly — `/vern:mediocre`, `/vern:hole`, `/vern:discovery`, etc.

## Parse Arguments

Parse the first argument to determine which Vern persona to invoke:

### Core Personas
| Alias | Skill | Description |
|-------|-------|-------------|
| `mediocre` / `med` / `m` | `/vern:mediocre` | Sonnet, fast and scrappy |
| `great` / `vernile` / `g` | `/vern:great` | Opus excellence |
| `nyquil` / `nq` / `n` | `/vern:nyquil` | Haiku, brilliant but brief |
| `ketamine` / `ket` / `k` | `/vern:ketamine` | Multi-pass planning |
| `yolo` / `y` | `/vern:yolo` | Gemini chaos mode |
| `mighty` / `codex` / `c` | `/vern:mighty` | Codex power |

### Specialist Personas
| Alias | Skill | Description |
|-------|-------|-------------|
| `academic` / `acad` / `a` | `/vern:academic` | Needs more research, cites sources |
| `architect` / `arch` / `ar` | `/vern:architect` | Systems design, blueprints before builds |
| `enterprise` / `ent` / `e` | `/vern:enterprise` | Needs 6 meetings first |
| `historian` / `his` | `/vern:historian` | The one who actually reads the whole thing |
| `inverse` / `inv` / `i` | `/vern:inverse` | Contrarian takes only |
| `optimist` / `opt` / `o` | `/vern:optimist` | Everything will be fine! |
| `oracle` / `orc` / `ora` | `/vern:oracle` | Reads the council's wisdom, recommends VTS changes |
| `paranoid` / `para` / `p` | `/vern:paranoid` | What could go wrong? Everything. |
| `retro` / `ret` / `r` | `/vern:retro` | We solved this with cron in 2004 |
| `startup` / `su` / `s` | `/vern:startup` | MVP or die trying |
| `ux` / `u` | `/vern:ux` | Cool architecture, but can the user find the button? |

### Workflows & Pipelines
| Alias | Skill | Description |
|-------|-------|-------------|
| `setup` | `/vern:setup` | Configure LLMs, pipeline personas, preferences |
| `new-idea` / `new` / `ni` | `/vern:new-idea` | Create a discovery folder with input/output structure |
| `discovery` / `disco` / `d` | `/vern:discovery` | Full discovery pipeline |
| `rerun` / `rr` | `/vern:rerun-discovery` | Rerun discovery on existing project |
| `generate` / `gen` | `/vern:generate` | Generate a new Vern persona using AI |
| `hole` / `khole` / `vh` | `/vern:hole` | VernHole - random Vern passes |
| `hole-existing` / `vhe` | `/vern:vernhole-existing` | VernHole on existing output |

If no persona specified, default to `/vern:great`.

## Routing

Route the alias to the matching skill and pass the remaining arguments as the task:

- `academic` / `acad` / `a` → invoke `/vern:academic`
- `architect` / `arch` / `ar` → invoke `/vern:architect`
- `discovery` / `disco` / `d` → invoke `/vern:discovery`
- `enterprise` / `ent` / `e` → invoke `/vern:enterprise`
- `generate` / `gen` → invoke `/vern:generate`
- `great` / `vernile` / `g` → invoke `/vern:great`
- `help` / `h` / `?` → invoke `/vern:help`
- `historian` / `his` → invoke `/vern:historian`
- `hole` / `khole` / `vh` → invoke `/vern:hole`
- `hole-existing` / `vhe` → invoke `/vern:vernhole-existing`
- `inverse` / `inv` / `i` → invoke `/vern:inverse`
- `ketamine` / `ket` / `k` → invoke `/vern:ketamine`
- `mediocre` / `med` / `m` → invoke `/vern:mediocre`
- `mighty` / `codex` / `c` → invoke `/vern:mighty`
- `new-idea` / `new` / `ni` → invoke `/vern:new-idea`
- `nyquil` / `nq` / `n` → invoke `/vern:nyquil`
- `optimist` / `opt` / `o` → invoke `/vern:optimist`
- `oracle` / `orc` / `ora` → invoke `/vern:oracle`
- `rerun` / `rr` → invoke `/vern:rerun-discovery`
- `paranoid` / `para` / `p` → invoke `/vern:paranoid`
- `retro` / `ret` / `r` → invoke `/vern:retro`
- `setup` → invoke `/vern:setup`
- `startup` / `su` / `s` → invoke `/vern:startup`
- `ux` / `u` → invoke `/vern:ux`
- `yolo` / `y` → invoke `/vern:yolo`

Begin with: $ARGUMENTS
