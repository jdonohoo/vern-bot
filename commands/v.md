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
| `architect` / `arch` / `ar` | `/vern:architect` | Systems design, blueprints before builds |
| `inverse` / `inv` / `i` | `/vern:inverse` | Contrarian takes only |
| `paranoid` / `para` / `p` | `/vern:paranoid` | What could go wrong? Everything. |
| `optimist` / `opt` / `o` | `/vern:optimist` | Everything will be fine! |
| `academic` / `acad` / `a` | `/vern:academic` | Needs more research, cites sources |
| `startup` / `su` / `s` | `/vern:startup` | MVP or die trying |
| `enterprise` / `ent` / `e` | `/vern:enterprise` | Needs 6 meetings first |
| `ux` / `u` | `/vern:ux` | Cool architecture, but can the user find the button? |
| `retro` / `ret` / `r` | `/vern:retro` | We solved this with cron in 2004 |

### Workflows & Pipelines
| Alias | Skill | Description |
|-------|-------|-------------|
| `setup` | `/vern:setup` | Configure LLMs, pipeline personas, preferences |
| `new-idea` / `new` / `ni` | `/vern:new-idea` | Create a discovery folder with input/output structure |
| `discovery` / `disco` / `d` | `/vern:discovery` | Full discovery pipeline |
| `hole` / `khole` / `vh` | `/vern:hole` | VernHole - random Vern passes |

If no persona specified, default to `/vern:great`.

## Routing

Route the alias to the matching skill and pass the remaining arguments as the task:

- `mediocre` / `med` / `m` → invoke `/vern:mediocre`
- `great` / `vernile` / `g` → invoke `/vern:great`
- `nyquil` / `nq` / `n` → invoke `/vern:nyquil`
- `ketamine` / `ket` / `k` → invoke `/vern:ketamine`
- `yolo` / `y` → invoke `/vern:yolo`
- `mighty` / `codex` / `c` → invoke `/vern:mighty`
- `architect` / `arch` / `ar` → invoke `/vern:architect`
- `inverse` / `inv` / `i` → invoke `/vern:inverse`
- `paranoid` / `para` / `p` → invoke `/vern:paranoid`
- `optimist` / `opt` / `o` → invoke `/vern:optimist`
- `academic` / `acad` / `a` → invoke `/vern:academic`
- `startup` / `su` / `s` → invoke `/vern:startup`
- `enterprise` / `ent` / `e` → invoke `/vern:enterprise`
- `ux` / `u` → invoke `/vern:ux`
- `retro` / `ret` / `r` → invoke `/vern:retro`
- `hole` / `khole` / `vh` → invoke `/vern:hole`
- `new-idea` / `new` / `ni` → invoke `/vern:new-idea`
- `discovery` / `disco` / `d` → invoke `/vern:discovery`
- `setup` → invoke `/vern:setup`
- `help` / `h` / `?` → invoke `/vern:help`

Begin with: $ARGUMENTS
