# Agent Instructions

This is a standalone offline Go research project for BTCUSDT 5m range-strategy work.
The active trading target is Binance USDT-M futures candles, not spot candles.

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
- Record candle source, market type, coverage, and row counts whenever a run or
  verdict depends on a CSV. A spot/futures source change invalidates promotion
  conclusions until a futures impact review is done.

Context budget:

- Treat tracked memory as a compact working index, not a full transcript.
- Do not make future sessions read every historical doc by default; use
  `README.md` as an index and open only docs relevant to the current task.
- Treat always-read memory file size targets as soft judgment bands, not hard
  triggers. Around `300-350` lines can be fine when the extra detail is
  genuinely useful; compact or split memory once an always-read file starts
  feeling bulky or repetitive.
- Put detailed evidence in focused docs, generated artifacts under `results/`,
  or git history instead of expanding always-read memory files.

Hard boundaries:

- Offline research only.
- No live orders, exchange API keys, deploy scripts, grid, martingale, averaging down, or two-exchange execution.
- BTCUSDT only until the project explicitly changes scope.
- Use Binance USDT-M futures 5m data for current research; legacy spot-based
  evidence is historical context only until rerun/reviewed on futures data.
- Use confirmed closed-candle decisions only.
- Keep results explainable and reproducible.
