# LLM working area

This directory defines how AI coding agents should work with local context.

Stable documentation can be versioned here.

Local agent scratch files, run logs, transcripts and private context must not be
versioned.

Ignored local directories:

```text
docs/llm/scratch/
docs/llm/runs/
docs/llm/transcripts/
docs/llm/context/
```

Rules:

- Do not store secrets here.
- Do not store API keys here.
- Do not paste private tokens here.
- Do not commit agent transcripts by default.
- Summarize durable decisions into `PROJECT_STATE.md` or `docs/DECISIONS.md`.
- Summarize durable architecture choices into `docs/ARCHITECTURE.md`.
- Keep temporary agent work in ignored subdirectories.
