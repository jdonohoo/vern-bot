#!/bin/bash
#
# Install git hooks for vern-bot development
#
# Run this after cloning the repo:
#   ./install-hooks.sh
#

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
HOOKS_SRC="$REPO_ROOT/hooks"
HOOKS_DST="$REPO_ROOT/.git/hooks"

if [ ! -d "$HOOKS_DST" ]; then
    echo "Error: Not a git repository (no .git/hooks found)"
    exit 1
fi

for hook in "$HOOKS_SRC"/*; do
    hook_name="$(basename "$hook")"
    target="$HOOKS_DST/$hook_name"

    if [ -L "$target" ] && [ "$(readlink "$target")" = "$hook" ]; then
        echo "  $hook_name — already installed"
        continue
    fi

    if [ -e "$target" ] && [ ! -L "$target" ]; then
        echo "  $hook_name — backing up existing hook to ${hook_name}.bak"
        mv "$target" "${target}.bak"
    fi

    ln -sf "$hook" "$target"
    chmod +x "$hook"
    echo "  $hook_name — installed"
done

echo ""
echo "Done. Hooks are symlinked so edits to hooks/ take effect immediately."
