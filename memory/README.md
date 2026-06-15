# Project Memory

This folder is the tracked handoff memory for `range-strategy-lab`.

Use it before nontrivial work so each Codex session starts from the current project state rather than guessing.

## Files

- `PROGRESS.md`: current milestone status, verification, result paths, and next step.
- `DECISIONS.md`: durable constraints and design decisions.
- `NEXT_CODEX_BRIEF.md`: ready prompt for the next implementation session.

## Maintenance Rules

- Keep notes short, dated, and factual.
- Keep always-read memory files compact; `PROGRESS.md` should be a rolling
  snapshot plus milestone index, not a transcript.
- Record commands and result paths instead of copying large generated outputs.
- Keep generated CSV/JSON under `results/`; `results/` is ignored by Git.
- Update `PROGRESS.md` after each completed milestone.
- Update `DECISIONS.md` only when a durable decision changes or is added.
- Remove stale next-step text when a newer milestone supersedes it.
- When `PROGRESS.md` grows past roughly 300 lines, compact older milestones to
  date, artifact/doc paths, row counts or verdict, and commit.
- Keep `NEXT_CODEX_BRIEF.md` focused on one next task and name only the docs
  that task actually needs.
