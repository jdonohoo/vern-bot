---
description: Create a new discovery idea folder with standardized input/output structure. Use when the user wants to prepare an idea for discovery.
argument-hint: <idea-name>
---

# Vern New Idea

Set up a standardized discovery folder so the user can prepare input materials before running the pipeline.

## Step 1: Get the Idea Name

The idea name comes from `$ARGUMENTS`. If not provided, ask the user for one.

Slugify the name for the directory: lowercase, replace non-alphanumeric with hyphens, collapse multiple hyphens, truncate to 50 chars.

## Step 2: Choose Location

Ask the user using AskUserQuestion:

> "Where should I create the discovery folder?"

Options:
- **Current directory** (Recommended) - creates `./discovery/{name}/`
- **Choose a path** - let them type a custom path

## Step 3: Create the Folder Structure

```
{location}/discovery/{name}/
├── input/
│   └── prompt.md       # Created with placeholder
└── output/             # Empty, pipeline writes here
```

Create `input/prompt.md` with:
```markdown
# Discovery: {name}

## Prompt
<!-- Describe your idea here. This will be fed to the discovery pipeline. -->


## Additional Context
<!-- Add any extra context, constraints, or goals. -->


## Reference Files
<!-- List any files in this input/ folder and what they contain. -->
```

## Step 4: Tell the User What to Do Next

Report:
- Folder created at `{location}/discovery/{name}/`
- They can now:
  1. Edit `input/prompt.md` with their idea description
  2. Drop any reference files into `input/` (specs, diagrams, existing code, etc.)
  3. Run `/vern:discovery {name}` to execute the pipeline
- Or skip prep and just run `/vern:discovery {prompt}` directly

Begin setup for: $ARGUMENTS
