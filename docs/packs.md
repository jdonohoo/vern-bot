# Packs: share personas and pipelines

A pack adds personas and discovery pipelines to vern without forking it. It's a folder:

```
my-pack/
├── pack.json            # optional: {"name": "my-pack", "description": "...", "version": "1.0.0"}
├── agents/              # personas, same format as vern's agents/*.md
│   └── editor.md
└── pipelines/           # discovery pipelines
    └── story.json
```

## Where vern looks

In this order. The first pack that has a persona or pipeline name wins:

1. `--pack <dir>` on any command (repeatable)
2. The project's `.vern/` folder (the nearest one above the current directory), then the packs its `config.json` lists
3. `VERN_PATH`: folders separated by `:` (`;` on Windows)
4. `~/.config/vern/packs/<name>/`, one folder per installed pack

vern's own personas come after every pack, so a pack can replace one (a project-specific `mighty`, say). Pack pipelines never replace a pipeline your config already defines: give them their own names.

`vern packs` shows what was found, `vern pipelines` lists the pipelines you can pass to `vern discovery --pipeline <name>`, and `vern models` shows each persona's tier.

## Writing a persona

A persona is a markdown file with frontmatter, then the prompt:

```markdown
---
name: editor
description: Editor Vern - cuts everything that doesn't serve the story. Use for tightening drafts and outlines.
model: sonnet
tier: balanced
color: teal
---

You are Editor Vern. You have read ten thousand drafts and loved maybe forty of them...

SIGN-OFF:
End with a dad joke about red pens.
```

- `name` must match the file name.
- `description` starts with the display name, then " - ", then when to use the persona. The VernHole council shows it.
- `tier` is `deep`, `balanced` or `fast`; it picks the model on whichever LLM runs the persona (see the README's "Models, effort and tiers"). Without it, `model: opus|sonnet|haiku` maps to deep, balanced and fast.
- `model` is also what Claude Code uses if the pack is installed as plugin agents.

Pack personas join the VernHole roster like vern's own.

## Writing a pipeline

```json
{
  "name": "story",
  "description": "Outline, stress-test and break a story idea into writing tasks",
  "kind": "software",
  "steps": [
    {"step": 1, "name": "Premise", "persona": "editor", "llm": "claude", "context_mode": "prompt_only",
     "prompt_prefix": "Develop this premise: stakes, protagonist, the question the story asks."},
    {"step": 2, "name": "Architect Breakdown", "persona": "architect", "llm": "claude", "context_mode": "previous",
     "prompt_prefix": "Break this into tasks with ### TASK N: headers..."}
  ]
}
```

Each step takes the same fields as `discovery_pipelines` in the config: `persona`, `llm`, `context_mode` (`prompt_only`, `previous`, `all_previous`, `consolidation`), `prompt_prefix`, and optionally `model`, `effort` and `tier`. The last step's output is split into tasks, so it should use `### TASK N: Title` headers.

## Project config: `.vern/config.json`

A project's `.vern/` folder is a pack, plus an optional `config.json` that sets defaults for runs inside that project:

```json
{
  "pipeline": "story",
  "sink": "exec:./scripts/file-tasks.sh",
  "packs": ["../shared-vern-pack"],
  "llm_mode": "mixed_claude_fallback",
  "allow_writes": false,
  "tiers": {"deep": {"models": {"codex": "gpt-6-astra"}, "effort": "high"}}
}
```

Every field is optional. `packs` paths are relative to `.vern/`. `tiers` merges over your user config.

## Task sinks

After a discovery run, `--sink` (or `sink` in `.vern/config.json`) sends the tasks on:

| Sink | What it does |
|---|---|
| `vts` | Nothing more: the VTS files in `output/vts/` (the default) |
| `json` or `json:<path>` | `output/tasks.json`, copied to `<path>` when given |
| `beads` | Creates Beads issues with the `br` CLI, like `vern tobeads --apply` |
| `exec:<command>` | Runs `<command>` through the shell with the task graph on stdin |

An `exec` sink is how vern feeds tools it knows nothing about, such as an issue tracker or a team's work board. The command:

- gets `output/tasks.json` ([schema v1](../schemas/tasks.v1.json)) on stdin
- has `VERN_TASKS_FILE`, `VERN_RUN_DIR` and `VERN_MANIFEST_FILE` set
- runs in the current directory, with a 10 minute limit
- has its output saved to `output/sink.log`, and its result recorded in `output/run.json`

`vern sink <discovery-dir> --sink <sink>` resends a finished run's tasks, which is handy while writing a sink.

A minimal sink that files one issue per task with the GitHub CLI:

```sh
#!/bin/sh
# file-tasks.sh: create a GitHub issue for each vern task
jq -c '.tasks[]' | while read -r task; do
  gh issue create \
    --title "$(echo "$task" | jq -r '.id + ": " + .title')" \
    --body "$(echo "$task" | jq -r '.body')"
done
```

## Sharing a pack

Put the folder in a git repo. Others clone it into `~/.config/vern/packs/<name>/` (for everything they run), or into a project and list it under `packs` in `.vern/config.json` (for that project only). Keep personas self-contained: a persona should read well on its own, since vern pastes it into the prompt as is.
