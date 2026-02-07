---
description: Shorthand router for Vern personas. Prefer invoking skills directly (e.g. /vern-bot:vern-mediocre, /vern-bot:vernhole, /vern-bot:vern-discovery).
argument-hint: <persona> <task>
allowed-tools: [Read, Glob, Grep, Bash, Write, Edit, Task, AskUserQuestion]
---

# Vern-Bot Router

Shorthand router that maps `/vern-bot:vern <alias> <task>` to the correct skill. Users can also invoke skills directly — `/vern-bot:vern-mediocre`, `/vern-bot:vernhole`, `/vern-bot:vern-discovery`, etc.

## Parse Arguments

Parse the first argument to determine which Vern persona to invoke:

### Core Personas
| Alias | Skill | Description |
|-------|-------|-------------|
| `mediocre` / `med` / `m` | `/vern-bot:vern-mediocre` | Sonnet, fast and scrappy |
| `great` / `vernile` / `g` | `/vern-bot:vernile-great` | Opus excellence |
| `nyquil` / `nq` / `n` | `/vern-bot:nyquil-vern` | Haiku, brilliant but brief |
| `ketamine` / `ket` / `k` | `/vern-bot:ketamine-vern` | Multi-pass planning |
| `yolo` / `y` | `/vern-bot:yolo-vern` | Gemini chaos mode |
| `mighty` / `codex` / `c` | `/vern-bot:mighty-vern` | Codex power |

### Specialist Personas
| Alias | Skill | Description |
|-------|-------|-------------|
| `architect` / `arch` / `ar` | `/vern-bot:architect-vern` | Systems design, blueprints before builds |
| `inverse` / `inv` / `i` | `/vern-bot:inverse-vern` | Contrarian takes only |
| `paranoid` / `para` / `p` | `/vern-bot:paranoid-vern` | What could go wrong? Everything. |
| `optimist` / `opt` / `o` | `/vern-bot:optimist-vern` | Everything will be fine! |
| `academic` / `acad` / `a` | `/vern-bot:academic-vern` | Needs more research, cites sources |
| `startup` / `su` / `s` | `/vern-bot:startup-vern` | MVP or die trying |
| `enterprise` / `ent` / `e` | `/vern-bot:enterprise-vern` | Needs 6 meetings first |
| `ux` / `u` | `/vern-bot:ux-vern` | Cool architecture, but can the user find the button? |
| `retro` / `ret` / `r` | `/vern-bot:retro-vern` | We solved this with cron in 2004 |

### Workflows & Pipelines
| Alias | Skill | Description |
|-------|-------|-------------|
| `setup` | `/vern-bot:setup` | Configure LLMs, pipeline personas, preferences |
| `new-idea` / `new` / `ni` | `/vern-bot:vern-new-idea` | Create a discovery folder with input/output structure |
| `discovery` / `disco` / `d` | `/vern-bot:vern-discovery` | Full discovery pipeline |
| `hole` / `khole` / `vh` | `/vern-bot:vernhole` | VernHole - random Vern passes |

If no persona specified, default to `/vern-bot:vernile-great`.

## Routing

Route the alias to the matching skill and pass the remaining arguments as the task:

- `mediocre` / `med` / `m` → invoke `/vern-bot:vern-mediocre`
- `great` / `vernile` / `g` → invoke `/vern-bot:vernile-great`
- `nyquil` / `nq` / `n` → invoke `/vern-bot:nyquil-vern`
- `ketamine` / `ket` / `k` → invoke `/vern-bot:ketamine-vern`
- `yolo` / `y` → invoke `/vern-bot:yolo-vern`
- `mighty` / `codex` / `c` → invoke `/vern-bot:mighty-vern`
- `architect` / `arch` / `ar` → invoke `/vern-bot:architect-vern`
- `inverse` / `inv` / `i` → invoke `/vern-bot:inverse-vern`
- `paranoid` / `para` / `p` → invoke `/vern-bot:paranoid-vern`
- `optimist` / `opt` / `o` → invoke `/vern-bot:optimist-vern`
- `academic` / `acad` / `a` → invoke `/vern-bot:academic-vern`
- `startup` / `su` / `s` → invoke `/vern-bot:startup-vern`
- `enterprise` / `ent` / `e` → invoke `/vern-bot:enterprise-vern`
- `ux` / `u` → invoke `/vern-bot:ux-vern`
- `retro` / `ret` / `r` → invoke `/vern-bot:retro-vern`
- `hole` / `khole` / `vh` → invoke `/vern-bot:vernhole`
- `new-idea` / `new` / `ni` → invoke `/vern-bot:vern-new-idea`
- `discovery` / `disco` / `d` → invoke `/vern-bot:vern-discovery`
- `setup` → invoke `/vern-bot:setup`
- `help` / `h` / `?` → invoke `/vern-bot:help`

Begin with: $ARGUMENTS
