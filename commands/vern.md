---
description: Shorthand router for Vern personas. Prefer invoking skills directly (e.g. /vern-mediocre, /vernhole, /vern-discovery).
argument-hint: <persona> <task>
allowed-tools: [Read, Glob, Grep, Bash, Write, Edit, Task, AskUserQuestion]
---

# Vern-Bot Router

Shorthand router that maps `/vern <alias> <task>` to the correct skill. Users can also invoke skills directly — `/vern-mediocre`, `/vernhole`, `/vern-discovery`, etc.

## Parse Arguments

Parse the first argument to determine which Vern persona to invoke:

### Core Personas
| Alias | Skill | Description |
|-------|-------|-------------|
| `mediocre` / `med` / `m` | `/vern-mediocre` | Sonnet, fast and scrappy |
| `great` / `vernile` / `g` | `/vernile-great` | Opus excellence |
| `nyquil` / `nq` / `n` | `/nyquil-vern` | Haiku, brilliant but brief |
| `ketamine` / `ket` / `k` | `/ketamine-vern` | Multi-pass planning |
| `yolo` / `y` | `/yolo-vern` | Gemini chaos mode |
| `mighty` / `codex` / `c` | `/mighty-vern` | Codex power |

### Specialist Personas
| Alias | Skill | Description |
|-------|-------|-------------|
| `architect` / `arch` / `ar` | `/architect-vern` | Systems design, blueprints before builds |
| `inverse` / `inv` / `i` | `/inverse-vern` | Contrarian takes only |
| `paranoid` / `para` / `p` | `/paranoid-vern` | What could go wrong? Everything. |
| `optimist` / `opt` / `o` | `/optimist-vern` | Everything will be fine! |
| `academic` / `acad` / `a` | `/academic-vern` | Needs more research, cites sources |
| `startup` / `su` / `s` | `/startup-vern` | MVP or die trying |
| `enterprise` / `ent` / `e` | `/enterprise-vern` | Needs 6 meetings first |

### Workflows & Pipelines
| Alias | Skill | Description |
|-------|-------|-------------|
| `setup` | `/setup` | Configure LLMs, pipeline personas, preferences |
| `new-idea` / `new` / `ni` | `/vern-new-idea` | Create a discovery folder with input/output structure |
| `discovery` / `disco` / `d` | `/vern-discovery` | Full discovery pipeline |
| `hole` / `khole` / `vh` | `/vernhole` | VernHole - 5-12 random Vern passes |

If no persona specified, default to `/vernile-great`.

## Routing

Route the alias to the matching skill and pass the remaining arguments as the task:

- `mediocre` / `med` / `m` → invoke `/vern-mediocre`
- `great` / `vernile` / `g` → invoke `/vernile-great`
- `nyquil` / `nq` / `n` → invoke `/nyquil-vern`
- `ketamine` / `ket` / `k` → invoke `/ketamine-vern`
- `yolo` / `y` → invoke `/yolo-vern`
- `mighty` / `codex` / `c` → invoke `/mighty-vern`
- `architect` / `arch` / `ar` → invoke `/architect-vern`
- `inverse` / `inv` / `i` → invoke `/inverse-vern`
- `paranoid` / `para` / `p` → invoke `/paranoid-vern`
- `optimist` / `opt` / `o` → invoke `/optimist-vern`
- `academic` / `acad` / `a` → invoke `/academic-vern`
- `startup` / `su` / `s` → invoke `/startup-vern`
- `enterprise` / `ent` / `e` → invoke `/enterprise-vern`
- `hole` / `khole` / `vh` → invoke `/vernhole`
- `new-idea` / `new` / `ni` → invoke `/vern-new-idea`
- `discovery` / `disco` / `d` → invoke `/vern-discovery`
- `setup` → invoke `/setup`

Begin with: $ARGUMENTS
