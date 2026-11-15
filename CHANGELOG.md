# Changelog

All notable changes to vern-bot. Versions follow [Semantic Versioning](https://semver.org/).

## Unreleased

### Changed

- **Read-only by default.** Every LLM CLI now runs without write access: Claude with `--permission-mode dontAsk`, Codex with `--sandbox read-only`, Antigravity in plan mode with its terminal sandbox. vern still writes every output file itself. Use `--allow-writes`, `VERN_ALLOW_WRITES=1` or `"allow_writes": true` to restore the old behaviour. Grok has no read-only headless mode yet and keeps write access, with a warning.
- **Personas have tiers** (`deep`, `balanced`, `fast`) instead of vendor models. A tier comes from `tier:` in the agent file or from its Claude model, and config maps each tier to a model per LLM. By default deep personas run Claude on Opus, balanced on Sonnet and fast on Haiku; other CLIs keep their own default model.

### Added

- **Model and effort per step**: `model`, `effort` and `tier` on discovery pipeline steps, and `vern run --model/--effort/--tier`. A step's model goes only to its own LLM, never to a fallback.
- **`vern doctor`**: checks each CLI's install, version, sign-in and whether it still accepts the flags vern passes. `--json` for scripts.
- **`vern models`**: shows the tier mapping, personas by tier, and the models Codex and Antigravity offer.
- Antigravity runs that end with "no output produced" (a tool denied in headless mode) now count as failures, so the step retries or falls back instead of keeping an error as its output.
- The run log records the model, effort and whether the run was read-only.

## 3.0.0

### Breaking

- **Copilot support removed.** Configs, pipeline steps and flags that still say `copilot` fall back to `claude` with a warning. `mixed_copilot_fallback` becomes `mixed_claude_fallback`.
- **Gemini CLI replaced by the Antigravity CLI (`agy`).** Google retired the Gemini CLI in 2026. The provider is now `antigravity`; `gemini` and `g` are aliases that run `agy`. `mixed_gemini_fallback` becomes `mixed_antigravity_fallback`.
- **Go 1.25** is required to build from source.

### Added

- **Grok** provider for xAI's Grok Build CLI (`grok`, aliases `k`, `xai`) and a `mixed_grok_fallback` mode.
- **Provider registry** (`go/internal/llm/provider.go`): each LLM CLI is one entry, and setup, the TUI and the fallback modes read from it.
- **CI on pull requests**: build, vet and tests on Linux, macOS and Windows.
- `LICENSE` (MIT) and this changelog.

### Fixed

- **Large prompts no longer fail.** Claude and Codex receive the prompt on stdin instead of as one command-line argument, which broke above 128 KB on Linux and 32 KB on Windows. Antigravity and Grok take the prompt as an argument, so prompts too big for the command line are written to a temporary file the CLI is told to read.
- Antigravity's own print timeout now matches vern's step timeout.

### Upgrading from 2.x

Nothing to do in most cases: v2 config files (`~/.config/vern/config.json`, `~/.claude/vern-bot-config.json`) are migrated when they load. To make it permanent, run `vern setup` once, which rewrites the file with the new provider names. Install the Antigravity CLI (`agy`) if you used Gemini.

## 2.9.1 and earlier

See the [GitHub releases](https://github.com/jdonohoo/vern-bot/releases).
