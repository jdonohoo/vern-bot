---
description: Show all available Vern-Bot commands, personas, and workflows.
---

# Vern-Bot Help

Display the following help text to the user:

```
=== VERN-BOT ===
16 AI personas. Multi-LLM discovery pipelines. VernHole chaos mode.

CORE PERSONAS
  /vern:mediocre <task>    Sonnet  - Quick fixes, "just make it work"
  /vern:great <task>       Opus    - Architecture, elegant solutions
  /vern:nyquil <task>      Haiku   - Quick answers, max brevity
  /vern:ketamine <task>    Opus    - Deep exploration, multi-pass planning
  /vern:yolo <task>        Gemini  - No guardrails, full send
  /vern:mighty <task>      Codex   - Comprehensive code gen, raw power

SPECIALIST PERSONAS
  /vern:architect <task>   Opus    - Systems design, scalable architecture
  /vern:inverse <task>     Sonnet  - Devil's advocate, contrarian takes
  /vern:paranoid <task>    Sonnet  - Risk assessment, failure modes
  /vern:optimist <task>    Haiku   - Encouragement, can-do energy
  /vern:academic <task>    Opus    - Evidence-based, cites sources
  /vern:startup <task>     Sonnet  - MVP thinking, lean approach
  /vern:enterprise <task>  Opus    - Governance, compliance, process
  /vern:ux <task>          Opus    - User experience, empathy-driven design
  /vern:retro <task>       Sonnet  - Historical perspective, proven tech
  /vern:oracle <task>      Opus    - Reads the council's wisdom, recommends VTS changes

WORKFLOWS & PIPELINES
  /vern:discovery <prompt> Full multi-LLM discovery pipeline (5 or 7 steps)
  /vern:new-idea <name>    Create discovery folder with input/output structure
  /vern:hole <idea>        Summon random Verns to brainstorm your idea

CONFIGURATION
  /vern:setup              Configure LLMs, pipeline personas, preferences

SHORTHAND ROUTER
  /vern:v <alias> <task>   Route to any persona by short alias

  Aliases: med, great, nq, ket, yolo, mighty, arch, inv, para, opt,
           acad, su, ent, ux, ret, orc, hole, disco, new, setup

EXAMPLES
  /vern:mediocre fix the broken login redirect
  /vern:great design a REST API for auth
  /vern:architect design a scalable notification system
  /vern:hole should we use GraphQL or REST?
  /vern:discovery build a SaaS for freelancers
  /vern:v med write a quick bash script
  /vern:v hole monolith vs microservices?
```

That's it. Just display the help text above. Do not execute any commands.

Begin: $ARGUMENTS
