# Agent Instructions

This is a standalone offline Go research project for BTCUSDT 5m range-strategy work.

Before nontrivial work:

1. Read `memory/README.md`.
2. Read `memory/PROGRESS.md`.
3. Read `memory/DECISIONS.md`.
4. Read `README.md` and the relevant docs under `docs/`.

Keep project memory current:

- Update `memory/PROGRESS.md` after each completed milestone.
- Update `memory/DECISIONS.md` only for durable decisions or constraints.
- Keep generated CSV/JSON outputs under `results/`; do not paste bulky generated output into memory.
- Record verification commands, result paths, and short factual outcomes.
- Keep notes dated and concise so future sessions can trust them quickly.

Hard boundaries:

- Offline research only.
- No live orders, exchange API keys, deploy scripts, grid, martingale, averaging down, or two-exchange execution.
- BTCUSDT only until the project explicitly changes scope.
- Use confirmed closed-candle decisions only.
- Keep results explainable and reproducible.
