# Agent Instructions

This is a standalone offline Go research project for BTCUSDT 5m range-strategy work.

Before nontrivial work:

1. Read `memory/README.md`.
2. Read `memory/PROGRESS.md`.
3. Read `memory/DECISIONS.md`.
4. Read `README.md` and only the relevant docs under `docs/`.

Keep project memory current:

- Update `memory/PROGRESS.md` after each completed milestone.
- Update `memory/DECISIONS.md` only for durable decisions or constraints.
- Keep generated CSV/JSON outputs under `results/`; do not paste bulky generated output into memory.
- Record verification commands, result paths, and short factual outcomes.
- Keep notes dated and concise so future sessions can trust them quickly.

Context budget:

- Treat tracked memory as a compact working index, not a full transcript.
- Do not make future sessions read every historical doc by default; use
  `README.md` as an index and open only docs relevant to the current task.
- When `memory/PROGRESS.md` grows past roughly 300 lines, compact older
  milestones to date, artifact/doc paths, row counts or verdict, and commit.
- Put detailed evidence in focused docs, generated artifacts under `results/`,
  or git history instead of expanding always-read memory files.

Hard boundaries:

- Offline research only.
- No live orders, exchange API keys, deploy scripts, grid, martingale, averaging down, or two-exchange execution.
- BTCUSDT only until the project explicitly changes scope.
- Use confirmed closed-candle decisions only.
- Keep results explainable and reproducible.
