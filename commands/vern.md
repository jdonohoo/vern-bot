---
description: Summon a Vern persona for your task. 12 personas, VernHole chaos, and full discovery pipeline.
argument-hint: <persona> <task>
allowed-tools: [Read, Glob, Grep, Bash, Write, Edit, Task, AskUserQuestion]
---

# Vern-Bot: Your Chaotic Discovery Companion

## First: Confirm Output Location

Before executing, ask the user:
- **Output directory**: Default to current working directory (`$CWD`)
- Let them specify a different path if desired

## Parse Arguments

Parse the first argument to determine which Vern persona to invoke:

### Core Personas
| Persona | Aliases | Description |
|---------|---------|-------------|
| `mediocre` | `med`, `m` | Vern the Mediocre - Sonnet, fast and scrappy |
| `great` | `vernile`, `g` | Vernile the Great - Opus excellence |
| `nyquil` | `nq`, `n` | Nyquil Vern - Haiku, brilliant but brief |
| `ketamine` | `ket`, `k` | Ketamine Vern - Permission bypass, multi-pass planning |
| `yolo` | `y` | YOLO Vern - Gemini chaos mode |
| `mighty` | `codex`, `c` | MightyVern - Codex power |

### Specialist Personas
| Persona | Aliases | Description |
|---------|---------|-------------|
| `inverse` | `inv`, `i` | Inverse Vern - Contrarian takes only |
| `paranoid` | `para`, `p` | Paranoid Vern - What could go wrong? Everything. |
| `optimist` | `opt`, `o` | Optimist Vern - Everything will be fine! |
| `academic` | `acad`, `a` | Academic Vern - Needs more research, cites sources |
| `startup` | `su`, `s` | Startup Vern - MVP or die trying |
| `enterprise` | `ent`, `e` | Enterprise Vern - Needs 6 meetings first |

### Workflows & Pipelines
| Command | Aliases | Description |
|---------|---------|-------------|
| `setup` | | Configure LLMs, pipeline personas, preferences |
| `new-idea` | `new`, `ni` | Create a discovery folder with input/output structure |
| `discovery` | `disco`, `d` | Full discovery pipeline: Codex->Claude->Gemini->Consolidate |
| `hole` | `khole`, `vh` | VernHole - 5-12 random Vern discovery passes |

If no persona specified, default to `vernile` (Vernile the Great).

### Discovery Workflows

**Prepared discovery** (recommended for complex ideas):
```
/vern new-idea my-saas-app        # Creates discovery/my-saas-app/input/ + output/
# User adds files to input/, edits input/prompt.md
/vern discovery my-saas-app       # Runs pipeline using that folder
```

**Quick discovery** (skip the prep):
```
/vern discovery build a SaaS for freelancers
# Prompts for name + location, creates folder, runs pipeline
```

## Execution

1. Ask user to confirm output location (default: current working directory)
2. Route to the appropriate sub-skill based on persona
3. Execute with the remaining arguments as the task

The sub-skills will create output in:
- `{output_dir}/vern-output/` for single-pass tasks
- `{output_dir}/discovery/` for discovery pipeline
- `{output_dir}/vernhole/` for vernhole passes

## Sub-skill Routing

- `vern-mediocre` for mediocre/med/m
- `vernile-great` for great/vernile/g
- `nyquil-vern` for nyquil/nq/n
- `ketamine-vern` for ketamine/ket/k
- `yolo-vern` for yolo/y
- `mighty-vern` for mighty/codex/c
- `inverse-vern` for inverse/inv/i
- `paranoid-vern` for paranoid/para/p
- `optimist-vern` for optimist/opt/o
- `academic-vern` for academic/acad/a
- `startup-vern` for startup/su/s
- `enterprise-vern` for enterprise/ent/e
- `vernhole` for hole/khole/vh
- `vern-new-idea` for new-idea/new/ni
- `vern-discovery` for discovery/disco/d
- Route `setup` to the `/vern-bot:setup` command

Begin with: $ARGUMENTS
