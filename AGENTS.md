# AGENTS.md

Notes for coding agents (Claude Code, Codex, Antigravity, Grok and others), in two parts: driving `vern` from an agent, and working on this repo.

## Driving vern from an agent

Every command runs without prompts when it gets its arguments, and the commands below print JSON on stdout with `--json` (progress goes to stderr).

| Goal | Command | JSON on stdout |
|---|---|---|
| Plan an idea into tasks | `vern discovery --json "<idea>" [dir]` | The run manifest (`output/run.json`) |
| Get many perspectives | `vern hole --json -d <dir> "<idea>"` | Council files and the synthesis |
| One prompt to one LLM | `vern run --json <llm> "<prompt>"` | Output, LLM, model, exit code, usage |
| Check the setup | `vern doctor --json` | Each CLI's install, version, sign-in, missing flags |
| See models and tiers | `vern models --json` | Tier mapping, personas by tier, models per CLI |
| Render a run for people | `vern report <dir>` | (writes `output/report.html`, prints its path) |

What a discovery run leaves in its folder:

- `output/run.json`: status (`ok`, `partial`, `failed`), every step with its persona, LLM (and fallback), model, effort, duration and usage.
- `output/tasks.json`: the tasks as JSON, [schema v1](schemas/tasks.v1.json). Read this instead of parsing the markdown.
- `output/vts/vts-NNN-*.md`: the same tasks as Vern Task Spec files.
- `output/report.html`: the whole run as one page.

Exit codes:

| Command | 0 | 1 | 3 | 124 |
|---|---|---|---|---|
| `discovery --json` | ok | could not start, or every step failed | partial | |
| `run` | ok | usage error or the CLI failed | | timeout |
| `hole` | ok | every Vern failed | | |
| `doctor` | claude ready | claude missing, signed out or changed | | |

Things to know:

- Runs are read-only: the LLM CLIs cannot edit files or run commands unless you pass `--allow-writes`. vern writes its own output files either way.
- Discovery takes minutes (one LLM call per step, 20 minutes timeout each by default). Run it in the background and read `output/run.json`, which is updated after every step.
- `--single-llm claude` (or another LLM) runs every step on one CLI when others are missing; `vern doctor` tells you which are installed.
- `--resume-from N` reruns a discovery from step N after a failure.

## Working on this repo

- The Go CLI lives in `go/` (Go 1.25+): `cd go && go build ./... && go test ./...`.
- Personas are `agents/*.md`. After changing one, run `cd go && go generate ./internal/embedded/`, which refreshes the embedded copies and the persona text in `docs/index.html`; CI fails if they drift.
- Each LLM CLI is one entry in `go/internal/llm/provider.go`: its flags, read-only mode, model and effort flags, and the checks `vern doctor` runs. Verify flags against the real CLI's `--help` before changing them.
- Tests must pass on Linux, macOS and Windows: build paths with `filepath.Join`, and set `USERPROFILE` alongside `HOME` in tests.
- Keep vern generic: no code for one company's tools or services. Integrations belong behind task sinks and packs.
- Record user-visible changes in `CHANGELOG.md` under "Unreleased".
